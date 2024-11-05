package shelves

import (
	"basement/main/internal/common"
	"basement/main/internal/items"
	"net/http"

	"github.com/gofrs/uuid/v5"
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

type shelfTemplateDetailsData struct {
	Shelf *Shelf
	Edit  bool
}

type ShelfDB interface {
	CreateShelf(shelf *Shelf) error
	Shelf(id uuid.UUID) (*Shelf, error)
	UpdateShelf(shelf *Shelf) error
	DeleteShelf(id uuid.UUID) error
	ShelfListRows(searchString string, limit int, pageNr int) (shelfRows []items.ListRow, err error)
	ShelfListCounter(queryString string) (count int, err error)
}

//	type PaginationData struct {
//		PageNumber int
//		Offset     int
//	}
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

// return the Shelf in type map
func (s *Shelf) Map() map[string]interface{} {
	shelfMap := map[string]interface{}{
		"ID":             s.ID,
		"Label":          s.Label,
		"Description":    s.Description,
		"Picture":        s.Picture,
		"PreviewPicture": s.PreviewPicture,
		"Height":         s.Height,
		"Width":          s.Width,
		"Depth":          s.Depth,
		"Rows":           s.Rows,
		"Cols":           s.Cols,
	}

	return shelfMap
}

// return new Shelf with default Values
func newShelf() *Shelf {
	s := items.NewBasicInfoWithLabel("Shelf")
	return &Shelf{
		BasicInfo: s,
		Height:    2.0,
		Width:     1.0,
		Depth:     0.5,
		Rows:      5,
		Cols:      10,
	}
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
