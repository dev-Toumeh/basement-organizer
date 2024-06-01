package auth

import (
	"basement/main/internal/templates"
	"basement/main/internal/util"
	"fmt"
	"log"
	"net/http"
	"text/template"
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

	user, _ := db.User(username)

	if !util.CheckPasswordHash(password, user.PasswordHash) {
		log.Println("pw hash doesnt match")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	log.Println("login successful")
	fmt.Fprintf(w, "Welcome %v\n", username)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/templates/index.html", "internal/templates/login.html")
	if err != nil {
		log.Printf("loginPage: %v\n", err)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	templateData := templates.IndexTemplate{Title: "login"}

	if err := tmpl.ExecuteTemplate(w, "index", templateData); err != nil {
		log.Printf("loginPage: %v\n", err)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
	}
}
