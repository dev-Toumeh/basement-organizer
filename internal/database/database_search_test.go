package database

import (
	"basement/main/internal/items"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gofrs/uuid/v5"
)

var uuidValue1 uuid.UUID = uuid.Must(uuid.FromString("623e4567-e89b-12d3-a456-426614174000"))
var boxId1 = &uuidValue1

var uuidValue2 uuid.UUID = uuid.Must(uuid.FromString("323e4567-e89b-12d3-a456-426614174000"))
var boxId2 = &uuidValue2

var uuidValue3 uuid.UUID = uuid.Must(uuid.FromString("423e4567-e89b-12d3-a456-426614174000"))
var boxId3 = &uuidValue3

var uuidValue4 uuid.UUID = uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174111"))
var boxId4 = &uuidValue4

// Clone 4
var box1 = &items.Box{
	ID:          *boxId1,
	Label:       "box 1",
	Description: "This is the sixth box",
	Picture:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HwAFAgL/uXBuZwAAAABJRU5ErkJggg==",
	QRcode:      "uvwxyzabcdefg",
	OuterBoxID:  uuid.Nil,
}

var box3 = &items.Box{
	ID:          *boxId2,
	Label:       "box 3",
	Description: "This is the third box",
	Picture:     box1.Picture,
	QRcode:      "abababababcd",
	OuterBoxID:  *boxId1,
}

var box4 = &items.Box{
	ID:          *boxId3,
	Label:       "box 4",
	Description: "This is the fourth box",
	Picture:     box1.Picture,
	QRcode:      "efghefghefgh",
	OuterBoxID:  *boxId1,
}

var box5 = &items.Box{
	ID:          *boxId4,
	Label:       "box 5",
	Description: "This is the fifth box",
	Picture:     box1.Picture,
	QRcode:      "ijklmnopqrst",
	OuterBoxID:  *boxId1,
}

func TestVirtualBoxInsirt(t *testing.T) {
	// Setup
	boxList, _ := testData()
	testBox := boxList[0]

	// Create the outerbox
	_, err := dbTest.CreateBox(testBox)
	if err != nil {
		t.Fatalf("Failed to create outer box: %v", err)
	}

	// Check if the outerbox exists in box_fts
	exist := dbTest.VirtualBoxExist(testBox.ID)
	assert.Equal(t, exist, true)

	virtualBox, err := dbTest.VirtualBoxById(testBox.ID)
	if err != nil {
		t.Fatalf("Failed to create outer box: %v", err)
	}
	assert.Equal(t, testBox.ID, virtualBox.BoxID)
	assert.Equal(t, testBox.Label, virtualBox.Label)
	assert.Equal(t, testBox.OuterBoxID, virtualBox.OuterBoxID)

	EmptyTestDatabase()

}

func TestVirtualBoxUpdate(t *testing.T) {
	defer EmptyTestDatabase()

	boxList, _ := testData()
	testbox := boxList[0]
	outerBox := boxList[3]

	// Create the outerbox
	_, err := dbTest.CreateBox(outerBox)
	if err != nil {
		t.Fatalf("Failed to create outer box: %v", err)
	}

	// Create the testBox
	_, err = dbTest.CreateBox(testbox)
	if err != nil {
		t.Fatalf("Failed to create outer box while checking the BoxTriger: %v", err)
	}

	//change the testBoxLabel
	testboxClone := *testbox
	testbox.Label = "newTestBoxLable"
	dbTest.UpdateBox(testboxClone)

	// change the Outerbox label
	outerBoxClone := *outerBox
	outerBoxClone.Label = "newLable"
	dbTest.UpdateBox(outerBoxClone)

	// Get the box_fts to check if the outerbox_label  was updated
	afterUpdate, err := dbTest.VirtualBoxById(testbox.ID)
	if err != nil {
		t.Fatalf("Failed to fetch the testbox while checking the BoxTriger: %v", err)
	}

	assert.Equal(t, afterUpdate.OuterBoxLabel, outerBoxClone.Label)
	assert.Equal(t, afterUpdate.Label, testboxClone.Label)
}

func TestVirtualBoxDelete(t *testing.T) {
	defer EmptyTestDatabase()

	boxList, _ := testData()
	testbox := boxList[0]

	// Create the testBox
	_, err := dbTest.CreateBox(testbox)
	if err != nil {
		t.Fatalf("Failed to create outer box while checking the BoxTriger: %v", err)
	}

	dbTest.DeleteBox(testbox.ID)
	if err != nil {
		t.Fatalf("error while Deleting the the box: %v", err)
	}

	// Check if the outerbox exists in box_fts
	exist := dbTest.VirtualBoxExist(testbox.ID)
	assert.NotEqual(t, exist, true)
}

func TestBoxFuzzyFinder(t *testing.T) {

	boxesToInsert := []items.Box{*box1, *box3, *box4, *box5}

	// Insert the new boxes
	for _, box := range boxesToInsert {
		_, err := dbTest.insertNewBox(&box)
		if err != nil {
			t.Fatalf("insertNewBox failed while testing the boxFuzzyFinder: %v", err)
		}
	}

	testCases := []struct {
		name     string
		query    string
		expected int
	}{
		{"Basic Search", "b", 4},
		{"Partial Match", "bo", 4},
		{"Specific Box", "box 3", 1},
		{"Case Insensitivity", "BOX 1", 1},
		{"No Matches", "nonexistent", 0},
		{"Empty Query", "", 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			virtualBoxes, err := dbTest.BoxFuzzyFinder(tc.query, 10, 1)
			if err != nil {
				t.Fatalf("error occurred while testing boxFuzzyFinder(): %v", err)
			}

			assert.Equal(t, tc.expected, len(virtualBoxes))

		})
	}
}
