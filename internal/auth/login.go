package auth

import (
	"context"
	"fmt"
	"net/http"

	"basement/main/internal/logg"
	"basement/main/internal/templates"
)

const (
	LOGIN_FAILED_MESSAGE string = "Login failed"
)

func LoginHandler(db AuthDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			loginUser(w, r, db)
		}
		if r.Method == http.MethodGet {
			loginPage(w, r)
		}
		w.Header().Add("Allow", http.MethodGet)
		w.Header().Add("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func loginUser(w http.ResponseWriter, r *http.Request, db AuthDatabase) {
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
	user, err := db.UserByField(ctx, "username", username)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		logg.Info(LOGIN_FAILED_MESSAGE)
		logg.Debug("User", username, "doesn't exist")
		return
	}

	if user.Username != username {
		w.WriteHeader(http.StatusInternalServerError)
		logg.Info(LOGIN_FAILED_MESSAGE)
		logg.Errf(`username "%s" does not match user.Username from database "%s". This should not happen!`, username, user.Username)
		return
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, LOGIN_FAILED_MESSAGE)
		logg.Info(LOGIN_FAILED_MESSAGE)
		logg.Debug("pw hash doesnt match")
		return
	}

	saveSession(w, r, user)

	logg.Info("login successful")

	// https://htmx.org/headers/hx-location/
	w.Header().Add("HX-Location", "/items")
	fmt.Fprintf(w, "Welcome %v\n", username)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := Authenticated(r)
	data := templates.NewPageTemplate()
	data.Title = "login"
	data.Authenticated = authenticated
	logg.Debug(data)
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

	err := templates.SafeRender(w, templates.TEMPLATE_LOGIN_FORM, data)
	if err != nil {
		logg.Debug(http.StatusText(http.StatusInternalServerError))
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

func saveSession(w http.ResponseWriter, r *http.Request, user User) {
	session, _ := store.Get(r, COOKIE_NAME)
	session.Values["authenticated"] = true
	session.Values["username"] = user.Username
	session.Values["id"] = user.Id.String()
	session.Save(r, w)
}

func UserSessionData(r *http.Request) (string, string) {

	session, _ := store.Get(r, COOKIE_NAME)
	username, ok1 := session.Values["username"].(string)
	id, ok2 := session.Values["id"].(string)
	if !ok1 || !ok2 {
		logg.Err("corrupted session, check UserSessionData function")
		http.RedirectHandler("/logout", http.StatusUnauthorized)
	}
	return username, id
}
