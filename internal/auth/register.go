package auth

import (
	"basement/main/internal/database"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	REGISTER_FAILED_MESSAGE string = "register failed"
	COOKIE_NAME             string = "mycookie"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// this function will check the type of the request
// if it is from type post it will register the user otherwise it will generate the register template
func RegisterHandler(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			registerUser(w, r, db)
			return
		}
		if r.Method == http.MethodGet {
			generateRegisterPage(w, r)
			return
		}
		w.Header().Add("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		logg.Debug(w, "Method:'", r.Method, "' not allowed")
	}
}

func generateRegisterPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := Authenticated(r)
	data := templates.PageTemplate{
		Title:         "Register",
		Authenticated: authenticated,
		User:          Username(r),
	}

	if err := templates.Render(w, templates.TEMPLATE_REGISTER_PAGE, data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

func registerUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	passwordConfirm := r.PostFormValue("password-confirm")

	if username == "" {
		logg.Debug("Missing username form input")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, REGISTER_FAILED_MESSAGE)
		return
	}
	if password == "" {
		logg.Debug("Missing password form input")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, REGISTER_FAILED_MESSAGE)
		return
	}
	if passwordConfirm == "" {
		logg.Debug("Missing password-confirm form input")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, REGISTER_FAILED_MESSAGE)
		return
	}
	if password != passwordConfirm {
		logg.Debugf(`Mismatch between password: "%v" and password-confirm: "%v" form input`, password, passwordConfirm)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, REGISTER_FAILED_MESSAGE)
		return
	}

	ctx := context.TODO() // i don't now which kind of context we need to use so i keep it todo for now

	user, err := db.User(ctx, username)
	if (user.Username == username) || (err == nil) {
		logg.Debugf(`User already exists: "%v"`, user)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, REGISTER_FAILED_MESSAGE)
		return
	}

	newHashedPassword, err := hashPassword(password)
	if err != nil {
		logg.Fatal(err)
	}

	if err := db.CreateNewUser(ctx, username, newHashedPassword); err != nil {
		if err == sql.ErrNoRows {
			templates.Render(w, templates.TEMPLATE_REGISTER_PAGE, "")
		} else {
			logg.Debug(http.StatusText(http.StatusInternalServerError))
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		// https://htmx.org/headers/hx-location/
		http.RedirectHandler("/login-form", http.StatusOK)
		w.Header().Add("HX-Location", "/login")
		logg.Debugf("User %s registered successfully:", username)

		return
	}
}
