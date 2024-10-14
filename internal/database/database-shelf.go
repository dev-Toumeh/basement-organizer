package database

import (
	"basement/main/internal/logg"
	"basement/main/internal/shelves"
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type SQLShelf struct {
	Id             sql.NullString
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
	logg.Debugf("SQL: %s", stmt)
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

// type ShelfDB interface {
// 	CreateShelf(shelf *shelves.Shelf) error
// 	Shelf(id uuid.UUID) (*shelves.Shelf, error)
// 	UpdateShelf(shelf *shelves.Shelf) error
// 	DeleteShelf(id uuid.UUID) error
// }

// CreateShelf creates a shelf entry in database from the provided shelf.
func (db *DB) CreateShelf(shelf *shelves.Shelf) error {
	var (
		picture64        string
		previewPicture64 string
		err              error
	)
	// picture := bytes.Buffer{}
	// previewPicture := bytes.Buffer{}

	// fmt.Println(shelf.Picture)
	if shelf.Picture != "" {
		// picture64 = ByteArrayToBase64ByteArray(shelf.Picture)
		// picture, err = base64.StdEncoding.DecodeString(shelf.Picture)
		// enc := base64.NewEncoder(base64.StdEncoding, &picture64)
		// enc.Write(shelf.Picture)
		// enc.Close()
		// base64.StdEncoding.Encode(picture.Bytes(), shelf.Picture)
		// picture64 = shelf.Picture
		_, err = base64.StdEncoding.DecodeString(shelf.Picture)
		if err != nil {
			return logg.Errorf("CreateShelf: invalid base64 in Picture %w", err)
		}
		picture64 = shelf.Picture
	}

	if shelf.PreviewPicture != "" {
		// previewPicture64 = ByteArrayToBase64ByteArray(shelf.PreviewPicture)
		// previewPicture = shelf.PreviewPicture
		// enc := base64.NewEncoder(base64.StdEncoding, &previewPicture64)
		// enc.Write(shelf.PreviewPicture)
		// enc.Close()
		// base64.StdEncoding.Encode(previewPicture, shelf.PreviewPicture)
		_, err = base64.StdEncoding.DecodeString(shelf.PreviewPicture)
		if err != nil {
			return logg.Errorf("CreateShelf: invalid base64 in PreviewPicture %w", err)
		}
		previewPicture64 = shelf.PreviewPicture
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
		shelf.Id.String(),
		shelf.Label,
		shelf.Description,
		picture64,
		previewPicture64,
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
	var (
		label          string
		description    sql.NullString
		picture        sql.NullString
		previewPicture sql.NullString
		qrcode         sql.NullString
		height         sql.NullFloat64
		width          sql.NullFloat64
		depth          sql.NullFloat64
		rows           sql.NullInt64
		cols           sql.NullInt64
	)

	stmt := `
        SELECT
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
		return nil, logg.WrapErr(err)
	}

	// var pictureBase64, previewPictureBase64 string
	// if len(picture) > 0 {
	// 	pictureBase64 = base64.StdEncoding.EncodeToString(picture)
	// }
	// if len(previewPicture) > 0 {
	// 	previewPictureBase64 = base64.StdEncoding.EncodeToString(previewPicture)
	// }
	shelf := &shelves.Shelf{
		Id:             id,
		Label:          label,
		Description:    description.String,
		Picture:        picture.String,
		PreviewPicture: previewPicture.String,
		QRcode:         qrcode.String,
		Height:         float32(height.Float64),
		Width:          float32(width.Float64),
		Depth:          float32(depth.Float64),
		Rows:           int(rows.Int64),
		Cols:           int(cols.Int64),
		Items:          nil,
		Boxes:          nil,
	}

	// f, err := os.ReadFile("./internal/static/pen.png")
	// if err != nil {
	// 	panic(err)
	// }
	// shelf2 := &shelves.Shelf{
	// Id:             uuid.FromStringOrNil("1cad0cb3-3307-43cc-b005-3f9f29eec8b4"),
	// shelf.PreviewPicture = f
	// base64.StdEncoding.Decode(shelf.PreviewPicture, shelf.PreviewPicture)
	// shelf.PreviewPicture = base64.StdEncoding.EncodeToString(f)
	// PreviewPicture: f,
	// }
	// _, err = base64.StdEncoding.Decode(shelf.PreviewPicture, previewPicture)
	// if err != nil {
	// 	panic(err)
	// }
	return shelf, nil
}

func ByteArrayToBase64ByteArray(src []byte) []byte {
	buf := bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	enc.Write(src)
	enc.Close()
	return buf.Bytes()
}

func (db *DB) UpdateShelf(shelf *shelves.Shelf) error {
	var (
		picture64        string
		previewPicture64 string
		err              error
	)

	if shelf.Picture != "" {
		// picture64 = ByteArrayToBase64ByteArray(shelf.Picture)
		// picture = shelf.Picture
		_, err = base64.StdEncoding.DecodeString(shelf.Picture)
		if err != nil {
			return logg.Errorf("UpdateShelf: invalid base64 in Picture %w", err)
		}
		picture64 = shelf.Picture
	}

	if shelf.PreviewPicture != "" {
		// previewPicture64 = ByteArrayToBase64ByteArray(shelf.PreviewPicture)
		_, err = base64.StdEncoding.DecodeString(shelf.PreviewPicture)
		if err != nil {
			return logg.Errorf("UpdateShelf: invalid base64 in PreviewPicture %w", err)
		}
		previewPicture64 = shelf.PreviewPicture
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
		picture64,
		previewPicture64,
		shelf.QRcode,
		shelf.Height,
		shelf.Width,
		shelf.Depth,
		shelf.Rows,
		shelf.Cols,
		shelf.Id.String(),
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
