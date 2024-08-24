package routes

import (
	"basement/main/internal/items"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type handlerInput struct {
	R *http.Request
	W httptest.ResponseRecorder
}

// boxDatabaseError returns errors on every function.
type boxDatabaseError struct{}

func (db *boxDatabaseError) CreateBox() (string, error) {
	return "", errors.New("AAAAAAAA")
}

func (db *boxDatabaseError) Box(id string) (items.Box, error) {
	return items.Box{}, errors.New("AAAAAAAA")
}

// boxDatabaseSuccess never returns errors.
type boxDatabaseSuccess struct{}

const BOX_ID = "fa2e3db6-fcf8-49c6-ac9c-54ce5855bf0b"

func (db *boxDatabaseSuccess) CreateBox() (string, error) {
	return BOX_ID, nil
}

func (db *boxDatabaseSuccess) Box(id string) (items.Box, error) {
	return items.Box{}, nil
}

func TestBoxHandlerDBErrors(t *testing.T) {
	// logg.EnableDebugLoggerS()

	dbErr := &boxDatabaseError{}

	// Add mux handler, without it r.PathValue("id") will not work.
	mux := http.NewServeMux()
	mux.Handle("/box", BoxHandler(WriteFprint, dbErr))
	mux.Handle("/box/", BoxHandler(WriteFprint, dbErr))
	mux.Handle("/api/v2/box/{id}", BoxHandler(WriteFprint, dbErr))
	mux.Handle("/api/v2/box/", BoxHandler(WriteFprint, dbErr))

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
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Create box DB error",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPost, "/box", nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		// @TODO
		// {
		// 	name: "Can't delete box, not found",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodDelete, "/box", nil),
		// 		W: *httptest.NewRecorder(),
		// 	},
		// 	expectedStatusCode: http.StatusNotFound,
		// },
		// {
		// 	name: "Can't update box, not found",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodPut, "/box", nil),
		// 		W: *httptest.NewRecorder(),
		// 	},
		// 	expectedStatusCode: http.StatusNotFound,
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mux.ServeHTTP(&tc.input.W, tc.input.R)
			assert.Equal(t, tc.expectedStatusCode, tc.input.W.Result().StatusCode, "URL: "+tc.input.R.URL.String())
		})
	}
}

func TestBoxHandlerInputErrors(t *testing.T) {
	// logg.EnableDebugLoggerS()

	dbOk := boxDatabaseSuccess{}

	// Add mux handler, without it r.PathValue("id") will not work.
	mux := http.NewServeMux()
	mux.Handle("/box", BoxHandler(WriteFprint, &dbOk))
	mux.Handle("/box/", BoxHandler(WriteFprint, &dbOk))
	mux.Handle("/api/v2/box/{id}", BoxHandler(WriteFprint, &dbOk))
	mux.Handle("/api/v2/box/", BoxHandler(WriteFprint, &dbOk))

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
				R: httptest.NewRequest(http.MethodGet, "/api/v2/box/333", nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Box not found, empty path id",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/api/v2/box/", nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Box not found, empty query param id",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/box?id=", nil),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
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
	mux.Handle("/box", BoxHandler(WriteFprint, &dbOk))
	mux.Handle("/box/", BoxHandler(WriteFprint, &dbOk))
	mux.Handle("/api/v2/box/{id}", BoxHandler(WriteFprint, &dbOk))
	mux.Handle("/api/v2/box/", BoxHandler(WriteFprint, &dbOk))

	testCases := []struct {
		name               string
		input              handlerInput
		expectedStatusCode int
	}{
		// @TODO
		// {
		// 	name: "Create box ok",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodPost, "/box", nil),
		// 		W: *httptest.NewRecorder(),
		// 	},
		// 	expectedStatusCode: http.StatusOK,
		// },
		// @TODO
		// {
		// 	name: "Should use query param value /box?id={id}",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodGet, "/box?id="+BOX_ID, nil),
		// 		W: *httptest.NewRecorder(),
		// 	},
		// 	expectedStatusCode: http.StatusOK,
		// },
		// {
		// 	name: "Should use path value id /box/{id}",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodGet, "/api/v2/box/"+BOX_ID, nil),
		// 		W: *httptest.NewRecorder(),
		// 	},
		// 	expectedStatusCode: http.StatusOK,
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mux.ServeHTTP(&tc.input.W, tc.input.R)
			assert.Equal(t, tc.expectedStatusCode, tc.input.W.Result().StatusCode, "URL: "+tc.input.R.URL.String())
		})
	}
}
