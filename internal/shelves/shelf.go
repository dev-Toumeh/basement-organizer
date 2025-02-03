package shelves

import (
	"basement/main/internal/common"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type Shelf struct {
	common.BasicInfo
	Items          []*common.ListRow `json:"items"`
	Boxes          []*common.ListRow `json:"boxes"`
	InnerItemsList common.ListTemplate
	InnerBoxesList common.ListTemplate
	Height         float32
	Width          float32
	Depth          float32
	Rows           int
	Cols           int
	AreaId         uuid.UUID
	AreaLabel      string
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
	Shelf(id uuid.UUID) (*Shelf, error)
	CreateShelf(shelf *Shelf) error
	UpdateShelf(shelf *Shelf) error
	DeleteShelf(id uuid.UUID) (label string, err error)
	ShelfListRows(searchString string, limit int, pageNr int) (shelfRows []common.ListRow, err error)
	ShelfListCounter(queryString string) (count int, err error)
	ErrorNotEmpty() error

	// required in common.Database interface
	InnerListRowsFrom2(belongsToTable string, belongsToTableID uuid.UUID, listRowsTable string) ([]common.ListRow, error)
	InnerListRowsPaginatedFrom(belongsToTable string, belongsToTableID uuid.UUID, listRowsTable string, searchQuery string, limit int, page int) (listRows []common.ListRow, err error)
	InnerBoxInBoxListCounter(searchString string, inTable string, inTableID uuid.UUID) (count int, err error)
	InnerShelfInTableListCounter(searchString string, inTable string, inTableID uuid.UUID) (count int, err error)
	InnerThingInTableListCounter(searchString string, thing int, inTable string, inTableID uuid.UUID) (count int, err error)
	MoveShelfToArea(shelfID uuid.UUID, toAreaID uuid.UUID) error
	BoxListCounter(searchQuery string) (count int, err error)
	AreaListCounter(searchQuery string) (count int, err error)
	BoxListRows(searchQuery string, limit int, page int) ([]common.ListRow, error)
	AreaListRows(searchQuery string, limit int, page int) (areaRows []common.ListRow, err error)
	DeleteItem(itemID uuid.UUID) error
	DeleteBox(boxID uuid.UUID) error
	DeleteShelf2(id uuid.UUID) error
	DeleteArea(areaID uuid.UUID) error
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
		"AreaID":         s.AreaId,
		"AreaLabel":      s.AreaLabel,
		"InnerBoxesList": s.InnerBoxesList,
		"InnerItemsList": s.InnerItemsList,
	}

	return shelfMap
}

// return new Shelf with default Values
func newShelf() *Shelf {
	s := common.NewBasicInfoWithLabel("Shelf")
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
		BasicInfo: common.BasicInfo{
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
