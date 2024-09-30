package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"time"

	"github.com/gofrs/uuid/v5"
)

type BoxListItem struct {
	Box_Id         uuid.UUID
	Label          string
	OuterBox_label string
	OuterBox_id    uuid.UUID
	Shelve_label   string
	Area_label     string
	PreviewPicture string
}

func (box *BoxListItem) Map() map[string]any {
	return map[string]interface{}{
		"Box_Id":         box.Box_Id,
		"Label":          box.Label,
		"OuterBox_label": box.OuterBox_label,
		"OuterBox_id":    box.OuterBox_id,
		"Shelve_label":   box.Shelve_label,
		"Area_label":     box.Area_label,
		"PreviewPicture": box.PreviewPicture,
	}
}

type AnotherItem struct {
	Quantity    int64     `json:"quantity"    validate:"omitempty,numeric,gte=1"`
	Weight      string    `json:"weight"      validate:"omitempty,numeric"`
	Id          uuid.UUID `json:"id"`
	Label       string    `json:"label"       validate:"required,lte=128"`
	Description string    `json:"description" validate:"omitempty,lte=256"`
	Picture     string    `json:"picture"     validate:"omitempty,base64"`
	QRcode      string    `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	Items       []*AnotherItem
	InnerBoxes  []*Box
	OuterBox    *Box
}

type Box struct {
	Id          uuid.UUID          `json:"id"`
	Label       string             `json:"label"       validate:"required,lte=128"`
	Description string             `json:"description" validate:"omitempty,lte=256"`
	Picture     string             `json:"picture"     validate:"omitempty,base64"`
	QRcode      string             `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	OuterBoxId  uuid.UUID          `json:"outerboxId"`
	Items       []*Item            `json:"items"`
	InnerBoxes  []*Box             `json:"innerboxes"`
	OuterBox    *Box               `json:"outerbox" `
	Shelve      *ShelveCoordinates `json:"shelveinfo" `
}

func (box *Box) Map() map[string]any {
	return map[string]interface{}{
		"Id":          box.Id,
		"Label":       box.Label,
		"Description": box.Description,
		"Picture":     box.Picture,
		"Qrcode":      box.QRcode,
		"OuterboxId":  box.OuterBoxId,
		"Items":       box.Items,
		"Innerboxes":  box.InnerBoxes,
		"Outerbox":    box.OuterBox,
		"Shelveinfo":  box.Shelve,
	}
}

type ShelveCoordinates struct {
	Id    uuid.UUID `json:"id"`
	Label string    `json:"label"       validate:"required,lte=128"`
	Rows  int
	Cols  int
}

