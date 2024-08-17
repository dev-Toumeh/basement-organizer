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

	"basement/main/internal/auth"
	"basement/main/internal/database"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
)

const (
	ID         string = "id"
	LABEL      string = "label"
	DESCRIPTIO string = "description"
	PICTURE    string = "picture"
	QUANTITY   string = "quantity"
	WEIGHT     string = "weight"
	QRCODE     string = "qrcode"
)

var validate *validator.Validate

// this function will check the type of the request
// if it is from type post it will create the item otherwise it will generate the  create title  template
func CreateItemHandler(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createNewItem(w, r, db)
		} else if r.Method == http.MethodGet {
			generateAddItemForm(w, r)
		}
	}
}

func createNewItem(w http.ResponseWriter, r *http.Request, db *database.DB) {
	logg.Debug(r.URL)
	var responseMessage []string
	newItem := item(r)

	if valiedItem, err := validateItem(newItem, &responseMessage); err != nil {
		responseGenerator(w, responseMessage, false)
	} else {
		ctx := context.TODO()
		if err := db.CreateNewItem(ctx, valiedItem); err != nil {
			if err == database.ErrExist {
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

// this function will generate a new from if the request was from type GET
func generateAddItemForm(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	data := templates.PageTemplate{
		Title:         "Personal",
		Authenticated: authenticated,
		User:          auth.Username(r),
	}

	if err := templates.Render(w, templates.TEMPLATE_CREATE_ITEM_PAGE, data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

// this function will validate the new Item  and will return ether true will a Struct full of data
// or false with an empty Struct
func validateItem(newItem database.Item, errorMessages *[]string) (database.Item, error) {

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
		return database.Item{}, err

	} else {

		log.Println("create Item Validation succeeded")
		return newItem, nil
	}
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
// @TODO: Need error return in case something wrong happens
func item(r *http.Request) database.Item {
	var id uuid.UUID
	if updatedId, err := uuid.FromString(r.PostFormValue(ID)); err != nil {
		id, _ = uuid.NewV4()
	} else {
		id = updatedId
	}
	logg.Debug("Creating item id:", id)
	logg.Debug("Content-Type:", r.Header.Get("Content-Type"))

	b64encodedPictureString := parsePicture(r)

	newItem := database.Item{
		Id:          id,
		Label:       r.PostFormValue(LABEL),
		Description: r.PostFormValue(DESCRIPTIO),
		Picture:     b64encodedPictureString,
		Quantity:    parseQuantity(r.PostFormValue(QUANTITY)),
		Weight:      r.PostFormValue(WEIGHT),
		QRcode:      r.PostFormValue(QRCODE),
	}
	return newItem
}

// parsePicture returns base64 encoded string of picture uploaded if there is any
func parsePicture(r *http.Request) string {
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
		logg.Info("Could not parse quantity. Set quantity to 1")
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
