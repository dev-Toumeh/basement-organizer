package auth

import (
	"basement/main/internal/database"
	"basement/main/internal/templates"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	REGISTER_FAILED_MESSAGE string = "register failed"
	USERNAME                string = "username"
	PASSWORD                string = "password"
	COOKIE_NAME             string = "mycookie"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// this function will check the type of the request
// if it is from type post it will register the user otherwise it will generate the register template
func RegisterHandler(db *database.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			registerUser(w, r, db)
		}

		generateRegisterPage(w, r)
	}
}

func generateRegisterPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := Authenticated(r)
	data := templates.PageTemplate{
		Title:         "Register",
		Authenticated: authenticated,
		User:          Username(r),
	}

	if err := templates.ApplyPageTemplate(w, templates.REGISTER_TEMPLATE_FILE_WITH_PATH, data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

func registerUser(w http.ResponseWriter, r *http.Request, db *database.JsonDB) {
	NewUsername := r.PostFormValue(USERNAME)
	NewPassword := r.PostFormValue(PASSWORD)

	if NewUsername == "" {
		log.Println("Missing username")
		fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
		return
	}
	if NewPassword == "" {
		log.Println("Missing password")
		fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
		return
	}

	//check the if the username is exist
	_, exist := db.User(NewUsername)

	if exist {
		log.Printf("the user %s is already exist", NewUsername)
		fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
		return
	}

	// hash the password
	NewHashedPassword, err := hashPassword(NewPassword)
	if err != nil {
		log.Fatal(err)
		fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
	}

	// add the new user to the Databse
	err = database.AddUser(NewUsername, NewHashedPassword, db)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
		return
	}



	 log.Printf("User %s registered successfully:", NewUsername)
	fmt.Fprintln(w, "User registered successfully")

	return
}
