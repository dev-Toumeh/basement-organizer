package items

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"

	"basement/main/internal/auth"
	"basement/main/internal/templates"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

var validate *validator.Validate

// this function will check the type of the request
// if it is from type post it will create the item otherwise it will generate the  create title  template
func CreateItemHandler(db *auth.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			 createNewItem(w, r, db)
		}

    		generateAddItemForm(w, r)
	}
}

func  createNewItem(w http.ResponseWriter, r *http.Request, db *auth.JsonDB) {
    validateItem(r) 
		return
}

func generateAddItemForm(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	data := templates.PageTemplate{
		Title:         "Personal",
		Authenticated: authenticated,
    User: auth.Username(r),
	}

	if err := templates.ApplyPageTemplate(w, templates.CREATE_ITEM_TEMPLATE_FILE_WITH_PATH, data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}}

// this function will validate the item request and will return ether true will a Struct full of data 
// or false with an empty Struct
func validateItem( r *http.Request) auth.Item {

	validate = validator.New(validator.WithRequiredStructEnabled())
	newItem := auth.Item{
		Id:          uuid.New(),
		Label:       "ExampleItem",
		Description: "ExampleDescription",
		Picture:     "iVBORw0KGgoAAAANSUhEUgAAAAUA",
		Quantity:    json.Number("10"),
		Weight:      "5",
		QRcode:      "QRCODE123",
	}

	// Validate the item
	if err := validate.Struct(newItem); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation succeeded")
	}
	return newItem

}
