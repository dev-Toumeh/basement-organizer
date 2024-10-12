package routes

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

type handlerInput struct {
	R *http.Request
	W httptest.ResponseRecorder
}

const BOX_ID_VALID string = "fa2e3db6-fcf8-49c6-ac9c-54ce5855bf0b"
const BOX_ID_NOT_FOUND string = "da2e3db6-fcf8-49c6-ac9c-54ce5855bf0b"
const BOX_ID_INVALID_EMPTY string = ""
const BOX_ID_INVALID_2 string = "aaaa"
const BOX_ID_INVALID_UUID_FORMAT string = "ac9c-54ce5855bf0b"

// boxDatabaseError returns errors on every function.
type boxDatabaseError struct{}

func (db *boxDatabaseError) CreateBox(newBox *items.Box) (uuid.UUID, error) {
	return uuid.Nil, errors.New("AAAAAAAA")
}

func (db *boxDatabaseError) BoxById(id uuid.UUID) (items.Box, error) {
	return items.Box{BasicInfo: items.BasicInfo{ID: uuid.Nil}}, errors.New("AAAAAAAA")
}

func (db *boxDatabaseError) BoxIDs() ([]string, error) {
	return nil, errors.New("AAAAAAAA")
}

func (db *boxDatabaseError) MoveBoxToBox(id1 uuid.UUID, id2 uuid.UUID) error {
	return errors.New("AAAAAAAA")
}

func (db *boxDatabaseError) BoxByField(field string, value string) (*items.Box, error) {
	return &items.Box{}, errors.New("AAAAAAAA")
}

func (db *boxDatabaseError) BoxExistById(id uuid.UUID) bool {
	return false
}

func (db *boxDatabaseError) ErrorExist() error {
	return errors.New("AAAAAAAA")
}

func (db *boxDatabaseError) UpdateBox(box items.Box) error {
	return errors.New("AAAAA")
}

func (db *boxDatabaseError) DeleteBox(boxId uuid.UUID) error {
	return errors.New("AAAAA")
}

func (db *boxDatabaseError) BoxFuzzyFinder(query string, limit int, page int) ([]items.ListRow, error) {
	return make([]items.ListRow, 0), errors.New("AAAAAAAA")
}

func (db *boxDatabaseError) BoxListRowByID(id uuid.UUID) (items.ListRow, error) {
	return items.ListRow{}, errors.New("AAAAAAAA")
}

// boxDatabaseSuccess never returns errors.
type boxDatabaseSuccess struct{}

func (db *boxDatabaseSuccess) CreateBox(newBox *items.Box) (uuid.UUID, error) {
	return uuid.Must(uuid.FromString(BOX_ID_VALID)), nil
}

func (db *boxDatabaseSuccess) BoxById(id uuid.UUID) (items.Box, error) {
	return items.Box{BasicInfo: items.BasicInfo{ID: uuid.Must(uuid.FromString(BOX_ID_VALID))}}, nil
}

func (db *boxDatabaseSuccess) BoxIDs() ([]string, error) {
	return []string{"id1", "id2", "id3"}, nil
}

func (db *boxDatabaseSuccess) MoveBoxToBox(id1 uuid.UUID, id2 uuid.UUID) error {
	return nil
}

func (db *boxDatabaseSuccess) BoxExistById(id uuid.UUID) bool {
	return true
}

func (db *boxDatabaseSuccess) ErrorExist() error {
	return nil
}

func (db *boxDatabaseSuccess) UpdateBox(box items.Box) error {
	return nil
}

func (db *boxDatabaseSuccess) DeleteBox(boxId uuid.UUID) error {
	return nil
}

func (db *boxDatabaseSuccess) BoxFuzzyFinder(query string, limit int, page int) ([]items.ListRow, error) {
	return make([]items.ListRow, 0), nil
}

func (db *boxDatabaseSuccess) BoxListRowByID(id uuid.UUID) (items.ListRow, error) {
	return items.ListRow{}, nil
}

