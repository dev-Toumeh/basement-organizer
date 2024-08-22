package database

import (
	"basement/main/internal/auth"
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
)

// Test if the UserByField works properly
func TestUserByField(t *testing.T) {
	ctx := context.Background()
	_, err := dbTest.Sql.ExecContext(ctx, "INSERT INTO user (id, username, passwordhash) VALUES (?, ?, ?)", "123e4567-e89b-12d3-a456-426614174000", "testuser", "hash")
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Test the User function for existing user
	user, err := dbTest.UserByField(ctx, "username", "testuser")
	if err != nil {
		t.Errorf("Error fetching user: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	// test for non-existing user
	_, err = dbTest.UserByField(context.Background(), "username", "nonexistent")
	if err == nil {
		t.Errorf("Expected an error for non-existing user, got none")
	}

	// test for non-existing user
	_, err = dbTest.UserByField(context.Background(), "id", "123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Errorf("Expected id '123e4567-e89b-12d3-a456-426614174000', got '%s'", user.Id)
	}
}

// Test if the TestUpdateUser works properly
func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	// Step 1: Insert a test user into the database
	testUserId := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174040"))
	_, err := dbTest.Sql.ExecContext(ctx, "INSERT INTO user (id, username, passwordhash) VALUES (?, ?, ?)", "123e4567-e89b-12d3-a456-426614174040", "originaluser", "originalhash")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Step 2: Update the test user
	updatedUser := auth.User{
		Id:           testUserId,
		Username:     "updateduser",
		PasswordHash: "updatedhash",
	}

	err = dbTest.UpdateUser(ctx, updatedUser)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	// Step 3: Verify that the update was successful
	user, err := dbTest.UserByField(ctx, "id", "123e4567-e89b-12d3-a456-426614174040")
	if err != nil {
		t.Fatalf("Failed to fetch updated user: %v", err)
	}
	if user.Username != "updateduser" {
		t.Errorf("Expected username 'updateduser', got '%s'", user.Username)
	}
	if user.PasswordHash != "updatedhash" {
		t.Errorf("Expected password hash 'updatedhash', got '%s'", user.PasswordHash)
	}

	// Step 4: Try to update a non-existent user and expect an error
	nonExistentUser := auth.User{
		Id:           uuid.Must(uuid.FromString("223e4567-e89b-12d3-a456-426614174001")),
		Username:     "nonexistentuser",
		PasswordHash: "somehash",
	}

	err = dbTest.UpdateUser(ctx, nonExistentUser)
	if err == nil {
		t.Errorf("Expected an error when updating a non-existent user, got none")
	}
}
