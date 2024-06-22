package auth

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/google/uuid"
)

// JsonDB handles a JSON file as a simple storage solution to hold user information
type JsonDB struct {
	File       *os.File
	Users      map[string]DBUser2
	FileReader io.Reader
	Items      map[string]Item
}

// DBUser2 represents user entries in a database
type DBUser2 struct {
	Uuid         uuid.UUID `json:"uuid"`
	PasswordHash string    `json:"passwordhash"`
}

type Item struct {
	Id          uuid.UUID   `json:"id"`
	Label       string      `json:"label"`
	Description string      `json:"description"`
	Picture     string      `json:"picture"`
	Quantity    json.Number `json:"quantity"`
	Weight      string      `json:"weight"`
	QRcode      string      `json:"qrcode"`
}

// AuthDatabase is for authentication handler functions that need database access
// type AuthDatabase interface {
// 	User(string) (DBUser2, bool)
// 	LoginHandler(w http.ResponseWriter, r *http.Request)
// 	RegisterHandler(w http.ResponseWriter, r *http.Request)
// }

// CreateJsonDB an object from a JSON file to be used as simple storage
func CreateJsonDB() (*JsonDB, error) {
	db := JsonDB{}
	err := db.connect("internal/auth/users2.json")
	if err != nil {
		log.Println("createDB() error", err)
		return &JsonDB{}, err
	}

	return &db, nil
}

func (db *JsonDB) connect(filepath string) error { // @TODO: Change filepath string to io.Reader for more flexibility
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

// InitFieldFromFile reads JSON file from disk to populate field.
//
// `field` must be an internal field of the database instance.
//
// Example: InitFieldFromFile("file.json", &db.Items)
func (db *JsonDB) InitFieldFromFile(filepath string, field interface{}) error {
	// reflect
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
	db.File = file
	if err != nil {
		log.Printf("Error opening file '%v': %v", file.Name(), err)
		return err
	}

	_ = db.InitField(file, field)
	return nil
}

// InitField reads data to populate Items field.
//
// `field` must be an internal field of the database instance.
//
// Example: InitFieldFromFile("file.json", &db.Items)
func (db *JsonDB) InitField(data io.Reader, field any) error {
	db.FileReader = data
	err := json.NewDecoder(data).Decode(field)
	if err != nil {
		log.Printf("Error decoding JSON from file '%v': %v", data, err)
		return err
	}
	log.Printf("InitField: %v\n", reflect.TypeOf(field))

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
