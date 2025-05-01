package items

import (
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/validate"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type Item struct {
	common.BasicInfo
	Quantity   int64
	Weight     float64
	BoxID      uuid.UUID
	BoxLabel   string
	ShelfID    uuid.UUID
	ShelfLabel string
	AreaID     uuid.UUID
	AreaLabel  string
}

func (i Item) String() string {
	return fmt.Sprintf("Item[ID=%s, Label=%s, Quantity=%d, Weight=%f, BoxID=%s, BoxLabel=%s, ShelfID=%s, ShelfLabel=%s, AreaID=%s, AreaLabel=%s]",
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
	UpdateItem(item Item, ignorePicture bool, pictureFormat string) error
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
func (s *Item) Map() map[string]any {
	shelfMap := map[string]any{
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
		Weight:    1.00,
		BoxID:     uuid.Nil,
		ShelfID:   uuid.Nil,
		AreaID:    uuid.Nil,
	}
}

// Converts a validated ItemValidate struct into a clean Item struct used for storage
// or business logic (e.g., database insertion).
func ToItem(validatedItem validate.ItemValidate) Item {
	item := Item{
		BasicInfo: common.BasicInfo{
			ID:          validatedItem.ID.UUID(),
			Label:       validatedItem.Label.String(),
			Description: validatedItem.Description.String(),
			Picture:     validatedItem.Picture.String(),
		},
		Quantity: validatedItem.Quantity.Int(),
		Weight:   validatedItem.Weight.Float64(),
		BoxID:    validatedItem.BoxID.UUID(),
		ShelfID:  validatedItem.ShelfID.UUID(),
		AreaID:   validatedItem.AreaID.UUID(),
	}
	return item
}

// Parses form input from an HTTP request, builds a validate.ItemValidate struct, and runs field-level validation.
// Returns the validator with error messages if any validations fail.
func ValidateItem(r *http.Request, w http.ResponseWriter) (validate.Validate, error) {
	item := validate.ItemValidate{
		BasicInfoValidate: validate.BasicInfoValidate{
			ID:          validate.NewUUIDField(r.PostFormValue(ID)),
			Label:       validate.NewStringField(r.PostFormValue(LABEL)),
			Description: validate.NewStringField(r.PostFormValue(DESCRIPTION)),
			Picture:     validate.NewStringField(common.ParsePicture(r)),
		},
		Quantity: validate.NewIntField(r.PostFormValue(QUANTITY)),
		Weight:   validate.NewFloatField(r.PostFormValue(WEIGHT)),
		BoxID:    validate.NewUUIDField(r.PostFormValue(BOX_ID)),
		ShelfID:  validate.NewUUIDField(r.PostFormValue(SHELF_ID)),
		AreaID:   validate.NewUUIDField(r.PostFormValue(AREA_ID)),
	}
	logg.DebugJSONLite(item.Map(), 50)

	validator := validate.Validate{Item: item}

	if err := validator.ValidateItem(w, item); err != nil {
		return validator, err
	}
	if validator.HasValidateErrors() {
		return validator, validator.Err()
	}

	return validator, nil
}
