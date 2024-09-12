package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type SqlVertualBox struct {
	BoxID          sql.NullString
	Label          sql.NullString
	OuterBoxLabel  sql.NullString
	OuterBoxID     sql.NullString
	ShelveLabel    sql.NullString
	AreaLabel      sql.NullString
	PreviewPicture sql.NullString
}

type SqlVirtualItem struct {
	ItemID         sql.NullString
	Label          sql.NullString
	OuterBoxLabel  sql.NullString
	OuterBoxID     sql.NullString
	ShelveLabel    sql.NullString
	AreaLabel      sql.NullString
	PreviewPicture sql.NullString
}

func (db *DB) ItemFuzzyFinder(query string) ([]items.VirtualItem, error) {
	rows, err := db.Sql.Query(` SELECT item_id, label FROM item_fts WHERE label LIKE ? ORDER BY item_id; `, query+"%")
	if err != nil {
		return nil, fmt.Errorf("error while fetching virtual items: %w", err)
	}
	defer rows.Close()

	var virtualItems []items.VirtualItem
	var sqlItem SqlVirtualItem

	for rows.Next() {
		err = rows.Scan(&sqlItem.ItemID, &sqlItem.Label)
		if err != nil {
			return nil, fmt.Errorf("error while assigning the Data to the VirtualItem struct: %w", err)
		}
		vItem, err := mapSqlVertualItemToVertualItem(sqlItem)
		if err != nil {
			return nil, err
		}
		virtualItems = append(virtualItems, vItem)
	}

	return virtualItems, nil
}

// BoxFuzzyFinder retrieves virtual boxes by label.
// If the query is empty or contains only spaces, it returns 10 default results.
func (db *DB) BoxFuzzyFinder(query string) ([]items.VirtualBox, error) {
	var rows *sql.Rows
	var err error

	if strings.TrimSpace(query) == "" {
		rows, err = db.Sql.Query(`SELECT box_id, label, outerbox_id, outerbox_label 
                                  FROM box_fts 
                                  ORDER BY label ASC 
                                  LIMIT 10;`)
	} else {
		rows, err = db.Sql.Query(`SELECT box_id, label, outerbox_id, outerbox_label 
                                  FROM box_fts 
                                  WHERE label LIKE ? 
                                  ORDER BY label ASC;`, query+"%")
	}

	if err != nil {
		return []items.VirtualBox{}, fmt.Errorf("error while fetching the virtualBox from box_fts: %w", err)
	}
	defer rows.Close()

	var sqlVertualBox SqlVertualBox
	var virtualBoxes []items.VirtualBox

	for rows.Next() {
		err := rows.Scan(
			&sqlVertualBox.BoxID,
			&sqlVertualBox.Label,
			&sqlVertualBox.OuterBoxID,
			&sqlVertualBox.OuterBoxLabel,
		)
		if err != nil {
			return []items.VirtualBox{}, fmt.Errorf("error while assigning the Data to the Virtualbox struct: %w", err)
		}
		vBox, err := mapSqlVertualBoxToVertualBox(sqlVertualBox)
		if err != nil {
			return []items.VirtualBox{}, err
		}
		virtualBoxes = append(virtualBoxes, vBox)
	}

	return virtualBoxes, nil
}

// check if the virtual box is empty
func (db *DB) VirtualBoxExist(id uuid.UUID) bool {
	query := fmt.Sprint("SELECT COUNT(*) FROM box_fts WHERE box_id = ?; ")
	var count int
	err := db.Sql.QueryRow(query, id.String()).Scan(&count)
	if err != nil {
		logg.Errf("Error checking item existence %v:", err)
		return false
	}
	return count > 0
}

// Get the virtual Box based on his ID
func (db *DB) VirtualBoxById(id uuid.UUID) (items.VirtualBox, error) {

	if !db.VirtualBoxExist(id) {
		return items.VirtualBox{}, fmt.Errorf("the Box Id does not exsist in the virtual table")
	}

	query := fmt.Sprintf("SELECT box_id, label, outerbox_id, outerbox_label FROM box_fts WHERE box_id = ?")
	row, err := db.Sql.Query(query, id.String())
	if err != nil {
		return items.VirtualBox{}, fmt.Errorf("error while fetching the virtual box: %w", err)
	}

	var sqlVertualBox SqlVertualBox
	for row.Next() {
		err := row.Scan(
			&sqlVertualBox.BoxID,
			&sqlVertualBox.Label,
			&sqlVertualBox.OuterBoxID,
			&sqlVertualBox.OuterBoxLabel,
		)
		if err != nil {
			return items.VirtualBox{}, fmt.Errorf("error while assigning the Data to the Virtualbox struct : %w", err)
		}
	}

	vBox, err := mapSqlVertualBoxToVertualBox(sqlVertualBox)
	if err != nil {
		return items.VirtualBox{}, err
	}
	return vBox, nil
}

// private function to map the sql virtual box into normal virtual box
func mapSqlVertualBoxToVertualBox(sqlBox SqlVertualBox) (items.VirtualBox, error) {
	id, err := UUIDFromSqlString(sqlBox.BoxID)
	if err != nil {
		return items.VirtualBox{}, err
	}

	return items.VirtualBox{
		Box_Id:         id,
		Label:          ifNullString(sqlBox.Label),
		OuterBox_label: ifNullString(sqlBox.OuterBoxLabel),
		OuterBox_id:    ifNullUUID(sqlBox.OuterBoxID),
		Shelve_label:   "Shelve 1",
		Area_label:     "Area 1",
		PreviewPicture: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HwAFAgL/uXBuZwAAAABJRU5ErkJggg==",
	}, nil
}

func mapSqlVertualItemToVertualItem(sqlItem SqlVirtualItem) (items.VirtualItem, error) {
	id, err := UUIDFromSqlString(sqlItem.ItemID)
	if err != nil {
		return items.VirtualItem{}, err
	}
	return items.VirtualItem{
		Item_Id:        id,
		Label:          ifNullString(sqlItem.Label),
		Box_label:      "box 1",
		Box_id:         uuid.Nil,
		Shelve_label:   "shelve 1",
		Area_label:     "shelve 1",
		PreviewPicture: "pic1",
	}, nil
}

// fmt.Printf("virtuaLBox: \n box_id: %s  box label: %s \n outerbox id; %s outerbox label %s \n ",
//             sqlVertualBox.BoxID.String(), sqlVertualBox.Label.String,
//             sqlVertualBox.OuterBoxID.String, sqlVertualBox.OuterBoxLabel.String)
