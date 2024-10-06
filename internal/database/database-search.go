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
	ShelfLabel     sql.NullString
	AreaLabel      sql.NullString
	PreviewPicture sql.NullString
}

type SqlVirtualItem struct {
	ItemID         sql.NullString
	Label          sql.NullString
	OuterBoxLabel  sql.NullString
	OuterBoxID     sql.NullString
	ShelfLabel     sql.NullString
	AreaLabel      sql.NullString
	PreviewPicture sql.NullString
}

// Search items based on search query, return array of virtualItems
func (db *DB) ItemFuzzyFinder(query string) ([]items.ItemListRow, error) {
	rows, err := db.Sql.Query(` SELECT item_id, label FROM item_fts WHERE label LIKE ? ORDER BY item_id; `, query+"%")
	if err != nil {
		return nil, fmt.Errorf("error while fetching virtual items: %w", err)
	}
	defer rows.Close()

	var virtualItems []items.ItemListRow
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

// Search items based on search query, return limited number of results used to generate pagination
func (db *DB) ItemFuzzyFinderWithPagination(query string, limit, offset int) ([]items.ItemListRow, error) {
	rows, err := db.Sql.Query(`
        SELECT item_id, label 
        FROM item_fts 
        WHERE label LIKE ? 
        ORDER BY item_id 
        LIMIT ? OFFSET ?; 
    `, query+"%", limit, offset)

	if err != nil {
		return nil, fmt.Errorf("error while fetching virtual items: %w", err)
	}
	defer rows.Close()

	var virtualItems []items.ItemListRow
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
func (db *DB) BoxFuzzyFinder(query string, limit int, page int) ([]items.BoxListRow, error) {
	var rows *sql.Rows
	var err error
	if page == 0 {
		panic("page starts at 1, cant be 0")
	}

	if strings.TrimSpace(query) == "" {
		rows, err = db.Sql.Query(`SELECT box_id, label, outerbox_id, outerbox_label 
                                  FROM box_fts 
                                  ORDER BY label ASC
				LIMIT ? OFFSET ?;`, limit, (page-1)*limit)

	} else {
		rows, err = db.Sql.Query(`SELECT box_id, label, outerbox_id, outerbox_label 
                                  FROM box_fts 
                                  WHERE label LIKE ? 
                                  ORDER BY label ASC
				LIMIT ? OFFSET ?;`, query+"%", limit, (page-1)*limit)
	}

	if err != nil {
		return []items.BoxListRow{}, fmt.Errorf("error while fetching the virtualBox from box_fts: %w", err)
	}
	defer rows.Close()

	var sqlVertualBox SqlVertualBox
	var virtualBoxes []items.BoxListRow

	for rows.Next() {
		err := rows.Scan(
			&sqlVertualBox.BoxID,
			&sqlVertualBox.Label,
			&sqlVertualBox.OuterBoxID,
			&sqlVertualBox.OuterBoxLabel,
		)
		if err != nil {
			return []items.BoxListRow{}, logg.Errorf("error while assigning the Data to the Virtualbox struct %w", err)
		}
		vBox, err := mapSqlVertualBoxToVertualBox(sqlVertualBox)
		if err != nil {
			return []items.BoxListRow{}, err
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
func (db *DB) VirtualBoxById(id uuid.UUID) (items.BoxListRow, error) {

	if !db.VirtualBoxExist(id) {
		return items.BoxListRow{}, fmt.Errorf("the Box Id does not exsist in the virtual table")
	}

	query := fmt.Sprintf("SELECT box_id, label, outerbox_id, outerbox_label FROM box_fts WHERE box_id = ?")
	row, err := db.Sql.Query(query, id.String())
	if err != nil {
		return items.BoxListRow{}, fmt.Errorf("error while fetching the virtual box: %w", err)
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
			return items.BoxListRow{}, fmt.Errorf("error while assigning the Data to the Virtualbox struct : %w", err)
		}
	}

	vBox, err := mapSqlVertualBoxToVertualBox(sqlVertualBox)
	if err != nil {
		return items.BoxListRow{}, err
	}
	return vBox, nil
}

// private function to map the sql virtual box into normal virtual box
func mapSqlVertualBoxToVertualBox(sqlBox SqlVertualBox) (items.BoxListRow, error) {
	id, err := UUIDFromSqlString(sqlBox.BoxID)
	if err != nil {
		return items.BoxListRow{}, err
	}

	return items.BoxListRow{
		BoxID:          id,
		Label:          ifNullString(sqlBox.Label),
		OuterBoxLabel:  ifNullString(sqlBox.OuterBoxLabel),
		OuterBoxID:     ifNullUUID(sqlBox.OuterBoxID),
		ShelfLabel:     "Shelf 1",
		AreaLabel:      "Area 1",
		PreviewPicture: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HwAFAgL/uXBuZwAAAABJRU5ErkJggg==",
	}, nil
}

func mapSqlVertualItemToVertualItem(sqlItem SqlVirtualItem) (items.ItemListRow, error) {
	id, err := UUIDFromSqlString(sqlItem.ItemID)
	if err != nil {
		return items.ItemListRow{}, err
	}
	return items.ItemListRow{
		ItemID:         id,
		Label:          ifNullString(sqlItem.Label),
		BoxLabel:       "box 1",
		BoxID:          uuid.Nil,
		ShelfLabel:     "shelf 1",
		AreaLabel:      "shelf 1",
		PreviewPicture: "pic1",
	}, nil
}

func (db *DB) NumOfItemRecords(searchString string) (int, error) {
	searchString = strings.TrimSpace(searchString)

	var query string
	if searchString == "" {
		query = "SELECT COUNT(*) FROM item_fts;"
	} else {
		query = fmt.Sprintf("SELECT COUNT(*) FROM item_fts WHERE label LIKE ?;")
	}

	var count int
	var err error
	if searchString == "" {
		err = db.Sql.QueryRow(query).Scan(&count)
	} else {
		err = db.Sql.QueryRow(query, searchString+"%").Scan(&count)
	}

	if err != nil {
		return 0, fmt.Errorf("Error checking the number of records %v:", err)
	}
	logg.Debugf("count: %d \n", count)
	return count, nil
}
