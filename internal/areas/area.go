package areas

import (
	"basement/main/internal/common"
	"basement/main/internal/server"
	"basement/main/internal/validate"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type AreaDatabase interface {
	CreateArea(newArea Area) (uuid.UUID, error)
	UpdateArea(area Area, ignorePicture bool, pictureFormat string) error
	DeleteArea(id uuid.UUID) error
	AreaById(id uuid.UUID) (Area, error)
	AreaIDs() ([]uuid.UUID, error)
	AreaListRows(query string, limit int, page int) ([]common.ListRow, error)
	AreaListRowByID(id uuid.UUID) (common.ListRow, error)
	AreaListCounter(searchString string) (count int, err error)
	BoxListCounter(searchQuery string) (count int, err error)
	ShelfListCounter(searchQuery string) (count int, err error)
	BoxListRows(searchQuery string, limit int, page int) ([]common.ListRow, error)
	ShelfListRows(searchQuery string, limit int, page int) (shelfRows []common.ListRow, err error)
	InnerListRowsFrom2(belongsToTable string, belongsToTableID uuid.UUID, listRowsTable string) ([]common.ListRow, error)
}

const (
	ID             string = "id"
	LABEL          string = "label"
	DESCRIPTION    string = "description"
	PICTURE        string = "picture"
	PREVIEWPICTURE string = "previewpicture"
	QRCODE         string = "qrcode"
)

type Area struct {
	common.BasicInfo
}

func NewArea() Area {
	b := common.NewBasicInfoWithLabel("Area")
	return Area{BasicInfo: b}
}

func (area Area) Map() map[string]any {
	m := area.BasicInfo.Map()
	return m
}

func areaFromPostFormValue(id uuid.UUID, r *http.Request) (area Area, ignorePicture bool) {
	ignorePicture = server.ParseIgnorePicture(r)
	area.BasicInfo = common.BasicInfoFromPostFormValue(id, r, false)
	return area, ignorePicture
}

// ValidateShelf parses form input from an HTTP request,
// performs inline field-level validation directly within the handler,
// and returns both the validation result and a Shelf struct if valid.
// If any validation fails, it returns the validator with error details.
func ValidateArea(w http.ResponseWriter, r *http.Request) (area Area, validator validate.Validate, err error) {
	varea := validate.AreaValidate{
		BasicInfoValidate: validate.BasicInfoValidate{
			ID:             validate.NewUUIDField(r.PostFormValue(ID)),
			Label:          validate.NewStringField(r.PostFormValue(LABEL)),
			Description:    validate.NewStringField(r.PostFormValue(DESCRIPTION)),
			Picture:        validate.NewStringField(common.ParsePicture(r)),
			PreviewPicture: validate.NewStringField(common.ParsePicture(r)),
			QRCode:         validate.NewStringField(r.PostFormValue(QRCODE)),
		},
	}

	validator = validate.Validate{Area: varea}
	if err = validator.ValidateArea(w, varea); err != nil {
		return area, validator, err
	}

	if validator.HasValidateErrors() {
		return area, validator, validator.Err()
	}

	area = Area{
		BasicInfo: common.BasicInfo{
			ID:             varea.ID.UUID(),
			Label:          varea.Label.String(),
			Description:    varea.Description.String(),
			Picture:        varea.Picture.String(),
			PreviewPicture: varea.PreviewPicture.String(),
			QRCode:         varea.QRCode.String(),
		},
	}

	return area, validator, nil
}
