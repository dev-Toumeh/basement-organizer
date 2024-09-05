package database

import (
	"basement/main/internal/items"
	"database/sql"
	"fmt"
	"log"

	"github.com/gofrs/uuid/v5"
)

func (db *DB) PrintUserRecords() {
	query := "SELECT id, username, passwordhash FROM user;"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying user records: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("User records:")

	for rows.Next() {
		var id string
		var username, passwordhash string
		if err := rows.Scan(&id, &username, &passwordhash); err != nil {
			log.Printf("Error scanning user record: %v", err)
			continue // Log the error and continue with the next row
		}
		fmt.Printf("id: %s, username: %s, passwordhash: %s\n", id, username, passwordhash)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}
}

func (db *DB) PrintItemRecords() {
	query := "SELECT id, label, description, quantity, weight, qrcode , box_id FROM item;"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying user records: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("User records:")

	for rows.Next() {

		var item items.Item
		var idStr string
		var boxId sql.NullString
		err := rows.Scan(&idStr, &item.Label, &item.Description, &item.Quantity, &item.Weight, &item.QRcode, &boxId)
		if err != nil {
			log.Printf("Error scanning item record: %v", err)
			continue // Log the error and continue with the next row
		}
		fmt.Printf("id: %s, label: %s, description: %s, quantity: %d, weight: %s, qrcode: %s \n", idStr, item.Label, item.Description, item.Quantity, item.Weight, item.QRcode)
		if boxId.Valid {
			fmt.Printf("Box ID: %s\n", boxId.String)
		} else {
			fmt.Printf("Box ID is null\n")
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}
}

func (db *DB) PrintTables() {
	query := "SELECT name FROM sqlite_master WHERE type='table';"
	//	query := " SELECT name FROM pragma_table_info('user');"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Fatalf("Error querying tables: %v", err)
	}
	defer rows.Close()

	fmt.Println("Available tables:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Error scanning table name: %v", err)
		}
		fmt.Println(tableName)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error during rows iteration: %v", err)
	}
}

// need to be deleted, use it to fix the database after update
func (db *DB) DatabasePatcher() error {
	// Define the UUIDs for the boxes
	patchBoxId := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174111"))
	outerBoxId := uuid.Must(uuid.FromString("18c60ba9-ffac-48f1-8c7c-473bd35acbea"))

	// Define the box structures
	outerBox := &items.Box{
		Id:          outerBoxId,
		Label:       "OuterBox",
		Description: "This is the outer box",
		Picture:     "base64encodedouterbox",
		QRcode:      "QRcodeOuterBox",
		OuterBoxId:  uuid.Nil,
	}
	patchBox := &items.Box{
		Id:          patchBoxId,
		Label:       "PatchBox",
		Description: "This box will allow you to add items again",
		Picture:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==",
		QRcode:      "AB123CD",
		OuterBoxId:  outerBoxId,
	}

	// Create the boxes using the CreateNewBox function
	if _, err := db.CreateBox(outerBox); err != nil {
		return err
	}
	if _, err := db.CreateBox(patchBox); err != nil {
		return err
	}

	// Adds the box_id column
	alterStmt := `ALTER TABLE item ADD COLUMN box_id TEXT REFERENCES box(id);`
	_, err := db.Sql.Exec(alterStmt)
	if err != nil {
		return err
	}

	// Update the box_id for all rows in the item table with the patchBoxId
	updateStmt := `UPDATE item SET box_id = ?;`
	_, err = db.Sql.Exec(updateStmt, patchBoxId.String())
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) addBoxIDColumn() error {
	alterStmt := `ALTER TABLE item ADD COLUMN box_id TEXT REFERENCES box(id);`
	_, err := db.Sql.Exec(alterStmt)
	if err != nil {
		return err
	}
	return nil
}
