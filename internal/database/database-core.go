package database

import (
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/gofrs/uuid/v5"
	_ "modernc.org/sqlite"
)

const DATABASE_PROD_V1_FILE_PATH = "./internal/database/sqlite-database-prod-v1.db"

var ErrExist = errors.New("already exists")
var ErrNotExist = errors.New("does not exist")
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
}

type DB struct {
	Sql *sql.DB
}

// Connect creates the database file if it doesn't exist and opens it.
func (db *DB) Connect() {
	if env.Development() {
		logg.Info("using in-memory database")
		db.open(":memory:")
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
		shelfRows, err := db.ShelfListRowsPaginated(1, 3)
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

	if env.Production() {
		logg.Infof(`using "%s" database`, DATABASE_PROD_V1_FILE_PATH)
		db.createFile(DATABASE_PROD_V1_FILE_PATH)
		db.open(DATABASE_PROD_V1_FILE_PATH)
	}
}

// createFile creates only db file if it doesn't exist, no tables.
// If error occurs porgram will shut down with os.Exit(1).
func (db *DB) createFile(dbFile string) {
	_, err := os.Stat(dbFile)
	if err == nil {
		logg.Debugf(`"%s" exists`, dbFile)
		return
	}

	logg.Debugf(`creating "%s"`, dbFile)
	file, err := os.Create("./sqlite-database-test.db")
	if err != nil {
		logg.Fatalf("Failed to create database: %v", err)
	}
	defer file.Close()
	logg.Infof(`"%s" created`, dbFile)
}

// open the connection and create tables if they don't exist.
// If error occurs porgram will shut down with os.Exit(1).
func (db *DB) open(dbFile string) {
	var err error
	db.Sql, err = sql.Open("sqlite", dbFile)
	if err != nil {
		logg.Fatalf("Failed to open database: %v", err)
	}
	logg.Debugf("opened '%s'", dbFile)
	logg.Info("Database Connection established")

	db.createTable(*mainTables)
	db.createTable(*virtualTables)
	db.createTable(*triggers)
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
	err1 := ValidTable(entityType)
	err2 := ValidVirtualTable(entityType)
	// not a vaid table and not valid virtual table
	if err1 != nil && err2 != nil {
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

func (db *DB) deleteFrom(table string, id uuid.UUID) error {
	err := ValidTable(table)
	if err != nil {
		return logg.WrapErr(err)
	}

	stmt := fmt.Sprintf(`DELETE FROM %s WHERE id = ?;`, table)
	result, err := db.Sql.Exec(stmt, id.String())
	if err != nil {
		return logg.Errorf(`can't delete "%s" from "%s" %w`, id, table, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return logg.WrapErr(err)
	}
	if rowsAffected == 0 {
		return logg.NewError(fmt.Sprintf(`id "%s" not found in "%s"`, id, table))
	} else if rowsAffected != 1 {
		return logg.NewError(fmt.Sprintf(`unexpected number of rows affected (%d) while deleting "%s" from "%s"`, rowsAffected, id, table))
	}
	return nil
}

func ValidTable(table string) error {
	_, ok := (*mainTables)[table]
	if !ok {
		return logg.NewError(fmt.Sprintf(`"%s" is not a valid table`, table))
	}
	return nil
}

func ValidVirtualTable(table string) error {
	_, ok := (*virtualTables)[table]
	if !ok {
		return logg.NewError(fmt.Sprintf(`"%s" is not a valid virtual table`, table))
	}
	return nil
}

// MoveTo moves item/box/shelf to a box/shelf/area.
//
// Example move item to a box:
//
//	MoveTo("item", itemID, "box", boxID)
//
// To move things out set
//
//	toTableID = uuid.Nil
func (db *DB) MoveTo(table string, id uuid.UUID, toTable string, toTableID uuid.UUID) error {
	err := ValidTable(table)
	if err != nil {
		return logg.Errorf(`table: "%s" %w`, table, err)
	}
	err = ValidTable(toTable)
	if err != nil {
		return logg.Errorf(`toTable: "%s" %w`, toTable, err)
	}

	errMsg := fmt.Sprintf(`moving "%s" "%s" to "%s" "%s"`, table, id.String(), toTable, toTableID.String())

	exists, err := db.Exists(table, id)
	if err != nil {
		return logg.WrapErr(err)
	}
	if exists == false {
		return logg.Errorf("%s %w", errMsg, ErrNotExist)
	}

	// check if the table where the item is being moved to exists
	if toTableID != uuid.Nil {
		exists, err := db.Exists(toTable, toTableID)
		if err != nil {
			return logg.WrapErr(err)
		}
		if !exists {
			return logg.Errorf("%s %w", errMsg, ErrNotExist)
		}
	}

	// Update the item's shelf_id
	stmt := fmt.Sprintf(`UPDATE %s SET %s_id = ? WHERE id = ?`, table, toTable)
	result, err := db.Sql.Exec(stmt, toTableID.String(), id.String())
	if err != nil {
		return logg.WrapErr(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return logg.WrapErr(err)
	}
	if rows != 1 {
		return logg.NewError(fmt.Sprintf("rows should be != 1 but is %d", rows))
	}
	return nil
}
