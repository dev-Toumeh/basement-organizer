package database

import (
	"basement/main/internal/boxes"
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type SQLBox struct {
	SQLBasicInfo
	OuterBoxID    sql.NullString
	OuterBoxLabel sql.NullString
	ShelfID       sql.NullString
	ShelfLabel    sql.NullString
	AreaID        sql.NullString
	AreaLabel     sql.NullString
}

// RowsToScan returns list of pointers for *sql.Rows.Scan() method.
//
//	// example usage:
//	rows.Scan(listRow.RowsToScan()...)
func (b *SQLBox) RowsToScan() []any {
	s := append(b.SQLBasicInfo.RowsToScan(), &b.OuterBoxID, &b.OuterBoxLabel,
		&b.ShelfID, &b.ShelfLabel, &b.AreaID, &b.AreaLabel)
	return s
}

// Vals returns all scanned values as strings.
func (s SQLBox) Vals() []string {
	return append(s.SQLBasicInfo.Vals(), s.OuterBoxID.String, s.ShelfID.String, s.AreaID.String)
}

// this function used inside of BoxByField to convert the box sql struct into normal struct
func (s *SQLBox) ToBox() (*boxes.Box, error) {
	box := &boxes.Box{}
	info, err := s.ToBasicInfo()
	if err != nil {
		return box, logg.WrapErr(err)
	}

	return &boxes.Box{
		BasicInfo:     info,
		OuterBoxID:    ifNullUUID(s.OuterBoxID),
		OuterBoxLabel: ifNullString(s.OuterBoxLabel),
		ShelfID:       ifNullUUID(s.ShelfID),
		ShelfLabel:    ifNullString(s.ShelfLabel),
		AreaID:        ifNullUUID(s.AreaID),
		AreaLabel:     ifNullString(s.AreaLabel),
	}, nil
}

// Create New Item Record
func (db *DB) CreateBox(newBox *boxes.Box) (uuid.UUID, error) {
	if db.BoxExistById(newBox.ID) {
		return uuid.Nil, db.ErrorExist()
	}

	id, err := db.insertNewBox(newBox)
	if err != nil {
		return uuid.Nil, logg.Errorf("error while creating new Box: %v", err)
	}
	return id, nil
}

// check if the Box Exist based on Id
// wrapper function for boxExist,
func (db *DB) BoxExistById(id uuid.UUID) bool {
	return db.BoxExist("id", id.String())
}

// check if the Box Exist based on given Field
func (db *DB) BoxExist(field string, value string) bool {
	query := "SELECT COUNT(*) FROM box WHERE " + field + " = ?"
	var count int
	err := db.Sql.QueryRow(query, value).Scan(&count)
	if err != nil {
		logg.Errf("Error checking item existence: %v", err)
		return false
	}
	return count > 0
}

// BoxIDs returns IDs of all boxes.
func (db *DB) BoxIDs() (ids []uuid.UUID, err error) {
	sqlStatement := `SELECT id FROM BOX`
	rows, err := db.Sql.Query(sqlStatement)
	if err != nil {
		return ids, logg.Errorf("Error while executing Box ids: %w", err)
	}
	for rows.Next() {
		var idStr string
		err := rows.Scan(&idStr)
		if err != nil {
			return ids, logg.Errorf("Error scanning Box ids: %v", err)
		}
		ids = append(ids, uuid.FromStringOrNil(idStr))
	}

	return ids, nil
}

// update box data
func (db *DB) UpdateBox(box boxes.Box, ignorePicture bool) error {
	exist := db.BoxExistById(box.ID)
	if !exist {
		return logg.Errorf("the box does not exist")
	}
	err := identicalThing(box.OuterBoxID, box.ID)
	if err != nil {
		return logg.Errorf("Can't have \""+box.Label+"\" in itself %w", err)
	}

	var stmt string
	var result sql.Result
	if ignorePicture {
		stmt = "UPDATE box SET label = ?, description = ?, qrcode = ?, box_id = ?, shelf_id = ?, area_id = ? WHERE id = ?"
		result, err = db.Sql.Exec(stmt, box.Label, box.Description, box.QRCode, box.OuterBoxID, box.ShelfID, box.AreaID, box.ID)
	} else {
		stmt = "UPDATE box SET label = ?, description = ?, picture = ?, preview_picture = ?, qrcode = ?, box_id = ?, shelf_id = ?, area_id = ? WHERE id = ?"
		box.PreviewPicture, err = ResizePNG(box.Picture, 50)
		if err != nil {
			return logg.Errorf("Error while resizing picture of box '%s' to create a preview picture %w", box.Label, err)
		}
		result, err = db.Sql.Exec(stmt, box.Label, box.Description, box.Picture, box.PreviewPicture, box.QRCode, box.OuterBoxID, box.ShelfID, box.AreaID, box.ID)
	}

	if err != nil {
		return logg.Errorf("something wrong happened while runing the box update query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return logg.Errorf("error while finding the item %w", err)
	}
	if rowsAffected == 0 {
		return logg.Errorf("the Record with the id: %s was not found; this should not have happened while updating", box.ID.String())
	} else if rowsAffected != 1 {
		return logg.Errorf("the id: %s has an unexpected number of rows affected (more than one or less than 0)", box.ID.String())
	}
	return nil
}

// delete Box
func (db *DB) DeleteBox(boxId uuid.UUID) error {
	id := boxId.String()

	// check if box is not Empty
	itemExist := db.ItemExist("box_id", id)
	boxExist := db.BoxExist("box_id", boxId.String())
	if itemExist || boxExist {
		return logg.Errorf(`the box with id="%s" is not empty`, id)
	}

	err := db.deleteFrom("box", boxId)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// Get Box based on his ID
// Wrapper function for BoxByField
func (db *DB) BoxById(id uuid.UUID) (boxes.Box, error) {
	box := boxes.Box{}
	if !db.BoxExistById(id) {
		return box, logg.Errorf("box does not exist \n")
	}
	b, err := db.BoxByField("id", id.String())
	if err != nil {
		return box, logg.WrapErr(err)
	}
	box = *b
	return box, err
}

// Get Box  based on given Field
func (db *DB) BoxByField(field string, value string) (*boxes.Box, error) {
	var sqlBox SQLBox
	stmt := fetchBoxQuery(true, field)

	err := db.Sql.QueryRow(stmt, value).Scan(sqlBox.RowsToScan()...)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	box, err := sqlBox.ToBox()
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	items, err := db.InnerListRowsFrom2("box", box.ID, "item_fts")
	if err != nil {
		return nil, logg.WrapErr(err)
	}
	box.Items = items

	boxes, err := db.InnerListRowsFrom2("box", box.ID, "box_fts")
	if err != nil {
		return nil, logg.WrapErr(err)
	}
	box.InnerBoxes = boxes
	logg.Debug(boxes)

	if box.OuterBoxID != uuid.Nil {
		outerbox, err := db.BoxListRowByID(box.OuterBoxID)
		if err != nil {
			return nil, logg.WrapErr(err)
		}
		box.OuterBox = &outerbox
	}

	return box, nil
}

// insert new Box record in the Database
func (db *DB) insertNewBox(box *boxes.Box) (uuid.UUID, error) {
	if db.BoxExistById(box.ID) {
		return uuid.Nil, db.ErrorExist()
	}

	sqlStatement := "INSERT INTO box (" + ALL_BOX_COLS + ") VALUES (?,?,?,?,?,?,?,?,?)"

	updatePicture(&box.Picture, &box.PreviewPicture)

	result, err := db.Sql.Exec(sqlStatement, box.ID.String(), box.Label, box.Description,
		box.Picture, box.PreviewPicture, box.QRCode, box.OuterBoxID.String(),
		box.ShelfID.String(), box.AreaID.String())
	if err != nil {
		return uuid.Nil, logg.Errorf("Error while executing create new box statement: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return uuid.Nil, logg.Errorf("Error while executing create new box statement: %w", err)
	}
	if rowsAffected != 1 {
		return uuid.Nil, logg.Errorf("unexpected number of effected rows, check insirtNewBox")
	}

	return box.ID, nil
}

// MoveBoxToBox moves box1 to another box2.
// To move box out of box2 set
//
//	box1 = uuid.Nil
func (db *DB) MoveBoxToBox(box1 uuid.UUID, box2 uuid.UUID) error {
	// Check if toBoxID is inside boxID.
	// Can't move if if this is the case.
	stmt := "SELECT box_id FROM box WHERE id = ?;"
	var id sql.NullString
	db.Sql.QueryRow(stmt, box2.String()).Scan(&id)
	if id.Valid && id.String == box1.String() {
		return logg.NewError(
			"can't move box1 (" + box1.String() + ") to box2 (" + box2.String() +
				"). box2 is already in box1 and they can't be inside each other at the same time",
		)
	}

	err := db.moveTo("box", box1, "box", box2)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// MoveBoxToShelf moves box to a shelf.
// To move box out of a shelf set
//
//	toShelfID = uuid.Nil
func (db *DB) MoveBoxToShelf(boxID uuid.UUID, toShelfID uuid.UUID) error {
	err := db.moveTo("box", boxID, "shelf", toShelfID)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// MoveBoxToArea moves box to an area.
// To move box out of an area set
//
//	toAreaID = uuid.Nil
func (db *DB) MoveBoxToArea(boxID uuid.UUID, toAreaID uuid.UUID) error {
	err := db.moveTo("box", boxID, "area", toAreaID)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// BoxListRows retrieves virtual boxes by label.
// If the query is empty or contains only spaces, it returns default results.
func (db *DB) BoxListRows(searchQuery string, limit int, page int) (listRows []common.ListRow, err error) {
	listRows, err = db.listRowsPaginatedFrom("box_fts", searchQuery, limit, page)
	if err != nil {
		return listRows, logg.WrapErr(err)
	}
	return listRows, nil
}

// Get the virtual Box based on his ID
func (db *DB) BoxListRowByID(id uuid.UUID) (common.ListRow, error) {
	row, err := db.listRowByID("box_fts", id)
	if err != nil {
		return row, logg.WrapErr(err)
	}
	return row, nil
}

func (db *DB) InnerBoxListRows(id uuid.UUID) (listRows []common.ListRow, err error) {
	listRows, err = db.InnerListRowsFrom2("box", id, "box_fts")
	if err != nil {
		return listRows, logg.WrapErr(err)
	}
	return listRows, nil
}

// returns the count of rows in the box_fts table that match the specified searchString.
// If queryString is empty, it returns the count of all rows in the table.
func (db *DB) BoxListCounter(searchString string) (count int, err error) {
	countQuery := `SELECT COUNT(*) FROM box_fts;`

	if searchString != "" {
		countQuery = ` SELECT COUNT(*) FROM box_fts WHERE label MATCH '` + searchString + `*'`
	}

	err = db.Sql.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error while fetching the number of box from the database: %v", err)
	}
	return count, nil
}

// returns the count of rows in the box_fts table that match the specified searchString.
// If queryString is empty, it returns the count of all rows in the table.
func (db *DB) InnerBoxInBoxListCounter(searchString string, inTable string, inTableID uuid.UUID) (count int, err error) {
	countQuery := `SELECT COUNT(*) FROM box_fts WHERE ` + inTable + `_id = ?;`

	if searchString != "" {
		countQuery = ` SELECT COUNT(*) FROM box_fts WHERE label MATCH '` + searchString + `*' AND ` + inTable + `_id = ?`
	}

	err = db.Sql.QueryRow(countQuery, inTableID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error while fetching the number of box from the database: %v", err)
	}
	return count, nil
}

// check if the box row  exist
func (db *DB) BoxRowExist(id uuid.UUID) (bool, error) {
	return db.Exists("box_fts", id)
}
