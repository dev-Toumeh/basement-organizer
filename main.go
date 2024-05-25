package main

import (
	"basement/main/internal/auth"
	"basement/main/internal/util"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var db auth.AuthDatabaseHandler
	var err error
	db, err = createDB()
	if err != nil {
		log.Fatalf("Can't create DB, shutting server down")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	http.HandleFunc("/login", auth.LoginPage)
	http.HandleFunc("/login/user", db.LoginHandler)
	http.HandleFunc("/register", db.RegisterHandler)

	http.ListenAndServe("localhost:8000", nil)
}

func createDB() (auth.AuthDatabaseHandler, error) {
	var db util.DBWithFileConnector
	db = &util.JsonDB{}
	err := db.Connect("./internal/auth/users2.json")
	if err != nil {
		log.Println("createDB() error", err)
		return nil, err
	}

	var authDBHandler auth.AuthDatabaseHandler
	authDBHandler = &auth.AuthJsonDB{db.(*util.JsonDB)}

	return authDBHandler, nil
}
