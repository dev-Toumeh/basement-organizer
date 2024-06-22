package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
)

const (
	STATIC           string = "/static/"
	API_V1_READ_ITEM string = "/api/v1/read/item/{id}"
)

func RegisterRoutes(db *auth.JsonDB) {
	http.Handle(STATIC, http.StripPrefix("/static/", http.FileServer(http.Dir("internal/static"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login", db.LoginHandler)
	http.HandleFunc("/register", auth.RegisterHandler(db))
	http.HandleFunc("/logout", auth.LogoutHandler)

	http.HandleFunc("/api/v1/create/item", CreateItem)
	http.HandleFunc("/api/v1/read/items", ReadItems)
	http.HandleFunc(API_V1_READ_ITEM, ReadItem(db))
	http.HandleFunc("/api/v1/update/item/{id}", UpdateItem)
	http.HandleFunc("/api/v1/delete/item/{id}", DeleteItem)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func ReadItems(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func ReadItem(db *auth.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, db.Items[r.PathValue("id")])
	}
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
