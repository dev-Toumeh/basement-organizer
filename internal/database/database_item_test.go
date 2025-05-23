package database

import (
	"basement/main/internal/common"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gofrs/uuid/v5"
)

func TestInsertNewItem(t *testing.T) {
	EmptyTestDatabase()
	resetTestItems()
	resetTestBoxes()

	item := ITEM_1

	err := dbTest.insertNewItem(*item)
	assert.Equal(t, err, nil)

	retrievedItem, err := dbTest.ItemById(item.ID)
	assert.Equal(t, err, nil)

	assert.Equal(t, item.ID, retrievedItem.ID)
	assert.Equal(t, item.Label, retrievedItem.Label)
	assert.Equal(t, item.Description, retrievedItem.Description)
	assert.Equal(t, item.Quantity, retrievedItem.Quantity)
	assert.Equal(t, item.QRCode, retrievedItem.QRCode)
	assert.NotEqual(t, "", retrievedItem.PreviewPicture)

	// ListRow
	var retrievedItemRow *common.ListRow
	retrievedItemRow, err = dbTest.ItemListRowByID(item.ID)
	assert.Equal(t, err, nil)

	assert.Equal(t, item.ID, retrievedItemRow.ID)
	assert.Equal(t, item.Label, retrievedItemRow.Label)
}

func TestUpdateItem(t *testing.T) {
	EmptyTestDatabase()
	resetTestItems()
	resetTestBoxes()

	item := ITEM_1

	err := dbTest.insertNewItem(*item)
	assert.Equal(t, err, nil)

	item.Label = "Updated Item Label"
	item.Description = "Updated Description"

	err = dbTest.UpdateItem(*item, true, "image/png")
	assert.Equal(t, err, nil)
	retrievedItem, err := dbTest.ItemById(item.ID)
	assert.Equal(t, err, nil)

	assert.Equal(t, item.Label, retrievedItem.Label)
	assert.Equal(t, item.Description, retrievedItem.Description)

	// ListRow
	var retrievedItemRow *common.ListRow
	retrievedItemRow, err = dbTest.ItemListRowByID(item.ID)
	assert.Equal(t, err, nil)

	assert.Equal(t, item.ID, retrievedItemRow.ID)
	assert.Equal(t, item.Label, retrievedItemRow.Label)
}

func TestDeleteItem(t *testing.T) {
	EmptyTestDatabase()
	resetTestItems()
	resetTestBoxes()

	item := ITEM_1

	err := dbTest.insertNewItem(*item)
	assert.Equal(t, err, nil)

	err = dbTest.DeleteItem(item.ID)
	assert.Equal(t, err, nil)

	_, err = dbTest.ItemById(item.ID)
	assert.NotEqual(t, nil, err)
}

func TestItemExist(t *testing.T) {
	EmptyTestDatabase()
	resetTestItems()
	resetTestBoxes()

	item := ITEM_1

	exists := dbTest.ItemExist("id", item.ID.String())
	assert.Equal(t, false, exists)

	err := dbTest.insertNewItem(*item)
	assert.Equal(t, err, nil)

	exists = dbTest.ItemExist("id", item.ID.String())
	assert.Equal(t, true, exists)
}

// Test retrieving item IDs
func TestItemIDs(t *testing.T) {
	EmptyTestDatabase()
	resetTestItems()
	resetTestBoxes()

	items := testItems()
	for _, item := range items {
		err := dbTest.insertNewItem(item)
		assert.Equal(t, err, nil)
	}

	ids, err := dbTest.ItemIDs()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(items), len(ids))

	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id.String()] = true
	}
	for _, item := range items {
		assert.Equal(t, true, idMap[item.ID.String()])
	}
}

