package database

import (
	"basement/main/internal/common"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestInnerListRowsFrom2(t *testing.T) {
	var err error
	EmptyTestDatabase()
	resetTestBoxes()
	resetShelves()
	dbTest.CreateNewItem(*ITEM_1)
	dbTest.CreateBox(BOX_1)
	dbTest.CreateBox(BOX_2)
	dbTest.CreateShelf(SHELF_1)
	dbTest.CreateArea(*AREA_1)
	err = dbTest.MoveItemToBox(ITEM_1.ID, BOX_1.ID)
	assert.Equal(t, err, nil)
	err = dbTest.MoveItemToShelf(ITEM_1.ID, SHELF_1.ID)
	assert.Equal(t, err, nil)
	err = dbTest.MoveItemToArea(ITEM_1.ID, AREA_1.ID)
	assert.Equal(t, err, nil)

	var rows []common.ListRow
	rows, err = dbTest.InnerListRowsFrom2("box", BOX_1.ID, "item_fts")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(rows), 1)
	assert.Equal(t, rows[0].ID, ITEM_1.ID)

	rows, err = dbTest.InnerListRowsFrom2("shelf", SHELF_1.ID, "item_fts")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(rows), 1)
	assert.Equal(t, rows[0].ID, ITEM_1.ID)

	rows, err = dbTest.InnerListRowsFrom2("area", AREA_1.ID, "item_fts")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(rows), 1)
	assert.Equal(t, rows[0].ID, ITEM_1.ID)

	err = dbTest.MoveBoxToBox(BOX_1.ID, BOX_2.ID)
	assert.Equal(t, err, nil)
	err = dbTest.MoveBoxToShelf(BOX_1.ID, SHELF_1.ID)
	assert.Equal(t, err, nil)
	err = dbTest.MoveBoxToArea(BOX_1.ID, AREA_1.ID)
	assert.Equal(t, err, nil)

	rows, err = dbTest.InnerListRowsFrom2("box", BOX_2.ID, "box_fts")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(rows), 1)
	assert.Equal(t, rows[0].ID, BOX_1.ID)

	rows, err = dbTest.InnerListRowsFrom2("shelf", SHELF_1.ID, "box_fts")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(rows), 1)
	assert.Equal(t, rows[0].ID, BOX_1.ID)

	rows, err = dbTest.InnerListRowsFrom2("area", AREA_1.ID, "box_fts")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(rows), 1)
	assert.Equal(t, rows[0].ID, BOX_1.ID)

	err = dbTest.MoveShelfToArea(SHELF_1.ID, AREA_1.ID)
	assert.Equal(t, err, nil)

	rows, err = dbTest.InnerListRowsFrom2("area", AREA_1.ID, "shelf_fts")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(rows), 1)
	assert.Equal(t, rows[0].ID, SHELF_1.ID)

	// inner shelves in area
	dbTest.CreateShelf(SHELF_2)
	dbTest.MoveShelfToArea(SHELF_2.ID, AREA_1.ID)
	rows, err = dbTest.InnerListRowsPaginatedFrom("area_fts", AREA_1.ID, "shelf", "", 1, 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(rows), 1)
	assert.Equal(t, rows[0].ID, SHELF_1.ID)

}
