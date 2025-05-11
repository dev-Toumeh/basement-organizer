package boxes

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"basement/main/internal/validate"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

const (
	ID          string = "id"
	LABEL       string = "label"
	DESCRIPTION string = "description"
	PICTURE     string = "picture"
	QUANTITY    string = "quantity"
	WEIGHT      string = "weight"
	QRCODE      string = "qrcode"
	BOX_ID      string = "box_id"
	BOX_LABEL   string = "box_label"
	SHELF_ID    string = "shelf_id"
	SHELF_LABEL string = "shelf_label"
	AREA_ID     string = "area_id"
	AREA_LABEL  string = "area_label"
)

type BoxDatabase interface {
	CreateBox(newBox *Box) (uuid.UUID, error)
	MoveBoxToBox(box1 uuid.UUID, box2 uuid.UUID) error
	MoveBoxToShelf(boxID uuid.UUID, toShelfID uuid.UUID) error
	MoveBoxToArea(boxID uuid.UUID, toAreaID uuid.UUID) error
	UpdateBox(box Box, ignorePicture bool, pictureFormat string) error
	DeleteBox(boxId uuid.UUID) error
	BoxById(id uuid.UUID) (Box, error)
	BoxIDs() ([]uuid.UUID, error)
	BoxListRows(searchQuery string, limit int, page int) ([]common.ListRow, error)
	BoxListRowByID(id uuid.UUID) (common.ListRow, error)
	// InnerListRowsFrom2(belongsToTable string, belongsToTableID uuid.UUID, listRowsTable string) ([]common.ListRow, error)
	BoxListCounter(searchQuery string) (count int, err error)
	ShelfListCounter(searchQuery string) (count int, err error)
	ShelfListRows(searchQuery string, limit int, page int) (shelfRows []common.ListRow, err error)
	AreaListCounter(searchQuery string) (count int, err error)
	AreaListRows(searchQuery string, limit int, page int) (rows []common.ListRow, err error)
}

type Box struct {
	common.BasicInfo
	Items            []common.ListRow
	InnerBoxes       []common.ListRow
	OuterBox         *common.ListRow
	OuterBoxLabel    string
	ShelfID          uuid.UUID
	OuterBoxID       uuid.UUID
	ShelfLabel       string
	AreaID           uuid.UUID
	AreaLabel        string
	ShelfCoordinates *ShelfCoordinates
}

func (box *Box) Map() map[string]any {
	m := box.BasicInfo.Map()
	m["Items"] = templates.SliceToSliceMaps(box.Items)
	m["InnerBoxes"] = templates.SliceToSliceMaps(box.InnerBoxes)
	m["OuterBox"] = box.OuterBox
	m["OuterBoxID"] = box.OuterBoxID
	m["OuterBoxLabel"] = box.OuterBoxLabel
	m["ShelfID"] = box.ShelfID
	m["ShelfLabel"] = box.ShelfLabel
	m["AreaID"] = box.AreaID
	m["AreaLabel"] = box.AreaLabel
	return m
}

type ShelfCoordinates struct {
	ID      uuid.UUID
	ShelfID uuid.UUID
	Label   string
	Row     int
	Col     int
}

// type BoxTemplateData struct {
// 	Box
// 	Edit   bool
// 	Create bool
// }
//
// func (tmpl BoxTemplateData) Map() map[string]any {
// 	data := make(map[string]any, 0)
// 	maps.Copy(data, tmpl.Box.Map())
// 	data["Edit"] = tmpl.Edit
// 	data["Create"] = tmpl.Create
// 	return data
// }

// type BoxListTemplateData struct {
// 	Boxes []common.ListRow
// }
//
// func (tmpl BoxListTemplateData) Map() map[string]any {
// 	data := make([]map[string]any, 0)
// 	for i := range tmpl.Boxes {
// 		data = append(data, tmpl.Boxes[i].Map())
// 	}
// 	return map[string]any{"Boxes": data}
// }

// NewBox returns an empty box with a new uuid.
func NewBox() Box {
	b := common.NewBasicInfoWithLabel("Box")
	return Box{BasicInfo: b}
}

