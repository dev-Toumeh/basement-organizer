package database

import (
	"basement/main/internal/common"
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

	var err error

	// should keep same ID
	EmptyTestDatabase()
	resetShelves()
	shelf := SHELF_1
	err = dbTest.CreateShelf(shelf)
	createdShelf, err := dbTest.Shelf(shelf.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, shelf.ID, createdShelf.ID)

	// item does not exist and should not be created
	shelf.Items = append(shelf.Items, &common.ListRow{ID: ITEM_1.ID})
	assert.Equal(t, len(shelf.Items), 1)
	err = dbTest.CreateShelf(shelf)
	assert.NotEqual(t, err, nil)
	shelf.Items = nil

	// box does not exist and should not be created
	shelf.Boxes = append(shelf.Boxes, &common.ListRow{ID: BOX_1.ID})
	assert.Equal(t, len(shelf.Boxes), 1)
	err = dbTest.CreateShelf(shelf)
	assert.NotEqual(t, err, nil)
	shelf.Items = nil

	err = dbTest.CreateNewItem(*ITEM_1)
	assert.Equal(t, err, nil)
	_, err = dbTest.CreateBox(BOX_1)
	assert.Equal(t, err, nil)
	shelf.Items = append(shelf.Items, &common.ListRow{ID: ITEM_1.ID})
	shelf.Boxes = append(shelf.Boxes, &common.ListRow{ID: BOX_1.ID})
	err = dbTest.CreateShelf(shelf)
	assert.NotEqual(t, err, nil)

	shelf.Items = nil
	shelf.Boxes = nil
	err = dbTest.CreateShelf(shelf)
	createdShelf, err = dbTest.Shelf(shelf.ID)
	assert.Equal(t, err, nil)

	assert.Equal(t, shelf.Label, createdShelf.Label)
	assert.Equal(t, shelf.Description, createdShelf.Description)
	assert.Equal(t, VALID_BASE64_PNG, createdShelf.Picture)
	assert.NotEqual(t, "", createdShelf.PreviewPicture)
	assert.Equal(t, shelf.QRCode, createdShelf.QRCode)
	assert.Equal(t, shelf.Height, createdShelf.Height)
	assert.Equal(t, shelf.Width, createdShelf.Width)
	assert.Equal(t, shelf.Depth, createdShelf.Depth)
	assert.Equal(t, shelf.Rows, createdShelf.Rows)
	assert.Equal(t, shelf.Cols, createdShelf.Cols)

	EmptyTestDatabase()
	resetShelves()
	shelf.Picture = INVALID_BASE64_PNG

	// Expected error log converting picture
	// but NO error returned!
	err = dbTest.CreateShelf(shelf)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, uuid.Nil, shelf.ID)

	createdShelf, err = dbTest.Shelf(shelf.ID)
	assert.Equal(t, "", createdShelf.Picture)
}

func TestDeleteShelf(t *testing.T) {
	EmptyTestDatabase()
	resetShelves()
	resetTestItems()
	var err error
	err = dbTest.CreateShelf(SHELF_1)
	assert.Equal(t, err, nil)
	err = dbTest.DeleteShelf(SHELF_1.ID)
	assert.Equal(t, err, nil)

	// should not delete shelf with an item
	dbTest.CreateShelf(SHELF_1)
	dbTest.CreateNewItem(*ITEM_1)
	dbTest.MoveItemToShelf(ITEM_1.ID, SHELF_1.ID)

	err = dbTest.DeleteShelf(SHELF_1.ID)
	assert.NotEqual(t, err, nil)
}

