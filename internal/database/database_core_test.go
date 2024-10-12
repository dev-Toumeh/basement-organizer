package database

import (
	"basement/main/internal/logg"
	"fmt"
	"os"
	"testing"

	_ "github.com/gofrs/uuid/v5"
	_ "modernc.org/sqlite"
)

const DATABASE_TEST_V1_FILE_PATH = "./internal/database/sqlite-database-test-v1.db"

var dbTest = &DB{}

func TestMain(m *testing.M) {

	setup()
	defer teardown()

	code := m.Run()

	os.Exit(code)
}

func setup() {
	// setup db with real file
	// dbTest.createFile(DATABASE_PROD_V1_FILE_PATH)
	// dbTest.open(DATABASE_PROD_V1_FILE_PATH)

	// setup in-memory db
	dbTest.open(":memory:")

	// dbTest.PrintTables()
}

func teardown() {

	EmptyTestDatabase()
	logg.Info("Testing Database Package was finished, Tables was cleared")
	dbTest.Sql.Close()
}

func EmptyTestDatabase() {
	for tableName := range *mainTables {
		sqlStatement := fmt.Sprintf("DELETE FROM %s;", tableName)
		_, err := dbTest.Sql.Exec(sqlStatement)
		if err != nil {
			logg.Fatalf("Failed to delete from table \"%s\"\n\tquery: \"%s\"\n\t%s", tableName, sqlStatement, err)
			return
		}
	}
	for tableName := range *virtualTables {
		sqlStatement := fmt.Sprintf("DELETE FROM %s;", tableName)
		logg.Debug(sqlStatement)
		_, err := dbTest.Sql.Exec(sqlStatement)
		if err != nil {
			logg.Fatalf("Failed to delete from table %s: %s", tableName, err)
			return
		}
	}
}
