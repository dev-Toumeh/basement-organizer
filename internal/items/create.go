package items

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"

	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
)

// BasicInfo is present in item, box, shelf and area
type BasicInfo struct {
	ID             uuid.UUID
	Label          string
	Description    string
	Picture        string
	PreviewPicture string
	QRCode         string
}

func (b BasicInfo) Map() map[string]any {
	return map[string]interface{}{
		"ID":             b.ID,
		"Label":          b.Label,
		"Description":    b.Description,
		"Picture":        b.Picture,
		"PreviewPicture": b.PreviewPicture,
		"QRCode":         b.QRCode,
	}
}

func NewBasicInfo() BasicInfo {
	return BasicInfo{ID: uuid.Must(uuid.NewV4())}.MakeLabelWithTime("thing")
}

func NewBasicInfoWithLabel(label string) BasicInfo {
	return BasicInfo{ID: uuid.Must(uuid.NewV4())}.MakeLabelWithTime(label)
}

func (b BasicInfo) MakeLabelWithTime(label string) BasicInfo {
	t := time.Now().Format("2006-01-02_15_04_05")
	b.Label = fmt.Sprintf("%s_%s", label, t)
	return b
}

// ListRow is a single row entry used for list templates.
type ListRow struct {
	ID             uuid.UUID
	Label          string
	BoxID          uuid.UUID
	BoxLabel       string
	ShelfID        uuid.UUID
	ShelfLabel     string
	AreaID         uuid.UUID
	AreaLabel      string
	PreviewPicture string
}

func (row *ListRow) Map() map[string]any {
	return map[string]interface{}{
		"ID":             row.ID,
		"Label":          row.Label,
		"BoxID":          row.BoxID,
		"BoxLabel":       row.BoxLabel,
		"ShelfID":        row.ShelfID,
		"ShelfLabel":     row.ShelfLabel,
		"AreaID":         row.AreaID,
		"AreaLabel":      row.AreaLabel,
		"PreviewPicture": row.PreviewPicture,
	}
}

type Item struct {
	BasicInfo
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
	ItemListRowByID(id uuid.UUID) (*ListRow, error)
	Item(id string) (Item, error)
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
	ItemFuzzyFinder(query string) ([]ListRow, error)
	ItemFuzzyFinderWithPagination(query string, limit, offset int) ([]ListRow, error)
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

var validate *validator.Validate

// this function will check the type of the request
// if it is from type post it will create the item otherwise it will generate the  create title  template
func CreateItemHandler(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createNewItem(w, r, db)
		} else {

		}
	}
}

func createNewItem(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	var responseMessage []string
	newItem, err := item(r)
	if err != nil {
		logg.Err(err)
		templates.RenderErrorNotification(w, "Error while generating the User please comeback later")
	}
	if valiedItem, err := validateItem(newItem, &responseMessage); err != nil {
		responseGenerator(w, responseMessage, false)
	} else {
		if err := db.CreateNewItem(valiedItem); err != nil {
			if err == db.ErrorExist() {
				responseMessage = append(responseMessage, "the Label is already token please choice another one")
				responseGenerator(w, responseMessage, false)
			} else {
				responseMessage = append(responseMessage, "Unable to add new item due to technical issues. Please try again later.")
				responseGenerator(w, responseMessage, false)
			}
		} else {
			responseMessage = append(responseMessage, "your Item was added successfully")
			responseGenerator(w, responseMessage, true)
		}
	}
}

// this function will validate the new Item  and will return ether true will a Struct full of data
// or false with an empty Struct
func validateItem(newItem Item, errorMessages *[]string) (Item, error) {

	if !env.Development() {
		validate = validator.New(validator.WithRequiredStructEnabled())

		if err := validate.Struct(newItem); err != nil {

			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, validationErr := range validationErrors {

					switch validationErr.Field() {
					case "Label":
						if validationErr.Tag() == "required" {
							*errorMessages = append(
								*errorMessages,
								"<div>The Label field is required but missing.</div>",
							)
						} else if validationErr.Tag() == "lte" {
							*errorMessages = append(*errorMessages, "<div>The Label field must be less than or equal to 15 characters.</div>")
						}
					case "Description":
						if validationErr.Tag() == "alphanum" {
							*errorMessages = append(
								*errorMessages,
								"<div>The Description must contain only alphanumeric characters.</div>",
							)
						} else if validationErr.Tag() == "lte" {
							*errorMessages = append(*errorMessages, "<div>The Description must be less than or equal to 255 characters.</div>")
						}
					case "Picture":
						if validationErr.Tag() == "base64" {
							*errorMessages = append(
								*errorMessages,
								"<div>The Picture must be a valid Base64 string.</div>",
							)
						}
					case "Quantity":
						if validationErr.Tag() == "numeric" {
							*errorMessages = append(
								*errorMessages,
								"<div>The Quantity must be a number.</div>",
							)
						} else if validationErr.Tag() == "gte" {
							*errorMessages = append(*errorMessages, "<div>The Quantity must be greater than or equal to 1.</div>")
						} else if validationErr.Tag() == "lte" {
							*errorMessages = append(*errorMessages, "<div>The Quantity must be less than or equal to 15.</div>")
						}
					case "Weight":
						if validationErr.Tag() == "numeric" {
							*errorMessages = append(
								*errorMessages,
								"<div>The Weight must be a number.</div>",
							)
						}
					case "QRCode":
						if validationErr.Tag() == "alphanumunicode" {
							*errorMessages = append(
								*errorMessages,
								"<div>The QRCode must contain only alphanumeric characters and unicode.</div>",
							)
						}
					default:
						*errorMessages = append(
							*errorMessages,
							fmt.Sprintf(
								"Field '%s' is invalid: %s",
								validationErr.Field(),
								validationErr.Tag(),
							),
						)
					}
				}
			} else {
				// Other errors
				*errorMessages = append(*errorMessages, err.Error())
			}

			log.Println("create Item Validation failed")
			err := errors.New("validation failed")
			return Item{}, err

		} else {

			log.Println("create Item Validation succeeded")
			return newItem, nil
		}
	}
	return newItem, nil
}

func responseGenerator(w http.ResponseWriter, responseMessage []string, success bool) {
	var htmlstring string
	if success {
		htmlstring = strings.Join(responseMessage, "")

	} else {
		htmlstring = fmt.Sprintf(
			"please check the following errors and try again: </br> %s",
			strings.Join(responseMessage, " "),
		)
	}
	tmp, err := template.New("div").Parse(htmlstring)

	if err != nil {
		fmt.Fprintln(w, "something wrong happened please comeback later")
	}
	tmp.Execute(w, nil)
}

// this function will pack the request into struct from type Item, so it will be easier to handle it
func item(r *http.Request) (Item, error) {

	id, boxId, err := common.CheckIDs(r.PostFormValue(ID), r.PostFormValue(BOX_ID))
	if err != nil {
		return Item{}, err
	}

	newItem := Item{
		BasicInfo: BasicInfo{
			ID:          id,
			Label:       r.PostFormValue(LABEL),
			Description: r.PostFormValue(DESCRIPTION),
			Picture:     common.ParsePicture(r),
			QRCode:      r.PostFormValue(QRCODE),
		},
		Quantity: common.ParseQuantity(r.PostFormValue(QUANTITY)),
		Weight:   r.PostFormValue(WEIGHT),
		BoxID:    boxId,
		ShelfID:  uuid.FromStringOrNil(r.PostFormValue(SHELF_ID)),
		AreaID:   uuid.FromStringOrNil(r.PostFormValue(AREA_ID)),
	}
	return newItem, nil
}
