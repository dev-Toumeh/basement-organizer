package items

import (
	"basement/main/internal/database"
	"fmt"
	"net/http"
)

func UpdateItemHandler(db *database.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateItem(w, r, db)
		} else if r.Method == http.MethodGet {
			generateAddItemForm(w, r)
		}
	}
}

func updateItem(w http.ResponseWriter, r *http.Request, db *database.JsonDB) {

	var errorMessages []string
	updatedItem := item(r)

	if valiedItem, err := validateItem(updatedItem, &errorMessages); err != nil {
		responseGenerator(w, errorMessages, false)
	} else {
		if responseMessage, err := db.UpdateItem(valiedItem); err != nil {
			fmt.Println(err)
			responseGenerator(w, responseMessage, false)
		} else {
			responseGenerator(w, responseMessage, true)
		}
	}
	return
}
