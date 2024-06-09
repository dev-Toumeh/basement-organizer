package auth

import (
	"basement/main/internal/templates"
	"fmt"
	"log"
	"net/http"
)

const LOGIN_FAILED_MESSAGE string = "Login failed"

func (db *AuthJsonDB) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db.loginUser(w, r)
	}
	if r.Method == http.MethodGet {
		db.loginPage(w, r)
	}
}

func (db *AuthJsonDB) loginUser(w http.ResponseWriter, r *http.Request) {
	authenticated, ok := Authenticated(r)

	if ok {
		if authenticated {
			log.Printf("LoginHandler - ok: %v authenticated: %v", ok, authenticated)
			fmt.Fprint(w, "already logged in")
			return
		}
	}

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

	user, _ := db.User(username)

	if !checkPasswordHash(password, user.PasswordHash) {
		log.Println("pw hash doesnt match")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	saveSession(w, r)

	log.Println("login successful")

	// https://htmx.org/headers/hx-location/
	w.Header().Add("HX-Location", "/")
	fmt.Fprintf(w, "Welcome %v\n", username)
}

func (db *AuthJsonDB) loginPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := Authenticated(r)
	data := templates.PageTemplate{
		Title:         "login",
		Authenticated: authenticated,
	}
	if err := templates.ApplyPageTemplate(w, "internal/templates/login.html", data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

func Authenticated(r *http.Request) (bool, bool) {
	session, _ := store.Get(r, COOKIE_NAME)
	authenticated, ok := session.Values["authenticated"].(bool)
	log.Println("session authenticated", session.Values["authenticated"])
	return authenticated, ok
}

func saveSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, COOKIE_NAME)
	session.Values["authenticated"] = true
	session.Save(r, w)
}
