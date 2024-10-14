package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// Search items based on search query, return array of virtualItems
func (db *DB) ItemFuzzyFinder(query string) ([]items.ListRow, error) {
	rows, err := db.Sql.Query(` SELECT item_id, label FROM item_fts WHERE label LIKE ? ORDER BY item_id; `, query+"%")
	if err != nil {
		return nil, fmt.Errorf("error while fetching virtual items: %w", err)
	}
	defer rows.Close()

	var virtualItems []items.ListRow
	var sqlItem SQLListRow

	for rows.Next() {
		err = rows.Scan(&sqlItem.ID, &sqlItem.Label)
		if err != nil {
			return nil, fmt.Errorf("error while assigning the Data to the VirtualItem struct: %w", err)
		}
		vItem, err := sqlItem.ToListRow()
		if err != nil {
			return nil, err
		}
		virtualItems = append(virtualItems, *vItem)
	}

	return virtualItems, nil
}

// Search items based on search query, return limited number of results used to generate pagination
func (db *DB) ItemFuzzyFinderWithPagination(query string, limit, offset int) ([]items.ListRow, error) {
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

	var virtualItems []items.ListRow
	var sqlItem SQLListRow

	for rows.Next() {
		err = rows.Scan(&sqlItem.ID, &sqlItem.Label)
		if err != nil {
			return nil, fmt.Errorf("error while assigning the Data to the VirtualItem struct: %w", err)
		}
		vItem, err := sqlItem.ToListRow()
		if err != nil {
			return nil, err
		}
		virtualItems = append(virtualItems, *vItem)
	}

	return virtualItems, nil
}

// BoxFuzzyFinder retrieves virtual boxes by label.
// If the query is empty or contains only spaces, it returns 10 default results.
func (db *DB) BoxFuzzyFinder(query string, limit int, page int) ([]items.ListRow, error) {
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
		return []items.ListRow{}, logg.Errorf("error while fetching the virtualBox from box_fts: %w", err)
	}
	defer rows.Close()

	var sqlBoxListRow SQLListRow
	var virtualBoxes []items.ListRow

	for rows.Next() {
		err := rows.Scan(
			&sqlBoxListRow.ID,
			&sqlBoxListRow.Label,
			&sqlBoxListRow.BoxID,
			&sqlBoxListRow.BoxLabel,
			&sqlBoxListRow.PreviewPicture,
		)
		if err != nil {
			return []items.ListRow{}, logg.Errorf("error while assigning the Data to the Virtualbox struct %w", err)
		}
		vBox, err := sqlBoxListRow.ToListRow()
		if err != nil {
			return []items.ListRow{}, logg.WrapErr(err)
		}
		virtualBoxes = append(virtualBoxes, *vBox)
	}

	return virtualBoxes, nil
}

// Search shelves based on search query, return array of virtueShelves
func (db *DB) ShelfFuzzyFinder(query string) ([]items.ListRow, error) {
	rows, err := db.Sql.Query(` SELECT id, label, area_id, area_label, preview_picture tokenize 
                              FROM shelf_fts WHERE label LIKE ? ORDER BY id; `, query+"%")
	if err != nil {
		return nil, fmt.Errorf("error while fetching virtual shelves: %w", err)
	}
	defer rows.Close()

	var virtuaShelves []items.ListRow
	var sqlShelf SQLListRow

	for rows.Next() {
		err = rows.Scan(&sqlShelf.ID, &sqlShelf.Label, &sqlShelf.ID, sqlShelf.AreaLabel, sqlShelf.PreviewPicture)
		if err != nil {
			return nil, fmt.Errorf("error while assigning the Data to the VirtualItem struct: %w", err)
		}
		vShelf, err := sqlShelf.ToListRow()
		if err != nil {
			return nil, err
		}
		virtuaShelves = append(virtuaShelves, *vShelf)
	}

	return virtuaShelves, nil
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
func (db *DB) BoxListRowByID(id uuid.UUID) (items.ListRow, error) {

	if !db.VirtualBoxExist(id) {
		return items.ListRow{}, fmt.Errorf("the Box Id does not exsist in the virtual table")
	}

	query := fmt.Sprintf("SELECT box_id, label, outerbox_id, outerbox_label, shelf_id, shelf_label, area_id, area_label FROM box_fts WHERE box_id = ?")
	row, err := db.Sql.Query(query, id.String())
	if err != nil {
		return items.ListRow{}, fmt.Errorf("error while fetching the virtual box: %w", err)
	}

	var sqlVertualBox SQLListRow
	for row.Next() {
		err := row.Scan(
			&sqlVertualBox.ID,
			&sqlVertualBox.Label,
			&sqlVertualBox.BoxID,
			&sqlVertualBox.BoxLabel,
			&sqlVertualBox.ShelfID,
			&sqlVertualBox.ShelfLabel,
			&sqlVertualBox.AreaID,
			&sqlVertualBox.AreaLabel,
		)
		if err != nil {
			return items.ListRow{}, fmt.Errorf("error while assigning the Data to the Virtualbox struct : %w", err)
		}
	}

	vBox, err := sqlVertualBox.ToListRow()
	if err != nil {
		return items.ListRow{}, err
	}
	return *vBox, nil
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

func (db *DB) ShelfSearchListRowsPaginated(page int, rows int, search string) (shelfRows []*items.ListRow, found int, err error) {

	if page < 1 {
		return shelfRows, found, logg.NewError(fmt.Sprintf("invalid page '%d', only positive page numbers starting from 1 are valid", page))
	}

	if rows < 1 {
		return shelfRows, found, logg.NewError(fmt.Sprintf("invalid rows '%d', needs at least 1 row", rows))
	}

	shelfRows = make([]*items.ListRow, rows)

	limit := rows
	offset := (page - 1) * rows

	searchTrimmed := strings.TrimSpace(search)
	re := regexp.MustCompile(`\s+`)
	searchModified := re.ReplaceAllString(searchTrimmed, "*")
	querySearch := fmt.Sprintf(`
		SELECT
			id, label, area_id, area_label, preview_picture
		FROM shelf_fts
		WHERE label MATCH '%s*'
		LIMIT ? OFFSET ?;`, searchModified)

	results, err := db.Sql.Query(querySearch, limit, offset)
	if err != nil {
		return shelfRows, found, logg.WrapErr(err)
	}

	i := 0
	for results.Next() {
		listRow := SQLListRow{}
		err = results.Scan(&listRow.ID, &listRow.Label, &listRow.AreaID, &listRow.AreaLabel, &listRow.PreviewPicture)
		if err != nil {
			return shelfRows, found, logg.WrapErr(err)
		}
		shelfRows[i], err = listRow.ToListRow()
		if err != nil {
			return shelfRows, found, logg.WrapErr(err)
		}
		i += 1
		found += 1
	}
	if results.Err() != nil {
		return shelfRows, found, logg.WrapErr(err)
	}

	return shelfRows, found, nil
}