type Shelve struct {
	Id             uuid.UUID      `json:"id"`
	Label          string         `json:"label"       validate:"required,lte=128"`
	Description    string         `json:"description" validate:"omitempty,lte=256"`
	Picture        string         `json:"picture"     validate:"omitempty,base64"`
	PreviewPicture string         `json:"previewpicture"     validate:"omitempty,base64"`
	QRcode         string         `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	Items          []*VirtualItem `json:"items"`
	Boxes          []*BoxListItem `json:"boxes"`
	Height         float32
	Width          float32
	Depth          float32
	Rows           int
	Cols           int
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
	// Quantity    int64     `json:"quantity"    validate:"omitempty,numeric,gte=1"`
	// Weight      string    `json:"weight"      validate:"omitempty,numeric"`
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
	// rand.Seed(uint64(time.Now().UnixNano()))
	// num := rand.Intn(10000)
	return Box{
		Id:          uuid.Must(uuid.NewV4()),
		Label:       fmt.Sprintf("Box_%s", label),
		Description: fmt.Sprintf("Box description %s", label),
	}
}

func (b *Box) MarshalJSON() ([]byte, error) {
	c := BoxC{}
	for _, item := range b.Items {
		it := *item
		c.Items = append(c.Items, it)
	}

	for _, innerb := range b.InnerBoxes {
		c.InnerBoxes = append(c.InnerBoxes, *innerb)
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
			b.OuterBox.Picture = shortenPictureForLogs(b.OuterBox.Picture)
		}
		for i := range b.InnerBoxes {
			b.InnerBoxes[i].Picture = shortenPictureForLogs(b.InnerBoxes[i].Picture)
		}
		for i := range b.Items {
			b.Items[i].Picture = shortenPictureForLogs(b.Items[i].Picture)
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

func (box *Box) MoveTo(other any) error {
	switch v := other.(type) {
	case *Box:
		logg.Debug("Moving '", box.Label, "' to '", v.Label, "'")
		if box == v {
			return fmt.Errorf("Can't move this box \"%v\" to itself \"%v\"", box.Label, v.Label)
		}
		v.InnerBoxes = append(v.InnerBoxes, box)
		box.OuterBox = v

		if box.OuterBox == nil {
			return nil
		}
	case *AnotherItem:
		// @TODO: Implement
		return fmt.Errorf("MoveTo AnotherItem is not implemented.")
	default:
		return fmt.Errorf("Can't move this box \"%v\" to \"%v\"", box.Label, other)
	}
	return nil
}

func (box *Box) MoveOutOfOtherBox() error {
	logg.Debug("Moving '", box.Label, "' out of '", box.OuterBox.Label, "'")
	outer := box.OuterBox

	if box == outer {
		return fmt.Errorf("Can't move this box \"%v\" out of itself \"%v\"", box.Label, outer.Label)
	}
	if outer == nil {
		return fmt.Errorf("Trying to move box '%v' out of OuterBox but OuterBox is nil", box.Label)
	}
	idx := slices.Index(outer.InnerBoxes, box)
	if idx == -1 {
		return fmt.Errorf("Trying to move box '%v' out of OuterBox but it doesn't have this box. '%v'", box.Label, outer.InnerBoxes)
	}
	outer.InnerBoxes = slices.Delete(outer.InnerBoxes, idx, idx+1)
	// Make slice nil that was the last element
	if len(outer.InnerBoxes) == 0 {
		outer.InnerBoxes = nil
	}

	box.OuterBox = nil

	return nil
}

type Movable interface {
	// MoveTo moves this instance inside another instance.
	// Currently 'other' works only with Box structs.
	// Returns an error if moving is not successful.
	MoveTo(other any) error
}

type BoxNode interface {
	// InnerBoxes represent boxes that are inside this current box.
	//
	// Returns a slice of pointers to (inner) boxes if it has other boxes inside, else it returns nil.
	// InnerBoxes() []*Box

	// OuterBox is the box where this current box is inside of.
	//
	// Returns a Box pointer to the (outer) box if it is inside that other box, else it returns nil.
	// OuterBox() *Box

	// Self returns this Box bnstance that implements the BoxNode interface.
	Self() *Box

	// Movable interface implements MoveTo functions.
	Movable

	// MoveOutOfOtherBox moves this box out.
	// Returns an error if moving out is not succesful.
	MoveOutOfOtherBox() error
}

type BoxDatabase interface {
	// CreateBox returns id of box if successful, otherwise error.
	CreateBox() (string, error)
	Box(id string) (Box, error)
	// BoxIDs returns IDs of all boxes.
	BoxIDs() ([]string, error)
	// MoveBox moves box with id1 into box with id2.
	MoveBox(id1 string, id2 string) error
}

// func (box *Box) OuterBox() *Box {
// 	return box.outerBox
// }
//
// func (box *Box) InnerBoxes() []*Box {
// 	return box.innerBoxes
// }
//
// func (box *Box) Self() *Box {
// 	return box
// }
//
// func MoveThisBoxToThatBox(this *Box, that *Box) error {
// 	if this == that {
// 		return fmt.Errorf("Can't move this box \"%v\" to that box \"%v\"", this.Label, that.Label)
// 	}
// 	that.InnerBoxes = append(that.InnerBoxes, this)
// 	return nil
// }
