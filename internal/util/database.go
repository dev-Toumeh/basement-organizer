package util

import (
	"encoding/json"
	"log"
	"os"
)

type DBUser2 struct {
	Uuid         string `json:"uuid"`
	PasswordHash string `json:"passwordhash"`
}

type DBWithFileConnector interface {
	Connect(string) error
}

type JsonDB struct {
	File  *os.File
	Users map[string]DBUser2
}

func (db *JsonDB) Connect(filepath string) error {
	var err error

	db.File, err = os.Open(filepath)
	if err != nil {
		log.Printf("Error happened while opening '%v'", filepath)
		log.Println(err)
		return err
	}
	log.Printf("Opened JsonDB: %v\n", filepath)

	err = json.NewDecoder(db.File).Decode(&db.Users)
	if err != nil {
		log.Printf("Error happened while opening %v file: %v\n", filepath, err)
		return err
	}
	return nil
}

func (db *JsonDB) User(username string) DBUser2 {
	u, ok := db.Users[username]
	if !ok {
		return DBUser2{}
	}
	log.Printf("%v: %v", username, u)
	return u
}
