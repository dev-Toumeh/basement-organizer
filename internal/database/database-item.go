package database

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type SQLBasicInfo struct {
	ID             sql.NullString
	Label          sql.NullString
	Description    sql.NullString
	Picture        sql.NullString
	PreviewPicture sql.NullString
	QRCode         sql.NullString
}

// Vals returns all scanned values as strings.
func (s SQLBasicInfo) Vals() []string {
	return []string{s.ID.String, s.Label.String, s.Description.String, s.Picture.String, s.PreviewPicture.String, s.QRCode.String}
}

// RowsToScan returns list of pointers for *sql.Rows.Scan() method.
//
//	// example usage:
//	rows.Scan(listRow.RowsToScan()...)
func (s *SQLBasicInfo) RowsToScan() []any {
	return []any{&s.ID, &s.Label, &s.Description, &s.Picture, &s.PreviewPicture, &s.QRCode}
}

func (s SQLBasicInfo) ToBasicInfo() (common.BasicInfo, error) {
	id, err := uuid.FromString(s.ID.String)
	if err != nil {
		return common.BasicInfo{}, logg.WrapErr(err)
	}

	return common.BasicInfo{
		ID:             id,
		Label:          ifNullString(s.Label),
		Description:    ifNullString(s.Description),
		Picture:        ifNullString(s.Picture),
		PreviewPicture: ifNullString(s.PreviewPicture),
		QRCode:         ifNullString(s.QRCode),
	}, nil
}

type SQLItem struct {
	SQLBasicInfo
	Quantity   sql.NullInt64
	Weight     sql.NullString
	BoxID      sql.NullString
	BoxLabel   sql.NullString
	ShelfID    sql.NullString
	ShelfLabel sql.NullString
	AreaID     sql.NullString
	AreaLabel  sql.NullString
}

func (i SQLItem) String() string {
	return fmt.Sprintf("SQLItem[ID=%s, Label=%s, QRCode=%s, Quantity=%d, Weight=%s, BoxID=%s, BoxLabel=%s, ShelfID=%s, ShelfLabel=%s, AreaID=%s, AreaLabel=%s]",
		i.SQLBasicInfo.ID.String, i.SQLBasicInfo.Label.String, i.SQLBasicInfo.QRCode.String, i.Quantity.Int64, i.Weight.String,
		i.BoxID.String, i.BoxLabel.String, i.ShelfID.String, i.ShelfLabel.String, i.AreaID.String, i.AreaLabel.String)
}

// this function used inside of BoxByField to convert the sql Item struct into normal struct
func (s *SQLItem) ToItem() (*items.Item, error) {
	// var err error
	info, err := s.SQLBasicInfo.ToBasicInfo()
	if err != nil {
		return nil, logg.WrapErr(err)
	}
	item := &items.Item{BasicInfo: info}

	// Convert and assign the ID
	if s.BoxID.Valid {
		item.BoxID, err = uuid.FromString(s.BoxID.String)
		if err != nil {
			return nil, logg.Errorf("Error parsing UUID for box ID: '%v' %w", s.BoxID, err)
		}
	} else {
		return nil, logg.NewError(fmt.Sprintf("box ID is required but was null in item '%v'", s))
	}

	if s.Quantity.Valid {
		item.Quantity = s.Quantity.Int64
	} else {
		item.Quantity = 1
	}

	if s.Weight.Valid {
		item.Weight = s.Weight.String
	} else {
		item.Weight = ""
	}

	item.ShelfID = ifNullUUID(s.ShelfID)
	item.AreaID = ifNullUUID(s.AreaID)

	return item, nil
}

type SQLListRow struct {
	ID             sql.NullString
	Label          sql.NullString
	Description    sql.NullString
	PreviewPicture sql.NullString
	BoxID          sql.NullString
	BoxLabel       sql.NullString
	ShelfID        sql.NullString
	ShelfLabel     sql.NullString
	AreaID         sql.NullString
	AreaLabel      sql.NullString
}

func (s SQLListRow) ToListRow() (*common.ListRow, error) {
	id, err := uuid.FromString(s.ID.String)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	return &common.ListRow{
		ID:             id,
		Label:          ifNullString(s.Label),
		Description:    ifNullString(s.Description),
		PreviewPicture: ifNullString(s.PreviewPicture),
		BoxID:          ifNullUUID(s.BoxID),
		BoxLabel:       ifNullString(s.BoxLabel),
		ShelfID:        ifNullUUID(s.ShelfID),
		ShelfLabel:     ifNullString(s.ShelfLabel),
		AreaID:         ifNullUUID(s.AreaID),
		AreaLabel:      ifNullString(s.AreaLabel),
	}, nil

}

