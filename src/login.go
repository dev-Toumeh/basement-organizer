package main

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("Welcome to my login")

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" {
		fmt.Println("Missing username")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}
	if password == "" {
		fmt.Println("Missing password")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	file, err := os.Open("users.json")
	if err != nil {
		fmt.Printf("Error happened while opening users.json file: %v", err)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	var users_file DBUser
	json.NewDecoder(file).Decode(&users_file)

	if !CheckPasswordHash(password, users_file.PassHash) {
		fmt.Println("pw hash doesnt match")
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		return
	}

	fmt.Println("login successful")
	fmt.Fprintf(w, "Welcome %v\n", username)
}
