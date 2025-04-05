package database

import (
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const DATABASE_FILE_PATH = "./internal/database/sqlite-database.db"

var ErrExist = errors.New("already exists")
var ErrNotExist = errors.New("does not exist")
var ErrNotEmpty = errors.New("not empty")
var ErrNotImplemented = errors.New("is not implemented")
var ErrIdenticalThing = errors.New("Thing IDs are the same")

// add statement to create new table
var mainTables = &map[string]string{
	"user":  CREATE_USER_TABLE_STMT,
	"item":  CREATE_ITEM_TABLE_STMT,
	"box":   CREATE_BOX_TABLE_STMT,
	"shelf": CREATE_SHELF_TABLE_STMT,
	"area":  CREATE_AREA_TABLE_STMT,
}

var virtualTables = &map[string]string{
	"item_fts":  CREATE_ITEM_TABLE_STMT_FTS,
	"box_fts":   CREATE_BOX_TABLE_STMT_FTS,
	"shelf_fts": CREATE_SHELF_TABLE_STMT_FTS,
	"area_fts":  CREATE_AREA_TABLE_STMT_FTS,
}

var triggers = &map[string]string{
	"item_fts_trigger_insert":  CREATE_ITEM_INSERT_TRIGGER,
	"item_fts_trigger_update":  CREATE_ITEM_UPDATE_TRIGGER,
	"item_fts_trigger_delete":  CREATE_ITEM_DELETE_TRIGGER,
	"box_fts_trigger_insert":   CREATE_BOX_INSERT_TRIGGER,
	"box_fts_trigger_update":   CREATE_BOX_UPDATE_TRIGGER,
	"box_fts_trigger_delete":   CREATE_BOX_DELETE_TRIGGER,
	"shelf_fts_trigger_insert": CREATE_SHELF_INSERT_TRIGGER,
	"shelf_fts_trigger_update": CREATE_SHELF_UPDATE_TRIGGER,
	"shelf_fts_trigger_delete": CREATE_SHELF_DELETE_TRIGGER,
	"area_fts_trigger_insert":  CREATE_AREA_INSERT_TRIGGER,
	"area_fts_trigger_update":  CREATE_AREA_UPDATE_TRIGGER,
	"area_fts_trigger_delete":  CREATE_AREA_DELETE_TRIGGER,
}

type DB struct {
	Sql       *sql.DB
	fileExist bool
}

// Connect creates the database file if it doesn't exist and opens it.
func (db *DB) Connect() {
	if !env.CurrentConfig().UseMemoryDB() {
		// create the database File and open it
		db.createFile(env.CurrentConfig().DbPath())
	}
	db.open(env.CurrentConfig().DbPath())

	// create the necessary Tables
	db.createTable(*mainTables)
	db.createTable(*virtualTables)
	db.createTable(*triggers)

	db.PrintItemRecords()
	// add dummy data
	if !db.fileExist && env.Development() {
		db.insertDummyData()
	}
}

// createFile creates only db file if it doesn't exist, no tables.
// If error occurs program will shut down with os.Exit(1).
func (db *DB) createFile(dbFile string) {
	_, err := os.Stat(dbFile)
	if err == nil {
		db.fileExist = true
		return
	}

	// Ensure the directory exists
	dir := filepath.Dir(dbFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logg.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	logg.Debugf(`creating "%s"`, dbFile)
	file, err := os.Create(dbFile)
	if err != nil {
		logg.Fatalf("Failed to create database: %v", err)
	}
	defer file.Close()
	logg.Infof(`"%s" created`, dbFile)
}

// open the connection and create tables if they don't exist.
// If error occurs program will shut down with os.Exit(1).
func (db *DB) open(dbFile string) {
	var err error
	db.Sql, err = sql.Open("sqlite", dbFile)
	if err != nil {
		logg.Fatalf("Failed to open database: %v", err)
	}
	logg.Debugf("opened '%s'", dbFile)
	logg.Info("Database Connection established")
}

func (db *DB) createTable(statements map[string]string) {
	for tableName, createStatement := range statements {
		row, err := db.Sql.Exec(createStatement)
		if err != nil {
			logg.Fatalf("Failed to create table \"%s\"\nSQL statement:\n\"%s\"\n%v", tableName, createStatement, err)
		}
		numEffectedRows, err := row.RowsAffected()
		if err != nil {
			logg.Debug("SQL statement: " + createStatement)
			logg.Fatalf("Failed to check the number of effected rows while creating the %s table: %v ", tableName, err)
		}
		if numEffectedRows != 0 {
			// logg.Debugf("Table '%s' created successfully\n", tableName)
		}
	}
}
