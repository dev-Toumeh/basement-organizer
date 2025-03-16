package database

import (
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"fmt"
	"slices"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gofrs/uuid/v5"
)

func TestInsertNewBox(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testBox := BOX_1

	// Step 2: Insert boxes
	for _, box := range testBoxes() {
		_, err := dbTest.insertNewBox(box)
		if err != nil {
			t.Fatalf("insertNewBox failed: %v", err)
		}
	}

	// Step 3: Insert items
	for _, item := range testItems() {
		// fmt.Println(item)
		err := dbTest.insertNewItem(item)
		if err != nil {
			t.Fatalf("insertNewItem failed: %v", err)
		}
	}

	//	Step 4: Verify that the insertion of items was successful
	for _, item := range testItems() {
		_, err := dbTest.ItemByField("id", item.ID.String())
		if err != nil {
			t.Fatalf("get item error: %v", err)
		}
	}

	fetchedBox, err := dbTest.BoxById(testBox.ID)
	if err != nil {
		t.Fatalf(" the function BoxByfield not working properly : %v %v", err.Error(), testBox)
	}

	//Compare the fetched box with the original test box
	assert.Equal(t, testBox.Label, fetchedBox.Label)
	assert.Equal(t, testBox.Description, fetchedBox.Description)
	assert.NotEqual(t, "", fetchedBox.PreviewPicture)
	assert.Equal(t, testBox.OuterBox, nil)
	assert.Equal(t, testBox.OuterBoxID, uuid.Nil)
	assert.Equal(t, testBox.InnerBoxes, nil)
	assert.Equal(t, testBox.ShelfID, uuid.Nil)
	assert.Equal(t, testBox.AreaID, uuid.Nil)

	duplicateBox := *BOX_1

	_, err = dbTest.insertNewBox(&duplicateBox)
	if err == nil {
		t.Errorf("Expected an error when inserting a box with an existing ID, got none")
	}
}

func TestInsertNewBoxWithOuterBox(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()
	var err error
	outerBox := BOX_1
	innerBox := BOX_2
	innerBox.OuterBoxID = outerBox.ID
	_, err = dbTest.insertNewBox(innerBox)
	assert.Equal(t, err, nil)
	_, err = dbTest.insertNewBox(outerBox)
	assert.Equal(t, err, nil)

	fetchedOuterBox, err := dbTest.BoxById(outerBox.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(fetchedOuterBox.InnerBoxes), 1)
	fetchedInnerBox, err := dbTest.BoxById(innerBox.ID)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, fetchedInnerBox.OuterBox, nil)
}

func TestBoxByField(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testBox := BOX_1
	dbTest.insertNewBox(testBox)

	// Testing retrieval by a field that should exist
	fetchedBox, err := dbTest.BoxById(testBox.ID)
	assert.Equal(t, err, nil)
	if err != nil {
		t.Fatalf("Failed to retrieve box by id: %v", err)
	}
	assert.Equal(t, fetchedBox.ID.String(), testBox.ID.String())

	// Testing retrieval by a non-existent field
	_, err = dbTest.BoxByField("non_existent_field", "some_value")
	assert.NotEqual(t, err, nil)

}

func TestCreateNewBox(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testBox := BOX_1

	// Testing creation of a new box that does not already exist
	_, err := dbTest.CreateBox(testBox)
	assert.Equal(t, nil, err)
	if err != nil {
		t.Fatalf("Failed to create new box: %v", err)
	}

	// Verify box was created
	exists := dbTest.BoxExistById(testBox.ID)
	assert.Equal(t, true, exists)

	// Test creating the same box again to trigger an error
	_, err = dbTest.CreateBox(testBox)
	assert.NotEqual(t, nil, err)

}

