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
	fmt.Println("Welcome to my login")

	// fmt.Println("username: ", r.FormValue("username"))
	// fmt.Println("pass_hash: ", r.FormValue("pass_hash"))

	reader := bufio.NewReader(r.Body)
	msg, _ := reader.ReadString('\n')
	msgs := strings.Split(msg, "&")

	var user User

	for _, msg_pair := range msgs {
		msg_split := strings.Split(msg_pair, "=")
		if len(msg_split) != 2 {
			fmt.Println("malformed key value pair", msg_split)
			continue
		}
		fmt.Println(msg_split)
		k := msg_split[0]
		v := msg_split[1]

		// fmt.Printf("%v: %v:%v", i, k, v)
		// fmt.Println()

		switch k {
		case "username":
			user.Username = v
		case "pass_hash":
			user.PassHash = v
		}
	}

	if user.Username == "" || user.PassHash == "" {
		fmt.Println("user missing information")
		fmt.Printf("user.username: %v\n, user.PassHash: %v", user.Username, user.PassHash)
		return
	}

	var users_file User

	file, _ := os.Open("users.json")

	json.NewDecoder(file).Decode(&users_file)

	// decoder := json.NewDecoder(r.Body)
	// decoder.Decode(&user)
	// user2 = user
	// encoder := json.NewEncoder(w)
	// encoder.Encode(&user2)

	fmt.Println()
	// pwhash, _ := HashPassword(user.PassHash)
	if !CheckPasswordHash(user.PassHash, users_file.PassHash) {
		fmt.Println("pw hash doesnt match")
		return
	} else {
		fmt.Println("login successful")
	}
	// if pwhash != users_file.PassHash {
	// 	fmt.Println("pw hash doesnt match")
	//        fmt.Printf("user pwhash: %v\nfile pwhash: %v", pwhash, users_file.PassHash)
	//        return
	// }

	// user2.PassHash = pwhash
	// fmt.Printf("user %s, %s, %s", user.Id, user.Username, user.PassHash)
	fmt.Println()
	// fmt.Printf("user2 %s, %s, %s", user2.Id, user2.Username, user2.PassHash)
	// fmt.Println()
	fmt.Printf("users_file: Id:%s\n, Username:%s\n, PassHash:%s\n", users_file.Id, users_file.Username, users_file.PassHash)
	fmt.Println()
	// fmt.Fprintf(w, "%s, %s, %s", user.Id, user.Username, user.PassHash)

	// json.NewDecoder(r.Body).Decode()
}
