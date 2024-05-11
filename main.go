package main

import (
	"fmt"
	"internal/auth"
	"internal/util"
	"log"
	"net/http"
)

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

func main() {
	var db auth.AuthDatabaseHandler
	var err error
	db, err = createDB()
	if err != nil {
		log.Fatalf("Can't create DB, shutting server down")
	}

	fmt.Println(db.User("alex"))
	fmt.Println(db.User("alx"))
	// fmt.Println(uuid.New())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	http.HandleFunc("/login", db.LoginHandler)
	// Change auth.Register to db.RegisterHandler
	http.HandleFunc("/register", auth.Register)

	http.ListenAndServe(":8000", nil)
}
