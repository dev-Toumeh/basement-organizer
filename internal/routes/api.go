package routes

import (
	"basement/main/internal/database"
	"fmt"
	"net/http"
)

func ReadItems(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func ApiReadItemHandler(db *database.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := r.PathValue("id")
			fmt.Fprint(w, db.Items[id])
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
