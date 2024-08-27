package database

import (
	"basement/main/internal/logg"
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/gofrs/uuid/v5"
	_ "modernc.org/sqlite"
)

const (
	CREATE_USER_TABLE_STMT = `CREATE TABLE IF NOT EXISTS user (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT UNIQUE,
    passwordhash TEXT);`

	CREATE_ITEM_TABLE_STMT = `CREATE TABLE IF NOT EXISTS item (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL,
    description TEXT,
    picture TEXT,
    quantity INTEGER,
    weight TEXT,
    qrcode TEXT,
    box_id TEXT REFERENCES box(id));`

	CREATE_BOX_TABLE_STMT = `CREATE TABLE IF NOT EXISTS box (
    id TEXT PRIMARY KEY,
    label TEXT NOT NULL, 
    description TEXT,
    picture TEXT,
    qrcode TEXT,
    outerbox_id TEXT REFERENCES box(id));`

	DATABASE_FILE_PATH = "./internal/database/sqlite-database.db"
)

var ErrExist = errors.New("the Record is already exist")

// add statement to create new table
var statements = &map[string]string{
	"user": CREATE_USER_TABLE_STMT,
	"item": CREATE_ITEM_TABLE_STMT,
	"box":  CREATE_BOX_TABLE_STMT,
}

type DB struct {
	Sql *sql.DB
}

// create the Database file if it was not exist and establish the connection with it
func (db *DB) Connect() {

	// create the sqlite database File it it wasn't exist
	if _, err := os.Stat(DATABASE_FILE_PATH); err != nil {
		log.Println("Creating sqlite-database.db...")
		file, internErr := os.Create(DATABASE_FILE_PATH)
		defer file.Close()
		if internErr != nil {
			logg.Fatal(internErr)
		}
		log.Println("sqlite-database.db created")
	}
	// open the connection
	var err error
	db.Sql, err = sql.Open("sqlite", DATABASE_FILE_PATH)
	if err != nil {
		logg.Fatalf("Failed to open database: %v", err)
	}
	log.Printf("Database Connection established")

	// create the Tables if were not exist
	db.createTable()
}

func (db *DB) createTable() {
	for tableName, createStatement := range *statements {
		// First, check if the table exists
		var exists bool
		err := db.Sql.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)", tableName).Scan(&exists)
		if err != nil {
			logg.Fatalf("Failed to check if table exists: %s", err)
		}

		// If the table doesn't exist, create it
		if !exists {
			_, err := db.Sql.Exec(createStatement)
			if err != nil {
				logg.Fatalf("Failed to create table: %s", err)
			}
			log.Printf("Table '%s' created successfully\n", tableName)
		}
	}
}

// ErrorExist returns a predefined error indicating that the requested SQL data insertion failed
// due to a duplicate record already existing in the database.
func (db *DB) ErrorExist() error {
	return ErrExist
}
