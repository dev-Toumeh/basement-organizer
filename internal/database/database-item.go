package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid/v5"
)

// Create New Item Record
func (db *DB) CreateNewItem(ctx context.Context, newItem items.Item) error {

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
func (db *DB) ItemByField(ctx context.Context, field string, value string) (items.Item, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	query := fmt.Sprintf("SELECT id, label, description, picture, quantity, weight, qrcode FROM item WHERE %s = ? \n", field)
	logg.Debug(query)
	row := db.Sql.QueryRowContext(ctx, query, value)

	var item items.Item
	var idStr string
	err := row.Scan(&idStr, &item.Label, &item.Description, &item.Picture, &item.Quantity, &item.Weight, &item.QRcode)

	if err != nil {
		// @TODO: should error not be returned?
		if err == sql.ErrNoRows {
			return items.Item{}, nil
		}
		log.Println("Error while checking if the Item is available:", err)
		return items.Item{}, err
	}
	item.Id = uuid.Must(uuid.FromString(idStr))
	// @TODO: Why is error returned after item is found?
	return item, db.ErrorExist()
}

// Item returns new Item struct if id matches.
func (db *DB) Item(id string) (items.Item, error) {
	ctx := context.Background()
	return db.ItemByField(ctx, "id", id)
}

// return items id's in array from type string
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
	return ids, nil
}

// here we run the insert new Item query separate from the public function
// it make the code more readable
func (db *DB) insertNewItem(ctx context.Context, item items.Item) error {
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

// update the item based on the id
func (db *DB) UpdateItem(ctx context.Context, item items.Item) error {
	sqlStatement := fmt.Sprintf(`UPDATE item Set label = "%s", description = "%s", picture = "%s",
    quantity = "%d", weight = "%s", qrcode = "%s" WHERE id = ?`,
		item.Label, item.Description, item.Picture, item.Quantity, item.Weight, item.QRcode)
	result, err := db.Sql.ExecContext(ctx, sqlStatement, item.Id.String())
	if err != nil {
		logg.Err(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logg.Err(err)
		return err
	}
	if rowsAffected == 0 {
		err := errors.New(fmt.Sprintf("the Record with the id: %s was not found that should not happened while updating", item.Id.String()))
		logg.Debug(err)
		return err
	} else if rowsAffected != 1 {
		err := errors.New(fmt.Sprintf("the id: %s has unexpected effected number of rows (more than one or less than 0)", item.Id.String()))
		logg.Err(err)
		return err
	}
	return nil
}

// Delete Item by Id
func (db *DB) DeleteItem(ctx context.Context, itemId uuid.UUID) error {
	id := itemId.String()
	sqlStatement := `DELETE FROM item WHERE id = ?;`
	result, err := db.Sql.ExecContext(ctx, sqlStatement, id)
	if err != nil {
		logg.Err(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logg.Err(err)
		return err
	}
	if rowsAffected == 0 {
		err := errors.New(fmt.Sprintf("the Record with the id: %s was not found that should not happened while deleting", id))
		logg.Debug(err)
		return err
	} else if rowsAffected != 1 {
		err := errors.New(fmt.Sprintf("the id: %s has unexpected effected number of rows (more than one or less than 0)", id))
		logg.Err(err)
		return err
	}
	return nil
}

// return all the available Items
func (db *DB) Items() ([][]string, error) {
	query := "SELECT id, label, description, picture, quantity, weight, qrcode FROM item;"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying user records: %v", err)
		return [][]string{}, err
	}
	defer rows.Close()

	var itemsArray [][]string
	var item items.Item
	for rows.Next() {
		var idStr string
		err := rows.Scan(&idStr, &item.Label, &item.Description, &item.Picture, &item.Quantity, &item.Weight, &item.QRcode)
		if err != nil {
			log.Printf("Error scanning item record: %v", err)
			return [][]string{}, err
		}

		formatted := fmt.Sprintf("id: %s, label: %s, description: %s, picture: %s, quantity: %d, weight: %s, qrcode: %s \n",
			idStr, item.Label, item.Description, item.Picture, item.Quantity, item.Weight, item.QRcode)
		itemsArray = append(itemsArray, []string{formatted})
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}
	return itemsArray, nil
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
