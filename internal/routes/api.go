package routes

import (
	"basement/main/internal/database"
	"fmt"
	"net/http"
)

func ReadItems(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func ReadItem(w http.ResponseWriter, r *http.Request) {
	db, _ := CreateItemsJsonDB()
	// get "Water Bottle"
	fmt.Fprint(w, db.Items["123e4567-e89b-12d3-a456-426614174002"])
	return
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
