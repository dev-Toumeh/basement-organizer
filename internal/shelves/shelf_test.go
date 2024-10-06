package shelves

import (
	"errors"

	"github.com/gofrs/uuid/v5"
)

var shelf1 *Shelf = &Shelf{
	ID:             uuid.Must(uuid.FromString("111e4567-e89b-12d3-a456-426614174000")),
	Label:          "Storage Shelf 1",
	Description:    "This is the first dummy shelf",
	Picture:        "base64PictureData1",
	PreviewPicture: "base64PreviewPictureData1",
	QRcode:         "QR1234ABC",
	Height:         250.0,
	Width:          120.0,
	Depth:          60.0,
	Rows:           4,
	Cols:           3,
	AreaID:         uuid.Must(uuid.FromString("222e4567-e89b-12d3-a456-426614174001")),
}

var shelf2 *Shelf = &Shelf{
	ID:             uuid.Must(uuid.FromString("333e4567-e89b-12d3-a456-426614174002")),
	Label:          "Storage Shelf 2",
	Description:    "This is the second dummy shelf",
	Picture:        "base64PictureData2",
	PreviewPicture: "base64PreviewPictureData2",
	QRcode:         "QR5678XYZ",
	Height:         300.0,
	Width:          150.0,
	Depth:          70.0,
	Rows:           5,
	Cols:           4,
	AreaID:         uuid.Must(uuid.FromString("444e4567-e89b-12d3-a456-426614174003")),
}

type ShelfDatabaseError struct{}

func (db *ShelfDatabaseError) CreateShelf(shelf *Shelf) error {
	return errors.New("unable to create shelf")
}

func (db *ShelfDatabaseError) Shelf(id uuid.UUID) (*Shelf, error) {
	return nil, errors.New("shelf not found")
}

func (db *ShelfDatabaseError) UpdateShelf(shelf *Shelf) error {
	return errors.New("unable to update shelf")
}

func (db *ShelfDatabaseError) DeleteShelf(id uuid.UUID) error {
	return errors.New("unable to delete shelf")
}

// ShelfDatabaseSuccess implements ShelfDB interface without errors for success testing.
type ShelfDatabaseSuccess struct{}

func (db *ShelfDatabaseSuccess) CreateShelf(shelf *Shelf) error {
	return nil
}

func (db *ShelfDatabaseSuccess) Shelf(id uuid.UUID) (*Shelf, error) {
	return &Shelf{ID: id}, nil
}

func (db *ShelfDatabaseSuccess) UpdateShelf(shelf *Shelf) error {
	return nil
}

func (db *ShelfDatabaseSuccess) DeleteShelf(id uuid.UUID) error {
	return nil
}
