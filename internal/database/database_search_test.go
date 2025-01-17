package database

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestVirtualBoxInsert(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testbox := BOX_1

	// Create the outerbox
	_, err := dbTest.CreateBox(testbox)
	if err != nil {
		t.Fatalf("Failed to create outer box: %v", err)
	}

	// Check if the outerbox exists in box_fts
	exist, err := dbTest.VirtualBoxExist(testbox.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, exist, true)

	boxListRow, err := dbTest.BoxListRowByID(testbox.ID)
	if err != nil {
		t.Fatalf("Failed to create outer box: %v", err)
	}
	assert.Equal(t, testbox.ID, boxListRow.ID)
	assert.Equal(t, testbox.Label, boxListRow.Label)
	assert.Equal(t, testbox.OuterBoxID, boxListRow.BoxID)
}

func TestVirtualBoxUpdate(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testbox := BOX_1
	outerBox := BOX_2
	testbox.OuterBoxID = BOX_2.ID

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

	testbox.Label = "new testbox label"
	dbTest.UpdateBox(*testbox, true)

	outerBox.Label = "new outerbox label"
	dbTest.UpdateBox(*outerBox, true)

	// Get the box_fts to check if the outerbox_label  was updated
	afterUpdate, err := dbTest.BoxListRowByID(testbox.ID)
	if err != nil {
		t.Fatalf("Failed to fetch the testbox while checking the BoxTriger: %v", err)
	}

	assert.Equal(t, afterUpdate.BoxLabel, outerBox.Label)
	assert.Equal(t, afterUpdate.Label, testbox.Label)
}

func TestVirtualBoxUpdateIgnorePicture(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testbox := BOX_1
	outerBox := BOX_2
	testbox.OuterBoxID = BOX_2.ID

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

	beforeUpdate, err := dbTest.BoxListRowByID(testbox.ID)
	assert.NotEqual(t, beforeUpdate.PreviewPicture, "")

	testbox.Label = "new testbox label"
	testbox.Picture = ""
	dbTest.UpdateBox(*testbox, true)

	outerBox.Label = "new outerbox label"
	outerBox.Picture = ""
	dbTest.UpdateBox(*outerBox, true)

	// Get the box_fts to check if the outerbox_label  was updated
	afterUpdate, err := dbTest.BoxListRowByID(testbox.ID)
	if err != nil {
		t.Fatalf("Failed to fetch the testbox while checking the BoxTriger: %v", err)
	}

	assert.Equal(t, afterUpdate.BoxLabel, outerBox.Label)
	assert.Equal(t, afterUpdate.Label, testbox.Label)
	assert.Equal(t, afterUpdate.PreviewPicture, beforeUpdate.PreviewPicture)
}

func TestVirtualBoxDelete(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()

	testbox := BOX_1

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
	exist, err := dbTest.VirtualBoxExist(testbox.ID)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, exist, true)
}

func TestBoxListRows(t *testing.T) {
	EmptyTestDatabase()
	resetTestBoxes()
	// Insert the new boxes
	for _, box := range testBoxes() {
		_, err := dbTest.insertNewBox(box)
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
			virtualBoxes, err := dbTest.BoxListRows(tc.query, 10, 1)
			if err != nil {
				t.Fatalf("error occurred while testing boxFuzzyFinder(): %v", err)
			}

			assert.Equal(t, tc.expected, len(virtualBoxes))

		})
	}
}
