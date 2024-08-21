package database

import (
	"basement/main/internal/auth"
	"basement/main/internal/logg"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid/v5"
)

// CreateNewUser inserts a new user into the database with the given username and passwordHash
// The function accepts also a variable of type userExist, which can be obtained by executing the UserExist function.
func (db *DB) CreateNewUser(ctx context.Context, username string, passwordhash string) error {
	err := db.insertNewUser(ctx, username, passwordhash)
	if err != nil {
		return fmt.Errorf("was not able to insert the user %w", err)
	}
	return nil
}

// check if the username is available
// if the user exist it will return user struct with nil
// if not it will return empty user struct with err
func (db *DB) User(ctx context.Context, username string) (auth.User, error) {
	var user auth.User
	var userId string

	if ctx == nil {
		ctx = context.Background()
	}

	query := "SELECT id, username, passwordhash FROM user WHERE username=?"
	row := db.Sql.QueryRowContext(ctx, query, username)

	err := row.Scan(&userId, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth.User{}, err
		}
		logg.Err("row.Scan error while checking if the username is available:", err)
		return auth.User{}, err
	}

	// Convert the string representation to a UUID
	user.Id, err = uuid.FromString(userId)
	if err != nil {
		logg.Err("Error parsing UUID:", err)
		return auth.User{}, err
	}

	return user, nil
}

func (db *DB) UserExist(ctx context.Context, username string) bool {
	if ctx == nil {
		ctx = context.Background()
	}

	query := "SELECT EXISTS(SELECT 1 FROM user WHERE username=?)"
	var exists bool
	row := db.Sql.QueryRowContext(ctx, query, username)

	err := row.Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		return true
	}
	return false
}

// here we run the insert new User query separate from the public function
// it make the code more readable
func (db *DB) insertNewUser(ctx context.Context, username string, passwordhash string) error {
	id, err := uuid.NewV4()
	if err != nil {
		logg.Err("Error generating UUID:", err)
		return err
	}

	sqlStatement := `INSERT INTO user (id, username, passwordhash) VALUES (?, ?, ?)`
	result, err := db.Sql.ExecContext(ctx, sqlStatement, id.String(), username, passwordhash) // Using ExecContext
	if err != nil {
		logg.Err("Error while executing create new user statement:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logg.Err("Error checking rows affected while executing create new user statement:", err)
		return err
	}
	if rowsAffected != 1 {
		logg.Err("No rows affected, user not added")
		return errors.New("user not added")
	}
	return nil
}
