package auth

import (
	"fmt"
	"log"
	"net/http"

	"basement/main/internal/templates"
)

const (
	LOGIN_FAILED_MESSAGE string = "Login failed"
)

func (db *JsonDB) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db.loginUser(w, r)
	}
	if r.Method == http.MethodGet {
		db.loginPage(w, r)
	}
}

func (db *JsonDB) loginUser(w http.ResponseWriter, r *http.Request) {
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
  http.RedirectHandler("/personal-page", http.StatusOK)

	// https://htmx.org/headers/hx-location/
	w.Header().Add("HX-Location", "/personal-page")
	fmt.Fprintf(w, "Welcome %v\n", username)
}

func (db *JsonDB) loginPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := Authenticated(r)
	data := templates.PageTemplate{
		Title:         "login",
		Authenticated: authenticated,
	}
	if err := templates.ApplyPageTemplate(w, templates.LOGIN_TEMPLATE_FILE_WITH_PATH, data); err != nil {
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
	session.Values["username"] = r.FormValue("username")
	session.Save(r, w)
}

func Username(r *http.Request) string {

	session, _ := store.Get(r, COOKIE_NAME)
	usernmae, ok := session.Values["username"].(string)
	if !ok {
		log.Println("Username not found in session or not a string")
		return ""
	}
	return usernmae
}
