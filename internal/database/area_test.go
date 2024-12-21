package database

import (
	"slices"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/gofrs/uuid/v5"
)

func TestCreateNewArea(t *testing.T) {
	EmptyTestDatabase()
	resetAreas()

	testArea := AREA_1

	// Testing creation of a new area that does not already exist
	_, err := dbTest.CreateArea(*testArea)
	assert.Equal(t, nil, err)
	if err != nil {
		t.Fatalf("Failed to create new area: %v", err)
	}

	// Verify area was created
	exists := dbTest.AreaExists(testArea.ID)
	assert.Equal(t, true, exists)

	// Test creating the same area again to trigger an error
	_, err = dbTest.CreateArea(*testArea)
	assert.NotEqual(t, nil, err)
}

func TestInsertNewArea(t *testing.T) {
	EmptyTestDatabase()
	resetAreas()

	testArea := AREA_1

	// Step 2: Insert boxes
	for _, area := range testAreas() {
		_, err := dbTest.insertNewArea(area)
		if err != nil {
			t.Fatalf("insertNewArea failed: %v", err)
		}
	}

	fetchedArea, err := dbTest.AreaById(testArea.ID)
	if err != nil {
		t.Fatalf(" the function AreaByfield not working properly : %v %v", err.Error(), testArea)
	}

	//Compare the fetched area with the original test area
	assert.Equal(t, testArea.Label, fetchedArea.Label)
	assert.Equal(t, testArea.Description, fetchedArea.Description)
	assert.NotEqual(t, "", fetchedArea.PreviewPicture)

	duplicateArea := *AREA_1

	_, err = dbTest.insertNewArea(duplicateArea)
	if err == nil {
		t.Errorf("Expected an error when inserting a area with an existing ID, got none")
	}
}

func TestAreaByField(t *testing.T) {
	EmptyTestDatabase()
	resetAreas()

	testArea := AREA_1
	dbTest.insertNewArea(*testArea)

	// Testing retrieval by a field that should exist
	fetchedArea, err := dbTest.AreaById(testArea.ID)
	assert.Equal(t, err, nil)
	if err != nil {
		t.Fatalf("Failed to retrieve area by id: %v", err)
	}
	assert.Equal(t, fetchedArea.ID.String(), testArea.ID.String())

	// Testing retrieval by a non-existent field
	_, err = dbTest.areaByField("non_existent_field", "some_value")
	assert.NotEqual(t, err, nil)
}

func TestAreaIDs(t *testing.T) {
	EmptyTestDatabase()
	resetAreas()

	expectedIDs := []uuid.UUID{AREA_1.ID, AREA_2.ID, AREA_3.ID, AREA_4.ID}

	// Insert test boxes into the database
	for _, testArea := range testAreas() {
		_, err := dbTest.insertNewArea(testArea)
		if err != nil {
			t.Fatalf("Failed to insert test area: %v", err)
		}
	}

	// Call the AreaIDs function
	actualIDs, err := dbTest.AreaIDs()
	if err != nil {
		t.Fatalf("AreaIDs function returned an error: %v", err)
	}

	// Verify the results
	for _, v := range expectedIDs {
		assert.Equal(t, slices.Contains(actualIDs, v), true)
	}
}

func TestUpdateArea(t *testing.T) {
	EmptyTestDatabase()
	resetAreas()

	testArea := AREA_1
	_, err := dbTest.insertNewArea(*testArea)
	if err != nil {
		t.Fatalf("error while inserting the area: %v", err)
	}

	oldLabel := testArea.Label
	oldDescr := testArea.Description

	testArea.Description = "updated"
	testArea.Label = "updated"
	assert.NotEqual(t, oldDescr, testArea.Description)
	assert.NotEqual(t, oldLabel, testArea.Label)

	err = dbTest.UpdateArea(*testArea)
	assert.Equal(t, err, nil)

	// Retrieve the updated area from the database
	updatedArea, err := dbTest.AreaById(testArea.ID)
	assert.Equal(t, err, nil)

	// Assert that the area was updated correctly (using individual asserts)
	assert.Equal(t, testArea.Label, updatedArea.Label)
	assert.Equal(t, testArea.Description, updatedArea.Description)
	assert.Equal(t, testArea.Picture, updatedArea.Picture)
	assert.Equal(t, testArea.QRCode, updatedArea.QRCode)

	assert.NotEqual(t, oldLabel, updatedArea.Label)
	assert.NotEqual(t, oldDescr, updatedArea.Description)

	EmptyTestDatabase()
}

func TestDeleteArea(t *testing.T) {
	EmptyTestDatabase()
	resetAreas()

	for _, area := range testAreas() {
		_, err := dbTest.insertNewArea(area)
		if err != nil {
			t.Fatalf("insertNewArea failed: %v", err)
		}
	}
	var err error

	err = dbTest.DeleteArea(AREA_1.ID)
	assert.NotEqual(t, err, nil)
}

func TestAreaListRowByID(t *testing.T) {
	EmptyTestDatabase()
	resetAreas()

	for _, area := range testAreas() {
		_, err := dbTest.insertNewArea(area)
		if err != nil {
			t.Fatalf("insertNewArea failed: %v", err)
		}
	}
	var err error

	area_1, err := dbTest.AreaListRowByID(AREA_1.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, area_1.ID, AREA_1.ID)
	assert.Equal(t, area_1.Label, AREA_1.Label)
	assert.Equal(t, area_1.Description, AREA_1.Description)

	area_3, err := dbTest.AreaListRowByID(AREA_3.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, area_3.ID, AREA_3.ID)
	assert.Equal(t, area_3.Label, AREA_3.Label)
	assert.Equal(t, area_3.Description, AREA_3.Description)
}

func TestAreaListCounter(t *testing.T) {
	EmptyTestDatabase()
	resetAreas()

	for _, area := range testAreas() {
		_, err := dbTest.insertNewArea(area)
		if err != nil {
			t.Fatalf("insertNewArea failed: %v", err)
		}
	}
	var err error

	count, err := dbTest.AreaListCounter("")
	assert.Equal(t, err, nil)
	assert.Equal(t, count, 6)

	count, err = dbTest.AreaListCounter("Area")
	assert.Equal(t, err, nil)
	assert.Equal(t, count, 4)

	count, err = dbTest.AreaListCounter("A")
	assert.Equal(t, err, nil)
	assert.Equal(t, count, 5)

	count, err = dbTest.AreaListCounter("B")
	assert.Equal(t, err, nil)
	assert.Equal(t, count, 1)

	count, err = dbTest.AreaListCounter("Test")
	assert.Equal(t, err, nil)
	assert.Equal(t, count, 3)

	// count, err = dbTest.AreaListCounter("Area A")
	// assert.Equal(t, err, nil)
	// assert.Equal(t, count, 1)

	count, err = dbTest.AreaListCounter("")
	assert.Equal(t, err, nil)
	assert.Equal(t, count, 6)

}
