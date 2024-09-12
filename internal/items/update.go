package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"context"
	"fmt"
	"net/http"
)

// update the item based on ID
func UpdateItemHandler(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateItem(w, r, db)
		} else if r.Method == http.MethodGet {
			generateAddItemForm(w, r)
		}
	}
}

func updateItem(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	logg.Debug(r.URL)
	var errorMessages []string
	updatedItem, err := item(r)
	if err != nil {
		templates.RenderErrorNotification(w, "Error while generating the User please comeback later")
	}

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
