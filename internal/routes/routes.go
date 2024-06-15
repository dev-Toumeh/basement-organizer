package routes

import (
	"net/http"

	"basement/main/internal/auth"
)

const (
	STATIC string = "/static/"
)

func RegisterRoutes(db auth.AuthDatabase) {
	http.Handle(STATIC, http.StripPrefix("/static/", http.FileServer(http.Dir("internal/static"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login", db.LoginHandler)
	http.HandleFunc("/register", db.RegisterHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/api/v1/create/item", CreateItem)
	http.HandleFunc("/api/v1/read/items", ReadItems)
	http.HandleFunc("/api/v1/read/item/id", ReadItem)
	http.HandleFunc("api/v1/update/item/id", UpdateItem)
	http.HandleFunc("api/v1/delete/itemh", DeleteItem)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func ReadItems(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func ReadItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

