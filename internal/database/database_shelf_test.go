package database

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gofrs/uuid/v5"
)

func init() {
	// logg.EnableInfoLogger()
	// logg.EnableDebugLogger()
	// logg.EnableErrorLogger()
}

func TestExistShelf(t *testing.T) {
	EmptyTestDatabase()
	exists, err := dbTest.Exists("shelf", SHELF_VALID_UUID_1)
	assert.Equal(t, err, nil)
	assert.Equal(t, exists, false)

	dbTest.createNewShelf(SHELF_VALID_UUID_1)
	exists, err = dbTest.Exists("shelf", SHELF_VALID_UUID_1)
	assert.Equal(t, err, nil)
	assert.Equal(t, exists, true)
}

func TestCreateNewShelf(t *testing.T) {
	EmptyTestDatabase()
	id, err := dbTest.CreateNewShelf()
	assert.Equal(t, err, nil)
	assert.NotEqual(t, id, uuid.Nil)
}

func TestCreateShelf(t *testing.T) {
	EmptyTestDatabase()
	resetShelves()

	shelf := SHELF_1

	err := dbTest.CreateShelf(shelf)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, uuid.Nil, SHELF_VALID_UUID_1)

	createdShelf, err := dbTest.Shelf(SHELF_VALID_UUID_1)
	assert.Equal(t, err, nil)

	assert.Equal(t, shelf.Label, createdShelf.Label)
	assert.Equal(t, shelf.Description, createdShelf.Description)
	assert.Equal(t, VALID_BASE64_PNG, createdShelf.Picture)
	assert.NotEqual(t, "", createdShelf.PreviewPicture)
	assert.Equal(t, shelf.QRcode, createdShelf.QRcode)
	assert.Equal(t, shelf.Height, createdShelf.Height)
	assert.Equal(t, shelf.Width, createdShelf.Width)
	assert.Equal(t, shelf.Depth, createdShelf.Depth)
	assert.Equal(t, shelf.Rows, createdShelf.Rows)
	assert.Equal(t, shelf.Cols, createdShelf.Cols)

	EmptyTestDatabase()
	shelf.Picture = INVALID_BASE64_PNG

	// Expected error log converting picture
	// but NO error returned!
	err = dbTest.CreateShelf(shelf)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, uuid.Nil, SHELF_VALID_UUID_1)

	createdShelf, err = dbTest.Shelf(SHELF_VALID_UUID_1)
	assert.Equal(t, "", createdShelf.Picture)
}

func TestDeleteShelf(t *testing.T) {
	EmptyTestDatabase()
	dbTest.createNewShelf(SHELF_VALID_UUID_1)
	id, err := dbTest.CreateNewShelf()
	assert.Equal(t, err, nil)
	assert.NotEqual(t, id, uuid.Nil)

	err = dbTest.DeleteShelf(SHELF_VALID_UUID_1)
	assert.Equal(t, err, nil)
}

func TestUpdateShelf(t *testing.T) {
	EmptyTestDatabase()
	resetShelves()

	shelf := SHELF_1

	err := dbTest.createNewShelf(shelf.Id)
	assert.Equal(t, err, nil)

	shelf.Label = "Updated Label"
	shelf.Description = "Updated Description"
	shelf.Height = 2.5
	shelf.Width = 1.8
	shelf.Depth = 0.7
	shelf.Rows = 4
	shelf.Cols = 5
	shelf.Picture = VALID_BASE64_PNG

	err = dbTest.UpdateShelf(shelf)
	assert.Equal(t, err, nil)

	updatedShelf, err := dbTest.Shelf(shelf.Id)
	assert.Equal(t, err, nil)

	assert.Equal(t, "Updated Label", updatedShelf.Label)
	assert.Equal(t, "Updated Description", updatedShelf.Description)
	assert.Equal(t, float32(2.5), updatedShelf.Height)
	assert.Equal(t, float32(1.8), updatedShelf.Width)
	assert.Equal(t, float32(0.7), updatedShelf.Depth)
	assert.Equal(t, 4, updatedShelf.Rows)
	assert.Equal(t, 5, updatedShelf.Cols)

	expectedPicture := VALID_BASE64_PNG
	assert.Equal(t, expectedPicture, updatedShelf.Picture)
	assert.NotEqual(t, "", updatedShelf.PreviewPicture)
}

func TestMoveItemToShelf(t *testing.T) {
	EmptyTestDatabase()
	resetTestItems()
	resetShelves()

	item := ITEM_1

	err := dbTest.CreateNewItem(*item)
	assert.Equal(t, err, nil)

	// Create a test shelf
	shelf := SHELF_1
	err = dbTest.CreateShelf(shelf)
	assert.Equal(t, err, nil)

	// Move the item to the shelf
	err = dbTest.MoveItemToShelf(item.ID, shelf.Id)
	assert.Equal(t, err, nil)

	// Verify item is associated with the shelf
	updatedItem, err := dbTest.ItemListRowByID(item.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, shelf.Id, updatedItem.ShelfID)

	// Move the item out of the shelf
	err = dbTest.MoveItemToShelf(item.ID, uuid.Nil)
	assert.Equal(t, err, nil)
	updatedItem, err = dbTest.ItemListRowByID(item.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, uuid.Nil, updatedItem.ShelfID)

	// Attempt to move a non-existent item
	err = dbTest.MoveItemToShelf(VALID_UUID_NOT_EXISTING, shelf.Id)
	// logg.Err(err)
	assert.NotEqual(t, err, nil)

	// Attempt to move the item to a non-existent shelf
	err = dbTest.MoveItemToShelf(item.ID, VALID_UUID_NOT_EXISTING)
	// logg.Err(err)
	assert.NotEqual(t, err, nil)
}
