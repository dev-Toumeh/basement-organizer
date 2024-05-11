package main

import (
	"fmt"
	// "github.com/google/uuid"
	"internal/util"
	"net/http"
)

func main() {
	db := util.JsonDB{}
	db.Connect("./users2.json")
	fmt.Println(db.User("alex"))
	fmt.Println(db.User("alx"))
	// fmt.Println(uuid.New())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	http.HandleFunc("/login", login)

	http.HandleFunc("/register", register)

	http.ListenAndServe(":8000", nil)
}
