package main

import (
	"fmt"
	"internal/auth"
	"internal/util"
	"net/http"
)

func createDB() *auth.AuthJsonDB {
	db := auth.AuthJsonDB{&util.JsonDB{}}
	db.Connect("./internal/auth/users2.json")
	return &db
}

func main() {
	db := createDB()

	fmt.Println(db.User("alex"))
	fmt.Println(db.User("alx"))
	// fmt.Println(uuid.New())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	http.HandleFunc("/login", db.LoginHandler)

	http.HandleFunc("/register", auth.Register)

	http.ListenAndServe(":8000", nil)
}
