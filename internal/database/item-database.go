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

func (db *DB) ItemIDs() ([]string, error) {
	query := "SELECT id FROM item;"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying item records: %v", err)
		return []string{}, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var idStr string
		err := rows.Scan(&idStr)
		if err != nil {
			log.Printf("Error scanning item record: %v", err)
			continue
		}
		ids = append(ids, idStr)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}

	// Print all IDs
	fmt.Println(ids)
	return ids, nil
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

func (db *DB) Items() ([][]string, error) {
	query := "SELECT id, label, description, picture, quantity, weight, qrcode FROM item;"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying user records: %v", err)
		return [][]string{}, err
	}
	defer rows.Close()

	var items [][]string
	for rows.Next() {
		var item Item
		var idStr string
		err := rows.Scan(&idStr, &item.Label, &item.Description, &item.Picture, &item.Quantity, &item.Weight, &item.QRcode)
		if err != nil {
			log.Printf("Error scanning item record: %v", err)
			return [][]string{}, err
		}

		formatted := fmt.Sprintf("id: %s, label: %s, description: %s, picture: %s, quantity: %d, weight: %s, qrcode: %s \n",
			idStr, item.Label, item.Description, item.Picture, item.Quantity, item.Weight, item.QRcode)
		items = append(items, []string{formatted})
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}
	return items, nil
}

// this is dynamic function but not ready
// am not really convinced from repeating the process every time i want to retrieve the data,
func (db *DB) ItemExperement(query string, refs []interface{}) {
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying user records: %v", err)
		return
	}
	defer rows.Close()

	var results [][]interface{}

	for rows.Next() {
		err := rows.Scan(refs...)
		if err != nil {
			log.Printf("Error scanning item: %v", err)
			continue
		}

		// Copy the data from refs to a new slice to store the results
		row := make([]interface{}, len(refs))
		for i, ref := range refs {
			row[i] = *ref.(*interface{})
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}
	fmt.Print(results)
}
