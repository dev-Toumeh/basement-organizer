package boxes

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type BoxDatabase interface {
	CreateBox(newBox *Box) (uuid.UUID, error)
	MoveBoxToBox(box1 uuid.UUID, box2 uuid.UUID) error
	MoveBoxToShelf(boxID uuid.UUID, toShelfID uuid.UUID) error
	MoveBoxToArea(boxID uuid.UUID, toAreaID uuid.UUID) error
	UpdateBox(box Box, ignorePicture bool) error
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
	Items            []common.ListRow `json:"items"`
	InnerBoxes       []common.ListRow `json:"innerboxes"`
	OuterBox         *common.ListRow  `json:"outerbox"`
	OuterBoxLabel    string           `json:"outer_box_label"`
	ShelfID          uuid.UUID
	OuterBoxID       uuid.UUID `json:"outerboxid"`
	ShelfLabel       string    `json:"shelf_label"`
	AreaID           uuid.UUID
	AreaLabel        string            `json:"area_label"`
	ShelfCoordinates *ShelfCoordinates `json:"shelfcoordinates"`
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
	ID      uuid.UUID `json:"id"`
	ShelfID uuid.UUID `json:"shelfid"`
	Label   string    `json:"label"       validate:"required,lte=128"`
	Row     int       `json:"row"`
	Col     int       `json:"col"`
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
	var shortenPicture bool
	if env.Development() {
		shortenPicture = true
	} else {
		shortenPicture = false
	}
	if shortenPicture {
		b.Picture = common.ShortenPictureForLogs(b.Picture)
		b.PreviewPicture = common.ShortenPictureForLogs(b.PreviewPicture)
		if b.OuterBox != nil {
			b.OuterBox.PreviewPicture = common.ShortenPictureForLogs(b.OuterBox.PreviewPicture)
		}
		for i := range b.InnerBoxes {
			b.InnerBoxes[i].PreviewPicture = common.ShortenPictureForLogs(b.InnerBoxes[i].PreviewPicture)
		}
		for i := range b.Items {
			b.Items[i].PreviewPicture = common.ShortenPictureForLogs(b.Items[i].PreviewPicture)
		}
	}

	data, err := json.Marshal(b)
	if err != nil {
		logg.Err("Can't JSON box to string:", err)
		return ""
	}
	s := fmt.Sprintf("%s", data)
	return s
}