// RowsToScan returns list of pointers for *sql.Rows.Scan() method.
//
//	// example usage:
//	rows.Scan(listRow.RowsToScan()...)
func (s *SQLListRow) RowsToScan() []any {
	return []any{
		&s.ID, &s.Label, &s.Description, &s.PreviewPicture, &s.BoxID, &s.BoxLabel, &s.ShelfID, &s.ShelfLabel, &s.AreaID, &s.AreaLabel,
	}
}

// Create New Item Record
func (db *DB) CreateNewItem(newItem items.Item) error {
	exist, err := db.Exists("item", newItem.ID)
	if exist {
		return db.ErrorExist()
	}
	if err != nil {
		return logg.WrapErr(err)
	}
	err = db.insertNewItem(newItem)
	if err != nil {
		return err
	}
	return nil
}

// Get Item Record based on given Field
func (db *DB) ItemByField(field string, value string) (items.Item, error) {

	if !db.ItemExist(field, value) {
		return items.Item{}, logg.WrapErr(sql.ErrNoRows)
	}

	query := fmt.Sprintf(`SELECT 
		id, label, description, picture, preview_picture, quantity, weight, qrcode, box_id, shelf_id, area_id
	FROM item WHERE %s = ?;`, field)
	row := db.Sql.QueryRow(query, value)

	sqlItem := &SQLItem{}
	err := row.Scan(&sqlItem.ID, &sqlItem.Label, &sqlItem.Description, &sqlItem.Picture, &sqlItem.PreviewPicture, &sqlItem.Quantity, &sqlItem.Weight, &sqlItem.QRCode, &sqlItem.BoxID, &sqlItem.ShelfID, &sqlItem.AreaID)

	if err != nil {
		return items.Item{}, logg.Errorf("Error while checking if the Item is available: %w ", err)
	}
	item, err := sqlItem.ToItem()
	if err != nil {
		return items.Item{}, logg.WrapErr(err)
	}

	return *item, nil
}

// check if the Item exist
func (db *DB) ItemExist(field string, value string) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM item WHERE %s = ?", field)
	var count int
	err := db.Sql.QueryRow(query, value).Scan(&count)
	if err != nil {
		log.Println("Error checking item existence:", err)
		return false
	}
	return count > 0
}

// Item returns new Item struct if id matches.
func (db *DB) ItemById(id uuid.UUID) (items.Item, error) {
	return db.ItemByField("id", id.String())
}

// ListItemById returns a single item with less information suitable for a list row.
func (db *DB) ItemListRowByID(id uuid.UUID) (*common.ListRow, error) {
	query := `
		SELECT 
            i.id, i.label, i.preview_picture,
            b.id, b.label,
			s.id, s.label,
			a.id, a.label
        FROM 
            item AS i
        LEFT JOIN 
            box AS b ON b.id = i.box_id 
        LEFT JOIN 
            shelf AS s ON s.id = i.shelf_id 
        LEFT JOIN 
            area AS a ON a.id = i.area_id 
        WHERE 
            i.id = ?;`
	queryRow := db.Sql.QueryRow(query, id.String())

	sqlListRow := SQLListRow{}

	err := queryRow.Scan(&sqlListRow.ID, &sqlListRow.Label, &sqlListRow.PreviewPicture, &sqlListRow.BoxID, &sqlListRow.BoxLabel, &sqlListRow.ShelfID, &sqlListRow.ShelfLabel, &sqlListRow.AreaID, &sqlListRow.AreaLabel)
	if err != nil {
		return nil, logg.Errorf("%s %w", query, err)
	}

	itemRow, err := sqlListRow.ToListRow()
	if err != nil {
		return itemRow, logg.WrapErr(err)
	}

	if env.Development() {
		b := bytes.Buffer{}
		server.WriteJSON(&b, itemRow)
		logg.Debugf("virtual item: %v", b.String())
	}

	return itemRow, nil
}

// return items id's in array from type string
func (db *DB) ItemIDs() (ids []uuid.UUID, err error) {
	query := "SELECT id FROM item;"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying item records: %v", err)
		return ids, err
	}
	defer rows.Close()

	for rows.Next() {
		var idStr string
		err := rows.Scan(&idStr)
		if err != nil {
			log.Printf("Error scanning item record: %v", err)
			continue
		}
		ids = append(ids, uuid.FromStringOrNil(idStr))
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}

	// Print all IDs
	return ids, nil
}

