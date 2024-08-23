package routes

import (
	"basement/main/internal/items"
	"context"
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

func TestBoxHandler(t *testing.T) {
	// logg.EnableDebugLoggerS()
	dbErrorCtx := context.WithValue(context.Background(), "db", &boxDatabaseError{})
	dbSuccessCtx := context.WithValue(context.Background(), "db", &boxDatabaseSuccess{})

	// Add mux handler, without it r.PathValue("id") will not work.
	mux := http.NewServeMux()
	mux.Handle("/box", BoxHandler(FprintWriteFunc))
	mux.Handle("/box/", BoxHandler(FprintWriteFunc))
	mux.Handle("/api/v2/box/{id}", BoxHandler(FprintWriteFunc))
	mux.Handle("/api/v2/box/", BoxHandler(FprintWriteFunc))

	testCases := []struct {
		name               string
		input              handlerInput
		expectedStatusCode int
	}{
		{
			name: "Box not found, incorrect id format",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/box?id=dfkjasdlk", nil).WithContext(dbErrorCtx),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Create box DB error",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPost, "/box", nil).WithContext(dbErrorCtx),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		// @TODO
		// {
		// 	name: "Create box ok",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodPost, "/box", nil).WithContext(dbSuccessCtx),
		// 		W: *httptest.NewRecorder(),
		// 	},
		// 	expectedStatusCode: http.StatusOK,
		// },
		{
			name: "Patch not allowed",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPatch, "/box", nil).WithContext(dbSuccessCtx),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		// @TODO
		// {
		// 	name: "Can't delete box, not found",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodDelete, "/box", nil).WithContext(dbErrorCtx),
		// 		W: *httptest.NewRecorder(),
		// 	},
		// 	expectedStatusCode: http.StatusNotFound,
		// },
		// {
		// 	name: "Can't update box, not found",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodPut, "/box", nil).WithContext(dbErrorCtx),
		// 		W: *httptest.NewRecorder(),
		// 	},
		// 	expectedStatusCode: http.StatusNotFound,
		// },
		{
			name: "Box not found, incorrect id path value /box/{id}",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/api/v2/box/333", nil).WithContext(dbSuccessCtx),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Box not found, empty path id",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/api/v2/box/", nil).WithContext(dbSuccessCtx),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Box not found, empty query param id",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/box?id=", nil).WithContext(dbSuccessCtx),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		// @TODO
		// {
		// 	name: "Should use query param value /box?id={id}",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodGet, "/box?id="+BOX_ID, nil).WithContext(dbSuccessCtx),
		// 		W: *httptest.NewRecorder(),
		// 	},
		// 	expectedStatusCode: http.StatusOK,
		// },
		// {
		// 	name: "Should use path value id /box/{id}",
		// 	input: handlerInput{
		// 		R: httptest.NewRequest(http.MethodGet, "/api/v2/box/"+BOX_ID, nil).WithContext(dbSuccessCtx),
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
