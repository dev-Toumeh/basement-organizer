package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
)

const (
	STATIC           string = "/static/"
	ITEMS_FILE_PATH  string = "internal/auth/items.json"
	USERS_FILE_PATH  string = "internal/auth/users2.json"
	API_V1_READ_ITEM string = "/api/v1/read/item/id"
)

func RegisterRoutes(db *auth.JsonDB) {
	http.Handle(STATIC, http.StripPrefix("/static/", http.FileServer(http.Dir("internal/static"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login", db.LoginHandler)
	http.HandleFunc("/register", auth.RegisterHandler(db))
	http.HandleFunc("/logout", auth.LogoutHandler)

	http.HandleFunc("/api/v1/create/item", CreateItem)
	http.HandleFunc("/api/v1/read/items", ReadItems)
	http.HandleFunc(API_V1_READ_ITEM, ReadItem)
	http.HandleFunc("/api/v1/update/item/id", UpdateItem)
	http.HandleFunc("/api/v1/delete/item", DeleteItem)
}

// CreateItemsJsonDB creates DB instance by reading or creating "items.json" file from disk.
func CreateItemsJsonDB() (*auth.JsonDB, error) {
	db := &auth.JsonDB{}
	db.InitFieldFromFile(ITEMS_FILE_PATH, &db.Items)
	db.InitFieldFromFile(USERS_FILE_PATH, &db.Users)
	// sss := reflect.VisibleFields(reflect.TypeOf(*db))
	// for _, v := range sss {
	// 	log.Println(v.Type)
	// }
	return db, nil
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

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
