package auth

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"internal/util"
)

const (
	REGISTER_FAILED_MESSAGE string = "register failed"
	Username                       = "username"
	Password                       = "password"
)

type registerPage struct {
	Title string
}

type DBUsers []DBUser

// Implement this method
func (db *AuthJsonDB) RegisterHandler(w http.ResponseWriter, r *http.Request) {}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		NewUsername := r.PostFormValue(Username)
		NewPassword := r.PostFormValue(Password)
		//		fmt.Printf("the usernmae is %s \n and the password is %s", username, password)

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
		// check the Database

		// open the json file
		file, err := os.Open("users.json")
		if err != nil {
			log.Println("Error opening users.json:", err)
			fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
			return
		}

		// turn into bytes
		byteValue, err := io.ReadAll(file)
		if err != nil {
			log.Println("Error happened while opening users.json file: ", err)
			fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
			return
		}

		defer file.Close()

		// convert json File to array of DBuser.
		var dbUsers DBUsers
		err = json.Unmarshal(byteValue, &dbUsers)
		if err != nil {
			log.Println("Error reading users.json:", err)
			fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
			return
		}

		//check the if the username is exist
		for _, dbUser := range dbUsers {
			if dbUser.Username == NewUsername {
				log.Printf("the username %s is already token", NewUsername)
				fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
				return
			}
		}

		// hash the password
		NewHashedPassword, err := util.HashPassword(NewPassword)
		if err != nil {
			log.Fatal(err)
			fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
		}
		// create new Record
		newUser := DBUser{
			Id:       uuid.New(),
			Username: NewUsername,
			PassHash: NewHashedPassword,
		}

		dbUsers = append(dbUsers, newUser)

		// Convert the updated users list back to JSON
		updatedJSON, err := json.Marshal(dbUsers)
		if err != nil {
			log.Println("Error marshalling new user data:", err)
			fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
			return
		}

		// Write the updated JSON back to the File
		err = os.WriteFile("users.json", updatedJSON, 0644)
		if err != nil {
			log.Println("Error writing to users.json:", err)
			fmt.Fprintln(w, REGISTER_FAILED_MESSAGE)
			return
		}

		log.Println("User registered successfully:", newUser.Username)
		fmt.Fprintln(w, "User registered successfully")

		return
	}

	tmpl, err := template.ParseFiles("./register.html")
	if err != nil {
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