func TestBoxIDs(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	expectedIDs := []uuid.UUID{BOX_1.ID, BOX_2.ID, BOX_3.ID, BOX_4.ID}

	// Insert test boxes into the database
	for _, testBox := range testBoxes() {
		_, err := dbTest.insertNewBox(testBox)
		if err != nil {
			t.Fatalf("Failed to insert test box: %v", err)
		}
	}

	// Call the BoxIDs function
	actualIDs, err := dbTest.BoxIDs()
	if err != nil {
		t.Fatalf("BoxIDs function returned an error: %v", err)
	}

	// Verify the results
	for _, v := range expectedIDs {
		assert.Equal(t, slices.Contains(actualIDs, v), true)
	}
}

func TestBoxUpdateO(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testBox := BOX_1
	assert.NotEqual(t, testBox.Picture, "")
	assert.Equal(t, testBox.Picture, VALID_BASE64_PNG)
	assert.Equal(t, testBox.PreviewPicture, "")
	_, err := dbTest.insertNewBox(testBox)
	if err != nil {
		t.Fatalf("error while inserting the box: %v", err)
	}

	oldLabel := testBox.Label
	oldDescr := testBox.Description
	oldPre := testBox.PreviewPicture

	testBox.Description = "updated"
	testBox.Label = "updated"
	testBox.Picture = VALID_BASE64_PNG_2
	assert.NotEqual(t, oldDescr, testBox.Description)
	assert.NotEqual(t, oldLabel, testBox.Label)
	assert.NotEqual(t, oldPre, "")

	logg.EnableDebugLogger()
	logg.EnableInfoLogger()
	logg.EnableErrorLogger()
	err = dbTest.UpdateBox(*testBox, false)
	if err != nil {
		t.Fatalf("error while updating the box: %v", err)
	}

	// Retrieve the updated box from the database
	updatedBox, err := dbTest.BoxById(testBox.ID)
	if err != nil {
		t.Fatalf("error while retrieving the updated box: %v", err)
	}

	// Assert that the box was updated correctly (using individual asserts)
	assert.Equal(t, testBox.Label, updatedBox.Label)
	assert.Equal(t, testBox.Description, updatedBox.Description)
	assert.Equal(t, testBox.Picture, updatedBox.Picture)
	assert.NotEqual(t, updatedBox.PreviewPicture, "")
	assert.NotEqual(t, oldPre, updatedBox.PreviewPicture)
	assert.Equal(t, testBox.QRCode, updatedBox.QRCode)
	assert.Equal(t, testBox.OuterBoxID, updatedBox.OuterBoxID)

	assert.NotEqual(t, oldLabel, updatedBox.Label)
	assert.NotEqual(t, oldDescr, updatedBox.Description)

	row, err := dbTest.listRowByID("box_fts", testBox.ID)
	assert.NotEqual(t, row.PreviewPicture, "")
	assert.Equal(t, row.PreviewPicture, updatedBox.PreviewPicture)
	rows, err := dbTest.BoxListRows("", 1, 1)
	assert.NotEqual(t, rows[0].PreviewPicture, "")
	fmt.Println(rows[0].PreviewPicture)
	frows, err := common.FilledRows(dbTest.BoxListRows, "", 1, 1, 1, common.ListRowTemplateOptions{})
	assert.NotEqual(t, frows[0].PreviewPicture, "")
	assert.NotEqual(t, frows[0].PreviewPicture, oldPre)
	assert.Equal(t, frows[0].PreviewPicture, updatedBox.PreviewPicture)

	EmptyTestDatabase()
}

func TestBoxUpdateBoxInBoxError(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()
	var err error

	_, err = dbTest.insertNewBox(BOX_1)
	if err != nil {
		t.Fatalf("error while inserting the box: %v", err)
	}

	BOX_1.OuterBoxID = BOX_1.ID

	err = dbTest.UpdateBox(*BOX_1, false)
	assert.NotEqual(t, err, nil)

	EmptyTestDatabase()
}

