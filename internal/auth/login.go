package auth

import (
	"fmt"
	"log"
	"net/http"

	"basement/main/internal/database"
	"basement/main/internal/templates"
)

const (
	LOGIN_FAILED_MESSAGE string = "Login failed"
)

func LoginHandler(db *database.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			loginUser(w, r, db)
		}
		if r.Method == http.MethodGet {
			loginPage(w, r, db)
		}
	}
}
func loginUser(w http.ResponseWriter, r *http.Request, db *database.JsonDB) {
	authenticated, ok := Authenticated(r)

	if ok {
		if authenticated {
			log.Printf("LoginHandler - ok: %v authenticated: %v", ok, authenticated)
			fmt.Fprint(w, "already logged in")
			return
		}
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Printf("%v %v ", r.URL, r.Form.Encode())

	if username == "" {
		log.Println("Missing username")
		w.WriteHeader(http.StatusBadRequest)
		templates.RenderErrorSnackbar(w, LOGIN_FAILED_MESSAGE)
		return
	}
	if password == "" {
		log.Println("Missing password")
		w.WriteHeader(http.StatusBadRequest)
		templates.RenderErrorSnackbar(w, LOGIN_FAILED_MESSAGE)
		return
	}

	user, _ := db.User(username)

	if !checkPasswordHash(password, user.PasswordHash) {
		log.Println("pw hash doesnt match")
		w.WriteHeader(http.StatusBadRequest)
		templates.RenderErrorSnackbar(w, LOGIN_FAILED_MESSAGE)
		return
	}

	saveSession(w, r)

	log.Println("login successful")
	http.RedirectHandler("/personal-page", http.StatusOK)

	// https://htmx.org/headers/hx-location/
	w.Header().Add("HX-Location", "/personal-page")
	fmt.Fprintf(w, "Welcome %v\n", username)
}

func loginPage(w http.ResponseWriter, r *http.Request, db *database.JsonDB) {
	authenticated, _ := Authenticated(r)
	data := templates.NewPageTemplate()
	data.Title = "login"
	data.Authenticated = authenticated

	err := templates.Render(w, templates.TEMPLATE_LOGIN_PAGE, data)
	if err != nil {
		templates.RenderErrorSnackbar(w, err.Error())
		log.Println("loginPage error:", err)
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
		fmt.Fprintln(w, "failed")
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
		log.Println("Username not found in session or not a string")
		return ""
	}
	return usernmae
}
