package routes

import (
	"basement/main/internal/database"
	"fmt"
	"net/http"
)

func ApiReadItemsHandler(db *database.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprint(w, db.Items)
			return
		}
		w.Header().Add("Allowed", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ApiReadItemHandler(db *database.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := r.PathValue("id")
			fmt.Fprint(w, db.Items[id])
			return
		}
		w.Header().Add("Allowed", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
