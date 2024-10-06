package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/templates"
	"fmt"
	"net/http"
)

func shelvesPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	user, _ := auth.UserSessionData(r)
	data := templates.PageTemplate{
		Title:         "Shelves",
		Authenticated: authenticated,
		User:          user,
	}

	err := templates.Render(w, "shelves-page", data.Map())
	if err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

func newShef(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	user, _ := auth.UserSessionData(r)
	data := templates.PageTemplate{
		Title:         "Shelves",
		Authenticated: authenticated,
		User:          user,
	}

	err := templates.Render(w, "new-shelf", data.Map())
	if err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}
