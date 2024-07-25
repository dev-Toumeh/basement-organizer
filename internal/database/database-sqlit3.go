package database

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/gofrs/uuid/v5"
	_ "modernc.org/sqlite"
)

const (
	CREATE_USER_TABLE_STMT = `CREATE TABLE IF NOT EXISTS user ( "id" BLOB NOT NULL PRIMARY KEY,
                              "username" TEXT UNIQUE, "passwordHash" TEXT);`
	CREATE_ITEM_TABLE_STMT = `CREATE TABLE IF NOT EXISTS Item ( Id TEXT PRIMARY KEY, Label TEXT NOT NULL,
                            Description TEXT, Picture TEXT, Quantity INTEGER, Weight TEXT, QRcode TEXT);`
	DATABASE_FILE_PATH = "./internal/database/sqlite-database.db"
)

// add statement to create new table
var statements = []string{CREATE_ITEM_TABLE_STMT, CREATE_USER_TABLE_STMT}

type DB struct {
	Sql   *sql.DB
	Users []User
	Items []Item
}

type User struct {
	Id           []byte
	Username     string
	PasswordHash string
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
		file.Close()
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
	for _, statement := range statements {
		db.createTable(statement)
	}
	return nil
}

func (db *DB) createTable(statment string) {
	statement, err := db.Sql.Prepare(statment)
	if err != nil {
		log.Fatalf("Failed to prepare the create table SQL statement: %s", err)
		return
	}
	defer statement.Close()

	result, err := statement.Exec()
	if err != nil {
		log.Fatalf("Failed to execute the create table statement: %s", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %s", err)
	} else {
		log.Printf("User table created, %d rows affected", rowsAffected)
	}
}

// close the Connection
func (db *DB) Close() {
	defer db.Sql.Close()
}

func (db *DB) TestConnection() {
	ctx := context.Background()

	var testOutput string
	_, err := db.Sql.ExecContext(ctx, "SELECT 'Test Connection'")
	if err != nil {
		log.Printf("Error during test query: %v", err)
	} else {
		log.Printf("Test query successful, output: %s", testOutput)
	}
}

func (db *DB) CheckDataPresence(tableName string) {
	ctx := context.Background()

	var count int
	err := db.Sql.QueryRowContext(ctx, "SELECT COUNT(*) FROM", tableName).Scan(&count)
	if err != nil {
		log.Printf("Error during data presence check: %v", err)
	} else if count > 0 {
		log.Printf("Data check successful, data is present")
	} else {
		log.Printf("Data check complete, no data found")
	}
}
