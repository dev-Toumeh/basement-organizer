package shelves

import (
	"net/http"

	"github.com/gofrs/uuid/v5"

	"basement/main/internal/common"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
)

type Shelf struct {
	ID             uuid.UUID            `json:"id"`
	Label          string               `json:"label"       validate:"required,lte=128"`
	Description    string               `json:"description" validate:"omitempty,lte=256"`
	Picture        string               `json:"picture"     validate:"omitempty,base64"`
	PreviewPicture string               `json:"previewpicture"     validate:"omitempty,base64"`
	QRcode         string               `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	Items          []*items.ItemListRow `json:"items"`
	Boxes          []*items.BoxListRow  `json:"boxes"`
	Height         float32
	Width          float32
	Depth          float32
	Rows           int
	Cols           int
	AreaID         uuid.UUID
}

// @TODO: Fix import cycle with items package.
type ShelfCoordinates struct {
	ID      uuid.UUID `json:"id"`
	ShelfID uuid.UUID `json:"shelfid"`
	Label   string    `json:"label"       validate:"required,lte=128"`
	Row     int       `json:"row"`
	Col     int       `json:"col"`
}

const (
	ID             string = "id"
	LABEL          string = "label"
	DESCRIPTION    string = "description"
	PICTURE        string = "picture"
	PREVIEWPICTURE string = "previewpicture"
	QRCODE         string = "qrcode"
	HEIGHT         string = "height"
	WIDTH          string = "width"
	DEPTH          string = "depth"
	ROWS           string = "rows"
	COLS           string = "cols"
	AREA_ID        string = "area_id"
)

type ShelfDB interface {
	CreateShelf(shelf *Shelf) error
	Shelf(id uuid.UUID) (*Shelf, error)
	UpdateShelf(shelf *Shelf) error
	DeleteShelf(id uuid.UUID) error
}

func CreateShelfHandler(db ShelfDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			shelve, err := shelf(r)
			if err != nil {
				logg.Errf("error while packing the shelve request data into %v", err)
			}
			// validate shelves request data @toDo
			err = db.CreateShelf(shelve)
			if err != nil {
				templates.RenderErrorNotification(w, "error while creating a new Shelf please try again later")
			}
			templates.RenderSuccessNotification(w, "the new shelf was created successfully")
		} else {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		}
	}
}

// this function will pack the request into struct from type Shelf, so it will be easier to handle it
func shelf(r *http.Request) (*Shelf, error) {
	id, areaId, err := common.CheckIDs(r.PostFormValue(ID), r.PostFormValue(AREA_ID))
	if err != nil {
		return &Shelf{}, err
	}

	newShelf := &Shelf{
		ID:             id,
		Label:          r.PostFormValue(LABEL),
		Description:    r.PostFormValue(DESCRIPTION),
		Picture:        "",
		PreviewPicture: "",
		QRcode:         r.PostFormValue(QRCODE),
		Height:         common.StringToFloat32(r.PostFormValue(HEIGHT)),
		Width:          common.StringToFloat32(r.PostFormValue(WIDTH)),
		Depth:          common.StringToFloat32(r.PostFormValue(DEPTH)),
		Rows:           common.StringToInt(r.PostFormValue(ROWS)),
		Cols:           common.StringToInt(r.PostFormValue(COLS)),
		AreaID:         areaId,
	}
	return newShelf, nil
}
