package auth

import (
	"fmt"
	"log"
	"net/http"

	"basement/main/internal/database"
	"basement/main/internal/templates"
)

const (
	LOGIN_FAILED_MESSAGE string = "Login failed"
)

func LoginHandler(db *database.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			loginUser(w, r, db)
		}
		if r.Method == http.MethodGet {
			loginPage(w, r, db)
		}
	}
}
func loginUser(w http.ResponseWriter, r *http.Request, db *database.JsonDB) {
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

func loginPage(w http.ResponseWriter, r *http.Request, db *database.JsonDB) {
	authenticated, _ := Authenticated(r)
	data := templates.PageTemplate{
		Title:         "login",
		Authenticated: authenticated,
	}
	if err := templates.ApplyPageTemplate(w, templates.LOGIN_TEMPLATE_FILE_WITH_PATH, data); err != nil {
		// t := templates.CreateTemplates()
		// t.ExecuteTemplate("")
		// if err := templates.ApplyPageTemplate(w, "internal/templates/login.html", data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

func Authenticated(r *http.Request) (bool, bool) {
	session, _ := store.Get(r, COOKIE_NAME)
	authenticated, ok := session.Values["authenticated"].(bool)
	// log.Println("session authenticated", session.Values["authenticated"])
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
