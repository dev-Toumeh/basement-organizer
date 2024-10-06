package database

import (
	"basement/main/internal/logg"
	"basement/main/internal/shelves"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type SQLShelf struct {
	ID             sql.NullString
	Label          sql.NullString
	Description    sql.NullString
	Picture        sql.NullString
	PreviewPicture sql.NullString
	QRcode         sql.NullString
	Height         sql.NullFloat64
	Width          sql.NullFloat64
	Depth          sql.NullFloat64
	Rows           sql.NullInt64
	Cols           sql.NullInt64
}

func (s SQLShelf) ToShelf() (*shelves.Shelf, error) {
	id, err := uuid.FromString(s.ID.String)
	if err != nil {
		return nil, logg.WrapErr(err)
	}
	return &shelves.Shelf{
		ID:             id,
		Label:          s.Label.String,
		Description:    s.Description.String,
		Picture:        s.Picture.String,
		PreviewPicture: s.PreviewPicture.String,
		QRcode:         s.QRcode.String,
		Height:         float32(s.Height.Float64),
		Width:          float32(s.Width.Float64),
		Depth:          float32(s.Depth.Float64),
		Rows:           int(s.Rows.Int64),
		Cols:           int(s.Cols.Int64),
		Items:          nil,
		Boxes:          nil,
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

	err := updatePicture(&shelf.Picture, &shelf.PreviewPicture)
	if err != nil {
		logg.Errf("Can't update picture %v", err.Error())
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

// Shelf returns shelf with provided id.
func (db *DB) Shelf(id uuid.UUID) (*shelves.Shelf, error) {
	var sqlShelf SQLShelf

	stmt := `
        SELECT
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
        FROM shelf WHERE id = ?`

	err := db.Sql.QueryRow(stmt, id.String()).Scan(
		&sqlShelf.ID,
		&sqlShelf.Label,
		&sqlShelf.Description,
		&sqlShelf.Picture,
		&sqlShelf.PreviewPicture,
		&sqlShelf.QRcode,
		&sqlShelf.Height,
		&sqlShelf.Width,
		&sqlShelf.Depth,
		&sqlShelf.Rows,
		&sqlShelf.Cols,
	)

	if err != nil {
		return nil, logg.WrapErr(err)
	}

	return sqlShelf.ToShelf()
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

// DeleteShelf deletes shelf.
func (db *DB) DeleteShelf(id uuid.UUID) error {
	stmt := `DELETE FROM shelf WHERE id = ?`
	_, err := db.Sql.Exec(stmt, id.String())
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// MoveItemToShelf moves item to a shelf.
// To move item out of a shelf set "toShelfID = uuid.Nil".
func (db *DB) MoveItemToShelf(itemID uuid.UUID, toShelfID uuid.UUID) error {
	errMsg := fmt.Sprintf(`moving item "%s" to shelf "%s"`, itemID.String(), toShelfID.String())

	exists, err := db.Exists("item", itemID)
	if err != nil {
		return logg.WrapErr(err)
	}

	if exists == false {
		return logg.Errorf(errMsg+" item: %w", ErrNotExist)
	}

	// If moving to a shelf, check if the shelf exists
	if toShelfID != uuid.Nil {
		var shelfExists int
		err = db.Sql.QueryRow(`SELECT EXISTS(SELECT 1 FROM shelf WHERE id = ?)`, toShelfID.String()).Scan(&shelfExists)
		if err != nil {
			return logg.WrapErr(err)
		}
		if shelfExists == 0 {
			return logg.Errorf(errMsg+" shelf: %w", ErrNotExist)
		}
	}

	// Update the item's shelf_id
	stmt := `UPDATE item SET shelf_id = ? WHERE id = ?`
	var shelfIDValue interface{}
	if toShelfID == uuid.Nil {
		shelfIDValue = nil
	} else {
		shelfIDValue = toShelfID.String()
	}
	_, err = db.Sql.Exec(stmt, shelfIDValue, itemID.String())
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// MoveBoxToShelf moves box to a shelf.
// To move box out of a shelf set "toShelfID = uuid.Nil".
func (db *DB) MoveBoxToShelf(boxID uuid.UUID, toShelfID uuid.UUID) error {
	return logg.WrapErr(ErrNotImplemented)
}

// MoveShelf moves shelf to an area.
// To move out of an area set "toAreaID = uuid.Nil".
func (db *DB) MoveShelf(shelfID uuid.UUID, toAreaID uuid.UUID) error {
	return logg.WrapErr(ErrNotImplemented)
}

func (db *DB) insertShelf(shelf *shelves.Shelf) error {
	return logg.WrapErr(ErrNotImplemented)
}
