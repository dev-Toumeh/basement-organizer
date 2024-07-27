package auth

import (
	"basement/main/internal/database"
	"basement/main/internal/templates"
	"context"
	"database/sql"
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
func RegisterHandler(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
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

	if err := templates.Render(w, templates.TEMPLATE_REGISTER_PAGE, data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

func registerUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
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
	NewHashedPassword, err := hashPassword(NewPassword)
	if err != nil {
		log.Fatal(err)
		fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
	}

	ctx := context.TODO() // i don't now which kind of context we need to use so i keep it todo for now
	if err := db.CreateNewUser(ctx, NewUsername, NewHashedPassword); err != nil {
		if err == sql.ErrNoRows {
			templates.Render(w, templates.TEMPLATE_REGISTER_PAGE, "")
		} else {
			templates.Render(w, "404.html", "")
		}
	} else {
		// https://htmx.org/headers/hx-location/
		http.RedirectHandler("/login-form", http.StatusOK)
		w.Header().Add("HX-Location", "/login")
		log.Printf("User %s registered successfully:", NewUsername)

		return
	}
}
