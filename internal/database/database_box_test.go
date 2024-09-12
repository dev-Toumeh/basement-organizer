package database

import (
	itemsPackage "basement/main/internal/items"
	"context"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gofrs/uuid/v5"
)

func TestInsertNewBox(t *testing.T) {
	EmptyTestDatabase()
	// Step 1: Get Test Data
	boxList, items := testData()
	boxToTest := boxList[0]

	// Step 2: Insert boxes
	for _, box := range boxList {
		_, err := dbTest.insertNewBox(box)
		if err != nil {
			t.Fatalf("insertNewBox failed: %v", err)
		}
	}

	// Step 3: Insert items
	for _, item := range *items {
		err := dbTest.insertNewItem(item)
		if err != nil {
			t.Fatalf("insertNewItem failed: %v", err)
		}
	}

	//	Step 4: Verify that the insertion of items was successful
	for _, item := range *items {
		_, err := dbTest.ItemByField("id", item.Id.String())
		if err != nil {
			t.Fatalf("get item error: %v", err)
		}
	}

	fetchedBox, err := dbTest.BoxByField("id", boxToTest.Id.String())
	if err != nil {
		t.Fatalf(" the function BoxByfield not working properly : %v", err)
	}

	//Compare the fetched box with the original test box
	assert.Equal(t, boxToTest.Label, fetchedBox.Label)
	assert.Equal(t, boxToTest.Description, fetchedBox.Description)

	duplicateBox := &itemsPackage.Box{
		Id:    boxToTest.Id,
		Label: "Duplicate Box",
	}

	_, err = dbTest.insertNewBox(duplicateBox)
	if err == nil {
		t.Errorf("Expected an error when inserting a box with an existing ID, got none")
	}
}

func TestBoxByField(t *testing.T) {
	defer EmptyTestDatabase()
	boxList, _ := testData()

	testBox := boxList[0] // Assuming you want to test the first box

	// Testing retrieval by a field that should exist
	fetchedBox, err := dbTest.BoxByField("id", testBox.Id.String())
	assert.Equal(t, err, nil)
	if err != nil {
		t.Fatalf("Failed to retrieve box by id: %v", err)
	}
	assert.Equal(t, fetchedBox.Id.String(), testBox.Id.String())

	// Testing retrieval by a non-existent field
	_, err = dbTest.BoxByField("non_existent_field", "some_value")
	assert.NotEqual(t, err, nil)

}

func TestCreateNewBox(t *testing.T) {
	defer EmptyTestDatabase()
	boxList, _ := testData()

	testBox := boxList[0]

	// Testing creation of a new box that does not already exist
	_, err := dbTest.CreateBox(testBox)
	assert.Equal(t, nil, err)
	if err != nil {
		t.Fatalf("Failed to create new box: %v", err)
	}

	// Verify box was created
	exists := dbTest.BoxExist("id", testBox.Id.String())
	assert.Equal(t, true, exists)

	// Test creating the same box again to trigger an error
	_, err = dbTest.CreateBox(testBox)
	assert.NotEqual(t, nil, err)

}

