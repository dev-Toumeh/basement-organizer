package database

import (
	"basement/main/internal/logg"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/gofrs/uuid/v5"
	_ "modernc.org/sqlite"
)

var dbTest = &DB{}

func TestMain(m *testing.M) {

	setup()
	defer teardown()

	code := m.Run()

	os.Exit(code)
}

func setup() {

	// 1. Create the sqlite database File it it wasn't exist
	if _, err := os.Stat("./sqlite-database-test.db"); err != nil {
		fmt.Print("Creating sqlite-database-test.db... \n")
		file, err := os.Create("./sqlite-database-test.db")
		if err != nil {
			logg.Fatalf("Failed to create database: %v", err)
		}
		defer file.Close()
		logg.Debug("sqlite-database-test.db was created")
	}

	//  2. Open the connection
	var err error
	if dbTest.Sql, err = sql.Open("sqlite", "./sqlite-database-test.db"); err != nil {
		logg.Fatalf("Failed to open database: %v", err)
	}

	// 3. Run our DDL statements to create the required tables if they do not exist
	createTestTables(*statementsMainTables)
	createTestTables(*statementVertualTabels)

	//dbTest.PrintTables()
}

func teardown() {

	EmptyTestDatabase()
	logg.Info("Testing Database Package was finished, Tables was cleared")
	dbTest.Sql.Close()
}

func EmptyTestDatabase() {
	statments := []string{"user", "item", "box", "shelf", "area", "item_fts", "box_fts"}
	for _, tableName := range statments {
		sqlStatement := fmt.Sprintf("DELETE FROM %s;", tableName)
		_, err := dbTest.Sql.Exec(sqlStatement)
		if err != nil {
			logg.Fatalf("Failed to delete from table %s: %s", tableName, err)
			return
		}
	}
}

func createTestTables(statements map[string]string) {
	for tableName, createStatement := range statements {
		var exists bool
		err := dbTest.Sql.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)", tableName).Scan(&exists)
		if err != nil {
			logg.Fatalf("Failed to check if table exists: %s", err)
		}
		if !exists {
			_, err := dbTest.Sql.Exec(createStatement)
			if err != nil {
				logg.Fatalf("Failed to create table: %s", err)
			}
			logg.Debugf("Table '%s' created successfully", tableName)
		}
	}
}
