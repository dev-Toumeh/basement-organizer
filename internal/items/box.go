package items

import (
	"basement/main/internal/logg"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/gofrs/uuid/v5"
)

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
	// Quantity    int64     `json:"quantity"    validate:"omitempty,numeric,gte=1"`
	// Weight      string    `json:"weight"      validate:"omitempty,numeric"`
	Id          uuid.UUID `json:"id"`
	Label       string    `json:"label"       validate:"required,lte=128"`
	Description string    `json:"description" validate:"omitempty,lte=256"`
	Picture     string    `json:"picture"     validate:"omitempty,base64"`
	QRcode      string    `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	Items       []*Item   `json:"items"`
	InnerBoxes  []*Box    `json:"innerboxes"`
	OuterBox    *Box      `json:"outerbox" `
}

type BoxTemplateData struct {
	// Quantity    int64     `json:"quantity"    validate:"omitempty,numeric,gte=1"`
	// Weight      string    `json:"weight"      validate:"omitempty,numeric"`
	Id          uuid.UUID `json:"id"`
	Label       string    `json:"label"       validate:"required,lte=128"`
	Description string    `json:"description" validate:"omitempty,lte=256"`
	Picture     string    `json:"picture"     validate:"omitempty,base64"`
	QRcode      string    `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	Items       []*Item   `json:"items"`
	InnerBoxes  []*Box    `json:"innerboxes"`
	OuterBox    *Box      `json:"outerbox" `
	Edit        bool
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
	return Box{
		Id:          uuid.Must(uuid.NewV4()),
		Label:       "Box",
		Description: "This box is empty.",
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
