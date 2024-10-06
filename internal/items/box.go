package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
)

type BoxListRow struct {
	BoxID          uuid.UUID
	Label          string
	OuterBoxID     uuid.UUID
	OuterBoxLabel  string
	ShelfID        uuid.UUID
	ShelfLabel     string
	AreaID         uuid.UUID
	AreaLabel      string
	PreviewPicture string
}

func (box *BoxListRow) Map() map[string]any {
	return map[string]interface{}{
		"BoxID":          box.BoxID,
		"Label":          box.Label,
		"OuterBoxID":     box.OuterBoxID,
		"OuterBoxLabel":  box.OuterBoxLabel,
		"ShelfID":        box.ShelfID,
		"ShelfLabel":     box.ShelfLabel,
		"AreaID":         box.AreaID,
		"AreaLabel":      box.AreaLabel,
		"PreviewPicture": box.PreviewPicture,
	}
}

type Box struct {
	ID             uuid.UUID      `json:"id"`
	Label          string         `json:"label"       validate:"required,lte=128"`
	Description    string         `json:"description" validate:"omitempty,lte=256"`
	Picture        string         `json:"picture"     validate:"omitempty,base64"`
	PreviewPicture string         `json:"previewpicture"     validate:"omitempty,base64"`
	QRcode         string         `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	OuterBoxID     uuid.UUID      `json:"outerboxid"`
	Items          []*ItemListRow `json:"items"`
	InnerBoxes     []*BoxListRow  `json:"innerboxes"`
	OuterBox       *BoxListRow    `json:"outerbox"`
	// @TODO: Fix import cycle with shelves package.
	// ShelfCoordinates *shelf         `json:"shelfcoordinates"`
}

func (box *Box) Map() map[string]any {
	return map[string]interface{}{
		"ID":             box.ID,
		"Label":          box.Label,
		"Description":    box.Description,
		"Picture":        box.Picture,
		"PreviewPicture": box.PreviewPicture,
		"QRcode":         box.QRcode,
		"OuterBoxID":     box.OuterBoxID,
		"Items":          box.Items,
		"InnerBoxes":     box.InnerBoxes,
		"OuterBox":       box.OuterBox,
		// @TODO: Fix import cycle with shelves package.
		// "ShelfCoordinates": box.ShelfCoordinates,
	}
}

type BoxTemplateData struct {
	*Box
	Edit bool
}

func (tmpl BoxTemplateData) Map() map[string]any {
	data := make(map[string]any, 0)
	maps.Copy(data, tmpl.Box.Map())
	data["Edit"] = tmpl.Edit
	return data
}

type BoxListTemplateData struct {
	Boxes []*Box
}

func (tmpl BoxListTemplateData) Map() map[string]any {
	data := make([]map[string]any, 0)
	for i := range tmpl.Boxes {
		data = append(data, tmpl.Boxes[i].Map())
	}
	return map[string]any{"Boxes": data}
}

func RenderBoxListItem(w http.ResponseWriter, data *Box) {
	templates.Render(w, templates.TEMPLATE_BOX_LIST_ITEM, data)
}

func RenderBoxList(w http.ResponseWriter, boxes []*Box) {
	var data any
	if boxes == nil {
		data = map[string]any{}
	} else {
		data = map[string][]*Box{"Boxes": boxes}
	}
	logg.Debug(data)
	templates.Render(w, templates.TEMPLATE_BOX_LIST, data)
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
	Id          uuid.UUID `json:"id"`
	Label       string    `json:"label"       validate:"required,lte=128"`
	Description string    `json:"description" validate:"omitempty,lte=256"`
	Picture     string    `json:"picture"     validate:"omitempty,base64"`
	QRcode      string    `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	Items       []Item    `json:"items"`
	InnerBoxes  []Box     `json:"innerboxes"`
	OuterBox    Box       `json:"outerbox" `
}

// NewBox returns an empty box with a new uuid.
func NewBox() Box {
	label := time.Now().Format("2006-01-02_15_04_05")
	return Box{
		ID:          uuid.Must(uuid.NewV4()),
		Label:       fmt.Sprintf("Box_%s", label),
		Description: fmt.Sprintf("Box description %s", label),
	}
}

func (b *Box) MarshalJSON() ([]byte, error) {
	c := Box{}
	for _, item := range b.Items {
		it := *item
		c.Items = append(c.Items, &ItemListRow{ItemID: it.ItemID, Label: it.Label, PreviewPicture: it.PreviewPicture})
	}

	for _, innerb := range b.InnerBoxes {
		c.InnerBoxes = append(c.InnerBoxes, &BoxListRow{BoxID: innerb.BoxID, Label: innerb.Label, PreviewPicture: innerb.PreviewPicture})
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
