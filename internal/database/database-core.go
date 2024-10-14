package database

import (
	"basement/main/internal/logg"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/gofrs/uuid/v5"
	_ "modernc.org/sqlite"
)

const DATABASE_FILE_PATH = "./internal/database/sqlite-database2.db"

var ErrExist = errors.New("already exists")
var ErrNotExist = errors.New("does not exist")
var ErrNotImplemented = errors.New("is not implemented")

// add statement to create new table
var statementsMainTables = &map[string]string{
	"user":  CREATE_USER_TABLE_STMT,
	"item":  CREATE_ITEM_TABLE_STMT,
	"box":   CREATE_BOX_TABLE_STMT,
	"shelf": CREATE_SHELF_TABLE_STMT,
	"area":  CREATE_AREA_TABLE_STMT,
}

var statementVirtualTabels = &map[string]string{
	"item_fts":                 CREATE_ITEM_TABLE_STMT_FTS,
	"box_fts":                  CREATE_BOX_TABLE_STMT_FTS,
	"shelf_fts":                CREATE_SHELF_TABLE_STMT_FTS,
	"Item_fts_trigger_insert":  CREATE_ITEM_INSERT_TRIGGER,
	"item_fts_trigger_update":  CREATE_ITEM_UPDATE_TRIGGER,
	"item_fts_trigger_delete":  CREATE_ITEM_DELETE_TRIGGER,
	"box_fts_trigger_insert":   CREATE_BOX_INSERT_TRIGGER,
	"box_fts_trigger_update":   CREATE_BOX_UPDATE_TRIGGER,
	"box_fts_trigger_delete":   CREATE_BOX_DELETE_TRIGGER,
	"shelf_fts_trigger_insert": CREATE_SHELF_INSERT_TRIGGER,
	"shelf_fts_trigger_update": CREATE_SHELF_UPDATE_TRIGGER,
	"shelf_fts_trigger_delete": CREATE_SHELF_DELETE_TRIGGER,
}

type DB struct {
	Sql *sql.DB
}

// create the Database file if it was not exist and establish the connection with it
func (db *DB) Connect() {

	// create the sqlite database File it it wasn't exist
	if _, err := os.Stat(DATABASE_FILE_PATH); err != nil {
		logg.Info("Creating sqlite-database.db...")
		file, internErr := os.Create(DATABASE_FILE_PATH)
		defer file.Close()
		if internErr != nil {
			logg.Fatal(internErr)
		}
		logg.Info("sqlite-database.db created")
	}
	// open the connection
	var err error
	db.Sql, err = sql.Open("sqlite", DATABASE_FILE_PATH)
	if err != nil {
		logg.Fatalf("Failed to open database: %v", err)
	}
	logg.Info("Database Connection established")

	// create the Tables if were not exist
	db.createTable(*statementsMainTables)
	db.createTable(*statementVirtualTabels)
}

func (db *DB) createTable(statements map[string]string) {
	for tableName, createStatement := range statements {
		row, err := db.Sql.Exec(createStatement)
		if err != nil {
			logg.Fatalf("Failed to create table %s: %v", tableName, err)
		}
		numEffectedRows, err := row.RowsAffected()
		if err != nil {
			logg.Fatalf("Failed to check the number of effected rows while creating the %s table: %v ", tableName, err)
		}
		if numEffectedRows != 0 {
			logg.Debugf("Table '%s' created successfully\n", tableName)
		}
	}
}

// ErrorExist returns a predefined error indicating that the requested SQL data insertion failed
func (db *DB) ErrorExist() error {
	return logg.WrapErrWithSkip(ErrExist, 2)
}

// Exists checks existence of entity (item, box, shelf, area).
// Returns error if other internal errors happen.
func (db *DB) Exists(entityType string, id uuid.UUID) (bool, error) {
	validEntities := []string{"item", "box", "shelf", "area"}
	if !slices.Contains(validEntities, entityType) {
		return false, logg.NewError(fmt.Sprintf("no entity with type: \"%s\"", entityType))
	}

	var itemExists int
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = ?)`, entityType)
	logg.Debug(query)
	err := db.Sql.QueryRow(query, id.String()).Scan(&itemExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, logg.WrapErr(err)
		}
	}
	if itemExists != 0 {
		return true, nil
	} else {
		return false, nil
	}
}
