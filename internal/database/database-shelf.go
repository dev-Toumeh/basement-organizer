package database

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/shelves"
	"bytes"
	"database/sql"
	"encoding/base64"
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
}

// convert struct from type SQLShelf to normal shelf struct
func (s SQLShelf) toShelf() (shelves.Shelf, error) {
	id, err := uuid.FromString(s.SQLBasicInfo.ID.String)
	if err != nil {
		return shelves.Shelf{}, logg.WrapErr(err)
	}

	return shelves.Shelf{
		BasicInfo: items.BasicInfo{
			ID:             id,
			Label:          ifNullString(s.Label),
			Description:    ifNullString(s.Description),
			Picture:        ifNullString(s.Picture),
			PreviewPicture: ifNullString(s.PreviewPicture),
			QRCode:         ifNullString(s.QRCode),
		},
		Height: float32(ifNullFloat64(s.Height)),
		Width:  float32(ifNullFloat64(s.Width)),
		Depth:  float32(ifNullFloat64(s.Depth)),
		Rows:   int(ifNullInt(s.Rows)),
		Cols:   int(ifNullInt(s.Cols)),
		Items:  nil,
		Boxes:  nil,
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
	if shelf.Items != nil || shelf.Boxes != nil {
		return logg.NewError(fmt.Sprintf(`shelf "%s" has %d items and %d boxes, they must be empty while creating a new shelf`,
			shelf.ID, len(shelf.Items), len(shelf.Boxes)))
	}

	err = updatePicture(&shelf.Picture, &shelf.PreviewPicture)
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
		picture64,
		previewPicture64,
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

	// var pictureBase64, previewPictureBase64 string
	// if len(picture) > 0 {
	// 	pictureBase64 = base64.StdEncoding.EncodeToString(picture)
	// }
	// if len(previewPicture) > 0 {
	// 	previewPictureBase64 = base64.StdEncoding.EncodeToString(previewPicture)
	// }

	// shelf := &shelves.Shelf{
	// 	BasicInfo: items.BasicInfo{
	// 		ID:             id,
	// 		Label:          label,
	// 		Description:    description.String,
	// 		Picture:        picture.String,
	// 		PreviewPicture: previewPicture.String,
	// 		QRCode:         qrcode.String,
	// 	},
	// 	Height: float32(height.Float64),
	// 	Width:  float32(width.Float64),
	// 	Depth:  float32(depth.Float64),
	// 	Rows:   int(rows.Int64),
	// 	Cols:   int(cols.Int64),
	// 	Items:  nil,
	// 	Boxes:  nil,
	// }

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

	shelf, err := sqlShelf.toShelf()
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
	return &shelf, nil
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

// MoveShelfToArea moves shelf to an area.
// To move out of an area set "toAreaID = uuid.Nil".
func (db *DB) MoveShelfToArea(shelfID uuid.UUID, toAreaID uuid.UUID) error {
	return logg.WrapErr(ErrNotImplemented)
}
