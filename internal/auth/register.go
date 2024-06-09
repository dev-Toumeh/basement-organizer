package auth

import (
	"basement/main/internal/templates"
	"basement/main/internal/util"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"text/template"
)

const (
	REGISTER_FAILED_MESSAGE string = "register failed"
	REGISTER_TEMPLATE_PATH  string = "internal/templates/register.html"
	USERNAME                string = "username"
	PASSWORD                string = "password"
	COOKIE_NAME             string = "mycookie"
)

// this function will check the type of the request 
type AuthDatabase interface {
	User(string) (util.DBUser2, bool)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	RegisterHandler(w http.ResponseWriter, r *http.Request)
}

type AuthJsonDB struct {
	util.JsonDB
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// this function will check the type of the request
// if it is from type post it will register the user otherwise it will generate the register template
func (db *AuthJsonDB) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db.registerUser(w, r)
	}

	generateRegisterPage(w)
}

func generateRegisterPage(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles(templates.PAGE_TEMPLATE_FILE, "internal/templates/register.html")
	if err != nil {
		log.Printf("%v or %v: %v\n", templates.PAGE_TEMPLATE, "register.html", err)
		fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
		return
	}

	templateData := templates.PageTemplate{Title: "Register"}

	if err := tmpl.ExecuteTemplate(w, templates.PAGE_TEMPLATE, templateData); err != nil {
		log.Printf("Error executing  register Template: %v", err)
		http.Error(w, "Error rendering  register page", http.StatusInternalServerError)
	}
}

func (db *AuthJsonDB) registerUser(w http.ResponseWriter, r *http.Request) {
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
	err = db.AddUser(NewUsername, NewHashedPassword)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
		fmt.Fprintln(w, "the problem is inside of the addUser function")
		return
	}

	log.Printf("User %s registered successfully:", NewUsername)
	fmt.Fprintln(w, "User registered successfully")

	return
}
