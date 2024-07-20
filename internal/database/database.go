package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/gofrs/uuid/v5"
)

// JsonDB handles a JSON file as a simple storage solution to hold user information
type JsonDB struct {
	UserFile   *os.File
	ItemFile   *os.File
	Users      map[string]DBUser2
	FileReader io.Reader
	Items      map[uuid.UUID]Item
}

// DBUser2 represents user entries in a database
type DBUser2 struct {
	Uuid         uuid.UUID `json:"uuid"`
	PasswordHash string    `json:"passwordhash"`
}

const (
	ITEMS_FILE_PATH            string = "internal/database/items.json"
	USERS_FILE_PATH            string = "internal/database/users2.json"
	ITEM_ERROR_GENERAL_MESSAGE string = "<div>we were not able to update the item, please come back later</div>"
)

type Item struct {
	Id          uuid.UUID   `json:"id"`
	Label       string      `json:"label"       validate:"required,lte=128"`
	Description string      `json:"description" validate:"omitempty,lte=256"`
	Picture     string      `json:"picture"     validate:"omitempty,base64"`
	Quantity    json.Number `json:"quantity"    validate:"omitempty,numeric,gte=1"`
	Weight      string      `json:"weight"      validate:"omitempty,numeric"`
	QRcode      string      `json:"qrcode"      validate:"omitempty,alphanumunicode"`
}

// CreateJsonDB an object from a JSON file to be used as simple storage
func CreateJsonDB() (*JsonDB, error) {
	db := JsonDB{}
	err1 := db.InitFieldFromItemFile(&db)
	err2 := db.InitFieldFromUserFile(&db)

	if err1 != nil || err2 != nil {
		log.Fatal("createDB() error")
	}

	return &db, nil
}

func (db *JsonDB) InitFieldFromUserFile(field *JsonDB) error {
	file, err := os.OpenFile(USERS_FILE_PATH, os.O_RDWR|os.O_CREATE, 0666)
	db.UserFile = file
	if err != nil {
		log.Printf("Error opening file '%v': %v", file.Name(), err)
		return err
	}

	_ = db.InitField(file, &field.Users)
	return nil
}

// InitFieldFromFile reads JSON file from disk to populate field.
// `field` must be an internal field of the database instance.
// Example: InitFieldFromFile("file.json", &db.Items)
func (db *JsonDB) InitFieldFromItemFile(field *JsonDB) error {
	file, err := os.OpenFile(ITEMS_FILE_PATH, os.O_RDWR|os.O_CREATE, 0666)
	db.ItemFile = file
	if err != nil {
		log.Printf("Error opening file '%v': %v", file.Name(), err)
		return err
	}

	_ = db.InitField(file, &field.Items)
	return nil
}

// InitField reads data to populate Items field.
// `field` must be an internal field of the database instance.
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

// this function will check if there is existing user withe same name and if not it will
// create new one at the end it will save it
func AddUser(username string, passwordHash string, db *JsonDB) error {
	if dbUser, exist := db.User(username); exist {
		return fmt.Errorf("user %s already exists", username)
	} else {
		dbUser.Uuid, _ = uuid.NewV4()
		dbUser.PasswordHash = passwordHash
		db.Users[username] = dbUser
	}

	return db.saveUser()
}

// this function is responsible of saving the new Record inside of the Database (user2.json)
func (db *JsonDB) saveUser() error {

	_, err := db.UserFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("Error seeking to start of file: %v", err)
		return err
	}

	encoder := json.NewEncoder(db.UserFile)

	err = encoder.Encode(db.Users)
	if err != nil {
		log.Printf("Error encoding users to JSON: %v", err)
		return err
	}

	currentPos, err := db.UserFile.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Printf("Error getting current file position: %v", err)
		return err
	}

	err = db.UserFile.Truncate(currentPos)
	if err != nil {
		log.Printf("Error truncating file: %v", err)
		return err
	}

	return nil
}

// ItemByLabel check if the Item Label exist
// If it was, the function will return the item with true.
// If not it will return empty item with false.
func (db *JsonDB) ItemByLabel(label string) (Item, bool) {
	for _, item := range db.Items {
		if label == item.Label {
			return item, true
		}
	}
	return Item{}, false
}

// ItemById check if the Item id exist
// If it was, the function will return the item with true.
// If not it will return empty item with false.
func (db *JsonDB) ItemById(id uuid.UUID) (Item, bool) {
	itemRecord, exist := db.Items[id]
	if exist {
		return itemRecord, true
	}

	return Item{}, false
}

// this function will add new Record to the database
func (db *JsonDB) AddItem(newItem Item) ([]string, error) {
	var responseMessage []string
	if _, exist := db.ItemByLabel(newItem.Label); exist {
		err := errors.New("the Label exist in the Database")
		responseMessage = append(responseMessage, "<div>the Label exist in the Database please choice another label</div>")
		return responseMessage, err
	} else {
		db.Items[newItem.Id] = newItem
		if err := db.saveItem(); err != nil {
			responseMessage = append(responseMessage, ITEM_ERROR_GENERAL_MESSAGE)
			return responseMessage, err
		}
		responseMessage = append(responseMessage, "<div>The Item has been added successfully</div>")
		return responseMessage, nil
	}
}

// this function will update the Item Record
func (db *JsonDB) UpdateItem(updatedItem Item) ([]string, error) {
	var responseMessage []string
	if _, exist := db.ItemById(updatedItem.Id); !exist {
		err := errors.New("the Id doesn't not exist in the Database")
		responseMessage = append(responseMessage, ITEM_ERROR_GENERAL_MESSAGE)
		return responseMessage, err
	} else {
		db.Items[updatedItem.Id] = updatedItem
		if err := db.saveItem(); err != nil {
			responseMessage = append(responseMessage, ITEM_ERROR_GENERAL_MESSAGE)
			return responseMessage, err
		}
		responseMessage = append(responseMessage, "<div>The Item has been updated successfully</div>")
		return responseMessage, nil
	}
}

func (db *JsonDB) saveItem() error {
	_, err := db.ItemFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("Error seeking to start of file: %v", err)
		return err
	}
	encoder := json.NewEncoder(db.ItemFile)

	err = encoder.Encode(db.Items)
	if err != nil {
		log.Printf("Error encoding users to JSON: %v", err)
		return err
	}

	currentPos, err := db.ItemFile.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Printf("Error getting current file position: %v", err)
		return err
	}

	err = db.ItemFile.Truncate(currentPos)
	if err != nil {
		log.Printf("Error truncating file: %v", err)
		return err
	}

	return nil
}
