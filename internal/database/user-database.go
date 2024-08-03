package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	Id           string
	Username     string
	PasswordHash string
}

var ctx context.Context

// CreateNewUser inserts a new user into the database with the given username and passwordHash
func (db *DB) CreateNewUser(ctx context.Context, username string, passwordhash string) error {
	if _, err := db.User(ctx, username); err != nil {
		// finding no rows means we can use the username to create new user
		if err == sql.ErrNoRows {
			err := db.inserNewUser(ctx, username, passwordhash)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		// otherwise return the real error
		return err
	}
	// if no error exist than you should choice another username
	return errors.New("username already exist")
}

// check if the username is available
// if the user exist it will return iuser struct with nil
// if not it will return empty user struct with err
func (db *DB) User(ctx context.Context, username string) (User, error) {
	var user User
	if ctx == nil {
		ctx = context.Background()
	}

	query := "SELECT id, username, passwordhash FROM user WHERE username=?"
	row := db.Sql.QueryRowContext(ctx, query, username)
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, err
		}
		log.Println("Error while checking if the username is available:", err)
		return User{}, err
	}

	return user, nil
}

// here we run the insert new User query separate from the public function
// it make the code more readable
func (db *DB) inserNewUser(ctx context.Context, username string, passwordhash string) error {
	id, err := uuid.NewV4() // Error handling for UUID generation
	if err != nil {
		log.Printf("Error generating UUID: %v", err)
		return err
	}

	sqlStatement := `INSERT INTO user (id, username, passwordhash) VALUES (?, ?, ?)`
	result, err := db.Sql.ExecContext(ctx, sqlStatement, id.String(), username, passwordhash) // Using ExecContext
	if err != nil {
		log.Printf("Error while executing create new user statement: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected while executing create new user statement: %v", err)
		return err
	}
	if rowsAffected != 1 {
		log.Println("No rows affected, user not added")
		return errors.New("user not added")
	}
	return nil
}
