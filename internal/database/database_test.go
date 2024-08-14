package database

import (
	"basement/main/internal/logg"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/gofrs/uuid/v5"
	_ "modernc.org/sqlite"
)

var dbTest = &DB{}

func TestMain(m *testing.M) {

	err := setup()
	if err != nil {
		logg.Fatalf("Can't create Test DB, shutting server down")
	}

	code := m.Run()

	teardown()
	os.Exit(code)
}

func setup() error {

	// 1. Create the sqlite database File it it wasn't exist
	if _, err := os.Stat("./sqlite-database-test.db"); err != nil {
		logg.Debug("Creating sqlite-database-test.db...")
		file, err := os.Create("./sqlite-database-test.db")
		if err != nil {
			logg.Fatalf("Failed to create database: %v", err)
			return err
		}
		defer file.Close()
		logg.Debug("sqlite-database-test.db was created")
		return err
	}

	//  2. Open the connection
	var err error
	if dbTest.Sql, err = sql.Open("sqlite", "./sqlite-database-test.db"); err != nil {
		logg.Fatalf("Failed to open database: %v", err)
		return err
	}

	// 3. Run our DDL statements to create the required tables if they do not exist
	for tableName, createStatement := range *statements {
		var exists bool
		err := dbTest.Sql.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)", tableName).Scan(&exists)
		if err != nil {
			logg.Fatalf("Failed to check if table exists: %s", err)
			return err
		}

		if !exists {
			_, err := dbTest.Sql.Exec(createStatement)
			if err != nil {
				logg.Fatalf("Failed to create table: %s", err)
				return err
			}
			logg.Debugf("Table '%s' created successfully", tableName)
		}
	}

	return nil
}

func teardown() {
	for tableName := range *statements {
		sqlStatement := fmt.Sprintf("DELETE FROM %s;", tableName)
		_, err := dbTest.Sql.Exec(sqlStatement)
		if err != nil {
			logg.Fatalf("Failed to delete from table %s: %s", tableName, err)
			return
		}
	}
	logg.Info("Testing Database Package was finished, Tables was cleared")
	dbTest.Sql.Close()
}

func TestUser(t *testing.T) {
	ctx := context.Background()
	_, err := dbTest.Sql.ExecContext(ctx, "INSERT INTO user (id, username, passwordhash) VALUES (?, ?, ?)", "123e4567-e89b-12d3-a456-426614174000", "testuser", "hash")
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Test the User function for existing user
	user, err := dbTest.User(ctx, "testuser")
	if err != nil {
		t.Errorf("Error fetching user: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	// test for non-existing user
	_, err = dbTest.User(context.Background(), "nonexistent")
	if err == nil {
		t.Errorf("Expected an error for non-existing user, got none")
	}
}
