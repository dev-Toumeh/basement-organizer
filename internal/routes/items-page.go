package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
)

func itemsPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	user, _ := auth.UserSessionData(r)
	data := templates.PageTemplate{
		Title:         "Personal",
		Authenticated: authenticated,
		User:          user,
	}

	err := templates.Render(w, "items-page", data)
	if err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

// this function will generate a new from
func itemTemp(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "create-item-form", ""); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

// generate Search Item Template, in case of get request
func searchItemTemp(w http.ResponseWriter, r *http.Request) {
	err := templates.Render(w, "search-item-form", "")
	if err != nil {
		logg.Debug(err)
		templates.RenderErrorSnackbar(w, "something wrong happened")
	}
}
