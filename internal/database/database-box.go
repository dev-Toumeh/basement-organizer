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
func (db *DB) CreateNewBox(newBox *items.Box) error {
	if db.BoxExist("id", newBox.Id.String()) {
		return db.ErrorExist()

	}
	_, err := db.insertNewBox(newBox)
	if err != nil {
		return fmt.Errorf("error whie creating new Box: %v", err)
	}
	return nil
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
			item.Quantity = 0 // Default or invalid value
		}

		if sqlItem.ItemWeight.Valid {
			item.Weight = sqlItem.ItemWeight.String
		} else {
			item.Weight = ""
		}

	}
	return item, nil
}
