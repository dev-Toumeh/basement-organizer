package items

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

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
	_, err := validateItem(r)
	if err != nil {
		htmlstring := fmt.Sprintf("<dev>%s</dev>", err)
		tmp, err := template.New("dev").Parse(htmlstring)
		if err != nil {
			fmt.Fprintln(w, "something wrong happened please comeback later")
		}
		tmp.Execute(w, nil)
	}

	return
}

func generateAddItemForm(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	data := templates.PageTemplate{
		Title:         "Personal",
		Authenticated: authenticated,
		User:          auth.Username(r),
	}

	if err := templates.ApplyPageTemplate(w, templates.CREATE_ITEM_TEMPLATE_FILE_WITH_PATH, data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

// this function will validate the item request and will return ether true will a Struct full of data
// or false with an empty Struct
func validateItem(r *http.Request) (database.Item, error) {

	validate = validator.New(validator.WithRequiredStructEnabled())
	newItem := database.Item{
		Id:          uuid.New(),
		Label:       r.PostFormValue(LABEL),
		Description: r.PostFormValue(DESCRIPTIO),
		Picture:     r.PostFormValue(PICTURE),
		Quantity:    json.Number(r.PostFormValue(QUANTITY)),
		Weight:      r.PostFormValue(WEIGHT),
		QRcode:      r.PostFormValue(QRCODE),
	}

	// Validate the item
	if err := validate.Struct(newItem); err != nil {
		log.Println("Validation failed:", err)
		return database.Item{}, err
	}

	log.Println("Validation succeeded")
	return newItem, nil

}
