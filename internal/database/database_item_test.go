package database

import (
	"basement/main/internal/items"
	"context"
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

	retrievedItem, err := dbTest.Item(item.ID.String())
	assert.Equal(t, err, nil)

	assert.Equal(t, item.ID, retrievedItem.ID)
	assert.Equal(t, item.Label, retrievedItem.Label)
	assert.Equal(t, item.Description, retrievedItem.Description)
	assert.Equal(t, item.Quantity, retrievedItem.Quantity)
	assert.Equal(t, item.QRCode, retrievedItem.QRCode)
	assert.NotEqual(t, "", retrievedItem.PreviewPicture)

	// ListRow
	var retrievedItemRow *items.ListRow
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

	err = dbTest.UpdateItem(context.Background(), *item)
	assert.Equal(t, err, nil)
	retrievedItem, err := dbTest.Item(item.ID.String())
	assert.Equal(t, err, nil)

	assert.Equal(t, item.Label, retrievedItem.Label)
	assert.Equal(t, item.Description, retrievedItem.Description)

	// ListRow
	var retrievedItemRow *items.ListRow
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

	_, err = dbTest.Item(item.ID.String())
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
		idMap[id] = true
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
		_, err := dbTest.Item(id.String())
		assert.NotEqual(t, nil, err)
	}
}

func TestMoveItem(t *testing.T) {
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
	retrievedItem, err := dbTest.Item(item.ID.String())
	assert.Equal(t, err, nil)
	retrievedBox, err := dbTest.BoxById(box2.ID)
	assert.Equal(t, err, nil)

	assert.Equal(t, box2.ID, retrievedItem.BoxID)
	assert.Equal(t, len(retrievedBox.Items), 1)
	assert.Equal(t, retrievedBox.Items[0].BoxID, box2.ID)

	// Move item out of box2
	err = dbTest.MoveItemToBox(item.ID, uuid.Nil)
	assert.Equal(t, err, nil)
	retrievedItem, err = dbTest.Item(item.ID.String())
	assert.Equal(t, err, nil)
	retrievedBox, err = dbTest.BoxById(box2.ID)
	assert.Equal(t, err, nil)

	assert.Equal(t, retrievedItem.BoxID, uuid.Nil)
	assert.Equal(t, len(retrievedBox.Items), 0)
}
