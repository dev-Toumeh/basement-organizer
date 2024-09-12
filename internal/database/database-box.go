package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"database/sql"
	"fmt"
	"log"

	"github.com/gofrs/uuid/v5"
)

type Exist func(filename string) bool

type SQLBox struct {
	Id          sql.NullString
	Label       sql.NullString
	Description sql.NullString
	Picture     sql.NullString
	QRcode      sql.NullString
	OuterboxId  sql.NullString
}

type SqlItem struct {
	ItemId          sql.NullString `json:"id"`
	ItemLabel       sql.NullString `json:"label"`
	ItemDescription sql.NullString `json:"description"`
	ItemPicture     sql.NullString `json:"picture"`
	ItemQRCode      sql.NullString `json:"qrcode"`
	ItemBoxId       sql.NullString `json:"box_id"`
	ItemQuantity    sql.NullInt64  `json:"quantity"`
	ItemWeight      sql.NullString `json:"weight"`
}

// Create New Item Record
func (db *DB) CreateBox(newBox *items.Box) (uuid.UUID, error) {
	if db.BoxExist("label", newBox.Label) {
		return uuid.Nil, db.ErrorExist()
	}

	id, err := db.insertNewBox(newBox)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while creating new Box: %v", err)
	}
	return id, nil
}

// check if the Box Exist based on given Field
func (db *DB) BoxExist(field string, value string) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM box WHERE %s = ?", field)
	var count int
	err := db.Sql.QueryRow(query, value).Scan(&count)
	if err != nil {
		log.Println("Error checking item existence:", err)
		return false
	}
	// fmt.Printf("count is %d", count)
	return count > 0
}

