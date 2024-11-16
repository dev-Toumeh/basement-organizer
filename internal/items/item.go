package items

import (
	"basement/main/internal/common"
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type Item struct {
	common.BasicInfo
	Quantity int64     `json:"quantity"    validate:"omitempty,numeric,gte=1"`
	Weight   string    `json:"weight"      validate:"omitempty,numeric"`
	BoxID    uuid.UUID `json:"box_id"`
	ShelfID  uuid.UUID `json:"shelf_id"`
	AreaID   uuid.UUID `json:"area_id"`
}

func (i Item) String() string {
	return fmt.Sprintf("Item[ID=%s, Label=%s, QRCode=%s, Quantity=%d, Weight=%s, BoxID=%s, ShelfID=%s, AreaID=%s]",
		i.BasicInfo.ID, i.BasicInfo.Label, i.BasicInfo.QRCode, i.Quantity, i.Weight, i.BoxID, i.ShelfID, i.AreaID)
}

type ItemDatabase interface {
	CreateNewItem(newItem Item) error
	ItemByField(field string, value string) (Item, error)
	ItemListRowByID(id uuid.UUID) (*common.ListRow, error)
	ItemById(id uuid.UUID) (Item, error)
	ItemIDs() ([]string, error)
	ItemExist(field string, value string) bool
	Items() ([][]string, error)
	UpdateItem(ctx context.Context, item Item) error
	DeleteItem(itemId uuid.UUID) error
	DeleteItems(itemId []uuid.UUID) error
	InsertSampleItems()
	ErrorExist() error
	MoveItemToBox(itemID uuid.UUID, boxID uuid.UUID) error

	// search functions
	ItemFuzzyFinder(query string) ([]common.ListRow, error)
	ItemFuzzyFinderWithPagination(query string, limit, offset int) ([]common.ListRow, error)
	NumOfItemRecords(searchString string) (int, error)
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
	SHELF_ID    string = "shelf_id"
	AREA_ID     string = "area_id"
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
		"Picture":        s.Picture,
		"PreviewPicture": s.PreviewPicture,
		"BoxID":          uuid.Nil,
		"ShelfID":        uuid.Nil,
		"AreaID":         uuid.Nil,
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
