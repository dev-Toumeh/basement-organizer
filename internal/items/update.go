package items

import (
	"basement/main/internal/database"
	"fmt"
	"net/http"
)

func UpdateItemHandler(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateItem(w, r, db)
		} else if r.Method == http.MethodGet {
			generateAddItemForm(w, r)
		}
	}
}

func updateItem(w http.ResponseWriter, r *http.Request, db *database.DB) {

	var errorMessages []string
	updatedItem := item(r)

	if valiedItem, err := validateItem(updatedItem, &errorMessages); err != nil {
		responseGenerator(w, errorMessages, false)
	} else {
		ctx := context.TODO()
		if err := db.UpdateItem(ctx, valiedItem); err != nil {
			fmt.Println(err)
			responseGenerator(w, []string{"we was not able to update the Item please comeback later"}, false)
		} else {
			responseGenerator(w, []string{"the Item has been updated Successfully"}, true)
		}
	}
	return
}