func TestBoxHandlerDBErrors(t *testing.T) {
	// logg.EnableDebugLoggerS()
	// defer logg.DisableDebugLoggerS()

	dbErr := &boxDatabaseError{}

	// Add mux handler, without it r.PathValue("id") will not work.
	mux := http.NewServeMux()
	mux.Handle("/box", boxHandler(dbErr))
	mux.Handle("/box/", boxHandler(dbErr))
	mux.Handle("/api/v1/box/{id}", boxHandler(dbErr))
	mux.Handle("/api/v1/box/", boxHandler(dbErr))

	testCases := []struct {
		name               string
		input              handlerInput
		expectedStatusCode int
	}{
		{
			name: "Create box DB error",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPost, "/box", nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Can't delete box, not found",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodDelete, "/box?id="+BOX_ID_NOT_FOUND, nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Can't update box, not found",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPut, "/box?id="+BOX_ID_NOT_FOUND, nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mux.ServeHTTP(&tc.input.W, tc.input.R)
			assert.Equal(t, tc.expectedStatusCode, tc.input.W.Result().StatusCode, "URL: "+tc.input.R.URL.String())
		})
	}
}

type methodTestCase struct {
	name               string
	input              handlerInput
	expectedStatusCode int
}

func TestBoxHandlerInputErrors(t *testing.T) {
	// logg.EnableDebugLoggerS()

	dbOk := boxDatabaseSuccess{}

	// Add mux handler, without it r.PathValue("id") will not work.
	mux := http.NewServeMux()
	mux.Handle("/box", boxHandler(&dbOk))
	mux.Handle("/box/", boxHandler(&dbOk))
	mux.Handle("/api/v1/box/{id}", boxHandler(&dbOk))
	mux.Handle("/api/v1/box/", boxHandler(&dbOk))

	testCases := []struct {
		name               string
		input              handlerInput
		expectedStatusCode int
	}{
		{
			name: "Patch not allowed",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPatch, "/box", nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name: "Box not found, incorrect id path value /box/{id}",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/api/v1/box/"+BOX_ID_INVALID_UUID_FORMAT, nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Box not found, empty path id",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/api/v1/box/"+BOX_ID_INVALID_EMPTY, nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Box not found, empty query param id",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/box?id="+BOX_ID_INVALID_EMPTY, nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
		// DELETE
		{
			name: "Can't delete box, invalid UUID format",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodDelete, "/box?id="+BOX_ID_INVALID_UUID_FORMAT, nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mux.ServeHTTP(&tc.input.W, tc.input.R)
			assert.Equal(t, tc.expectedStatusCode, tc.input.W.Result().StatusCode, "URL: "+tc.input.R.URL.String())
		})
	}
}

func TestBoxHandlerOK(t *testing.T) {
	// logg.EnableDebugLoggerS()

	dbOk := boxDatabaseSuccess{}

	// Add mux handler, without it r.PathValue("id") will not work.
	mux := http.NewServeMux()
	mux.Handle("/box", boxHandler(&dbOk))
	mux.Handle("/box/", boxHandler(&dbOk))
	mux.Handle("/api/v1/box/{id}", boxHandler(&dbOk))
	mux.Handle("/api/v1/box", boxHandler(&dbOk))

	testCases := []struct {
		name               string
		input              handlerInput
		expectedStatusCode int
		expectedTemplate   bool
	}{
		{
			name: "Create box ok template response",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPost, "/box", nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			expectedTemplate:   true,
		},
		{
			name: "Create box ok data only response",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPost, "/api/v1/box", nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			expectedTemplate:   false,
		},
		{
			name: "Should use query param value /box?id={id}",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/box?id="+BOX_ID_VALID, nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			expectedTemplate:   true,
		},
		{
			name: "Should use path value id /box/{id}",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/api/v1/box/"+BOX_ID_VALID, nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			expectedTemplate:   false,
		},
	}
	err := templates.InitTemplates("../templates")
	if err != nil {
		logg.Fatal(err)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mux.ServeHTTP(&tc.input.W, tc.input.R)
			url := "URL: " + tc.input.R.URL.String()
			assert.Equal(t, tc.expectedStatusCode, tc.input.W.Result().StatusCode, url)

			read, _ := io.ReadAll(tc.input.W.Result().Body)
			if tc.expectedTemplate {
				assert.Contains(t, string(read), "hx-", url)
			} else {
				assert.NotContains(t, string(read), "hx-", url)
			}
		})
	}
}
