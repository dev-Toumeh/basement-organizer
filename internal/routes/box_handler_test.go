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

func (db *boxDatabaseSuccess) CreateBox() (string, error) {
	return "asdfkaj", nil
}

func (db *boxDatabaseSuccess) Box(id string) (items.Box, error) {
	return items.Box{}, nil
}

func TestBoxHandler(t *testing.T) {
	dbErrorCtx := context.WithValue(context.Background(), "db", &boxDatabaseError{})
	dbSuccessCtx := context.WithValue(context.Background(), "db", &boxDatabaseSuccess{})

	testCases := []struct {
		name               string
		input              handlerInput
		expectedStatusCode int
	}{
		{
			name: "Get Box fails",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodGet, "/box?id=dfkjasdlk", nil).WithContext(dbErrorCtx),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "CreateBox fails",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPost, "/box", nil).WithContext(dbErrorCtx),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "CreateBox succedes",
			input: handlerInput{
				R: httptest.NewRequest(http.MethodPost, "/box", nil).WithContext(dbSuccessCtx),
				W: *httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := BoxHandler(FprintWriteFunc)
			h(&tc.input.W, tc.input.R)
			assert.Equal(t, tc.expectedStatusCode, tc.input.W.Result().StatusCode)
		})
	}
}
