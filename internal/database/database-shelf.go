package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/shelves"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type SQLShelf struct {
	SQLBasicInfo
	Height sql.NullFloat64
	Width  sql.NullFloat64
	Depth  sql.NullFloat64
	Rows   sql.NullInt64
	Cols   sql.NullInt64
	AreaID sql.NullString
}

func (s SQLShelf) ToShelf() (*shelves.Shelf, error) {
	info, err := s.SQLBasicInfo.ToBasicInfo()
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	return &shelves.Shelf{
		ID:             info.ID,
		Label:          info.Label,
		Description:    info.Description,
		Picture:        info.Picture,
		PreviewPicture: info.PreviewPicture,
		QRcode:         info.QRcode,
		Height:         float32(s.Height.Float64),
		Width:          float32(s.Width.Float64),
		Depth:          float32(s.Depth.Float64),
		Rows:           int(s.Rows.Int64),
		Cols:           int(s.Cols.Int64),
		Items:          nil,
		Boxes:          nil,
		AreaID:         ifNullUUID(s.AreaID),
	}, nil
}

// CreateNewShelf creates a new empty shelf in database.
func (db *DB) CreateNewShelf() (uuid.UUID, error) {
	nID, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, logg.WrapErr(err)
	}
	err = db.createNewShelf(nID)
	if err != nil {
		return uuid.Nil, logg.WrapErr(err)
	}
	return nID, err
}

