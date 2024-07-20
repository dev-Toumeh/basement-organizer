package items

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid/v5"

	"basement/main/internal/auth"
	"basement/main/internal/database"
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
func CreateItemHandler(db *database.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createNewItem(w, r, db)
		} else if r.Method == http.MethodGet {
			generateAddItemForm(w, r)
		}
	}
}

func createNewItem(w http.ResponseWriter, r *http.Request, db *database.JsonDB) {

	var errorMessages []string
	newItem := item(r)

	if _, exist := db.ItemByLabel(newItem.Label); exist {
		responseGenerator(
			w,
			[]string{"<div> the Label you have chosen is already selected</div>"},
			false,
		)
	}

	if valiedItem, err := validateItem(newItem, &errorMessages); err != nil {
		responseGenerator(w, errorMessages, false)
	} else {

		fmt.Println(valiedItem)
		if err := db.AddItem(valiedItem); err != nil {
			responseGenerator(w, []string{"<div>we was not able to add the new Item please comeback later</div>"}, false)
		} else {
			responseGenerator(w, []string{"<div>The Item has been added successfully</div>"}, true)
		}
	}
	return
}

// this function will generate a new from if the request was from type GET
func generateAddItemForm(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	data := templates.PageTemplate{
		Title:         "Personal",
		Authenticated: authenticated,
		User:          auth.Username(r),
	}

	if err := templates.Render(w,templates.TEMPLATE_CREATE_ITEM_PAGE, data); err != nil {
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
func item(r *http.Request) database.Item {

	newId, _ := uuid.NewV4()
	newItem := database.Item{
		Id:          newId,
		Label:       r.PostFormValue(LABEL),
		Description: r.PostFormValue(DESCRIPTIO),
		Picture:     r.PostFormValue(PICTURE),
		Quantity:    json.Number(r.PostFormValue(QUANTITY)),
		Weight:      r.PostFormValue(WEIGHT),
		QRcode:      r.PostFormValue(QRCODE),
	}
	return newItem
}
