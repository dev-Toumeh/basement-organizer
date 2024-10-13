package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/shelves"
	"database/sql"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
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
			continue
		}
		logg.Debugf("id: %s, username: %s, passwordhash: %s\n", id, username, passwordhash)
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
		err := rows.Scan(&idStr, &item.Label, &item.Description, &item.Quantity, &item.Weight, &item.QRCode, &boxId)
		if err != nil {
			log.Printf("Error scanning item record: %v", err)
			continue
		}
		logg.Debugf("id: %s, label: %s, description: %s, quantity: %d, weight: %s, qrcode: %s \n", idStr, item.Label, item.Description, item.Quantity, item.Weight, item.QRCode)
		if boxId.Valid {
			logg.Debugf("Box ID: %s\n", boxId.String)
		} else {
			logg.Debugf("Box ID is null\n")
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}
}

func (db *DB) PrintTables() {
	query := "SELECT name FROM sqlite_master WHERE type='table';"

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
		BasicInfo: items.BasicInfo{
			ID:          outerBoxId,
			Label:       "OuterBox",
			Description: "This is the outer box",
			Picture:     "base64encodedouterbox",
			QRcode:      "QRcodeOuterBox",
		},
		OuterBoxID: uuid.Nil,
	}
	patchBox := &items.Box{
		BasicInfo: items.BasicInfo{
			ID:          patchBoxId,
			Label:       "PatchBox",
			Description: "This box will allow you to add items again",
			Picture:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==",
			QRcode:      "AB123CD",
		},
		OuterBoxID: outerBoxId,
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

// this function will print all the records inside of  item_fts table
func (db *DB) CheckItemFTSData() error {
	rows, err := db.Sql.Query("SELECT id, label, description FROM item_fts;")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {

		var rowid, label, description string
		err := rows.Scan(&rowid, &label, &description)
		if err != nil {
			return err
		}
		logg.Debugf("Rowid: %s, Label: %s, Description: %s\n", rowid, label, description)
	}
	return nil
}

// Repopulate the item_fts table with data from the item table, but only if item_fts is currently empty.
// work only as path could be deleted when Alex adopt the changes
func (db *DB) RepopulateItemFTS() error {
	var count int
	err := db.Sql.QueryRow(`SELECT COUNT(*) FROM item_fts`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check item_fts: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, err = db.Sql.Exec(`
        INSERT INTO item_fts(id, label, description)
        SELECT id, label, description FROM item;
    `)
	if err != nil {
		return fmt.Errorf("failed to repopulate item_fts: %w", err)
	}

	fmt.Print("Item search data has been successfully repopulated from the item table to the item_fts table. \n")
	return nil
}

// 💀💀 this function will delete all the search relevant tables
// don't use it if you don't know how it works 💀💀
// func (db *DB) dropFTS5TablesAndTriggers() error {
// 	_, err := db.Sql.Exec(`
//         DROP TABLE IF EXISTS item_fts;
//         DROP TABLE IF EXISTS box_fts;
//         DROP TRIGGER IF EXISTS item_ai;
//         DROP TRIGGER IF EXISTS item_au;
//         DROP TRIGGER IF EXISTS item_ad;
//         DROP TRIGGER IF EXISTS box_ai;
//         DROP TRIGGER IF EXISTS box_au;
//         DROP TRIGGER IF EXISTS box_ad;
//     `)
// 	if err != nil {
// 		return fmt.Errorf("failed to drop FTS5 tables and triggers: %w", err)
// 	}
// 	return nil
// }

// 💀💀 this function will delete all the Records form the item_fts table
// don't use it if you don't know how it works 💀💀
// func (db *DB) clearItemFTS() error {
// 	_, err := db.Sql.Exec("DELETE FROM item_fts;")
// 	if err != nil {
// 		return fmt.Errorf("failed to clear item_fts: %w", err)
// 	}
// 	return nil
// }

const SEED = 1234

func (db *DB) InsertSampleItems() {
	gofakeit.Seed(SEED)

	for i := 0; i < 10; i++ {
		newItem := items.Item{
			BasicInfo: items.BasicInfo{
				ID:          uuid.Must(uuid.FromString(gofakeit.UUID())),
				Label:       gofakeit.ProductName(),
				Description: gofakeit.Sentence(5),
				Picture:     ByteToBase64String(gofakeit.ImagePng((i+1)*10, (10-i)*10)),
			},
			Quantity: int64(gofakeit.IntRange(0, 100)),
			Weight:   fmt.Sprintf("%.2f", gofakeit.Float32Range(0, 100)),
			QRCode:   gofakeit.HipsterWord(),
		}

		err := db.insertNewItem(newItem)
		if err != nil {
			logg.Errf("error while adding dummyData %v", err)
			return
		}
	}
	return
}

func (db *DB) InsertSampleBoxes() {
	gofakeit.Seed(SEED)

	for i := 0; i < 10; i++ {
		newBox := items.Box{
			BasicInfo: items.BasicInfo{
				ID:          uuid.Must(uuid.FromString(gofakeit.UUID())),
				Label:       gofakeit.ProductName(),
				Description: gofakeit.Sentence(5),
				Picture:     ByteToBase64String(gofakeit.ImagePng((i+1)*10, (10-i)*10)),
			},
		}

		_, err := db.CreateBox(&newBox)
		if err != nil {
			logg.Errf("error while adding dummyData %v", err)
			return
		}
	}
	return
}

func (db *DB) InsertSampleShelves() {
	gofakeit.Seed(SEED)

	for i := 0; i < 10; i++ {
		newShelf := shelves.Shelf{
			Id:          uuid.Must(uuid.FromString(gofakeit.UUID())),
			Label:       gofakeit.ProductName(),
			Description: gofakeit.Sentence(5),
			Picture:     ByteToBase64String(gofakeit.ImagePng((i+1)*10, (10-i)*10)),
		}

		err := db.CreateShelf(&newShelf)
		if err != nil {
			logg.Errf("error while adding dummyData %v", err)
			return
		}
	}
	return
}
