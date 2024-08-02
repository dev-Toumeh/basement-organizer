package database

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/gofrs/uuid/v5"
	_ "modernc.org/sqlite"
)

const (
	CREATE_USER_TABLE_STMT = `CREATE TABLE IF NOT EXISTS user ( id TEXT NOT NULL PRIMARY KEY, username TEXT UNIQUE, passwordhash TEXT);`
	CREATE_ITEM_TABLE_STMT = `CREATE TABLE IF NOT EXISTS item ( id TEXT PRIMARY KEY, label TEXT NOT NULL,
                            description TEXT, picture TEXT, quantity INTEGER, weight TEXT, qrcode TEXT);`
	DATABASE_FILE_PATH = "./internal/database/sqlite-database.db"
)

var ErrExist = errors.New("exist")

// add statement to create new table
var statements = &map[string]string{
	"user": CREATE_USER_TABLE_STMT,
	"item": CREATE_ITEM_TABLE_STMT,
}

type DB struct {
	Sql   *sql.DB
	Users []User
	Items []Item
}

// create the Database file if it was not exist and establish the connection with it
func (db *DB) Connect() error {

	// create the sqlite database File it it wasn't exist
	if _, err := os.Stat(DATABASE_FILE_PATH); err != nil {
		log.Println("Creating sqlite-database.db...")
		file, err := os.Create(DATABASE_FILE_PATH)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		defer file.Close()
		log.Println("sqlite-database.db created")
		return err
	}
	// open the connection
	var err error
	if db.Sql, err = sql.Open("sqlite", DATABASE_FILE_PATH); err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return err
	}

	// create the Tables if were not exist
	db.createTable()
	return nil
}

func (db *DB) createTable() {
	for tableName, createStatement := range *statements {
		// First, check if the table exists
		var exists bool
		err := db.Sql.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)", tableName).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check if table exists: %s", err)
			return
		}

		// If the table doesn't exist, create it
		if !exists {
			_, err := db.Sql.Exec(createStatement)
			if err != nil {
				log.Fatalf("Failed to create table: %s", err)
				return
			}
			log.Printf("Table '%s' created successfully", tableName)
		}
	}
	// If the table already exists, the function will exit silently
}
