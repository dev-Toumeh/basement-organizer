package auth

import (
	"fmt"
	"basement/main/internal/util"
	"log"
	"net/http"
)

const LOGIN_FAILED_MESSAGE string = "Login failed"

type AuthDatabaseHandler interface {
	User(string) (util.DBUser2, bool)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	RegisterHandler(w http.ResponseWriter, r *http.Request)
}

type AuthJsonDB struct {
	*util.JsonDB
}

func (db *AuthJsonDB) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Printf("%v %v ", r.URL, r.Form.Encode())

	if username == "" {
		log.Println("Missing username")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}
	if password == "" {
		log.Println("Missing password")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	user,_ := db.User(username)

	if !util.CheckPasswordHash(password, user.PasswordHash) {
		log.Println("pw hash doesnt match")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	log.Println("login successful")
	fmt.Fprintf(w, "Welcome %v\n", username)
}