func (b *Box) MarshalJSON() ([]byte, error) {
	c := Box{}
	for _, item := range b.Items {
		it := item
		c.Items = append(c.Items, common.ListRow{ID: it.ID, Label: it.Label, PreviewPicture: it.PreviewPicture})
	}

	for _, innerb := range b.InnerBoxes {
		c.InnerBoxes = append(c.InnerBoxes, common.ListRow{BoxID: innerb.BoxID, Label: innerb.Label, PreviewPicture: innerb.PreviewPicture})
	}

	// if b.OuterBox != nil {
	// 	c.OuterBox = *b.OuterBox
	// }
	return json.Marshal(c)
}

func (b Box) String() string {
	// @TODO: Shorteing picture to now blow up logs with base64 encoding.
	// A little dirty but is ok for now.
	bCopy := b
	shortenPicture := false
	if env.Development() {
		shortenPicture = true
	}
	if shortenPicture {
		bCopy.Picture = common.ShortenPictureForLogs(bCopy.Picture)
		bCopy.PreviewPicture = common.ShortenPictureForLogs(bCopy.PreviewPicture)
		if bCopy.OuterBox != nil {
			bCopy.OuterBox.PreviewPicture = common.ShortenPictureForLogs(bCopy.OuterBox.PreviewPicture)
		}
		for i := range bCopy.InnerBoxes {
			bCopy.InnerBoxes[i].PreviewPicture = common.ShortenPictureForLogs(bCopy.InnerBoxes[i].PreviewPicture)
		}
		for i := range bCopy.Items {
			bCopy.Items[i].PreviewPicture = common.ShortenPictureForLogs(bCopy.Items[i].PreviewPicture)
		}
	}

	data, err := json.Marshal(bCopy)
	if err != nil {
		logg.Err("Can't JSON box to string:", err)
		return ""
	}
	s := fmt.Sprintf("%s", data)
	return s
}

func (tmpl boxDetailsPageTemplate) Map() map[string]any {
	data := make(map[string]any, 0)
	maps.Copy(data, tmpl.PageTemplate.Map())
	maps.Copy(data, tmpl.Box.Map())
	data["Edit"] = tmpl.Edit
	data["Create"] = tmpl.Create
	return data
}

// ValidateBox parses form input from an HTTP request,
// performs field-level validation, and returns both the validation result and a Box struct if valid.
// If any validation fails, it returns the validator with error details
func ValidateBox(w http.ResponseWriter, r *http.Request) (box Box, validator validate.Validate, err error) {
	vbox := validate.BoxValidate{
		BasicInfoValidate: validate.BasicInfoValidate{
			ID:             validate.NewUUIDField(r.PostFormValue(ID)),
			Label:          validate.NewStringField(r.PostFormValue(LABEL)),
			Description:    validate.NewStringField(r.PostFormValue(DESCRIPTION)),
			Picture:        validate.NewStringField(common.ParsePicture(r)),
			PreviewPicture: validate.NewStringField(common.ParsePicture(r)),
		},
		ShelfID:    validate.NewUUIDField(r.PostFormValue(SHELF_ID)),
		OuterBoxID: validate.NewUUIDField(r.PostFormValue(BOX_ID)),
		AreaID:     validate.NewUUIDField(r.PostFormValue(AREA_ID)),
	}
	logg.DebugJSON(vbox, 50)

	validator = validate.Validate{Box: vbox}

	// used to return actual errors
	if err = validator.ValidateBox(w, vbox); err != nil {
		return box, validator, err
	}

	// return validation errors
	if validator.HasValidateErrors() {
		return box, validator, validator.Err()
	}

	// no error found return struct from type Item
	box = Box{
		BasicInfo: common.BasicInfo{
			ID:             vbox.ID.Value,
			Label:          vbox.Label.Value,
			Description:    vbox.Description.Value,
			Picture:        vbox.Picture.Value,
			PreviewPicture: vbox.PreviewPicture.Value,
		},
		ShelfID:    vbox.ShelfID.Value,
		OuterBoxID: vbox.OuterBoxID.Value,
		AreaID:     vbox.AreaID.Value,
	}

	return box, validator, nil
}
