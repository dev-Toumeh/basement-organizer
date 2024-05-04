package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type DBUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	PassHash string `json:"pass_hash"`
}

func login(w http.ResponseWriter, r *http.Request) {
	const LOGIN_FAILED_MESSAGE string = "Login failed"
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

	var users_file DBUser
	json.NewDecoder(file).Decode(&users_file)

	if !CheckPasswordHash(password, users_file.PassHash) {
		log.Println("pw hash doesnt match")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	log.Println("login successful")
	fmt.Fprintf(w, "Welcome %v\n", username)
}