// BoxIDs returns IDs of all boxes.
func (db *DB) BoxIDs() ([]string, error) {
	ids := []string{}
	sqlStatement := `SELECT id FROM BOX`
	rows, err := db.Sql.Query(sqlStatement)
	if err != nil {
		return ids, fmt.Errorf("Error while executing Box ids: %w", err)
	}
	for rows.Next() {
		var idStr string
		err := rows.Scan(&idStr)
		if err != nil {
			return []string{}, fmt.Errorf("Error scanning Box ids: %v", err)
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
	exist := db.BoxExist("id", box.Id.String())
	if !exist {
		return fmt.Errorf("the box is not exist")
	}
	sqlStatement := fmt.Sprintf(`UPDATE box SET label = '%s', description = '%s', picture = '%s', qrcode = '%s' WHERE id = ?`, box.Label, box.Description, box.Picture, box.QRcode)
	result, err := db.Sql.Exec(sqlStatement, box.Id)
	if err != nil {
		return fmt.Errorf("something wrong happened while runing the box update query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while finding the item %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("the Record with the id: %s was not found; this should not have happened while updating", box.Id.String())
	} else if rowsAffected != 1 {
		return fmt.Errorf("the id: %s has an unexpected number of rows affected (more than one or less than 0)", box.Id.String())
	}
	return nil
}

// delete Box
func (db *DB) DeleteBox(boxId uuid.UUID) error {
	id := boxId.String()

	// check if box is not Empty
	itemExist := db.ItemExist("box_id", id)
	boxExist := db.BoxExist("outerbox_id", id)
	if itemExist || boxExist {
		return fmt.Errorf("the box is not empty")
	}

	sqlStatement := `DELETE FROM box WHERE id = ?;`
	_, err := db.Sql.Exec(sqlStatement, id)
	if err != nil {
		return fmt.Errorf("error whiel deleting the box: %W", err)
	}
	return nil
}

// Get Box based on his ID
// Wrapper function for BoxByField
func (db *DB) BoxById(id uuid.UUID) (items.Box, error) {
	box := items.Box{}
	b, err := db.BoxByField("id", id.String())
	if err != nil {
		return box, fmt.Errorf("Box() error:\n\t%w", err)
	}
	box = *b
	return box, err
}

// Get Box  based on given Field
func (db *DB) BoxByField(field string, value string) (*items.Box, error) {
	var box *items.Box
	var outerBox = &items.Box{}
	var itemsList = []*items.Item{}
	var inBoxList = []*items.Box{}
	boxInitialized := false
	outerBoxInitialized := false
	addedInBoxes := make(map[string]bool)
	addedItems := make(map[string]*items.Item)

	query := fmt.Sprintf(
		`SELECT 
            b.id, b.label, b.description, b.picture, b.qrcode, b.outerbox_id,
            ob.id, ob.label, ob.description, ob.picture, ob.qrcode, ob.outerbox_id,
            ib.id, ib.label, ib.description, ib.picture, ib.qrcode, ib.outerbox_id,
            i.id, i.label, i.description, i.picture, i.quantity, i.weight, i.qrcode, i.box_id
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
		return nil, fmt.Errorf("error executing query: %w", err)
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
			&boxSQLBox.Id, &boxSQLBox.Label, &boxSQLBox.Description, &boxSQLBox.Picture, &boxSQLBox.QRcode, &boxSQLBox.OuterboxId,
			&outboxSQLBox.Id, &outboxSQLBox.Label, &outboxSQLBox.Description, &outboxSQLBox.Picture, &outboxSQLBox.QRcode, &outboxSQLBox.OuterboxId,
			&inboxSQLBox.Id, &inboxSQLBox.Label, &inboxSQLBox.Description, &inboxSQLBox.Picture, &inboxSQLBox.QRcode, &inboxSQLBox.OuterboxId,
			&itemSQLItem.ItemId, &itemSQLItem.ItemLabel, &itemSQLItem.ItemDescription, &itemSQLItem.ItemPicture, &itemSQLItem.ItemQuantity, &itemSQLItem.ItemWeight, &itemSQLItem.ItemQRCode, &itemSQLItem.ItemBoxId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Initialize the main Box only once
		if !boxInitialized {
			box, err = convertSQLBoxToBox(&boxSQLBox)
			if err != nil {
				return nil, fmt.Errorf("error assigning the box data: %w", err)
			}
			boxInitialized = true
		}

		// Initialize outer box only once
		if !outerBoxInitialized && outboxSQLBox.Id.Valid {
			outerBox, err = convertSQLBoxToBox(&outboxSQLBox)
			if err != nil {
				return &items.Box{}, fmt.Errorf("something wrong happened while assigning the outboxSQLBox data %w", err)
			}
			outerBoxInitialized = true
		}

		// Add the inboxes if its ID is valid and not added before
		if inboxSQLBox.Id.Valid && !addedInBoxes[inboxSQLBox.Id.String] {
			inbox, err = convertSQLBoxToBox(&inboxSQLBox)
			if err != nil {
				return nil, fmt.Errorf("error converting SQLBox to Box: %w", err)
			}
			inBoxList = append(inBoxList, inbox)
			addedInBoxes[inboxSQLBox.Id.String] = true
		}

		// Add the item to the itemsList if itemId is valid
		if itemSQLItem.ItemId.Valid {
			if _, exists := addedItems[itemSQLItem.ItemId.String]; !exists {
				item, err := convertSQLItemToItem(&itemSQLItem)
				if err != nil {
					return nil, fmt.Errorf("error converting SQLItem to Item: %w", err)
				}
				itemsList = append(itemsList, item)
				addedItems[itemSQLItem.ItemId.String] = item
			}
		}

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	// Assign the items list to the box
	box.OuterBox = outerBox
	box.Items = itemsList
	box.InnerBoxes = inBoxList

	return box, nil
}

// insert new Box record in the Database
func (db *DB) insertNewBox(box *items.Box) (uuid.UUID, error) {
	if db.BoxExist("id", box.Id.String()) {
		return uuid.Nil, db.ErrorExist()
	}

	sqlStatement := `INSERT INTO box (id, label, description, picture, qrcode, outerbox_id) VALUES (?, ?, ?, ?, ?, ?)`
	logg.Debugf("fuck %s", sqlStatement)
	result, err := db.Sql.Exec(sqlStatement, box.Id.String(), box.Label, box.Description, box.Picture, box.QRcode, box.OuterBoxId.String())
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error while executing create new box statement: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error while executing create new box statement: %w", err)
	}
	if rowsAffected != 1 {
		return uuid.Nil, fmt.Errorf("unexpected number of effected rows, check insirtNewBox")
	}

	return box.Id, nil
}

// this function used inside of BoxByField to convert the box sql struct into normal struct
func convertSQLBoxToBox(sqlBox *SQLBox) (*items.Box, error) {
	box := &items.Box{}

	// Convert and assign the ID
	if sqlBox.Id.Valid {
		var err error
		box.Id, err = uuid.FromString(sqlBox.Id.String)
		if err != nil {
			log.Printf("Error parsing UUID for box: %v", err)
			return box, err
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

		if sqlBox.QRcode.Valid {
			box.QRcode = sqlBox.QRcode.String
		} else {
			box.QRcode = ""
		}

		if sqlBox.OuterboxId.Valid {
			box.OuterBoxId, err = uuid.FromString(sqlBox.OuterboxId.String)
			if err != nil {
				log.Printf("Error parsing UUID for box: %v", err)
				return box, err
			}
		}

	}

	return box, nil
}

// this function used inside of BoxByField to convert the sql Item struct into normal struct
func convertSQLItemToItem(sqlItem *SqlItem) (*items.Item, error) {
	item := &items.Item{}

	// Convert and assign the ID
	if sqlItem.ItemId.Valid {
		var err error
		item.Id, err = uuid.FromString(sqlItem.ItemId.String)
		if err != nil {
			log.Printf("Error parsing UUID for item: %v", err)
			return nil, err
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

		if sqlItem.ItemBoxId.Valid {
			var err error
			item.BoxId, err = uuid.FromString(sqlItem.ItemBoxId.String)
			if err != nil {
				log.Printf("Error parsing UUID for box ID: %v", err)
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("box ID is required but was null")
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
			return uuid.Nil, fmt.Errorf("error while converting the string id into uuid: %w", err)
		}
		return id, nil
	}
	return uuid.Nil, fmt.Errorf("invalid VirtualItem Id string")
}
