package items

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"

	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
)

type Item struct {
	ID             uuid.UUID `json:"id"`
	Label          string    `json:"label"       validate:"required,lte=128"`
	Description    string    `json:"description" validate:"omitempty,lte=256"`
	Picture        string    `json:"picture"     validate:"omitempty,base64"`
	PreviewPicture string    `json:"preview_picture"     validate:"omitempty,base64"`
	Quantity       int64     `json:"quantity"    validate:"omitempty,numeric,gte=1"`
	Weight         string    `json:"weight"      validate:"omitempty,numeric"`
	QRcode         string    `json:"qrcode"      validate:"omitempty,alphanumunicode"`
	BoxID          uuid.UUID `json:"box_id"`
	ShelfID        uuid.UUID `json:"shelf_id"`
	AreaID         uuid.UUID `json:"area_id"`
}

type ItemDatabase interface {
	CreateNewItem(newItem Item) error
	ItemByField(field string, value string) (Item, error)
	ItemListRowByID(id uuid.UUID) (*ItemListRow, error)
	Item(id string) (Item, error)
	ItemIDs() ([]string, error)
	ItemExist(field string, value string) bool
	Items() ([][]string, error)
	UpdateItem(ctx context.Context, item Item) error
	DeleteItem(itemId uuid.UUID) error
	DeleteItems(itemId []uuid.UUID) error
	InsertDummyItems()
	ErrorExist() error
	MoveItem(id1 uuid.UUID, id2 uuid.UUID) error

	// search functions
	ItemFuzzyFinder(query string) ([]ItemListRow, error)
	ItemFuzzyFinderWithPagination(query string, limit, offset int) ([]ItemListRow, error)
	NumOfItemRecords(searchString string) (int, error)
}

const (
	ID         string = "id"
	LABEL      string = "label"
	DESCRIPTIO string = "description"
	PICTURE    string = "picture"
	QUANTITY   string = "quantity"
	WEIGHT     string = "weight"
	QRCODE     string = "qrcode"
	BOX_ID     string = "Box_id"
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
					case "QRcode":
						if validationErr.Tag() == "alphanumunicode" {
							*errorMessages = append(
								*errorMessages,
								"<div>The QRcode must contain only alphanumeric characters and unicode.</div>",
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

	id, boxId, err := checkIDs(r)
	if err != nil {
		return Item{}, err
	}
	// b64encodedPictureString := parsePicture(r)

	newItem := Item{
		ID:          id,
		Label:       r.PostFormValue(LABEL),
		Description: r.PostFormValue(DESCRIPTIO),
		Picture:     "", //b64encodedPictureString, // @TODO: Fix picture is not added while creating or updating item.
		Quantity:    parseQuantity(r.PostFormValue(QUANTITY)),
		Weight:      r.PostFormValue(WEIGHT),
		QRcode:      r.PostFormValue(QRCODE),
		BoxID:       boxId,
	}
	return newItem, nil
}

// ParsePicture returns base64 encoded string of picture uploaded if there is any
func ParsePicture(r *http.Request) string {
	logg.Info("Parsing multipart/form-data for picture")
	// 8 MB
	var maxSize int64 = 1000 * 1000 * 8
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		logg.Err(err)
		return ""
	}

	file, header, err := r.FormFile(PICTURE)
	if header != nil {
		logg.Debug("picture filename:", header.Filename)
	}
	if err != nil {
		logg.Err(err)
		return ""
	}

	readbytes, err := io.ReadAll(file)
	logg.Debug("picture size:", len(readbytes)/1000, "KB")
	if err != nil {
		logg.Err(err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(readbytes)
}

// parseQuantity returns by default at least 1
func parseQuantity(value string) int64 {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 1
	}
	return intValue
}

func stringToInt64(value string) int64 {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("Error converting string to int64: %v", err)
		return 0
	}
	return intValue
}

func checkIDs(r *http.Request) (uuid.UUID, uuid.UUID, error) {
	var err error
	id := uuid.Nil
	boxId := uuid.Nil

	lengID := len(r.PostFormValue(ID))
	lengBoxID := len(r.PostFormValue(BOX_ID))

	if lengBoxID != 0 {
		boxId, err = uuid.FromString(r.PostFormValue(BOX_ID))
		if err != nil {
			logg.Errf("error while converting the boxId to type uuid: %v", err)
		}
	}

	if lengID == 0 {
		id, err = uuid.NewV4()
		if err != nil {
			return id, boxId, fmt.Errorf("error while generating the new item uuid: %w", err)
		}
		return id, boxId, nil
	} else {
		id, err = uuid.FromString(r.PostFormValue(ID))
		if err != nil {
			return id, boxId, fmt.Errorf("error while converting the itemId to type uuid: %w", err)
		}
		return id, boxId, nil
	}
}
