package database

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"bytes"
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
	info, err := s.SQLBasicInfo.ToBasicInfo()
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	return &items.Item{
		BasicInfo:  info,
		Quantity:   ifNullInt64(s.Quantity),
		Weight:     ifNullString(s.Weight),
		BoxID:      ifNullUUID(s.BoxID),
		BoxLabel:   ifNullString(s.BoxLabel),
		ShelfID:    ifNullUUID(s.ShelfID),
		ShelfLabel: ifNullString(s.ShelfLabel),
		AreaID:     ifNullUUID(s.AreaID),
		AreaLabel:  ifNullString(s.AreaLabel),
	}, nil
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

func (s SQLListRow) ToListRow2() (row common.ListRow, err error) {
	id, err := uuid.FromString(s.ID.String)
	if err != nil {
		return row, logg.WrapErr(err)
	}

	return common.ListRow{
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

	query := fmt.Sprintf(`
        SELECT 
          i.id, i.label, i.description, i.picture, i.preview_picture, i.quantity, i.weight, i.qrcode, 
          COALESCE(i.box_id, '') AS box_id, COALESCE(b.label, '') AS box_label, COALESCE(i.shelf_id, '') AS shelf_id, 
          COALESCE(s.label, '') AS shelf_label, COALESCE(i.area_id, '') AS area_id, COALESCE(a.label, '') AS area_label
        FROM item as i
        LEFT JOIN box as b ON i.box_id = b.id
        LEFT JOIN shelf as s ON i.shelf_id = s.id
        LEFT JOIN area as a ON i.area_id = a.id
        WHERE i.%s = ?;
      `, field)

	row := db.Sql.QueryRow(query, value)

	sqlItem := &SQLItem{}
	err := row.Scan(
		&sqlItem.ID, &sqlItem.Label, &sqlItem.Description, &sqlItem.Picture, &sqlItem.PreviewPicture,
		&sqlItem.Quantity, &sqlItem.Weight, &sqlItem.QRCode, &sqlItem.BoxID, &sqlItem.BoxLabel,
		&sqlItem.ShelfID, &sqlItem.ShelfLabel, &sqlItem.AreaID, &sqlItem.AreaLabel)

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
func (db *DB) ItemById(id uuid.UUID) (*items.Item, error) {
	item, error := db.ItemByField("id", id.String())
	return &item, error
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
		logg.Debugf("virtual item: %v", itemRow.String())
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

	sqlStatement := `INSERT INTO item (id, label, description, picture, preview_picture, quantity, weight,
       qrcode, box_id, shelf_id, area_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
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
func (db *DB) UpdateItem(item items.Item, ignorePicture bool, pictureFormat string) error {
	var err error
	var sqlStatement string
	var result sql.Result
	if ignorePicture {
		sqlStatement = `UPDATE item SET 
			label = ?, description = ?, quantity = ?, weight = ?, 
			qrcode = ?, box_id = ?, shelf_id = ?, area_id = ? WHERE id = ?`

		result, err = db.Sql.Exec(sqlStatement,
			item.BasicInfo.Label, item.BasicInfo.Description, item.Quantity, item.Weight,
			item.BasicInfo.QRCode, item.BoxID.String(), item.ShelfID.String(), item.AreaID.String(), item.BasicInfo.ID.String())
	} else {
		item.PreviewPicture, err = ResizeImage(item.Picture, 50, pictureFormat)
		if err != nil {
			if errors.Is(err, UnsupportedImageFormat) {
				return logg.NewError(logg.CleanLastError(err) + err.Error())
			} else {
				return logg.Errorf("Error while resizing picture of item '%s' to create a preview picture %w", item.Label, err)
			}
		}

		sqlStatement = `UPDATE item SET 
			label = ?, description = ?, picture = ?, preview_picture = ?, quantity = ?, 
			weight = ?, qrcode = ?, box_id = ?, shelf_id = ?, area_id = ? WHERE id = ?`

		result, err = db.Sql.Exec(sqlStatement,
			item.BasicInfo.Label, item.BasicInfo.Description, item.BasicInfo.Picture,
			item.BasicInfo.PreviewPicture, item.Quantity, item.Weight, item.BasicInfo.QRCode,
			item.BoxID.String(), item.ShelfID.String(), item.AreaID.String(), item.BasicInfo.ID.String())
	}

	if err != nil {
		return logg.Errorf("Error while executing update item statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return logg.Errorf("Error checking rows affected while executing update item statement: %w", err)
	}
	if rowsAffected != 1 {
		return logg.Errorf("Unexpected number of rows affected during update: %d for ID %s", rowsAffected, item.BasicInfo.ID.String())
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
	err := db.moveTo("item", id1, "box", id2)
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
	err := db.moveTo("item", itemID, "shelf", toShelfID)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// MoveItemToArea moves item to a shelf.
// To move item out of a shelf set
//
//	toShelfID = uuid.Nil
func (db *DB) MoveItemToArea(itemID uuid.UUID, toAreaID uuid.UUID) error {
	err := db.moveTo("item", itemID, "area", toAreaID)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// ItemListRows retrieves items by label.
// If the query is empty or contains only spaces, it returns default results.
func (db *DB) ItemListRows(searchString string, limit int, pageNr int) (shelfRows []common.ListRow, err error) {
	shelfRows, err = db.listRowsPaginatedFrom("item_fts", searchString, limit, pageNr)
	if err != nil {
		return shelfRows, logg.WrapErr(err)
	}
	return shelfRows, nil
}

// ShelfCounter returns the count of rows in the shelf_fts table that match
// the specified queryString.
// If queryString is empty, it returns the count of all rows in the table.
func (db *DB) ItemListCounter(queryString string) (count int, err error) {
	countQuery := `SELECT COUNT(*) FROM item_fts;`

	if queryString != "" {
		countQuery = fmt.Sprintf(`
			SELECT COUNT(*)
			FROM item_fts
      WHERE label MATCH '%s*' `, queryString)
	}

	err = db.Sql.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error while fetching the number of Items from the Database: %v", err)
	}
	return count, nil
}

// @TODELETE after 01/03.2025
func (db *DB) AddItemToArea(itemID uuid.UUID, toAreaID uuid.UUID) error {
	item, err := db.ItemById(itemID)
	area, err := db.AreaById(toAreaID)
	if err != nil {
		return logg.WrapErr(err)
	}

	item.AreaID = area.ID
	item.AreaLabel = area.Label

	err = db.UpdateItem(*item, false, "image/png")
	if err != nil {
		return logg.WrapErr(err)
	}

	return nil
}

// @TODELETE after 01/03/2025
func (db *DB) AddItemToShelf(itemID uuid.UUID, toShelfID uuid.UUID) error {
	item, err := db.ItemById(itemID)
	shelf, err := db.Shelf(toShelfID)
	if err != nil {
		return logg.WrapErr(err)
	}

	item.ShelfID = shelf.ID
	item.ShelfLabel = shelf.Label

	err = db.UpdateItem(*item, false, "image/png")
	if err != nil {
		return logg.WrapErr(err)
	}

	return nil
}

// @TODELETE after 01/03/2025
func (db *DB) AddItemToBox(itemID uuid.UUID, boxID uuid.UUID) error {
	item, err := db.ItemById(itemID)
	box, err := db.BoxById(boxID)
	if err != nil {
		return logg.WrapErr(err)
	}
	item.BoxID = box.ID
	item.BoxLabel = box.Label

	err = db.UpdateItem(*item, false, "image/png")
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// @TODELETE after 01/03/2025
func (db *DB) MoveItemToObject(itemID uuid.UUID, objectID uuid.UUID, objectType string) error {
	item, err := db.ItemById(itemID)
	if err != nil {
		return logg.WrapErr(err)
	}

	switch objectType {
	case "area":
		area, err := db.AreaById(objectID)
		if err != nil {
			return logg.WrapErr(err)
		}
		item.AreaID = area.ID
		item.AreaLabel = area.Label

	case "shelf":
		shelf, err := db.Shelf(objectID)
		if err != nil {
			return logg.WrapErr(err)
		}
		item.ShelfID = shelf.ID
		item.ShelfLabel = shelf.Label

	case "box":
		box, err := db.BoxById(objectID)
		if err != nil {
			return logg.WrapErr(err)
		}
		item.BoxID = box.ID
		item.BoxLabel = box.Label

	default:
		return logg.WrapErr(fmt.Errorf("invalid object type: %s", objectType))
	}

	err = db.UpdateItem(*item, false, "image/png")
	if err != nil {
		return logg.WrapErr(err)
	}

	return nil
}
