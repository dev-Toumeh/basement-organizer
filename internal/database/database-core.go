package database

import (
	"basement/main/internal/logg"
	"database/sql"
	"errors"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

const DATABASE_FILE_PATH = "./internal/database/sqlite-database.db"

var ErrExist = errors.New("the Record is already exist")

// add statement to create new table
var statements = &map[string]string{
	"user":                   CREATE_USER_TABLE_STMT,
	"item":                   CREATE_ITEM_TABLE_STMT,
	"box":                    CREATE_BOX_TABLE_STMT,
	"item_fts":               CREATE_ITEM_TABLE_STMT_FTS,
	"box_fts":                CREATE_BOX_TABLE_STMT_FTS,
	"Item_fts_triger_insert": CREATE_ITEM_AI_TRIGGER,
	"item_fts_triger_update": CREATE_ITEM_AU_TRIGGER,
	"item_fts_triger_delete": CREATE_ITEM_AD_TRIGGER,
	"box_fts_triger_insert":  CREATE_BOX_AI_TRIGGER,
	"box_fts_triger_update":  CREATE_BOX_AU_TRIGGER,
	"box_fts_triger_delete":  CREATE_BOX_AD_TRIGGER,
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
	db.RepopulateItemFTS() // only patch take this function out after Alex adopt the changes
}

func (db *DB) createTable() {
	for tableName, createStatement := range *statements {
		row, err := db.Sql.Exec(createStatement)
		if err != nil {
			logg.Fatalf("Failed to create table %s: %v", tableName, err)
		}
		numEffectedRows, err := row.RowsAffected()
		if err != nil {
			logg.Fatalf("Failed to check the number of effected rows while creating the %s table: %v ", tableName, err)
		}
		if numEffectedRows != 0 {
			log.Printf("Table '%s' created successfully\n", tableName)
		}
	}
}

// ErrorExist returns a predefined error indicating that the requested SQL data insertion failed
func (db *DB) ErrorExist() error {
	return ErrExist
}
