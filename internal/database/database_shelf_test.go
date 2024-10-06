package database

import (
	"basement/main/internal/items"
	"basement/main/internal/shelves"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gofrs/uuid/v5"
)

func init() {
	// logg.EnableInfoLogger()
	// logg.EnableDebugLogger()
	// logg.EnableErrorLogger()
}

var SHELF_VALID_UUID uuid.UUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000"))
var ITEM_VALID_UUID uuid.UUID = uuid.Must(uuid.FromString("133e4567-e89b-12d3-a456-426614174000"))
var VALID_UUID_NOT_EXISTING uuid.UUID = uuid.Must(uuid.FromString("033e4567-e89b-12d3-a456-426614174000"))

const VALID_BASE64_PNG = "iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAIAAAACUFjqAAAAtUlEQVR4nGJp2XGEAQb+/P49J7cgY8ZUuAgTnDUjI7vUQf3m5e3/zxyakZENFW3ZcURGQf/r52cQBGfLKhq0bD/MqKBu+ufnL4jSm5e3QxjmtuEfPnyCGn7z8na4BAMDg7quZ2mia2thMAMDA0j31TMb4XJr5s2BMKr71zIwMLAwYIDq/rWMMDaLobs7mjTEWJC6CeuYjL08+o/eU9f1RDPgsbpTxvQpjMjBAvEucrAAAgAA//+Elk5AOfCu8QAAAABJRU5ErkJggg=="
const VALID_BASE64_PREVIEW_PNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGJaOLEJEAAA//8DkwG35JmAnAAAAABJRU5ErkJggg=="
const INVALID_BASE64_PNG = "invalid base 64"

func TestExistShelf(t *testing.T) {
	EmptyTestDatabase()
	exists, err := dbTest.Exists("shelf", SHELF_VALID_UUID)
	assert.Equal(t, err, nil)
	assert.Equal(t, exists, false)

	dbTest.createNewShelf(SHELF_VALID_UUID)
	exists, err = dbTest.Exists("shelf", SHELF_VALID_UUID)
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

	shelf := &shelves.Shelf{
		ID:             SHELF_VALID_UUID,
		Label:          "Test Shelf",
		Description:    "A shelf for testing",
		Picture:        VALID_BASE64_PNG,
		PreviewPicture: "",
		QRcode:         "testqrcode",
		Height:         2.0,
		Width:          1.5,
		Depth:          0.5,
		Rows:           3,
		Cols:           4,
	}

	err := dbTest.CreateShelf(shelf)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, uuid.Nil, SHELF_VALID_UUID)

	createdShelf, err := dbTest.Shelf(SHELF_VALID_UUID)
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
	assert.NotEqual(t, uuid.Nil, SHELF_VALID_UUID)

	createdShelf, err = dbTest.Shelf(SHELF_VALID_UUID)
	assert.Equal(t, "", createdShelf.Picture)
}

func TestDeleteShelf(t *testing.T) {
	EmptyTestDatabase()
	dbTest.createNewShelf(SHELF_VALID_UUID)
	id, err := dbTest.CreateNewShelf()
	assert.Equal(t, err, nil)
	assert.NotEqual(t, id, uuid.Nil)

	err = dbTest.DeleteShelf(SHELF_VALID_UUID)
	assert.Equal(t, err, nil)
}

func TestUpdateShelf(t *testing.T) {
	EmptyTestDatabase()

	shelf := &shelves.Shelf{
		ID:             SHELF_VALID_UUID,
		Label:          "Original Label",
		Description:    "Original Description",
		Picture:        "", // No picture initially
		PreviewPicture: "",
		QRcode:         "",
		Height:         2.0,
		Width:          1.5,
		Depth:          0.5,
		Rows:           3,
		Cols:           4,
	}

	err := dbTest.createNewShelf(SHELF_VALID_UUID)
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

func TestMoveItemToShelf(t *testing.T) {
	EmptyTestDatabase()

	// Create a test item
	item := items.Item{
		ID:          ITEM_VALID_UUID,
		Label:       "Test Item",
		Description: "A test item",
		Quantity:    1,
		Weight:      "1kg",
		QRcode:      "testitemqrcode",
	}
	err := dbTest.CreateNewItem(item)
	assert.Equal(t, err, nil)

	// Create a test shelf
	shelf := &shelves.Shelf{
		ID:          SHELF_VALID_UUID,
		Label:       "Test Shelf",
		Description: "A test shelf",
	}
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
	// logg.Err(err)
	assert.NotEqual(t, err, nil)
}

// func TestMoveItemToShelf(t *testing.T) {
// 	EmptyTestDatabase()
//
// Create a test item
// 	item := items.Item2{
// 		Id:          ITEM_VALID_UUID,
// 		Label:       "Test Item",
// 		Description: "A test item",
// 		Quantity:    1,
// 		Weight:      "1kg",
// 		QRcode:      "testitemqrcode",
// 	}
// 	err := dbTest.CreateNewItem2(item)
// 	assert.Equal(t, err, nil)
//
// 	// Create a test shelf
// 	shelf := &shelves.Shelf{
// 		Id:          SHELF_VALID_UUID,
// 		Label:       "Test Shelf",
// 		Description: "A test shelf",
// 	}
// 	err = dbTest.CreateShelf(shelf)
// 	assert.Equal(t, err, nil)
//
// 	// Move the item to the shelf
// 	err = dbTest.MoveItemToShelf(item.Id, shelf.Id)
// 	assert.Equal(t, err, nil)
//
// 	// Verify item is associated with the shelf
// 	updatedItem, err := dbTest.ListItemById2(item.Id)
// 	assert.Equal(t, err, nil)
// 	assert.Equal(t, shelf.Id, updatedItem.Shelf_Id)
//
// 	// Move the item out of the shelf
// 	err = dbTest.MoveItemToShelf(item.Id, uuid.Nil)
// 	assert.Equal(t, err, nil)
// 	updatedItem, err = dbTest.ListItemById2(item.Id)
// 	assert.Equal(t, err, nil)
// 	assert.Equal(t, uuid.Nil, updatedItem.Shelf_Id)
//
// 	// Attempt to move a non-existent item
// 	err = dbTest.MoveItemToShelf(VALID_UUID_NOT_EXISTING, shelf.Id)
// 	// logg.Err(err)
// 	assert.NotEqual(t, err, nil)
//
// 	// Attempt to move the item to a non-existent shelf
// 	err = dbTest.MoveItemToShelf(item.Id, VALID_UUID_NOT_EXISTING)
// 	// logg.Err(err)
// 	assert.NotEqual(t, err, nil)
// }