func TestBoxIDs(t *testing.T) {
	defer EmptyTestDatabase()
	// Prepare static test data with pre-defined UUIDs using uuid.Must
	testBox1Id := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000"))
	testBox2Id := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174001"))
	testBox3Id := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174002"))

	testBoxes := []*itemsPackage.Box{
		{Id: testBox1Id, Label: "Test Box 1", Description: "Test description 1", OuterBoxId: uuid.Nil},
		{Id: testBox2Id, Label: "Test Box 2", Description: "Test description 2", OuterBoxId: uuid.Nil},
		{Id: testBox3Id, Label: "Test Box 3", Description: "Test description 3", OuterBoxId: uuid.Nil},
	}

	expectedIDs := []string{testBox1Id.String(), testBox2Id.String(), testBox3Id.String()}

	// Insert test boxes into the database
	for _, testBox := range testBoxes {
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
	assert.Equal(t, expectedIDs, actualIDs)
}

func TestBoxUpdate(t *testing.T) {
	BoxId := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000"))

	BoxBeforeUpdate := &itemsPackage.Box{
		Id:          BoxId,
		Label:       "Test Box",
		Description: "Test description",
		Picture:     "my picture 1",
		QRcode:      "qrcode 1",
		OuterBoxId:  uuid.Nil,
	}

	BoxToUpdate := itemsPackage.Box{
		Id:          BoxId,
		Label:       "update Box",
		Description: "Update description",
		Picture:     "my picture 2",
		QRcode:      "qrcode 2",
		OuterBoxId:  uuid.Nil,
	}

	_, err := dbTest.insertNewBox(BoxBeforeUpdate)
	if err != nil {
		t.Fatalf("error while inserting the box: %v", err)
	}

	err = dbTest.UpdateBox(BoxToUpdate)
	if err != nil {
		t.Fatalf("error while updating the box: %v", err)
	}

	// Retrieve the updated box from the database
	updatedBox, err := dbTest.BoxByField("id", BoxId.String())
	if err != nil {
		t.Fatalf("error while retrieving the updated box: %v", err)
	}

	// Assert that the box was updated correctly (using individual asserts)
	assert.Equal(t, BoxToUpdate.Label, updatedBox.Label)
	assert.Equal(t, BoxToUpdate.Description, updatedBox.Description)
	assert.Equal(t, BoxToUpdate.Picture, updatedBox.Picture)
	assert.Equal(t, BoxToUpdate.QRcode, updatedBox.QRcode)
	assert.Equal(t, BoxToUpdate.OuterBoxId, BoxToUpdate.OuterBoxId)

	assert.NotEqual(t, BoxBeforeUpdate.Label, updatedBox.Label)
	assert.NotEqual(t, BoxBeforeUpdate.Description, updatedBox.Description)

	EmptyTestDatabase()
}

func TestDeleteBox(t *testing.T) {
	testBoxId := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000"))
	innerBoxId := uuid.Must(uuid.FromString("a0c201c2-5d5b-4587-938b-5a2b59c31e25"))
	itemId := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174002"))

	testBox := &itemsPackage.Box{
		Id:          testBoxId,
		Label:       "My Special Box",
		Description: "This box contains my precious items.",
		Picture:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==",
		QRcode:      "AB123CD",
		OuterBoxId:  uuid.Nil,
	}

	innerBox := &itemsPackage.Box{
		Id:          innerBoxId,
		Label:       "Inner Box 1",
		Description: "This is the first inner box",
		Picture:     "base64encodedinnerbox",
		QRcode:      "QRcodeInnerBox",
		OuterBoxId:  testBoxId,
	}

	item := itemsPackage.Item{
		Id:          itemId,
		Label:       "Item 1",
		Description: "Description for item 1",
		Picture:     "base64encodedstring1",
		Quantity:    10,
		Weight:      "5.5",
		QRcode:      "QRcode1",
		BoxId:       testBoxId,
	}

	boxList := []*itemsPackage.Box{innerBox, testBox}

	for _, box := range boxList {
		_, err := dbTest.insertNewBox(box)
		if err != nil {
			t.Fatalf("insertNewBox failed: %v", err)
		}
	}

	err := dbTest.insertNewItem(item)
	if err != nil {
		t.Fatalf("insertNewItem failed: %v", err)
	}

	err = dbTest.DeleteBox(testBox.Id)
	if err != nil && !strings.Contains(err.Error(), "the box is not empty") {
		t.Fatalf("the should not be deleted as the box is not empty")
	}
	err = dbTest.DeleteItem(item.Id)
	if err != nil {
		t.Fatalf("the item was not deleted")
	}
	err = dbTest.DeleteBox(innerBox.Id)
	if err != nil {
		t.Fatalf("deleting the innerbox was not succeed")
	}

	err = dbTest.DeleteBox(testBox.Id)
	if err != nil {
		t.Fatalf("delete the box after deleting the data inside of it was not succeed")
	}

	EmptyTestDatabase()
}

func TestMoveBox(t *testing.T) {
	// Prepare test data
	outerBox1Id := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000"))
	outerBox2Id := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174001"))
	innerBoxId := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174002"))

	testBoxes := []*itemsPackage.Box{
		{Id: outerBox1Id, Label: "Outer Box 1", Description: "This is the first outer box", OuterBoxId: uuid.Nil},
		{Id: outerBox2Id, Label: "Outer Box 2", Description: "This is the second outer box", OuterBoxId: uuid.Nil},
		{Id: innerBoxId, Label: "Inner Box", Description: "This is the inner box", OuterBoxId: outerBox1Id}, // Assign outerBox1Id by default
	}

	// Insert test boxes into the database using range
	for _, testBox := range testBoxes {
		_, err := dbTest.insertNewBox(testBox)
		if err != nil {
			t.Fatalf("Failed to insert test box: %v", err)
		}
	}

	// 1. Test successful move
	err := dbTest.MoveBox(innerBoxId, outerBox2Id)
	if err != nil {
		t.Fatalf("MoveBox function returned an error: %v", err)
	}

	updatedInnerBox, err := dbTest.BoxByField("id", innerBoxId.String())
	if err != nil {
		t.Fatalf("Failed to retrieve updated inner box: %v", err)
	}
	assert.Equal(t, outerBox2Id, updatedInnerBox.OuterBoxId)

	// 2. Test move to non-existent box (should return an error)
	nonExistentBoxId := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174003"))
	err = dbTest.MoveBox(innerBoxId, nonExistentBoxId)
	assert.Equal(t, err, err)

	EmptyTestDatabase()
}

// return data for testing Database
func testData() ([]*itemsPackage.Box, *[]itemsPackage.Item) {

	testBoxId := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174111"))
	outerBoxId := uuid.Must(uuid.FromString("18c60ba9-ffac-48f1-8c7c-473bd35acbea"))
	innerBoxId := uuid.Must(uuid.FromString("a0c201c2-5d5b-4587-938b-5a2b59c31e25"))
	innerBox2Id := uuid.Must(uuid.FromString("f47ac10b-58cc-4372-a567-0e02b2c3d479"))

	innerBox := &itemsPackage.Box{
		Id:          innerBoxId,
		Label:       "Inner Box 1",
		Description: "This is the first inner box",
		Picture:     "base64encodedinnerbox",
		QRcode:      "QRcodeInnerBox",
		OuterBoxId:  testBoxId,
	}

	innerBox2 := &itemsPackage.Box{
		Id:          innerBox2Id,
		Label:       "Inner Box 2",
		Description: "This is the second inner box",
		Picture:     "innerBox2Picture",
		QRcode:      "QR91011",
		OuterBoxId:  testBoxId,
	}

	outerBox := &itemsPackage.Box{
		Id:          outerBoxId,
		Label:       "OuterBox",
		Description: "This is the outer box",
		Picture:     "base64encodedouterbox",
		QRcode:      "QRcodeOuterBox",
		OuterBoxId:  uuid.Nil,
	}

	testBox := &itemsPackage.Box{
		Id:          testBoxId,
		Label:       "TestBox",
		Description: "This box contains my precious items.",
		Picture:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg==",
		QRcode:      "AB123CD",
		OuterBoxId:  outerBoxId,
		InnerBoxes:  []*itemsPackage.Box{innerBox, innerBox2},
		OuterBox:    outerBox,
	}

	item1 := itemsPackage.Item{
		Id:          uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000")),
		Label:       "Item 1",
		Description: "Description for item 1",
		Picture:     "base64encodedstring1",
		Quantity:    10,
		Weight:      "5.5",
		QRcode:      "QRcode1",
		BoxId:       testBoxId,
	}

	item2 := itemsPackage.Item{
		Id:          uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174001")),
		Label:       "Item 2",
		Description: "Description for item 2",
		Picture:     "base64encodedstring2",
		Quantity:    20,
		Weight:      "10.0",
		QRcode:      "QRcode2",
		BoxId:       testBoxId,
	}

	item3 := itemsPackage.Item{
		Id:          uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174002")),
		Label:       "Item 3",
		Description: "Description for item 3",
		Picture:     "base64encodedstring3",
		Quantity:    15,
		Weight:      "7.25",
		QRcode:      "QRcode3",
		BoxId:       testBoxId,
	}

	testBoxItemList := &[]itemsPackage.Item{item1, item2, item3}
	boxList := []*itemsPackage.Box{testBox, innerBox, innerBox2, outerBox}
	return boxList, testBoxItemList
}

// print the data that came from BoxByField()
// func priintData(t *testing.T) {
//
// 	boxList, _ := testData()
// 	fetchedBox, err := dbTest.BoxByField("id", boxList[0].Id.String())
// 	if err != nil {
// 		t.Fatalf("Failed to fetch inserted box: %v", err)
// 	}
//
// 	fmt.Print(" 1. Checking the items \n")
// 	for index, item := range fetchedBox.Items {
// 		fmt.Printf("item %d item %v \n", index, item)
// 	}
//
// 	fmt.Print(" 2. Checking the inner boxes \n")
// 	for index, item := range fetchedBox.InnerBoxes {
// 		fmt.Printf("item %d item %v \n", index, item)
// 	}
//
// 	fmt.Print(" 3. Checking the outer Box \n")
// 	fmt.Printf("the outerBox: %v \n", fetchedBox.OuterBox)
// }
