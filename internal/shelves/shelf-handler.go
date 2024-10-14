package shelves

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"basement/main/internal/common"
	"basement/main/internal/items"
	"basement/main/internal/server"
)

type Shelf struct {
	items.BasicInfo
	Items  []*items.ListRow `json:"items"`
	Boxes  []*items.ListRow `json:"boxes"`
	Height float32
	Width  float32
	Depth  float32
	Rows   int
	Cols   int
	AreaId uuid.UUID
}

type ShelfListRow struct {
	ID             uuid.UUID
	Label          string
	AreaID         uuid.UUID
	AreaLabel      string
	PreviewPicture string
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
	ShelfSearchListRowsPaginated(page int, rows int, search string) (shelfRows []*items.ListRow, found int, err error)
}

// handles read, create, update and delete for single shelf.
func ShelfHandler(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			const errMsgForUser = "Can't find box"

			id := server.ValidID(w, r, errMsgForUser)
			if id.IsNil() {
				return
			}

			shelf, err := db.Shelf(id)
			if err != nil {
				server.WriteNotFoundError(errMsgForUser, err, w, r)
				return
			}

			// Use API data writer
			if !server.WantsTemplateData(r) {
				server.WriteJSON(w, shelf)
				return
			}

			// Template writer
			renderShelfTemplate(shelf, w, r)
			break

		case http.MethodPost:
			createShelf(w, r, db)
			break

		case http.MethodDelete:
			deleteShelf(w, r, db)
			return

		case http.MethodPut:
			updateShelf(w, r, db)
			break

		default:
			// Other methods are not allowed.
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		}
	}
}

func renderShelfTemplate(box *Shelf, w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// this function will pack the request into struct from type Shelf, so it will be easier to handle it
func shelf(r *http.Request) (*Shelf, error) {
	id, areaId, err := common.CheckIDs(r.PostFormValue(ID), r.PostFormValue(AREA_ID))
	if err != nil {
		return &Shelf{}, err
	}

	newShelf := &Shelf{
		BasicInfo: items.BasicInfo{
			ID:             id,
			Label:          r.PostFormValue(LABEL),
			Description:    r.PostFormValue(DESCRIPTION),
			Picture:        common.ParsePicture(r),
			PreviewPicture: "",
		},
		Height: common.StringToFloat32(r.PostFormValue(HEIGHT)),
		Width:  common.StringToFloat32(r.PostFormValue(WIDTH)),
		Depth:  common.StringToFloat32(r.PostFormValue(DEPTH)),
		Rows:   common.StringToInt(r.PostFormValue(ROWS)),
		Cols:   common.StringToInt(r.PostFormValue(COLS)),
		AreaId: areaId,
	}
	return newShelf, nil
}