func TestBoxUpdateIgnorePicture(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testBox := BOX_1
	testBox.PreviewPicture = VALID_BASE64_PREVIEW_PNG
	_, err := dbTest.insertNewBox(testBox)
	if err != nil {
		t.Fatalf("error while inserting the box: %v", err)
	}

	oldLabel := testBox.Label
	oldDescr := testBox.Description
	oldPicture := testBox.Picture
	oldPreviewPicture := testBox.PreviewPicture

	testBox.Description = "updated"
	testBox.Label = "updated"
	testBox.Picture = ""
	testBox.PreviewPicture = ""
	assert.NotEqual(t, oldDescr, testBox.Description)
	assert.NotEqual(t, oldLabel, testBox.Label)
	assert.NotEqual(t, oldPicture, testBox.Picture)
	assert.NotEqual(t, oldPreviewPicture, testBox.PreviewPicture)

	err = dbTest.UpdateBox(*testBox, true)
	if err != nil {
		t.Fatalf("error while updating the box: %v", err)
	}

	// Retrieve the updated box from the database
	updatedBox, err := dbTest.BoxById(testBox.ID)
	if err != nil {
		t.Fatalf("error while retrieving the updated box: %v", err)
	}

	// Assert that the box was updated correctly (using individual asserts)
	assert.Equal(t, testBox.Label, updatedBox.Label)
	assert.Equal(t, testBox.Description, updatedBox.Description)
	assert.Equal(t, oldPicture, updatedBox.Picture)
	assert.Equal(t, oldPreviewPicture, updatedBox.PreviewPicture)
	assert.Equal(t, testBox.QRCode, updatedBox.QRCode)
	assert.Equal(t, testBox.OuterBoxID, updatedBox.OuterBoxID)

	assert.NotEqual(t, oldLabel, updatedBox.Label)
	assert.NotEqual(t, oldDescr, updatedBox.Description)

	EmptyTestDatabase()
}

func TestBoxUpdateRemovePicture(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testBox := BOX_1
	_, err := dbTest.insertNewBox(testBox)
	if err != nil {
		t.Fatalf("error while inserting the box: %v", err)
	}

	oldPicture := testBox.Picture
	oldPreviewPicture := testBox.PreviewPicture

	testBox.Picture = ""
	testBox.PreviewPicture = ""
	assert.NotEqual(t, oldPicture, testBox.Picture)
	assert.NotEqual(t, oldPreviewPicture, testBox.PreviewPicture)

	err = dbTest.UpdateBox(*testBox, false)
	if err != nil {
		t.Fatalf("error while updating the box: %v", err)
	}

	// Retrieve the updated box from the database
	updatedBox, err := dbTest.BoxById(testBox.ID)
	if err != nil {
		t.Fatalf("error while retrieving the updated box: %v", err)
	}

	// Assert that the box was updated correctly (using individual asserts)
	assert.Equal(t, "", updatedBox.Picture)
	assert.Equal(t, "", updatedBox.PreviewPicture)

	EmptyTestDatabase()
}

func TestBoxUpdateShelf(t *testing.T) {
	EmptyTestDatabase()
	resetShelves()

	shelf := SHELF_1
	box := BOX_1
	box.ShelfID = shelf.ID

	err := dbTest.CreateShelf(shelf)
	assert.Equal(t, err, nil)
	_, err = dbTest.CreateBox(box)
	assert.Equal(t, err, nil)
	boxrow, _ := dbTest.BoxListRowByID(box.ID)
	assert.Equal(t, boxrow.ShelfID, shelf.ID)
	assert.Equal(t, boxrow.ShelfLabel, shelf.Label)

	err = dbTest.CreateShelf(SHELF_2)
	assert.Equal(t, err, nil)
	box.ShelfID = SHELF_2.ID

	dbTest.UpdateBox(*box, true)
	boxrow, err = dbTest.BoxListRowByID(box.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, boxrow.ShelfID, SHELF_2.ID)
	assert.Equal(t, boxrow.ShelfLabel, SHELF_2.Label)
}

