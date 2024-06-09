package auth

// import "basement/main/internal/util"

import (
	"log"
)

func CreateJsonDB() (*AuthJsonDB, error) {
	db := AuthJsonDB{}
	err := db.Connect("./internal/auth/users2.json")
	if err != nil {
		log.Println("createDB() error", err)
		return &AuthJsonDB{}, err
	}

	return &db, nil
}
