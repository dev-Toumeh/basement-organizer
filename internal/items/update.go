package items

import (
	"basement/main/internal/database"
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
	panic("unimplemented")
}
