package auth

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"internal/util"
	"log"
	"net/http"
	"os"
)

type DBUser struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	PassHash string    `json:"pass_hash"`
}

const LOGIN_FAILED_MESSAGE string = "Login failed"

func Login(w http.ResponseWriter, r *http.Request) {
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

	file, err := os.Open("users.json")
	if err != nil {
		log.Println("Error happened while opening users.json file: ", err)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	var usersFile DBUser
	json.NewDecoder(file).Decode(&usersFile)

	if !util.CheckPasswordHash(password, usersFile.PassHash) {
		log.Println("pw hash doesnt match")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	log.Println("login successful")
	fmt.Fprintf(w, "Welcome %v\n", username)
}
