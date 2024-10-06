package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type SQLBoxListRow struct {
	BoxID          sql.NullString
	Label          sql.NullString
	OuterBoxID     sql.NullString
	OuterBoxLabel  sql.NullString
	ShelfID        sql.NullString
	ShelfLabel     sql.NullString
	AreaID         sql.NullString
	AreaLabel      sql.NullString
	PreviewPicture sql.NullString
}

type SqlVirtualItem struct {
	ItemID         sql.NullString
	Label          sql.NullString
	OuterBoxID     sql.NullString
	OuterBoxLabel  sql.NullString
	ShelfID        sql.NullString
	ShelfLabel     sql.NullString
	AreaID         sql.NullString
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
	queryNoSearch := `
		SELECT
			box_id, label, outerbox_id, outerbox_label, preview_picture
		FROM box_fts AS b_fts
		ORDER BY label ASC
		LIMIT ? OFFSET ?;`

	queryWithSearch := `
		SELECT 
			box_id, label, outerbox_id, outerbox_label, preview_picture
		FROM box_fts
		WHERE label LIKE ?
		ORDER BY label ASC
		LIMIT ? OFFSET ?;`

	if strings.TrimSpace(query) == "" {
		rows, err = db.Sql.Query(queryNoSearch, limit, (page-1)*limit)
	} else {
		rows, err = db.Sql.Query(queryWithSearch, query+"%", limit, (page-1)*limit)
	}

	if err != nil {
		return []items.BoxListRow{}, logg.Errorf("error while fetching the virtualBox from box_fts: %w", err)
	}
	defer rows.Close()

	var sqlBoxListRow SQLBoxListRow
	var virtualBoxes []items.BoxListRow

	for rows.Next() {
		err := rows.Scan(
			&sqlBoxListRow.BoxID,
			&sqlBoxListRow.Label,
			&sqlBoxListRow.OuterBoxID,
			&sqlBoxListRow.OuterBoxLabel,
			&sqlBoxListRow.PreviewPicture,
		)
		if err != nil {
			return []items.BoxListRow{}, logg.Errorf("error while assigning the Data to the Virtualbox struct %w", err)
		}
		vBox, err := mapSqlVertualBoxToVertualBox(sqlBoxListRow)
		if err != nil {
			return []items.BoxListRow{}, logg.WrapErr(err)
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
func (db *DB) BoxListRowByID(id uuid.UUID) (items.BoxListRow, error) {

	if !db.VirtualBoxExist(id) {
		return items.BoxListRow{}, fmt.Errorf("the Box Id does not exsist in the virtual table")
	}

	query := fmt.Sprintf("SELECT box_id, label, outerbox_id, outerbox_label FROM box_fts WHERE box_id = ?")
	row, err := db.Sql.Query(query, id.String())
	if err != nil {
		return items.BoxListRow{}, fmt.Errorf("error while fetching the virtual box: %w", err)
	}

	var sqlVertualBox SQLBoxListRow
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
func mapSqlVertualBoxToVertualBox(sqlBox SQLBoxListRow) (items.BoxListRow, error) {
	id, err := UUIDFromSqlString(sqlBox.BoxID)
	if err != nil {
		return items.BoxListRow{}, logg.WrapErr(err)
	}

	return items.BoxListRow{
		BoxID:          id,
		Label:          ifNullString(sqlBox.Label),
		OuterBoxID:     ifNullUUID(sqlBox.OuterBoxID),
		OuterBoxLabel:  ifNullString(sqlBox.OuterBoxLabel),
		ShelfID:        ifNullUUID(sqlBox.ShelfID),
		ShelfLabel:     ifNullString(sqlBox.ShelfLabel),
		AreaID:         ifNullUUID(sqlBox.AreaID),
		AreaLabel:      ifNullString(sqlBox.AreaLabel),
		PreviewPicture: ifNullString(sqlBox.PreviewPicture),
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