func TestDeleteBox(t *testing.T) {
	EmptyTestDatabase()
	resetTestItems()
	resetTestBoxes()

	ITEM_1.BoxID = BOX_1.ID
	BOX_2.OuterBoxID = BOX_1.ID

	for _, box := range testBoxes() {
		_, err := dbTest.insertNewBox(box)
		if err != nil {
			t.Fatalf("insertNewBox failed: %v", err)
		}
	}

	err := dbTest.insertNewItem(*ITEM_1)
	if err != nil {
		t.Fatalf("insertNewItem failed: %v", err)
	}

	err = dbTest.DeleteBox(BOX_1.ID)
	assert.NotEqual(t, err, nil) // err: can't delete, box not empty

	err = dbTest.DeleteItem(ITEM_1.ID)
	if err != nil {
		t.Fatalf("the item was not deleted: %v", err)
	}
	err = dbTest.DeleteBox(BOX_2.ID)
	if err != nil {
		t.Fatalf("deleting the innerbox was not succeed: %v", err)
	}

	err = dbTest.DeleteBox(BOX_1.ID)
	if err != nil {
		t.Fatalf("delete the box after deleting the data inside of it was not succeed")
	}
}

func TestMoveBoxToBox(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()
	resetTestItems()

	innerBox := BOX_1
	innerBox2 := BOX_3
	outerBox := BOX_2

	// Insert test boxes into the database using range
	for _, testBox := range testBoxes() {
		_, err := dbTest.insertNewBox(testBox)
		if err != nil {
			t.Fatalf("Failed to insert test box: %v", err)
		}
	}

	// 1. Test successful move
	err := dbTest.MoveBoxToBox(innerBox.ID, outerBox.ID)
	if err != nil {
		t.Fatalf("MoveBox function returned an error: %v", err)
	}
	err = dbTest.MoveBoxToBox(innerBox2.ID, outerBox.ID)
	assert.Equal(t, err, nil)

	// inner box
	updatedInnerBox, err := dbTest.BoxById(innerBox.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated inner box: %v", err)
	}
	assert.Equal(t, outerBox.ID, updatedInnerBox.OuterBoxID)
	assert.Equal(t, outerBox.ID, updatedInnerBox.OuterBox.ID)
	assert.Equal(t, outerBox.Label, updatedInnerBox.OuterBox.Label)
	updatedInnerBox2, err := dbTest.BoxById(innerBox.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, outerBox.ID, updatedInnerBox2.OuterBoxID)
	assert.Equal(t, outerBox.ID, updatedInnerBox2.OuterBox.ID)
	assert.Equal(t, outerBox.Label, updatedInnerBox2.OuterBox.Label)

	// outer box
	updatedOuterBox, err := dbTest.BoxById(outerBox.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(updatedOuterBox.InnerBoxes), 2)
	assert.Equal(t, updatedOuterBox.InnerBoxes[0].ID, innerBox.ID)
	assert.Equal(t, updatedOuterBox.InnerBoxes[0].Label, innerBox.Label)
	assert.Equal(t, updatedOuterBox.InnerBoxes[1].ID, innerBox2.ID)
	assert.Equal(t, updatedOuterBox.InnerBoxes[1].Label, innerBox2.Label)

	// 2. Test move to non-existent box (should return an error)
	nonExistentBoxId := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174003"))
	err = dbTest.MoveBoxToBox(innerBox.ID, nonExistentBoxId)
	assert.Equal(t, err, err)

	// Move innerbox out of outerbox
	err = dbTest.MoveBoxToBox(innerBox.ID, uuid.Nil)
	assert.Equal(t, err, nil)
	updatedInnerBox, err = dbTest.BoxById(innerBox.ID)
	assert.Equal(t, updatedInnerBox.OuterBoxID, uuid.Nil)
	assert.Equal(t, updatedInnerBox.OuterBox, nil)

	// Move to itself
	err = dbTest.MoveBoxToBox(innerBox.ID, innerBox.ID)
	assert.NotEqual(t, err, nil)

	// Inner and outerbox can't be inside eachother at the same time
	EmptyTestDatabase()
	resetTestBoxes()
	resetTestItems()
	innerBox = BOX_1
	outerBox = BOX_2
	for _, testBox := range testBoxes() {
		dbTest.insertNewBox(testBox)
	}
	err = dbTest.MoveBoxToBox(innerBox.ID, outerBox.ID)
	err = dbTest.MoveBoxToBox(outerBox.ID, innerBox.ID)
	assert.NotEqual(t, err, nil)
	updatedInnerBox, err = dbTest.BoxById(innerBox.ID)
	updatedOuterBox, err = dbTest.BoxById(outerBox.ID)
	assert.NotEqual(t, updatedOuterBox.OuterBoxID, updatedInnerBox.ID)
}

