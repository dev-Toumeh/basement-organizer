package auth

import (
	"basement/main/internal/templates"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
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
func RegisterHandler(db *JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			db.registerUser(w, r)
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

func (db *JsonDB) registerUser(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	log.Printf("User %s registered successfully:", NewUsername)
	fmt.Fprintln(w, "User registered successfully")

	return
}

// this function will check if there is existing user withe same name and if not it will
// create new one at the end it will save it
func (db *JsonDB) AddUser(username string, passwordHash string) error {
	if dbUser, exist := db.User(username); exist {
		return fmt.Errorf("user %s already exists", username)
	} else {
		dbUser.Uuid = uuid.New()
		dbUser.PasswordHash = passwordHash
		db.Users[username] = dbUser
	}

	return db.save()
}

// this function is responsible of saving the new Record inside of the Database (user2.json)
func (db *JsonDB) save() error {

	_, err := db.File.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("Error seeking to start of file: %v", err)
		return err
	}

	encoder := json.NewEncoder(db.File)

	err = encoder.Encode(db.Users)
	if err != nil {
		log.Printf("Error encoding users to JSON: %v", err)
		return err
	}

	currentPos, err := db.File.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Printf("Error getting current file position: %v", err)
		return err
	}

	err = db.File.Truncate(currentPos)
	if err != nil {
		log.Printf("Error truncating file: %v", err)
		return err
	}

	return nil
}