func (db *DB) createNewShelf(nID uuid.UUID) error {
	exists, err := db.Exists("shelf", nID)

	if err != nil {
		return logg.WrapErr(err)
	}
	if exists {
		return db.ErrorExist()
	}

	var (
		id             string = nID.String()
		label          string
		description    sql.NullString
		picture        []byte
		previewPicture []byte
		qrcode         sql.NullString
		height         sql.NullFloat64
		width          sql.NullFloat64
		depth          sql.NullFloat64
		rows           sql.NullInt64
		cols           sql.NullInt64
	)

	stmt := `INSERT INTO shelf (
		id,
		label,
		description,
		picture,
		preview_picture,
		qrcode,
		height,
		width,
		depth,
		rows,
		cols) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// logg.Debugf("SQL: %s", stmt)

	result, err := db.Sql.Exec(stmt,
		&id,
		&label,
		&description,
		&picture,
		&previewPicture,
		&qrcode,
		&height,
		&width,
		&depth,
		&rows,
		&cols,
	)

	if err != nil {
		return logg.WrapErr(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return logg.WrapErr(err)
	}
	if rowsAffected != 1 {
		return logg.Errorf("unexpected number of effected rows, check insirtNewBox")
	}

	return nil
}

// CreateShelf creates a shelf entry in database from the provided shelf.
func (db *DB) CreateShelf(shelf *shelves.Shelf) error {
	if shelf.ID == uuid.Nil {
		id := uuid.Must(uuid.NewV4())
		logg.Infof(`CreateShelf: Provided shelf "%s" had invalid uuid "%s". Creating a new one: "%s"`, shelf.Label, shelf.ID.String(), id.String())
		shelf.ID = id
	}

	if shelf.Items != nil || shelf.Boxes != nil {
		return logg.NewError(fmt.Sprintf(`shelf "%s" has %d items and %d boxes, they must be empty while creating a new shelf`, shelf.ID, len(shelf.Items), len(shelf.Boxes)))
	}

	err := updatePicture(&shelf.Picture, &shelf.PreviewPicture)
	if err != nil {
		logg.Infof("Can't update picture %v", err.Error())
	}

	stmt := `
        INSERT INTO shelf (
            id,
            label,
            description,
            picture,
            preview_picture,
            qrcode,
            height,
            width,
            depth,
            rows,
            cols
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err = db.Sql.Exec(stmt,
		shelf.ID.String(),
		shelf.Label,
		shelf.Description,
		shelf.Picture,
		shelf.PreviewPicture,
		shelf.QRcode,
		shelf.Height,
		shelf.Width,
		shelf.Depth,
		shelf.Rows,
		shelf.Cols,
	)
	if err != nil {
		return logg.Errorf("CreateShelf %w", err)
	}

	return nil
}

// innerListRowsFrom returns all items/boxes/shelves/etc belonging to another box, shelf or area.
//
// Example:
//
//	// get all items that belongs to a shelf.
//	innerListRowsFrom("shelf", shelf.ID, "item_fts")
//
// listRowsTable:
//
//	FROM "item_fts"
//	FROM "box_ftx"
//	...
//
// belongsToTable:
//
//	"item", "box", "shelf", ...
//
// belongsToTableID:
//
//	WHERE "item"_id = ID
//	WHERE "box"_id = ID
//	...
func (db *DB) innerListRowsFrom(belongsToTable string, belongsToTableID uuid.UUID, listRowsTable string) ([]*items.ListRow, error) {
	err := ValidVirtualTable(listRowsTable)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	err = ValidTable(belongsToTable)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	var listRows []*items.ListRow

	stmt := fmt.Sprintf(`
	SELECT id, label, box_id, box_label, shelf_id, shelf_label, area_id, area_label
	FROM %s	
	WHERE %s_id = ?;`, listRowsTable, belongsToTable)

	rows, err := db.Sql.Query(stmt, belongsToTableID.String())
	if err != nil {
		return nil, logg.Errorf("%s %w", stmt, err)
	}
	defer rows.Close()
	for rows.Next() {
		var sqlrow SQLListRow

		err = rows.Scan(
			&sqlrow.ID, &sqlrow.Label, &sqlrow.BoxID, &sqlrow.BoxLabel, &sqlrow.ShelfID, &sqlrow.ShelfLabel, &sqlrow.AreaID, &sqlrow.AreaLabel,
		)
		if err != nil {
			return nil, logg.WrapErr(err)
		}
		lrow, err := sqlrow.ToListRow()
		if err != nil {
			return nil, logg.WrapErr(err)
		}
		listRows = append(listRows, lrow)
	}
	return listRows, nil
}

// Shelf returns shelf with provided id.
func (db *DB) Shelf(id uuid.UUID) (*shelves.Shelf, error) {
	var sqlShelf SQLShelf
	stmt := `
	SELECT
		id, label, description, picture, preview_picture, qrcode, height, width, depth, rows, cols
	FROM 
		shelf
	WHERE 
		id = ?;`

	err := db.Sql.QueryRow(stmt, id.String()).Scan(
		&sqlShelf.SQLBasicInfo.ID, &sqlShelf.SQLBasicInfo.Label, &sqlShelf.SQLBasicInfo.Description, &sqlShelf.SQLBasicInfo.Picture, &sqlShelf.SQLBasicInfo.PreviewPicture, &sqlShelf.SQLBasicInfo.QRCode, &sqlShelf.Height, &sqlShelf.Width, &sqlShelf.Depth, &sqlShelf.Rows, &sqlShelf.Cols,
	)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	shelf, err := sqlShelf.ToShelf()
	if err != nil {
		return nil, logg.WrapErr(err)
	}
	items, err := db.innerListRowsFrom("shelf", shelf.ID, "item_fts")
	if err != nil {
		return nil, logg.WrapErr(err)
	}
	shelf.Items = items
	boxes, err := db.innerListRowsFrom("shelf", shelf.ID, "box_fts")
	if err != nil {
		return nil, logg.WrapErr(err)
	}
	shelf.Boxes = boxes
	return shelf, nil
}

func (db *DB) UpdateShelf(shelf *shelves.Shelf) error {
	err := updatePicture(&shelf.Picture, &shelf.PreviewPicture)
	if err != nil {
		logg.Errf("Can't update picture %v", err.Error())
	}

	stmt := `
        UPDATE shelf SET
            label = ?,
            description = ?,
            picture = ?,
            preview_picture = ?,
            qrcode = ?,
            height = ?,
            width = ?,
            depth = ?,
            rows = ?,
            cols = ?
        WHERE id = ?
    `
	_, err = db.Sql.Exec(stmt,
		shelf.Label,
		shelf.Description,
		shelf.Picture,
		shelf.PreviewPicture,
		shelf.QRcode,
		shelf.Height,
		shelf.Width,
		shelf.Depth,
		shelf.Rows,
		shelf.Cols,
		shelf.ID.String(),
	)
	if err != nil {
		return logg.Errorf("UpdateShelf: %w", err)
	}
	return nil
}

// DeleteShelf deletes a single shelf.
func (db *DB) DeleteShelf(id uuid.UUID) error {
	shelf, err := db.Shelf(id)
	if err != nil {
		return logg.WrapErr(err)
	}
	if shelf.Items != nil || shelf.Boxes != nil {
		return logg.NewError(fmt.Sprintf(`can't delete shelf "%s": has %d items and %d boxes`, shelf.ID, len(shelf.Items), len(shelf.Boxes)))
	}
	err = db.deleteFrom("shelf", id)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// ShelfListRowsPaginated returns always a slice with the number of rows specified.
// If rows=5 with 3 results, the last 2 rows will be nil.
// rows and page must be above 1.
func (db *DB) ShelfListRowsPaginated(page int, rows int) ([]*items.ListRow, error) {
	shelfRows := make([]*items.ListRow, 0)

	if page < 1 {
		return shelfRows, logg.NewError(fmt.Sprintf("invalid page '%d', only positive page numbers starting from 1 are valid", page))
	}

	if rows < 1 {
		return shelfRows, logg.NewError(fmt.Sprintf("invalid rows '%d', needs at least 1 row", rows))
	}

	shelfRows = make([]*items.ListRow, rows)

	limit := rows
	offset := (page - 1) * rows

	queryNoSearch := `
		SELECT
			id, label, area_id, area_label, preview_picture
		FROM shelf_fts
			ORDER BY label ASC
		LIMIT ? OFFSET ?;`
	results, err := db.Sql.Query(queryNoSearch, limit, offset)
	if err != nil {
		return shelfRows, logg.WrapErr(err)
	}

	i := 0
	for results.Next() {
		listRow := SQLListRow{}
		err = results.Scan(&listRow.ID, &listRow.Label, &listRow.AreaID, &listRow.AreaLabel, &listRow.PreviewPicture)
		if err != nil {
			return shelfRows, logg.WrapErr(err)
		}
		shelfRows[i], err = listRow.ToListRow()
		if err != nil {
			return shelfRows, logg.WrapErr(err)
		}
		i += 1
	}
	if results.Err() != nil {
		return shelfRows, logg.WrapErr(err)
	}

	return shelfRows, nil
}

func (db *DB) ShelfSearchListRowsPaginated(page int, rows int, search string) (shelfRows []*items.ListRow, found int, err error) {
	// shelfRows = make([]*items.ListRow, 0)
	// found = 0
	// err = nil
	if page < 1 {
		return shelfRows, found, logg.NewError(fmt.Sprintf("invalid page '%d', only positive page numbers starting from 1 are valid", page))
	}

	if rows < 1 {
		return shelfRows, found, logg.NewError(fmt.Sprintf("invalid rows '%d', needs at least 1 row", rows))
	}

	shelfRows = make([]*items.ListRow, rows)

	limit := rows
	offset := (page - 1) * rows

	// queryNoSearch := `
	// 	SELECT
	// 		id, label, area_id, area_label, preview_picture
	// 	FROM shelf_fts
	// 		ORDER BY label ASC
	// 	LIMIT ? OFFSET ?;`
	searchTrimmed := strings.TrimSpace(search)
	re := regexp.MustCompile(`\s+`)
	searchModified := re.ReplaceAllString(searchTrimmed, "*")
	querySearch := fmt.Sprintf(`
		SELECT
			id, label, area_id, area_label, preview_picture
		FROM shelf_fts
		WHERE label MATCH '%s*'
		LIMIT ? OFFSET ?;`, searchModified)

	// ORDER BY label ASC
	// results, err := db.Sql.Query(querySearch, fmt.Sprintf("'%s'", search), limit, offset)
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

// MoveShelfToArea moves shelf to an area.
// To move out of an area set "toAreaID = uuid.Nil".
func (db *DB) MoveShelfToArea(shelfID uuid.UUID, toAreaID uuid.UUID) error {
	return logg.WrapErr(ErrNotImplemented)
}
