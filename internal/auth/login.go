package auth

import (
	"context"
	"fmt"
	"net/http"

	"basement/main/internal/database"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
)

const (
	LOGIN_FAILED_MESSAGE string = "Login failed"
)

func LoginHandler(db *database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			loginUser(w, r, db)
		}
		if r.Method == http.MethodGet {
			loginPage(w, r)
		}
	}
}
func loginUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authenticated, ok := Authenticated(r)

	if ok {
		if authenticated {
			logg.Debugf("LoginHandler - ok: %v authenticated: %v", ok, authenticated)
			fmt.Fprint(w, "already logged in")
			return
		}
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	logg.Debugf("%v %v ", r.URL, r.Form.Encode())

	if username == "" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		logg.Info(LOGIN_FAILED_MESSAGE)
		logg.Debug("Missing username")
		return
	}
	if password == "" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		logg.Info(LOGIN_FAILED_MESSAGE)
		logg.Debug("Missing password")
		return
	}

	ctx := context.TODO()
	user, _ := db.User(ctx, username)

	if !checkPasswordHash(password, user.PasswordHash) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		logg.Info(LOGIN_FAILED_MESSAGE)
		logg.Debug("pw hash doesnt match")
		return
	}

	saveSession(w, r)

	logg.Info("login successful")

	// https://htmx.org/headers/hx-location/
	w.Header().Add("HX-Location", "/personal-page")
	fmt.Fprintf(w, "Welcome %v\n", username)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := Authenticated(r)
	data := templates.NewPageTemplate()
	data.Title = "login"
	data.Authenticated = authenticated

	err := templates.Render(w, templates.TEMPLATE_LOGIN_PAGE, data)
	if err != nil {
		templates.RenderErrorSnackbar(w, err.Error())
		logg.Err(err)
		return
	}
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := Authenticated(r)
	data := templates.NewPageTemplate()
	data.Title = "login"
	data.Authenticated = authenticated

	err := templates.Render(w, templates.TEMPLATE_LOGIN_FORM, data)
	if err != nil {
		logg.Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Authenticated shows if user is authenticated and has "authenticated" value in session cookie.
func Authenticated(r *http.Request) (authenticated bool, hasAuthenticatedCookieValue bool) {
	session, _ := store.Get(r, COOKIE_NAME)
	authenticated, hasAuthenticatedCookieValue = session.Values["authenticated"].(bool)
	// log.Println("session authenticated", session.Values["authenticated"])
	return
}

func saveSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, COOKIE_NAME)
	session.Values["authenticated"] = true
	session.Values["username"] = r.FormValue("username")
	session.Save(r, w)
}

func Username(r *http.Request) string {

	session, _ := store.Get(r, COOKIE_NAME)
	usernmae, ok := session.Values["username"].(string)
	if !ok {
		logg.Err("Username not found in session or not a string")
		return ""
	}
	return usernmae
}
