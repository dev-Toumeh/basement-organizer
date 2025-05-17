package shelves

import (
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/validate"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type Shelf struct {
	common.BasicInfo
	Items          []*common.ListRow `json:"items"`
	Boxes          []*common.ListRow `json:"boxes"`
	InnerItemsList common.ListTemplate
	InnerBoxesList common.ListTemplate
	Height         float64
	Width          float64
	Depth          float64
	Rows           int64
	Cols           int64
	AreaID         uuid.UUID
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
	UpdateShelf(shelf *Shelf, ignorePicture bool, pictureFormat string) error
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
func (s *Shelf) Map() map[string]any {
	shelfMap := map[string]any{
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
		"AreaID":         s.AreaID,
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
		Depth:     1.0,
		Rows:      5,
		Cols:      10,
	}
}

// ValidateShelf parses form input from an HTTP request,
// performs field-level validation, and returns both the validation result and a Shelf struct if valid.
// If any validation fails, it returns the validator with error details.
func ValidateShelf(w http.ResponseWriter, r *http.Request) (shelf *Shelf, validator validate.Validate, err error) {
	vshelf := validate.ShelfValidate{
		BasicInfoValidate: validate.BasicInfoValidate{
			ID:             validate.NewUUIDField(r.PostFormValue(ID)),
			Label:          validate.NewStringField(r.PostFormValue(LABEL)),
			Description:    validate.NewStringField(r.PostFormValue(DESCRIPTION)),
			Picture:        validate.NewStringField(common.ParsePicture(r)),
			PreviewPicture: validate.NewStringField(common.ParsePicture(r)),
			QRCode:         validate.NewStringField(r.PostFormValue(QRCODE)),
		},
		Height: validate.NewFloatField(r.PostFormValue(HEIGHT)),
		Width:  validate.NewFloatField(r.PostFormValue(WIDTH)),
		Depth:  validate.NewFloatField(r.PostFormValue(DEPTH)),
		Rows:   validate.NewIntField(r.PostFormValue(ROWS)),
		Cols:   validate.NewIntField(r.PostFormValue(COLS)),
		AreaID: validate.NewUUIDField(r.PostFormValue(AREA_ID)),
	}
	logg.Debug(vshelf)

	validator = validate.Validate{Shelf: vshelf}

	// validate fields
	if err = validator.ValidateShelf(w, vshelf); err != nil {
		return shelf, validator, err
	}

	if validator.HasValidateErrors() {
		return shelf, validator, validator.Err()
	}

	// pack into Shelf struct
	shelf = &Shelf{
		BasicInfo: common.BasicInfo{
			ID:             vshelf.ID.UUID(),
			Label:          vshelf.Label.String(),
			Description:    vshelf.Description.String(),
			Picture:        vshelf.Picture.String(),
			PreviewPicture: vshelf.PreviewPicture.String(),
			QRCode:         vshelf.QRCode.String(),
		},
		Height: vshelf.Height.Float64(),
		Width:  vshelf.Width.Float64(),
		Depth:  vshelf.Depth.Float64(),
		Rows:   vshelf.Rows.Int(),
		Cols:   vshelf.Cols.Int(),
		AreaID: vshelf.AreaID.UUID(),
	}

	return shelf, validator, nil
}