func TestDeleteItems(t *testing.T) {
	EmptyTestDatabase()
	resetTestItems()
	resetTestBoxes()

	items := testItems()
	var itemIDs []uuid.UUID
	for _, item := range items {
		err := dbTest.insertNewItem(item)
		assert.Equal(t, err, nil)
		itemIDs = append(itemIDs, item.ID)
	}

	err := dbTest.DeleteItems(itemIDs)
	assert.Equal(t, err, nil)

	for _, id := range itemIDs {
		_, err := dbTest.ItemById(id)
		assert.NotEqual(t, nil, err)
	}
}

func TestMoveItemToBox(t *testing.T) {
	EmptyTestDatabase()
	resetTestItems()
	resetTestBoxes()

	box1 := BOX_1
	box2 := BOX_2
	item := ITEM_1

	_, err := dbTest.insertNewBox(box1)
	assert.Equal(t, err, nil)
	_, err = dbTest.insertNewBox(box2)
	assert.Equal(t, err, nil)
	err = dbTest.insertNewItem(*item)
	assert.Equal(t, err, nil)

	// Move item in box2
	err = dbTest.MoveItemToBox(item.ID, box2.ID)
	assert.Equal(t, err, nil)
	retrievedItem, err := dbTest.ItemById(item.ID)
	assert.Equal(t, err, nil)
	retrievedBox, err := dbTest.BoxById(box2.ID)
	assert.Equal(t, err, nil)

	assert.Equal(t, box2.ID, retrievedItem.BoxID)
	assert.Equal(t, len(retrievedBox.Items), 1)
	assert.Equal(t, retrievedBox.Items[0].BoxID, box2.ID)

	// Move item out of box2
	err = dbTest.MoveItemToBox(item.ID, uuid.Nil)
	assert.Equal(t, err, nil)
	retrievedItem, err = dbTest.ItemById(item.ID)
	assert.Equal(t, err, nil)
	retrievedBox, err = dbTest.BoxById(box2.ID)
	assert.Equal(t, err, nil)

	assert.Equal(t, retrievedItem.BoxID, uuid.Nil)
	assert.Equal(t, len(retrievedBox.Items), 0)
}

func TestMoveItemToObject(t *testing.T) {
	// Reset the test database and objects
	EmptyTestDatabase()
	resetTestItems()
	resetAreas()
	resetShelves()
	resetTestBoxes()

	area := AREA_1
	shelf := SHELF_1
	box := BOX_1
	item := ITEM_1

	// Insert objects into the database
	_, err := dbTest.insertNewArea(*area)
	assert.Equal(t, nil, err)
	err = dbTest.CreateShelf(shelf)
	assert.Equal(t, nil, err)
	_, err = dbTest.insertNewBox(box)
	assert.Equal(t, nil, err)
	err = dbTest.insertNewItem(*item)
	assert.Equal(t, nil, err)

	// Move item to area
	err = dbTest.MoveItemToObject(item.ID, area.ID, "area")
	assert.Equal(t, nil, err)
	retrievedItem, err := dbTest.ItemById(item.ID)
	assert.Equal(t, nil, err)

	assert.Equal(t, area.ID, retrievedItem.AreaID)
	assert.Equal(t, area.Label, retrievedItem.AreaLabel)

	// Move item to shelf
	err = dbTest.MoveItemToObject(item.ID, shelf.ID, "shelf")
	assert.Equal(t, nil, err)
	retrievedItem, err = dbTest.ItemById(item.ID)
	assert.Equal(t, nil, err)

	assert.Equal(t, shelf.ID, retrievedItem.ShelfID)
	assert.Equal(t, shelf.Label, retrievedItem.ShelfLabel)

	// Move item to box
	err = dbTest.MoveItemToObject(item.ID, box.ID, "box")
	assert.Equal(t, nil, err)
	retrievedItem, err = dbTest.ItemById(item.ID)
	assert.Equal(t, nil, err)

	assert.Equal(t, box.ID, retrievedItem.BoxID)
	assert.Equal(t, box.Label, retrievedItem.BoxLabel)

	// Test invalid object type
	err = dbTest.MoveItemToObject(item.ID, box.ID, "invalid")
	assert.NotEqual(t, nil, err)
}
