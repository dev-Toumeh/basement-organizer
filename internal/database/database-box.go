package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type SQLBox struct {
	SQLBasicInfo
	OuterBoxID sql.NullString
	ShelfID    sql.NullString
	AreaID     sql.NullString
}

// RowsToScan returns list of pointers for Scan() method.
func (b *SQLBox) RowsToScan() []any {
	s := append(b.SQLBasicInfo.RowsToScan(), &b.OuterBoxID, &b.ShelfID, &b.AreaID)
	return s
}

// Vals returns all scanned values as strings.
func (s SQLBox) Vals() []string {
	return append(s.SQLBasicInfo.Vals(), s.OuterBoxID.String, s.ShelfID.String, s.AreaID.String)
}

// this function used inside of BoxByField to convert the box sql struct into normal struct
func (s *SQLBox) ToBox() (*items.Box, error) {
	box := &items.Box{}
	info, err := s.ToBasicInfo()
	if err != nil {
		return box, logg.WrapErr(err)
	}

	box.BasicInfo = info
	box.OuterBoxID = ifNullUUID(s.OuterBoxID)
	box.ShelfID = ifNullUUID(s.ShelfID)
	box.ShelfID = ifNullUUID(s.ShelfID)
	return box, nil
}

// Create New Item Record
func (db *DB) CreateBox(newBox *items.Box) (uuid.UUID, error) {
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
	query := fmt.Sprintf("SELECT COUNT(*) FROM box WHERE %s = ?", field)
	var count int
	err := db.Sql.QueryRow(query, value).Scan(&count)
	if err != nil {
		logg.Errf("Error checking item existence: %v", err)
		return false
	}
	return count > 0
}

// BoxIDs returns IDs of all boxes.
func (db *DB) BoxIDs() ([]string, error) {
	ids := []string{}
	sqlStatement := `SELECT id FROM BOX`
	rows, err := db.Sql.Query(sqlStatement)
	if err != nil {
		return ids, logg.Errorf("Error while executing Box ids: %w", err)
	}
	for rows.Next() {
		var idStr string
		err := rows.Scan(&idStr)
		if err != nil {
			return []string{}, logg.Errorf("Error scanning Box ids: %v", err)
		}
		ids = append(ids, idStr)
	}

	return ids, nil
}

// update box data
func (db *DB) UpdateBox(box items.Box) error {
	exist := db.BoxExistById(box.ID)
	if !exist {
		return logg.Errorf("the box does not exist")
	}

	var err error
	box.PreviewPicture, err = ResizePNG(box.Picture, 50)
	if err != nil {
		logg.Errorf("Error while resizing picture of box '%s' to create a preview picture %w", box.Label, err)
	}

	sqlStatement := "UPDATE box SET label = ?, description = ?, picture = ?, preview_picture = ?, qrcode = ?, shelf_id = ?, area_id = ? WHERE id = ?"
	result, err := db.Sql.Exec(sqlStatement, box.Label, box.Description, box.Picture, box.PreviewPicture, box.QRCode, box.ShelfID, box.AreaID, box.ID)

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
func (db *DB) BoxById(id uuid.UUID) (items.Box, error) {
	box := items.Box{}
	if !db.BoxExistById(id) {
		return box, logg.Errorf("box is not exist \n")
	}
	b, err := db.BoxByField("id", id.String())
	if err != nil {
		return box, logg.WrapErr(err)
	}
	box = *b
	return box, err
}

// Get Box  based on given Field
func (db *DB) BoxByField(field string, value string) (*items.Box, error) {
	var sqlBox SQLBox
	stmt := fmt.Sprintf(`
	SELECT
		id, label, description, picture, preview_picture, qrcode, box_id, shelf_id, area_id
	FROM 
		box
	WHERE 
		%s = ?;`, field)

	err := db.Sql.QueryRow(stmt, value).Scan(sqlBox.RowsToScan()...)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	box, err := sqlBox.ToBox()
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	items, err := db.innerListRowsFrom("box", box.ID, "item_fts")
	if err != nil {
		return nil, logg.WrapErr(err)
	}
	box.Items = items

	boxes, err := db.innerListRowsFrom("box", box.ID, "box_fts")
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
func (db *DB) insertNewBox(box *items.Box) (uuid.UUID, error) {
	if db.BoxExistById(box.ID) {
		return uuid.Nil, db.ErrorExist()
	}

	sqlStatement := `INSERT INTO box (id, label, description, picture, preview_picture, qrcode, box_id, shelf_id, area_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	updatePicture(&box.Picture, &box.PreviewPicture)

	result, err := db.Sql.Exec(sqlStatement, box.ID.String(), box.Label, box.Description, box.Picture, box.PreviewPicture, box.QRCode, box.OuterBoxID.String(), box.ShelfID.String(), box.AreaID.String())
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
		return logg.NewError(fmt.Sprintf("can't move box1 (%s) to box2 (%s). box2 is already in box1 and they can't be inside eachother at the same time", box1.String(), box2.String()))
	}

	err := db.MoveTo("box", box1, "box", box2)
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
	err := db.MoveTo("box", boxID, "shelf", toShelfID)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}
