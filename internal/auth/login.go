package auth

import (
	"basement/main/internal/templates"
	"basement/main/internal/util"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

const LOGIN_FAILED_MESSAGE string = "Login failed"
const COOKIE_NAME string = "mycookie"

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type AuthDatabaseHandler interface {
	User(string) (util.DBUser2, bool)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	RegisterHandler(w http.ResponseWriter, r *http.Request)
}

type AuthJsonDB struct {
	*util.JsonDB
}

func (db *AuthJsonDB) LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, COOKIE_NAME)
	authenticated, ok := session.Values["authenticated"].(bool)

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

	if !util.CheckPasswordHash(password, user.PasswordHash) {
		log.Println("pw hash doesnt match")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	session.Values["authenticated"] = true
	session.Save(r, w)
	log.Println("session authenticated", session.Values["authenticated"])

	log.Println("login successful")
	fmt.Fprintf(w, "Welcome %v\n", username)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(templates.ROOT_PAGE_TEMPLATE_FILE, "internal/templates/login.html")
	if err != nil {
		log.Printf("%v or %v: %v\n", templates.ROOT_PAGE_TEMPLATE, "login.html", err)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	templateData := templates.RootPageTemplate{Title: "login"}

	if err := tmpl.ExecuteTemplate(w, templates.ROOT_PAGE_TEMPLATE, templateData); err != nil {
		log.Printf("loginPage: %v\n", err)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
	}
}
