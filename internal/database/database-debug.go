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

// print the shelves Records in the console
func (db *DB) PrintShelvesRecords() {
	query := "SELECT id, label, description, picture, preview_picture, qrcode, height, width, depth, rows, cols FROM shelf"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying shelf records: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("Shelf records:")

	for rows.Next() {
		var (
			id             sql.NullString
			label          sql.NullString
			description    sql.NullString
			picture        sql.NullString
			previewPicture sql.NullString
			qrcode         sql.NullString
			height         sql.NullFloat64
			width          sql.NullFloat64
			depth          sql.NullFloat64
			rowsField      sql.NullInt64 // Renamed to avoid conflict
			cols           sql.NullInt64
		)
		err := rows.Scan(
			&id,
			&label,
			&description,
			&picture,
			&previewPicture,
			&qrcode,
			&height,
			&width,
			&depth,
			&rowsField,
			&cols,
		)
		if err != nil {
			log.Printf("Error scanning shelf record: %v", err)
			continue
		}

		// Extract values, handling NULLs
		idValue := "NULL"
		if id.Valid {
			idValue = id.String
		}
		labelValue := "NULL"
		if label.Valid {
			labelValue = label.String
		}
		descriptionValue := "NULL"
		if description.Valid {
			descriptionValue = description.String
		}
		pictureValue := "NULL"
		if picture.Valid {
			pictureValue = picture.String
		}
		previewPictureValue := "NULL"
		if previewPicture.Valid {
			previewPictureValue = previewPicture.String
		}
		qrcodeValue := "NULL"
		if qrcode.Valid {
			qrcodeValue = qrcode.String
		}
		heightValue := "NULL"
		if height.Valid {
			heightValue = fmt.Sprintf("%f", height.Float64)
		}
		widthValue := "NULL"
		if width.Valid {
			widthValue = fmt.Sprintf("%f", width.Float64)
		}
		depthValue := "NULL"
		if depth.Valid {
			depthValue = fmt.Sprintf("%f", depth.Float64)
		}
		rowsValue := "NULL"
		if rowsField.Valid {
			rowsValue = fmt.Sprintf("%d", rowsField.Int64)
		}
		colsValue := "NULL"
		if cols.Valid {
			colsValue = fmt.Sprintf("%d", cols.Int64)
		}

		// Corrected fmt.Printf with matching format specifiers and arguments
		fmt.Printf("id: %s, label: %s, description: %s, picture: %s, preview_picture: %s, qrcode: %s, height: %s, width: %s, depth: %s, rows: %s, cols: %s\n",
			idValue, labelValue, descriptionValue, pictureValue, previewPictureValue, qrcodeValue,
			heightValue, widthValue, depthValue, rowsValue, colsValue)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}
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
				QRCode:      gofakeit.HipsterWord(),
			},
			Quantity: int64(gofakeit.IntRange(0, 100)),
			Weight:   fmt.Sprintf("%.2f", gofakeit.Float32Range(0, 100)),
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
		newShelf := &shelves.Shelf{
			BasicInfo: items.BasicInfo{
				ID:          uuid.Must(uuid.FromString(gofakeit.UUID())),
				Label:       gofakeit.ProductName(),
				Description: gofakeit.Sentence(5),
				Picture:     ByteToBase64String(gofakeit.ImagePng((i+1)*10, (10-i)*10)),
			},
		}

		err := db.CreateShelf(newShelf)
		if err != nil {
			logg.Errf("error while adding dummyData %v", err)
			return
		}
	}
	return
}

// ShelfListRowsPaginated returns always a slice with the number of rows specified.
func (db *DB) ShelfListRowsPaginated(limit int, offset int) (shelfRows []*items.ListRow, count int, err error) {
	i := 0
	count, err = db.ShelfCounter("")
	shelfRows = make([]*items.ListRow, count)

	queryNoSearch := `
		SELECT
			id, label, area_id, area_label, preview_picture
		FROM shelf_fts
			ORDER BY label ASC
		LIMIT ? OFFSET ?;`
	results, err := db.Sql.Query(queryNoSearch, limit, offset)
	if err != nil {
		return shelfRows, count, logg.WrapErr(err)
	}

	for results.Next() {
		listRow := SQLListRow{}
		err = results.Scan(&listRow.ID, &listRow.Label, &listRow.AreaID, &listRow.AreaLabel, &listRow.PreviewPicture)
		if err != nil {
			return shelfRows, count, logg.WrapErr(err)
		}
		shelfRows[i], err = listRow.ToListRow()
		if err != nil {
			return shelfRows, count, logg.WrapErr(err)
		}
		i += 1
	}
	if results.Err() != nil {
		return shelfRows, count, logg.WrapErr(err)
	}

	return shelfRows, count, nil
}
