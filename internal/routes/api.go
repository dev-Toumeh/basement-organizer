package routes

import (
	"basement/main/internal/auth"
	"fmt"
	"net/http"
)

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