func TestUpdateShelf(t *testing.T) {
	EmptyTestDatabase()
	resetShelves()

	shelf := SHELF_1

	err := dbTest.createNewShelf(shelf.ID)
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

	updatedShelf, err := dbTest.Shelf(shelf.ID)
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

func TestShelfListRowsPaginated(t *testing.T) {
	EmptyTestDatabase()
	resetShelves()

	id1, _ := dbTest.CreateNewShelf()
	id2, _ := dbTest.CreateNewShelf()
	id3, _ := dbTest.CreateNewShelf()
	shelves, found, err := dbTest.ShelfListRowsPaginated(1, 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(shelves), 2)
	assert.Equal(t, found, 2)
	assert.Equal(t, id1, shelves[0].ID)
	assert.Equal(t, id2, shelves[1].ID)

	shelves, found, err = dbTest.ShelfListRowsPaginated(2, 2)
	assert.Equal(t, id3, shelves[0].ID)
	assert.Equal(t, len(shelves), 2)
	assert.Equal(t, found, 1)
	assert.Equal(t, nil, shelves[1])
}

func TestShelfSearchListRowsPaginated(t *testing.T) {
	EmptyTestDatabase()
	resetShelves()

	for _, shelf := range testShelves() {
		err := dbTest.CreateShelf(&shelf)
		if err != nil {
			t.Fatalf("create shelf setup failed: %v", err)
		}
	}

	// full word
	shelves, found, err := dbTest.ShelfSearchListRowsPaginated(1, 10, "shelf")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(shelves), 10)
	assert.Equal(t, found, 4)

	// part of a word
	shelves, found, err = dbTest.ShelfSearchListRowsPaginated(1, 10, "key")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(shelves), 10)
	assert.Equal(t, found, 2)

	// whitespace and single letter
	shelves, found, err = dbTest.ShelfSearchListRowsPaginated(1, 10, "            a")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(shelves), 10)
	assert.Equal(t, found, 2)

	// 2 parts of 2 different words
	shelves, found, err = dbTest.ShelfSearchListRowsPaginated(1, 10, "sh           3")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(shelves), 10)
	assert.Equal(t, found, 1)

	// 2 parts of 2 different words
	shelves, found, err = dbTest.ShelfSearchListRowsPaginated(1, 10, "Tes Sh ")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(shelves), 10)
	assert.Equal(t, found, 3)

	// 2 parts of 2 different words with pagination
	// page 1: SHELF_1, SHELF_2, page 2: SHELF_3
	shelves, found, err = dbTest.ShelfSearchListRowsPaginated(2, 2, "Tes Sh ")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(shelves), 2)
	assert.Equal(t, found, 1)
	assert.Equal(t, shelves[0].ID, SHELF_3.ID)

	shelves, found, err = dbTest.ShelfSearchListRowsPaginated(1, 10, "")
	assert.Equal(t, err, nil)
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
	err = dbTest.MoveItemToShelf(item.ID, shelf.ID)
	assert.Equal(t, err, nil)

	// Verify item is associated with the shelf
	updatedItem, err := dbTest.ItemListRowByID(item.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, shelf.ID, updatedItem.ShelfID)

	// Move the item out of the shelf
	err = dbTest.MoveItemToShelf(item.ID, uuid.Nil)
	assert.Equal(t, err, nil)
	updatedItem, err = dbTest.ItemListRowByID(item.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, uuid.Nil, updatedItem.ShelfID)

	// Attempt to move a non-existent item
	err = dbTest.MoveItemToShelf(VALID_UUID_NOT_EXISTING, shelf.ID)
	// logg.Err(err)
	assert.NotEqual(t, err, nil)

	// Attempt to move the item to a non-existent shelf
	err = dbTest.MoveItemToShelf(item.ID, VALID_UUID_NOT_EXISTING)
	assert.NotEqual(t, err, nil)

	// Move item to itself (makes no sense do to so)
	err = dbTest.MoveItemToBox(item.ID, item.ID)
	assert.NotEqual(t, err, nil)
	err = dbTest.MoveItemToShelf(item.ID, item.ID)
	assert.NotEqual(t, err, nil)
}

// ID is always checked for `uuid.Nil` in controller, this will never happen.
// func TestCreateShelfWithUUIDNil(t *testing.T) {
// 	shelf := SHELF_1
// 	shelf.ID = uuid.Nil
// 	err := dbTest.CreateShelf(shelf)
// 	_, err = dbTest.Shelf(shelf.ID)
// 	assert.NotEqual(t, err, nil)
// }
