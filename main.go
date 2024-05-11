package main

import (
	"fmt"
	"internal/auth"
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

	http.HandleFunc("/login", auth.Login)

	http.HandleFunc("/register", auth.Register)

	http.ListenAndServe(":8000", nil)
}
