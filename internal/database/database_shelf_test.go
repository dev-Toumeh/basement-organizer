package database

import (
	"basement/main/internal/items"
	"encoding/base64"
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

	// shelf, err := dbTest.Shelf(id)
	// assert.Equal(t, err, nil)
	// assert.Equal(t, len(shelf.Picture), 0)
	// assert.Equal(t, len(shelf.PreviewPicture), 0)
}

func TestCreateShelf(t *testing.T) {
	EmptyTestDatabase()
	picdata := []byte("test picture data")
	previewpicdata := []byte("test preview data")
	picdata64 := base64.StdEncoding.EncodeToString(picdata)
	previewpicdata64 := base64.StdEncoding.EncodeToString(previewpicdata)

	shelf := &items.Shelf{
		ID:             SHELF_VALID_UUID,
		Label:          "Test Shelf",
		Description:    "A shelf for testing",
		Picture:        picdata64,
		PreviewPicture: previewpicdata64,
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
	assert.Equal(t, picdata64, createdShelf.Picture)
	assert.Equal(t, previewpicdata64, createdShelf.PreviewPicture)
	assert.Equal(t, shelf.QRcode, createdShelf.QRcode)
	assert.Equal(t, shelf.Height, createdShelf.Height)
	assert.Equal(t, shelf.Width, createdShelf.Width)
	assert.Equal(t, shelf.Depth, createdShelf.Depth)
	assert.Equal(t, shelf.Rows, createdShelf.Rows)
	assert.Equal(t, shelf.Cols, createdShelf.Cols)

	// created pic data should be encoded in base64
	// pic64 := bytes.Buffer{}
	// pic64enc := base64.NewEncoder(base64.StdEncoding, &pic64)
	// pic64enc.Write(picdata)
	// pic64enc.Close()
	// assert.Equal(t, pic64.Bytes(), createdShelf.Picture)

	// created preview pic data should be encoded in base64
	// previewpic64 := bytes.Buffer{}
	// previewpic64enc := base64.NewEncoder(base64.StdEncoding, &previewpic64)
	// previewpic64enc.Write(previewpicdata)
	// previewpic64enc.Close()
	// assert.Equal(t, previewpic64.Bytes(), createdShelf.PreviewPicture)
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

	shelf := &items.Shelf{
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

	picdata := []byte("test picture data")
	previewpicdata := []byte("test preview data")

	picdata64 := base64.StdEncoding.EncodeToString(picdata)
	previewpicdata64 := base64.StdEncoding.EncodeToString(previewpicdata)

	err := dbTest.createNewShelf(SHELF_VALID_UUID)
	assert.Equal(t, err, nil)

	shelf.Label = "Updated Label"
	shelf.Description = "Updated Description"
	shelf.Height = 2.5
	shelf.Width = 1.8
	shelf.Depth = 0.7
	shelf.Rows = 4
	shelf.Cols = 5
	shelf.Picture = picdata64
	shelf.PreviewPicture = previewpicdata64

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

	// pictureData, err := base64.StdEncoding.DecodeString(updatedShelf.Picture)
	// base64.StdEncoding.Decode(picture, updatedShelf.Picture)

	// pic64 := bytes.Buffer{}
	// pic64enc := base64.NewEncoder(base64.StdEncoding, &pic64)
	// pic64enc.Write(pic)
	// pic64enc.Close()
	// expectedPicture := pic64.Bytes()
	// assert.Equal(t, err, nil)
	expectedPicture := picdata64
	assert.Equal(t, expectedPicture, updatedShelf.Picture)
	//
	// previewpic64 := bytes.Buffer{}
	// previewpic64enc := base64.NewEncoder(base64.StdEncoding, &previewpic64)
	// previewpic64enc.Write(previewpic)
	// previewpic64enc.Close()
	// expectedPreviewPicture := previewpic64.Bytes()

	expectedPreviewPicture := previewpicdata64
	// assert.Equal(t, err, nil)
	assert.Equal(t, expectedPreviewPicture, updatedShelf.PreviewPicture)
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
	shelf := &items.Shelf{
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
