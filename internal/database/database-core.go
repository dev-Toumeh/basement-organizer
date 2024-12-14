package database

import (
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"database/sql"
	"errors"
	"os"

	"github.com/gofrs/uuid/v5"
	_ "modernc.org/sqlite"
)

const DATABASE_FILE_PATH = "./internal/database/sqlite-database.db"

var ErrExist = errors.New("already exists")
var ErrNotExist = errors.New("does not exist")
var ErrNotEmpty = errors.New("already exists")
var ErrNotImplemented = errors.New("is not implemented")

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
	if !env.Config().UseMemoryDB() {
		// create the database File and open it
		db.createFile(env.Config().DBPath())
	}
	db.open(env.Config().DBPath())

	// create the necessary Tables
	db.createTable(*mainTables)
	db.createTable(*virtualTables)
	db.createTable(*triggers)

	// add dummy data
	if !db.fileExist && env.Development() {
		db.insirtDummyData()
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

// ErrorExist returns a predefined error indicating that the requested SQL data insertion failed
func (db *DB) ErrorExist() error {
	return logg.WrapErrWithSkip(ErrExist, 2)
}

// ErrorNotEmpty returns a predefined error indicating that the requested unit is not empty
func (db *DB) ErrorNotEmpty() error {
	return ErrNotEmpty
	// return logg.WrapErrWithSkip(ErrNotEmpty, 2)  @TODO (Alex) fix it so we can use it without colors
}

func (db *DB) insirtDummyData() {
	db.InsertSampleItems()
	db.InsertSampleBoxes()
	db.InsertSampleShelves()
	itemIDs, err := db.ItemIDs()
	if err != nil {
		logg.WrapErr(err)
	}
	boxIDs, err := db.BoxIDs()
	if err != nil {
		logg.WrapErr(err)
	}
	shelfRows, _, err := db.ShelfListRowsPaginated(1, 3)
	if err != nil {
		logg.WrapErr(err)
	}
	db.MoveItemToBox(uuid.FromStringOrNil(itemIDs[0]), uuid.FromStringOrNil(boxIDs[0]))
	db.MoveItemToBox(uuid.FromStringOrNil(itemIDs[1]), uuid.FromStringOrNil(boxIDs[0]))
	db.MoveItemToBox(uuid.FromStringOrNil(itemIDs[2]), uuid.FromStringOrNil(boxIDs[0]))
	db.MoveItemToBox(uuid.FromStringOrNil(itemIDs[3]), uuid.FromStringOrNil(boxIDs[1]))
	db.MoveItemToBox(uuid.FromStringOrNil(itemIDs[4]), uuid.FromStringOrNil(boxIDs[1]))
	db.MoveItemToShelf(uuid.FromStringOrNil(itemIDs[5]), shelfRows[0].ID)
	db.MoveItemToShelf(uuid.FromStringOrNil(itemIDs[6]), shelfRows[0].ID)
	db.MoveBoxToShelf(uuid.FromStringOrNil(boxIDs[0]), shelfRows[1].ID)
	db.MoveBoxToShelf(uuid.FromStringOrNil(boxIDs[2]), shelfRows[1].ID)
	db.MoveBoxToBox(uuid.FromStringOrNil(boxIDs[3]), uuid.FromStringOrNil(boxIDs[5]))
	db.MoveBoxToBox(uuid.FromStringOrNil(boxIDs[3]), uuid.FromStringOrNil(boxIDs[5]))
	db.MoveBoxToBox(uuid.FromStringOrNil(boxIDs[5]), uuid.FromStringOrNil(boxIDs[6]))
}
