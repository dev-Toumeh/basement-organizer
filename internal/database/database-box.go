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

// Moves box with id1 into box with id2.
func (db *DB) MoveBox(id1 uuid.UUID, id2 uuid.UUID) error {
	updateStmt := `UPDATE box SET outerbox_id = ? WHERE Id = ?;`
	_, err := db.Sql.Exec(updateStmt, id2, id1)
	if err != nil {
		return err
	}
	return nil
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
	result, err := db.Sql.Exec(sqlStatement, box.Label, box.Description, box.Picture, box.PreviewPicture, box.QRcode, box.ShelfID, box.AreaID, box.ID)

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
	boxExist := db.BoxExist("outerbox_id", boxId.String())
	if itemExist || boxExist {
		return logg.Errorf("the box is not empty")
	}

	sqlStatement := `DELETE FROM box WHERE id = ?;`
	_, err := db.Sql.Exec(sqlStatement, id)
	if err != nil {
		return logg.Errorf("error whiel deleting the box: %W", err)
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
	var box *items.Box
	var outerBox = &items.Box{}
	var itemListRows = []*items.ListRow{}
	var innerBoxListRows = []*items.ListRow{}
	boxInitialized := false
	outerBoxInitialized := false
	addedInBoxes := make(map[string]bool)
	addedItems := make(map[string]*items.Item)

	query := fmt.Sprintf(
		`SELECT 
            b.id, b.label, b.description, b.picture, b.preview_picture, b.qrcode, b.outerbox_id, b.shelf_id, b.area_id,
            ob.id, ob.label, ob.preview_picture,
            ib.id,
            i.id
        FROM 
            box AS b
        LEFT JOIN 
            box As ib ON b.id = ib.outerbox_id
        LEFT JOIN 
            box As ob ON b.outerbox_id = ob.id
        LEFT JOIN 
            item As i ON b.id = i.box_id 
        WHERE 
            b.%s = ?;`, field)
	rows, err := db.Sql.Query(query, value)
	if err != nil {
		return nil, logg.WrapErr(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			sqlBox      SQLBox
			sqlOuterBox SQLBox
			sqlInnerBox SQLBox
			sqlItem     SQLItem
			innerBox    *items.Box
		)

		err := rows.Scan(
			&sqlBox.ID, &sqlBox.Label, &sqlBox.Description, &sqlBox.Picture, &sqlBox.PreviewPicture, &sqlBox.QRCode, &sqlBox.OuterBoxID, &sqlBox.ShelfID, &sqlBox.AreaID,
			&sqlOuterBox.ID, &sqlOuterBox.Label, &sqlOuterBox.PreviewPicture,
			&sqlInnerBox.ID,
			&sqlItem.ID,
		)
		if err != nil {
			return nil, logg.Errorf("error scanning row: %w", err)
		}

		// Initialize the main Box only once
		if !boxInitialized {
			box, err = sqlBox.ToBox()
			if err != nil {
				return nil, logg.Errorf("error assigning the box data: %w", err)
			}
			boxInitialized = true
		}

		// Initialize outer box only once
		if !outerBoxInitialized && sqlOuterBox.ID.Valid {
			outerBox, err = sqlOuterBox.ToBox()
			if err != nil {
				return &items.Box{}, logg.Errorf("something wrong happened while assigning the outboxSQLBox data %w", err)
			}
			outerBoxInitialized = true
		}

		// Add the inboxes if its ID is valid and not added before
		if sqlInnerBox.ID.Valid && !addedInBoxes[sqlInnerBox.ID.String] {
			innerBox, err = sqlInnerBox.ToBox()
			if err != nil {
				return nil, logg.Errorf("error converting SQLBox to Box: %w", err)
			}
			innerBoxListRow, err := db.BoxListRowByID(innerBox.ID)
			if err != nil {
				return nil, logg.WrapErr(err)
			}
			innerBoxListRows = append(innerBoxListRows, &innerBoxListRow)
			addedInBoxes[sqlInnerBox.ID.String] = true
		}

		// Add the item to the itemsList if itemId is valid
		if sqlItem.ID.Valid {
			if _, exists := addedItems[sqlItem.ID.String]; !exists {
				item, err := sqlItem.ToItem()
				if err != nil {
					return nil, logg.Errorf("error converting SQLItem to Item: %w", err)
				}
				itemListRow, err := db.ItemListRowByID(item.ID)
				if err != nil {
					return nil, logg.WrapErr(err)
				}
				itemListRows = append(itemListRows, itemListRow)
				addedItems[sqlItem.ID.String] = item
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, logg.Errorf("error iterating through rows: %w", err)
	}

	// Assign the items list to the box
	box.OuterBox = &items.ListRow{BoxID: outerBox.ID, Label: outerBox.Label}
	box.Items = itemListRows
	box.InnerBoxes = innerBoxListRows

	return box, nil
}

// insert new Box record in the Database
func (db *DB) insertNewBox(box *items.Box) (uuid.UUID, error) {
	if db.BoxExistById(box.ID) {
		return uuid.Nil, db.ErrorExist()
	}

	sqlStatement := `INSERT INTO box (id, label, description, picture, preview_picture, qrcode, outerbox_id, shelf_id, area_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	logg.Debugf("SQL: %s", sqlStatement)
	result, err := db.Sql.Exec(sqlStatement, box.ID.String(), box.Label, box.Description, box.Picture, box.PreviewPicture, box.QRcode, box.OuterBoxID.String(), box.ShelfID.String(), box.AreaID.String())
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
