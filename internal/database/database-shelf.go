package database

import (
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/shelves"
	"database/sql"
	"fmt"

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
		BasicInfo: info,
		Height:    float32(ifNullFloat64(s.Height)),
		Width:     float32(ifNullFloat64(s.Width)),
		Depth:     float32(ifNullFloat64(s.Depth)),
		Rows:      int(ifNullInt(s.Rows)),
		Cols:      int(ifNullInt(s.Cols)),
		Items:     nil,
		Boxes:     nil,
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
	// will never happen, uuid is always checked for nil
	if shelf.ID == uuid.Nil {
		panic(fmt.Sprintf(`CreateShelf: Provided shelf "%s" had invalid uuid "%s". This will never happen, because uuid is always checked for nil.`, shelf.Label, shelf.ID.String()))
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
		shelf.QRCode,
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
		&sqlShelf.SQLBasicInfo.ID, &sqlShelf.SQLBasicInfo.Label, &sqlShelf.SQLBasicInfo.Description,
		&sqlShelf.SQLBasicInfo.Picture, &sqlShelf.SQLBasicInfo.PreviewPicture, &sqlShelf.SQLBasicInfo.QRCode,
		&sqlShelf.Height, &sqlShelf.Width, &sqlShelf.Depth, &sqlShelf.Rows, &sqlShelf.Cols,
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
		shelf.QRCode,
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

// MoveShelfToArea moves shelf to an area.
// To move out of an area set "toAreaID = uuid.Nil".
func (db *DB) MoveShelfToArea(shelfID uuid.UUID, toAreaID uuid.UUID) error {
	return logg.WrapErr(ErrNotImplemented)
}

// ShelfListRowsPaginated returns always a slice with the number of rows specified.
func (db *DB) SearchShelves(limit int, offset int, count int, searchString string) (shelfRows []*common.ListRow, err error) {
	i := 0
	if count > limit {
		shelfRows = make([]*common.ListRow, limit)
	} else {
		shelfRows = make([]*common.ListRow, count)
	}

	// by default use this query
	query := `
		SELECT
			id, label, area_id, area_label, preview_picture
		FROM shelf_fts
			ORDER BY label ASC
		LIMIT ? OFFSET ?;`

	// if the searchString is not empty use this query
	if searchString != "" {
		query = fmt.Sprintf(`
		SELECT
			id, label, area_id, area_label, preview_picture
		FROM shelf_fts
		WHERE label MATCH '%s*'
		LIMIT ? OFFSET ?;`, searchString)
	}

	results, err := db.Sql.Query(query, limit, offset)
	if err != nil {
		return shelfRows, fmt.Errorf("error while fetching shelves Data from the database: %v", err)
	}

	for results.Next() {
		listRow := SQLListRow{}
		err = results.Scan(&listRow.ID, &listRow.Label, &listRow.AreaID, &listRow.AreaLabel, &listRow.PreviewPicture)
		if err != nil {
			return shelfRows, fmt.Errorf("error while scanning the Shelf row during search in the database: %v", err)
		}
		shelfRows[i], err = listRow.ToListRow()
		if err != nil {
			return shelfRows, fmt.Errorf("error while mapping the row data to the shelfRows slice: %v", err)
		}
		i += 1
	}
	if results.Err() != nil {
		return shelfRows, fmt.Errorf("error while assigning the row to the variable after fetching it: %v", err)
	}

	return shelfRows, nil
}

// ShelfCounter returns the count of rows in the shelf_fts table that match
// the specified queryString.
// If queryString is empty, it returns the count of all rows in the table.
func (db *DB) ShelfCounter(queryString string) (count int, err error) {
	countQuery := `SELECT COUNT(*) FROM shelf_fts;`

	if queryString != "" {
		countQuery = fmt.Sprintf(`
			SELECT COUNT(*)
			FROM shelf_fts
      WHERE label MATCH '%s*' `, queryString)
	}

	err = db.Sql.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error while fetching the number of shelves from the database: %v", err)
	}
	return count, nil
}
