package items

import (
	"basement/main/internal/common"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type Item struct {
	common.BasicInfo
	Quantity   int64     `json:"quantity"    validate:"omitempty,numeric,gte=1"`
	Weight     string    `json:"weight"      validate:"omitempty,numeric"`
	BoxID      uuid.UUID `json:"box_id"`
	BoxLabel   string    `json:"box_label"`
	ShelfID    uuid.UUID `json:"shelf_id"`
	ShelfLabel string    `json:"shelf_label"`
	AreaID     uuid.UUID `json:"area_id"`
	AreaLabel  string    `json:"area_label"`
}

func (i Item) String() string {
	return fmt.Sprintf("Item[ID=%s, Label=%s, Quantity=%d, Weight=%s, BoxID=%s, BoxLabel=%s, ShelfID=%s, ShelfLabel=%s, AreaID=%s, AreaLabel=%s]",
		i.ID, i.Label, i.Quantity, i.Weight, i.BoxID, i.BoxLabel, i.ShelfID, i.ShelfLabel, i.AreaID, i.AreaLabel)
}

type ItemDatabase interface {
	CreateNewItem(newItem Item) error
	ItemByField(field string, value string) (Item, error)
	ItemListRowByID(id uuid.UUID) (*common.ListRow, error)
	ItemById(id uuid.UUID) (*Item, error)
	ItemIDs() ([]uuid.UUID, error)
	ItemExist(field string, value string) bool
	Items() ([][]string, error)
	UpdateItem(item Item, ignorePicture bool) error
	DeleteItem(itemId uuid.UUID) error
	DeleteItems(itemId []uuid.UUID) error
	InsertSampleItems()
	ErrorExist() error
	MoveItemToBox(itemID uuid.UUID, boxID uuid.UUID) error
	MoveItemToShelf(itemID uuid.UUID, shelfID uuid.UUID) error
	MoveItemToArea(itemID uuid.UUID, areaID uuid.UUID) error

	// search functions
	ItemListCounter(queryString string) (count int, err error)
	ItemListRows(searchString string, limit int, pageNr int) (shelfRows []common.ListRow, err error)

	// required in common.Database interface
	InnerListRowsFrom2(belongsToTable string, belongsToTableID uuid.UUID, listRowsTable string) ([]common.ListRow, error)
	InnerListRowsPaginatedFrom(belongsToTable string, belongsToTableID uuid.UUID, listRowsTable string, searchQuery string, limit int, page int) (listRows []common.ListRow, err error)
	InnerBoxInBoxListCounter(searchString string, inTable string, inTableID uuid.UUID) (count int, err error)
	InnerShelfInTableListCounter(searchString string, inTable string, inTableID uuid.UUID) (count int, err error)
	InnerThingInTableListCounter(searchString string, thing int, inTable string, inTableID uuid.UUID) (count int, err error)
	MoveShelfToArea(shelfID uuid.UUID, toAreaID uuid.UUID) error
	BoxListCounter(searchQuery string) (count int, err error)
	ShelfListCounter(searchQuery string) (count int, err error)
	ShelfListRows(searchQuery string, limit int, page int) (shelfRows []common.ListRow, err error)
	AreaListCounter(searchQuery string) (count int, err error)
	BoxListRows(searchQuery string, limit int, page int) ([]common.ListRow, error)
	AreaListRows(searchQuery string, limit int, page int) (areaRows []common.ListRow, err error)
	DeleteBox(boxID uuid.UUID) error
	DeleteShelf(id uuid.UUID) (label string, err error)
	DeleteShelf2(id uuid.UUID) error
	DeleteArea(areaID uuid.UUID) error
}

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

const (
	ITEM_PAGE_TEMPLATE   string = "item-page-template"
	ITEM_CREATE_TEMPLATE string = "item-create-template"
)

// return the Item in type map
func (s *Item) Map() map[string]interface{} {
	shelfMap := map[string]interface{}{
		"ID":             s.ID,
		"Label":          s.Label,
		"Description":    s.Description,
		"Weight":         s.Weight,
		"Quantity":       s.Quantity,
		"Picture":        s.Picture,
		"PreviewPicture": s.PreviewPicture,
		"BoxID":          s.BoxID,
		"BoxLabel":       s.BoxLabel,
		"ShelfID":        s.ShelfID,
		"ShelfLabel":     s.ShelfLabel,
		"AreaID":         s.AreaID,
		"AreaLabel":      s.AreaLabel,
	}

	return shelfMap
}

// return new Shelf with default Values
func newItem() *Item {
	s := common.NewBasicInfoWithLabel("Item")
	return &Item{
		BasicInfo: s,
		Quantity:  1,
		Weight:    "2",
		BoxID:     uuid.Nil,
		ShelfID:   uuid.Nil,
		AreaID:    uuid.Nil,
	}
}

// this function will pack the request into struct from type Item, so it will be easier to handle it
func item(r *http.Request) (Item, error) {
	newItem := Item{
		BasicInfo: common.BasicInfo{
			ID:          uuid.FromStringOrNil(r.PostFormValue(ID)),
			Label:       r.PostFormValue(LABEL),
			Description: r.PostFormValue(DESCRIPTION),
			Picture:     common.ParsePicture(r),
		},
		Quantity:   common.ParseQuantity(r.PostFormValue(QUANTITY)),
		Weight:     r.PostFormValue(WEIGHT),
		BoxID:      uuid.FromStringOrNil(r.PostFormValue(BOX_ID)),
		BoxLabel:   r.PostFormValue(BOX_LABEL),
		ShelfID:    uuid.FromStringOrNil(r.PostFormValue(SHELF_ID)),
		ShelfLabel: r.PostFormValue(SHELF_LABEL),
		AreaID:     uuid.FromStringOrNil(r.PostFormValue(AREA_ID)),
		AreaLabel:  r.PostFormValue(AREA_LABEL),
	}
	return newItem, nil
}