func TestMoveBoxToShelf(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()
	resetShelves()
	dbTest.CreateBox(BOX_1)
	dbTest.CreateShelf(SHELF_1)
	fetchedBox, _ := dbTest.BoxById(BOX_1.ID)
	fetchedShelf, _ := dbTest.Shelf(SHELF_1.ID)
	assert.Equal(t, fetchedBox.ShelfID, uuid.Nil)
	assert.Equal(t, fetchedShelf.Boxes, nil)

	// Move in
	err := dbTest.MoveBoxToShelf(BOX_1.ID, SHELF_1.ID)
	assert.Equal(t, err, nil)
	fetchedBox, _ = dbTest.BoxById(BOX_1.ID)
	fetchedShelf, _ = dbTest.Shelf(SHELF_1.ID)
	assert.Equal(t, fetchedBox.ShelfID, SHELF_1.ID)
	assert.NotEqual(t, fetchedShelf.Boxes, nil)
	assert.Equal(t, fetchedShelf.Boxes[0].ID, BOX_1.ID)

	// Move out
	err = dbTest.MoveBoxToShelf(BOX_1.ID, uuid.Nil)
	assert.Equal(t, err, nil)
	fetchedBox, _ = dbTest.BoxById(BOX_1.ID)
	fetchedShelf, _ = dbTest.Shelf(SHELF_1.ID)
	assert.Equal(t, fetchedBox.ShelfID, uuid.Nil)
	assert.Equal(t, fetchedShelf.Boxes, nil)

	// Move non existent ID
	err = dbTest.MoveBoxToShelf(BOX_1.ID, VALID_UUID_NOT_EXISTING)
	assert.NotEqual(t, err, nil)
}

func TestMoveBoxToArea(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()
	resetShelves()
	dbTest.CreateBox(BOX_1)
	dbTest.CreateArea(*AREA_1)
	fetchedBox, _ := dbTest.BoxById(BOX_1.ID)
	// fetchedArea, _ := dbTest.AreaById(AREA_1.ID)
	assert.Equal(t, fetchedBox.AreaID, uuid.Nil)
	// assert.Equal(t, fetchedArea.Boxes, nil)

	// Move in
	err := dbTest.MoveBoxToArea(BOX_1.ID, AREA_1.ID)
	assert.Equal(t, err, nil)
	fetchedBox, _ = dbTest.BoxById(BOX_1.ID)
	// fetchedArea, _ = dbTest.Area(*AREA_1.ID)
	assert.Equal(t, fetchedBox.AreaID, AREA_1.ID)
	// assert.NotEqual(t, fetchedArea.Boxes, nil)
	// assert.Equal(t, fetchedArea.Boxes[0].ID, BOX_1.ID)

	// Move out
	err = dbTest.MoveBoxToArea(BOX_1.ID, uuid.Nil)
	assert.Equal(t, err, nil)
	fetchedBox, _ = dbTest.BoxById(BOX_1.ID)
	// fetchedArea, _ = dbTest.Area(*AREA_1.ID)
	assert.Equal(t, fetchedBox.AreaID, uuid.Nil)
	// assert.Equal(t, fetchedArea.Boxes, nil)

	// Move non existent ID
	err = dbTest.MoveBoxToArea(BOX_1.ID, VALID_UUID_NOT_EXISTING)
	assert.NotEqual(t, err, nil)
}
