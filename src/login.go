package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	PassHash string `json:"pass_hash"`
}

func login(w http.ResponseWriter, r *http.Request) {
	// file, _ := os.OpenFile("users.json", os.O_RDONLY)
	reader := bufio.NewReader(r.Body)
	msg, _ := reader.ReadString('\n')
	msgs := strings.Split(msg, "&")
	fmt.Println("%v : %v", strings.Split(msgs[0], "="))
	// msg, _ = reader.ReadString('&')
	fmt.Println("2", msgs[1])

	var user User
	var users_file User

	file, _ := os.Open("users.json")

	fmt.Println("Welcome to my login")
	json.NewDecoder(file).Decode(&users_file)

	// decoder := json.NewDecoder(r.Body)
	// decoder.Decode(&user)
	// user2 = user
	// encoder := json.NewEncoder(w)
	// encoder.Encode(&user2)

	// pwhash, _ := HashPassword(user.PassHash)
	// user2.PassHash = pwhash
	fmt.Printf("user %s, %s, %s", user.Id, user.Username, user.PassHash)
	fmt.Println()
	// fmt.Printf("user2 %s, %s, %s", user2.Id, user2.Username, user2.PassHash)
	// fmt.Println()
	fmt.Printf("users_file %s, %s, %s", users_file.Id, users_file.Username, users_file.PassHash)
	fmt.Println()
	fmt.Fprintf(w, "%s, %s, %s", user.Id, user.Username, user.PassHash)

	// json.NewDecoder(r.Body).Decode()
}
