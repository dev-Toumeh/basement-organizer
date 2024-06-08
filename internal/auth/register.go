package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"basement/main/internal/util"
)

const (
	REGISTER_FAILED_MESSAGE string = "register failed"
	REGISTER_TEMPLATE_PATH  string = "internal/templates/register.html"
	USERNAME                string = "username"
	PASSWORD                string = "password"
)

type registerPage struct {
	Title string
}

func (db *AuthJsonDB) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
			log.Default().Fatalf("the user %s is already exist", NewUsername)
			fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
			return
		}

		// hash the password
		NewHashedPassword, err := util.HashPassword(NewPassword)
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

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
	}

	filePath := filepath.Join(pwd, REGISTER_TEMPLATE_PATH)
	tmpl, err := template.ParseFiles(filePath)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	data := registerPage{
		Title: "Register",
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
