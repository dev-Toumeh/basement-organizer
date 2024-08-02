package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid/v5"
)

type Item struct {
	Id          uuid.UUID `json:"id"`
	Label       string    `json:"label"       validate:"required,lte=128"`
	Description string    `json:"description" validate:"omitempty,lte=256"`
	Picture     string    `json:"picture"     validate:"omitempty,base64"`
	Quantity    int64     `json:"quantity"    validate:"omitempty,numeric,gte=1"`
	Weight      string    `json:"weight"      validate:"omitempty,numeric"`
	QRcode      string    `json:"qrcode"      validate:"omitempty,alphanumunicode"`
}

// Create New Item Record
func (db *DB) CreateNewItem(ctx context.Context, newItem Item) error {

	if _, err := db.ItemByField(ctx, "label", newItem.Label); err != nil {
		return err
	}
	err := db.insertNewItem(ctx, newItem)
	if err != nil {
		return err
	}
	return nil

}

// Get/Check if the Item exist
func (db *DB) ItemByField(ctx context.Context, field string, value string) (Item, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	query := fmt.Sprintf("SELECT id, label, description, picture, quantity, weight, qrcode FROM item WHERE %s = ? \n", field)
	row := db.Sql.QueryRowContext(ctx, query, value)

	var item Item
	var idStr string
	err := row.Scan(&idStr, &item.Label, &item.Description, &item.Picture, &item.Quantity, &item.Weight, &item.QRcode)

	if err != nil {
		if err == sql.ErrNoRows {
			return Item{}, nil
		}
		log.Println("Error while checking if the Item is available:", err)
		return Item{}, err
	}
	item.Id = uuid.Must(uuid.FromString(idStr))
	return item, ErrExist
}

// here we run the insert new Item query separate from the public function
// it make the code more readable
func (db *DB) insertNewItem(ctx context.Context, item Item) error {
	sqlStatement := `INSERT INTO item (id, label, description, picture, quantity, weight, qrcode) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Sql.ExecContext(ctx, sqlStatement, item.Id.String(), item.Label, item.Description, item.Picture, item.Quantity, item.Weight, item.QRcode)
	if err != nil {
		log.Printf("Error while executing create new item statement: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected while executing create new item statement: %v", err)
		return err
	}
	if rowsAffected != 1 {
		log.Println("No rows affected, item not added")
		return errors.New("item not added")
	}
	return nil
}