// here we run the insert new Item query separate from the public function
// it make the code more readable
func (db *DB) insertNewItem(item items.Item) error {
	updatePicture(&item.Picture, &item.PreviewPicture)

	sqlStatement := `INSERT INTO item (id, label, description, picture, preview_picture, quantity, weight, qrcode, box_id, shelf_id, area_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Sql.Exec(sqlStatement, item.BasicInfo.ID.String(),
		item.BasicInfo.Label, item.BasicInfo.Description, item.BasicInfo.Picture,
		item.BasicInfo.PreviewPicture, item.Quantity, item.Weight, item.BasicInfo.QRCode,
		item.BoxID.String(), item.ShelfID.String(), item.AreaID.String())
	if err != nil {
		return logg.Errorf("Error while executing create new item statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return logg.Errorf("Error checking rows affected while executing create new item statement: %w", err)
	}
	if rowsAffected != 1 {
		return logg.NewError("item not added")
	}
	return nil
}

// update the item based on the id
func (db *DB) UpdateItem(ctx context.Context, item items.Item) error {
	updatePicture(&item.Picture, &item.PreviewPicture)

	sqlStatement := `UPDATE item Set label = ?, description = ?, picture = ?, preview_picture = ?, quantity = ?, weight = ?, qrcode = ? WHERE id = ?`
	result, err := db.Sql.ExecContext(ctx, sqlStatement, item.Label, item.Description, item.Picture, item.PreviewPicture, item.Quantity, item.Weight, item.QRCode, item.ID.String())
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
		err := errors.New(fmt.Sprintf("the Record with the id: %s was not found that should not happened while updating", item.ID.String()))
		logg.Debug(err)
		return err
	} else if rowsAffected != 1 {
		err := errors.New(fmt.Sprintf("the id: %s has unexpected effected number of rows (more than one or less than 0)", item.ID.String()))
		logg.Err(err)
		return err
	}
	return nil
}

// Delete Item by Id
func (db *DB) DeleteItem(itemId uuid.UUID) error {
	err := db.deleteFrom("item", itemId)
	if err != nil {
		return logg.WrapErr(err)
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
		err := rows.Scan(&idStr, &item.Label, &item.Description, &item.Picture, &item.Quantity, &item.Weight, &item.QRCode)
		if err != nil {
			log.Printf("Error scanning item record: %v", err)
			return [][]string{}, err
		}

		formatted := fmt.Sprintf("id: %s, label: %s, description: %s, picture: %s, quantity: %d, weight: %s, qrcode: %s \n",
			idStr, item.Label, item.Description, item.Picture, item.Quantity, item.Weight, item.QRCode)
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

// delete one item or more
func (db *DB) DeleteItems(itemIds []uuid.UUID) error {
	if len(itemIds) == 0 {
		return nil
	}

	// Create placeholders and arguments
	placeholders := make([]string, len(itemIds))
	args := make([]interface{}, len(itemIds))
	for i, id := range itemIds {
		placeholders[i] = "?"
		args[i] = id
	}

	// Join the placeholders with commas
	sqlStatement := `DELETE FROM item WHERE id IN (` + strings.Join(placeholders, ",") + `);`

	// Execute the query with the item IDs as arguments
	result, err := db.Sql.Exec(sqlStatement, args...)
	if err != nil {
		logg.Err(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logg.Err(err)
		return err
	}
	if rowsAffected != int64(len(itemIds)) {
		err := logg.Errorf("unexpected number of rows affected while deleting. Expected: %d, Actual: %d", len(itemIds), rowsAffected)
		logg.Err(err)
		return err
	}

	return nil
}

// MoveItemToBox moves item to a box.
// To move item out of a box set
//
//	id2 = uuid.Nil
func (db *DB) MoveItemToBox(id1 uuid.UUID, id2 uuid.UUID) error {
	err := db.MoveTo("item", id1, "box", id2)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// MoveItemToShelf moves item to a shelf.
// To move item out of a shelf set
//
//	toShelfID = uuid.Nil
func (db *DB) MoveItemToShelf(itemID uuid.UUID, toShelfID uuid.UUID) error {
	err := db.MoveTo("item", itemID, "shelf", toShelfID)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}
