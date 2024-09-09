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
	BoxID          uuid.UUID
	Label          sql.NullString
	OuterBoxLabel  sql.NullString
	OuterBoxID     sql.NullString
	ShelveLabel    sql.NullString
	AreaLabel      sql.NullString
	PreviewPicture sql.NullString
}

func (db *DB) SearchItemsByLabel(query string) ([]struct {
	Id    string
	Label string
}, error) {
	rows, err := db.Sql.Query(` SELECT item_id, label FROM item_fts WHERE label LIKE ? ORDER BY item_id; `, query+"%")
	if err != nil {
		return nil, fmt.Errorf("error while fetching the items for search engine: %w", err)
	}
	defer rows.Close()

	var results []struct {
		Id    string
		Label string
	}
	for rows.Next() {
		var result struct {
			Id    string
			Label string
		}
		if err := rows.Scan(&result.Id, &result.Label); err != nil {
			return nil, fmt.Errorf("error while fetching the items for search engine: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
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
		virtualBoxes = append(virtualBoxes, mapSqlVertualBoxToVertualBox(sqlVertualBox))
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
	// fmt.Printf("virtuaLBox: \n box_id: %s  box label: %s \n outerbox id; %s outerbox label %s \n ",
	//             sqlVertualBox.BoxID.String(), sqlVertualBox.Label.String,
	//             sqlVertualBox.OuterBoxID.String, sqlVertualBox.OuterBoxLabel.String)

	return mapSqlVertualBoxToVertualBox(sqlVertualBox), nil
}

// private function to map the sql virtual box into normal virtual box
func mapSqlVertualBoxToVertualBox(sqlBox SqlVertualBox) items.VirtualBox {
	return items.VirtualBox{
		Box_Id:         sqlBox.BoxID,
		Label:          ifNullString(sqlBox.Label),
		OuterBox_label: ifNullString(sqlBox.OuterBoxLabel),
		OuterBox_id:    ifNullUUID(sqlBox.OuterBoxID),
		Shelve_label:   "Shelve 1",
		Area_label:     "Area 1",
		PreviewPicture: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HwAFAgL/uXBuZwAAAABJRU5ErkJggg==",
	}
}
