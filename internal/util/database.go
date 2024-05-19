package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"github.com/google/uuid"
)

type DBUser2 struct {
	Uuid         uuid.UUID `json:"uuid"`
	PasswordHash string    `json:"passwordhash"`
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

	db.File, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Error opening file '%v': %v", filepath, err)
		return err
	}
	log.Printf("Opened JsonDB: %v\n", filepath)

	err = json.NewDecoder(db.File).Decode(&db.Users)
	if err != nil {
		log.Printf("Error decoding JSON from file '%v': %v", filepath, err)
		return err
	}

	return nil
}


// User retrieves a DBUser2 from the JsonDB by username.
// If the user is found, it logs the user details and returns the user.
// If the user is not found, it returns an empty DBUser2 struct.
func (db *JsonDB) User(username string) (DBUser2, bool) {
	userRecord, exist := db.Users[username]
	if exist {
		log.Printf("%v: %v", username, userRecord)
		return userRecord, true
	}
	return DBUser2{}, false
}

// this function will check if there is existing user withe same name and if not it will
// create new one at the end it will save it
func (db *JsonDB) AddUser(username string, passwordHash string) error {
	if _, exist := db.Users[username]; exist {
      return fmt.Errorf("user %s already exists", username)
    }

	newUser := DBUser2{
		Uuid:         uuid.New(),
		PasswordHash: passwordHash,
	}
	db.Users[username] = newUser

	return db.save()
}


//this function is responsible of saving the new Record inside of the Database (user2.json) 
func (db *JsonDB) save() error {

	_, err := db.File.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("Error seeking to start of file: %v", err)
		return err
	}

	encoder := json.NewEncoder(db.File)

	err = encoder.Encode(db.Users)
	if err != nil {
		log.Printf("Error encoding users to JSON: %v", err)
		return err
	}

	currentPos, err := db.File.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Printf("Error getting current file position: %v", err)
		return err
	}

	err = db.File.Truncate(currentPos)
	if err != nil {
		log.Printf("Error truncating file: %v", err)
		return err
	}

	return nil
}
