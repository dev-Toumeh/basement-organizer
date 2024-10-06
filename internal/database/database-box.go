package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type SQLBox struct {
	ID             sql.NullString
	Label          sql.NullString
	Description    sql.NullString
	Picture        sql.NullString
	PreviewPicture sql.NullString
	QRcode         sql.NullString
	OuterboxID     sql.NullString
}

type SqlItem struct {
	ItemID             sql.NullString
	ItemLabel          sql.NullString
	ItemDescription    sql.NullString
	ItemPicture        sql.NullString
	ItemPreviewPicture sql.NullString
	ItemQRCode         sql.NullString
	ItemQuantity       sql.NullInt64
	ItemWeight         sql.NullString
	ItemBoxID          sql.NullString
	ItemShelfID        sql.NullString
	ItemAreaID         sql.NullString
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
		logg.Errorf2(fmt.Sprintf("Error while resizing picture of box '%s' to create a preview picture", box.Label), err)
	}

	sqlStatement := "UPDATE box SET label = ?, description = ?, picture = ?, preview_picture = ?, qrcode = ? WHERE id = ?"
	result, err := db.Sql.Exec(sqlStatement, box.Label, box.Description, box.Picture, box.PreviewPicture, box.QRcode, box.ID)

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
	var itemsList = []*items.ItemListRow{}
	var inBoxList = []*items.BoxListRow{}
	boxInitialized := false
	outerBoxInitialized := false
	addedInBoxes := make(map[string]bool)
	addedItems := make(map[string]*items.Item)

	query := fmt.Sprintf(
		`SELECT 
            b.id, b.label, b.description, b.picture, b.preview_picture, b.qrcode, b.outerbox_id,
            ob.id, ob.label, ob.preview_picture,
            ib.id, ib.label, ib.preview_picture,
            i.id, i.label, i.preview_picture, i.box_id, i.shelf_id, i.area_id
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
			boxSQLBox    SQLBox
			outboxSQLBox SQLBox
			inboxSQLBox  SQLBox
			itemSQLItem  SqlItem
			inbox        *items.Box
		)

		err := rows.Scan(
			&boxSQLBox.ID, &boxSQLBox.Label, &boxSQLBox.Description, &boxSQLBox.Picture, &boxSQLBox.PreviewPicture, &boxSQLBox.QRcode, &boxSQLBox.OuterboxID,
			&outboxSQLBox.ID, &outboxSQLBox.Label, &outboxSQLBox.PreviewPicture,
			&inboxSQLBox.ID, &inboxSQLBox.Label, &inboxSQLBox.PreviewPicture,
			&itemSQLItem.ItemID, &itemSQLItem.ItemLabel, &itemSQLItem.ItemPreviewPicture, &itemSQLItem.ItemBoxID, &itemSQLItem.ItemShelfID, &itemSQLItem.ItemAreaID,
		)
		if err != nil {
			return nil, logg.Errorf("error scanning row: %w", err)
		}

		// Initialize the main Box only once
		if !boxInitialized {
			box, err = convertSQLBoxToBox(&boxSQLBox)
			if err != nil {
				return nil, logg.Errorf("error assigning the box data: %w", err)
			}
			boxInitialized = true
		}

		// Initialize outer box only once
		if !outerBoxInitialized && outboxSQLBox.ID.Valid {
			outerBox, err = convertSQLBoxToBox(&outboxSQLBox)
			if err != nil {
				return &items.Box{}, logg.Errorf("something wrong happened while assigning the outboxSQLBox data %w", err)
			}
			outerBoxInitialized = true
		}

		// Add the inboxes if its ID is valid and not added before
		if inboxSQLBox.ID.Valid && !addedInBoxes[inboxSQLBox.ID.String] {
			inbox, err = convertSQLBoxToBox(&inboxSQLBox)
			if err != nil {
				return nil, logg.Errorf("error converting SQLBox to Box: %w", err)
			}
			inBoxList = append(inBoxList, &items.BoxListRow{BoxID: inbox.ID, Label: inbox.Label, PreviewPicture: inbox.PreviewPicture})
			addedInBoxes[inboxSQLBox.ID.String] = true
		}

		// Add the item to the itemsList if itemId is valid
		if itemSQLItem.ItemID.Valid {
			if _, exists := addedItems[itemSQLItem.ItemID.String]; !exists {
				item, err := convertSQLItemToItem(&itemSQLItem)
				if err != nil {
					return nil, logg.Errorf("error converting SQLItem to Item: %w", err)
				}
				itemsList = append(itemsList, &items.ItemListRow{ItemID: item.ID, Label: item.Label})
				addedItems[itemSQLItem.ItemID.String] = item
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, logg.Errorf("error iterating through rows: %w", err)
	}

	// Assign the items list to the box
	box.OuterBox = &items.BoxListRow{BoxID: outerBox.ID, Label: outerBox.Label}
	box.Items = itemsList
	box.InnerBoxes = inBoxList

	return box, nil
}

// insert new Box record in the Database
func (db *DB) insertNewBox(box *items.Box) (uuid.UUID, error) {
	if db.BoxExistById(box.ID) {
		return uuid.Nil, db.ErrorExist()
	}

	sqlStatement := `INSERT INTO box (id, label, description, picture, preview_picture, qrcode, outerbox_id) VALUES (?, ?, ?, ?, ?, ?, ?)`
	logg.Debugf("SQL: %s", sqlStatement)
	result, err := db.Sql.Exec(sqlStatement, box.ID.String(), box.Label, box.Description, box.Picture, box.PreviewPicture, box.QRcode, box.OuterBoxID.String())
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

// this function used inside of BoxByField to convert the box sql struct into normal struct
func convertSQLBoxToBox(sqlBox *SQLBox) (*items.Box, error) {
	box := &items.Box{}

	// Convert and assign the ID
	if sqlBox.ID.Valid {
		var err error
		box.ID, err = uuid.FromString(sqlBox.ID.String)
		if err != nil {
			return box, logg.Errorf("Error parsing UUID for box: %w", err)
		}

		// Assign other fields only if they are valid
		if sqlBox.Label.Valid {
			box.Label = sqlBox.Label.String
		} else {
			box.Label = ""
		}

		if sqlBox.Description.Valid {
			box.Description = sqlBox.Description.String
		} else {
			box.Description = ""
		}

		if sqlBox.Picture.Valid {
			box.Picture = sqlBox.Picture.String
		} else {
			box.Picture = ""
		}

		if sqlBox.PreviewPicture.Valid {
			box.PreviewPicture = sqlBox.PreviewPicture.String
		} else {
			box.PreviewPicture = ""
		}

		if sqlBox.QRcode.Valid {
			box.QRcode = sqlBox.QRcode.String
		} else {
			box.QRcode = ""
		}

		if sqlBox.OuterboxID.Valid {
			box.OuterBoxID, err = uuid.FromString(sqlBox.OuterboxID.String)
			if err != nil {
				return box, logg.Errorf("Error parsing UUID for box: %w", err)
			}
		}
	}

	return box, nil
}

// this function used inside of BoxByField to convert the sql Item struct into normal struct
func convertSQLItemToItem(sqlItem *SqlItem) (*items.Item, error) {
	item := &items.Item{}

	// Convert and assign the ID
	if sqlItem.ItemID.Valid {
		var err error
		item.ID, err = uuid.FromString(sqlItem.ItemID.String)
		if err != nil {
			return nil, logg.Errorf("Error parsing UUID for item: %w", err)
		}
		// Assign other fields only if they are valid
		if sqlItem.ItemLabel.Valid {
			item.Label = sqlItem.ItemLabel.String
		} else {
			item.Label = ""
		}

		if sqlItem.ItemDescription.Valid {
			item.Description = sqlItem.ItemDescription.String
		} else {
			item.Description = ""
		}

		if sqlItem.ItemPicture.Valid {
			item.Picture = sqlItem.ItemPicture.String
		} else {
			item.Picture = ""
		}

		if sqlItem.ItemQRCode.Valid {
			item.QRcode = sqlItem.ItemQRCode.String
		} else {
			item.QRcode = ""
		}

		if sqlItem.ItemBoxID.Valid {
			var err error
			item.BoxID, err = uuid.FromString(sqlItem.ItemBoxID.String)
			if err != nil {
				return nil, logg.Errorf("Error parsing UUID for box ID: '%v' %w", sqlItem.ItemBoxID, err)
			}
		} else {
			return nil, logg.NewError(fmt.Sprintf("box ID is required but was null in item %v", sqlItem.ItemBoxID.String))
		}

		if sqlItem.ItemQuantity.Valid {
			item.Quantity = sqlItem.ItemQuantity.Int64
		} else {
			item.Quantity = 1
		}

		if sqlItem.ItemWeight.Valid {
			item.Weight = sqlItem.ItemWeight.String
		} else {
			item.Weight = ""
		}

	} else {
		return item, logg.NewError("invalid")

	}
	return item, nil
}

// Helper function to check for null strings and return empty if null
func ifNullString(sqlStr sql.NullString) string {
	if sqlStr.Valid {
		return sqlStr.String
	}
	return ""
}

// Helper function to check for null UUIDs and return uuid.Nil if null
func ifNullUUID(sqlUUID sql.NullString) uuid.UUID {
	if sqlUUID.Valid {
		return uuid.FromStringOrNil(sqlUUID.String)
	}
	return uuid.Nil
}

func UUIDFromSqlString(boxID sql.NullString) (uuid.UUID, error) {
	if boxID.Valid {
		id, err := uuid.FromString(boxID.String)
		if err != nil {
			return uuid.Nil, logg.Errorf("error while converting the string id into uuid: %w", err)
		}
		return id, nil
	}
	return uuid.Nil, logg.Errorf("invalid Virtual Id string")
}
