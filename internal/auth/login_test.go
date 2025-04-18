package auth

import (
	"basement/main/internal/logg"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gofrs/uuid/v5"
)

// TestAuthDatabase mock implementation
type TestAuthDatabase struct{}

// ErrorExist implements AuthDatabase.
func (db *TestAuthDatabase) ErrorExist() error {
	return errors.New("the Record is already exist")
}

// CreateNewUser mock implementation
func (db *TestAuthDatabase) CreateNewUser(ctx context.Context, username string, passwordhash string) error {
	return nil
}

// UserByField mock implementation
func (db *TestAuthDatabase) UserByField(ctx context.Context, field string, value string) (User, error) {
	user := User{
		Id:           uuid.Must(uuid.FromString("18c60ba9-ccac-48f5-8c7c-473bd35acbea")),
		Username:     "testuser1",
		PasswordHash: "$2a$14$Lw/lCPdEm2JrmCgSuEIUN.rxZZYlHQbMSNbM/7zOLu8k5jZZ4pwPK", // "abc"
	}
	return user, nil
}

// UpdateUser mock implementation
func (db *TestAuthDatabase) UpdateUser(ctx context.Context, user User) error {
	return nil
}

// UserExist mock implementation
func (db *TestAuthDatabase) UserExist(ctx context.Context, username string) bool {
	return true
}

type TestAuthDatabaseError struct{}

// ErrorExist implements AuthDatabase.
func (db *TestAuthDatabaseError) ErrorExist() error {
	return errors.New("the Record is already exist")
}

// CreateNewUser mock implementation returns error
func (db *TestAuthDatabaseError) CreateNewUser(ctx context.Context, username string, passwordhash string) error {
	return errors.New("")
}

// UserByField mock implementation returns error
func (db *TestAuthDatabaseError) UserByField(ctx context.Context, field string, value string) (User, error) {
	return User{}, errors.New("")
}

// UpdateUser mock implementation
func (db *TestAuthDatabaseError) UpdateUser(ctx context.Context, user User) error {
	return errors.New("")
}

func (db *TestAuthDatabaseError) UserExist(ctx context.Context, username string) bool {
	return false
}

var testDB *TestAuthDatabase = &TestAuthDatabase{}
var testErrorDB *TestAuthDatabaseError = &TestAuthDatabaseError{}

func TestLoginHandlerMethodNotAllowed(t *testing.T) {
	methodsNotAllowed := []string{http.MethodConnect, http.MethodTrace, http.MethodPut, http.MethodOptions, http.MethodDelete, http.MethodPatch, http.MethodHead}
	for _, method := range methodsNotAllowed {
		runLoginHandlerMethodNotAllowed(method, t)
	}
}

func TestLoginMissingUsername(t *testing.T) {
	// body := strings.NewReader("username=testuser1")
	loginWithMalformedInputs("", t)
	loginWithMalformedInputs("?username=", t)
}

func TestLoginMissingPassword(t *testing.T) {
	loginWithMalformedInputs("?username=testuser1", t)
	loginWithMalformedInputs("?username=testuser1&password", t)
	loginWithMalformedInputs("?username=testuser1&password=", t)
}

func TestLoginIncorrectPassword(t *testing.T) {
	loginWithMalformedInputs("?username=testuser1&password=a", t)
}

func TestLoginUserDoesNotExist(t *testing.T) {
	urlParams := "?username=userdoesnotexist&password=abc"
	r := httptest.NewRequest(http.MethodPost, "/login"+urlParams, nil)
	w := &httptest.ResponseRecorder{}

	loginUser(w, r, testErrorDB)
	if w.Code != http.StatusForbidden {
		t.Log(logg.WantHave(http.StatusForbidden, w.Result().Status, "/login"+urlParams))
		t.Fail()
	}
}

func TestLoginUserDoesNotMatch(t *testing.T) {
	// Disable error logging because this test will throw an expected errror that should not be logged.
	logg.DisableErrorLoggerS()
	defer logg.EnableErrorLoggerS()

	urlParams := "?username=nomatchinguser&password=abc"
	r := httptest.NewRequest(http.MethodPost, "/login"+urlParams, nil)
	w := &httptest.ResponseRecorder{}
	loginUser(w, r, testDB)
	if w.Code != http.StatusInternalServerError {
		t.Log(logg.WantHave(http.StatusInternalServerError, w.Result().Status, "/login"+urlParams))
		t.Fail()
	}
}

func TestLoginCorrectPassword(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/login?username=testuser1&password=abc", nil)
	w := &httptest.ResponseRecorder{}
	loginUser(w, r, testDB)
	if w.Code != http.StatusOK {
		t.Log(logg.WantHave(http.StatusOK, w.Result().Status, "/login?username=testuser1&password=abc"))
		t.Fail()
	}
}

// func TestLoginHandlerMethodAllowed(t *testing.T) {
// 	methods := []string{http.MethodGet, http.MethodPost}
// 	for _, method := range methods {
// 		runLoginHandlerMethodNotAllowed(method, t)
// 	}
// }

func loginWithMalformedInputs(urlParams string, t *testing.T) {
	// body := strings.NewReader("username=testuser1")
	r := httptest.NewRequest(http.MethodPost, "/login"+urlParams, nil)
	w := &httptest.ResponseRecorder{}
	loginUser(w, r, testDB)
	if w.Code != http.StatusForbidden {
		t.Log(logg.WantHave(http.StatusForbidden, w.Result().Status, "/login"+urlParams))
		t.Fail()
	}
}

func runLoginHandlerMethodNotAllowed(method string, t *testing.T) {
	r := httptest.NewRequest(method, "/login", nil)
	w := &httptest.ResponseRecorder{}
	loginFunc := LoginHandler(testDB)
	loginFunc(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Log(logg.WantHave(http.StatusMethodNotAllowed, w.Result().Status, "Method="+method))
		t.Fail()
	}
	allowHeader := w.Result().Header.Values("allow")
	getAllowed := slices.Contains(allowHeader, http.MethodGet)
	postAllowed := slices.Contains(allowHeader, http.MethodPost)
	logg.Debug(allowHeader)
	logg.Debug(getAllowed)
	logg.Debug(postAllowed)
	if !getAllowed {
		t.Log(logg.WantHave("GET", allowHeader))
		t.Fail()
	}
	if !postAllowed {
		t.Log(logg.WantHave("POST", allowHeader))
		t.Fail()
	}
}

// func runLoginHandlerMethodAllowed(method string, t *testing.T) {
// 	r := httptest.NewRequest(method, "/login", nil)
// 	w := &httptest.ResponseRecorder{}
// 	loginFunc := LoginHandler(testDB)
// 	loginFunc(w, r)
// 	if w.Code != http.StatusMethodNotAllowed {
// 		t.Log("Method:", method, "=> Status:", w.Result().Status)
// 		t.Fail()
// 	}
// }
