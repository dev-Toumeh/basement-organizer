package boxes

import (
	"basement/main/internal/common"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type BoxDatabase interface {
	CreateBox(newBox *Box) (uuid.UUID, error)
	MoveBoxToBox(box1 uuid.UUID, box2 uuid.UUID) error
	MoveBoxToShelf(boxID uuid.UUID, toShelfID uuid.UUID) error
	UpdateBox(box Box) error
	DeleteBox(boxId uuid.UUID) error
	BoxById(id uuid.UUID) (Box, error)
	BoxIDs() ([]string, error) // @TODO: Change string to uuid.UUID
	BoxListRows(searchQuery string, limit int, page int) ([]common.ListRow, error)
	BoxListRowByID(id uuid.UUID) (common.ListRow, error)
	BoxListCounter(searchQuery string) (count int, err error)
	ShelfListCounter(searchQuery string) (count int, err error)
	ShelfListRows(searchQuery string, limit int, page int) (shelfRows []common.ListRow, err error)
}

type Box struct {
	common.BasicInfo
	OuterBoxID       uuid.UUID         `json:"outerboxid"`
	Items            []*common.ListRow `json:"items"`
	InnerBoxes       []*common.ListRow `json:"innerboxes"`
	OuterBox         *common.ListRow   `json:"outerbox"`
	ShelfID          uuid.UUID
	AreaID           uuid.UUID
	ShelfCoordinates *ShelfCoordinates `json:"shelfcoordinates"`
}

func (box *Box) Map() map[string]any {
	m := box.BasicInfo.Map()
	m["OuterBoxID"] = box.OuterBoxID
	m["Items"] = templates.SliceToSliceMaps(box.Items)
	m["InnerBoxes"] = templates.SliceToSliceMaps(box.InnerBoxes)
	m["OuterBox"] = box.OuterBox
	m["ShelfID"] = box.ShelfID
	m["AreaID"] = box.AreaID
	return m
}

type ShelfCoordinates struct {
	ID      uuid.UUID `json:"id"`
	ShelfID uuid.UUID `json:"shelfid"`
	Label   string    `json:"label"       validate:"required,lte=128"`
	Row     int       `json:"row"`
	Col     int       `json:"col"`
}

type BoxTemplateData struct {
	*Box
	Edit   bool
	Create bool
}

func (tmpl BoxTemplateData) Map() map[string]any {
	data := make(map[string]any, 0)
	maps.Copy(data, tmpl.Box.Map())
	data["Edit"] = tmpl.Edit
	data["Create"] = tmpl.Create
	return data
}

type BoxListTemplateData struct {
	Boxes []common.ListRow
}

func (tmpl BoxListTemplateData) Map() map[string]any {
	data := make([]map[string]any, 0)
	for i := range tmpl.Boxes {
		data = append(data, tmpl.Boxes[i].Map())
	}
	return map[string]any{"Boxes": data}
}

func RenderBoxListItem(w http.ResponseWriter, data *Box) {
	templates.Render(w, templates.TEMPLATE_BOX_LIST_ROW, data)
}

type boxPageTemplateData struct {
	*BoxTemplateData
	*templates.PageTemplate
}

// BoxPageTemplateData returns struct needed for "templates.TEMPLATE_BOX_DETAILS_PAGE" with default values.
func BoxPageTemplateData() *boxPageTemplateData {
	pageTmpl := templates.NewPageTemplate()
	boxTmpl := BoxTemplateData{}
	data := boxPageTemplateData{
		BoxTemplateData: &boxTmpl,
		PageTemplate:    &pageTmpl,
	}
	return &data
}

func (tmpl *boxPageTemplateData) Map() map[string]any {
	data := make(map[string]any, 0)
	maps.Copy(data, tmpl.BoxTemplateData.Map())
	maps.Copy(data, tmpl.PageTemplate.Map())
	return data
}

type BoxC struct {
	Id          uuid.UUID    `json:"id"`
	Label       string       `json:"label"       validate:"required,lte=128"`
	Description string       `json:"description" validate:"omitempty,lte=256"`
	Picture     string       `json:"picture"     validate:"omitempty,base64"`
	QRCode      string       `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	Items       []items.Item `json:"items"`
	InnerBoxes  []Box        `json:"innerboxes"`
	OuterBox    Box          `json:"outerbox" `
}

// NewBox returns an empty box with a new uuid.
func NewBox() Box {
	b := common.NewBasicInfoWithLabel("Box")
	return Box{BasicInfo: b}
}

func (b *Box) MarshalJSON() ([]byte, error) {
	c := Box{}
	for _, item := range b.Items {
		it := *item
		c.Items = append(c.Items, &common.ListRow{ID: it.ID, Label: it.Label, PreviewPicture: it.PreviewPicture})
	}

	for _, innerb := range b.InnerBoxes {
		c.InnerBoxes = append(c.InnerBoxes, &common.ListRow{BoxID: innerb.BoxID, Label: innerb.Label, PreviewPicture: innerb.PreviewPicture})
	}

	// if b.OuterBox != nil {
	// 	c.OuterBox = *b.OuterBox
	// }
	return json.Marshal(c)
}

func (b Box) String() string {
	// @TODO: Shorteing picture to now blow up logs with base64 encoding.
	// A little dirty but is ok for now.
	shortenPicture := true
	if shortenPicture {
		b.Picture = shortenPictureForLogs(b.Picture)
		if b.OuterBox != nil {
			b.OuterBox.PreviewPicture = shortenPictureForLogs(b.OuterBox.PreviewPicture)
		}
		for i := range b.InnerBoxes {
			b.InnerBoxes[i].PreviewPicture = shortenPictureForLogs(b.InnerBoxes[i].PreviewPicture)
		}
		for i := range b.Items {
			b.Items[i].PreviewPicture = shortenPictureForLogs(b.Items[i].PreviewPicture)
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

func shortenPictureForLogs(picture string) string {
	if len(picture) < 4 {
		return ""
	}
	return picture[0:3] + "...(shortened)"
}
