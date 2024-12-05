package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/templates"
)

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	if !authenticated {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
	username, _ := auth.UserSessionData(r)
	data := templates.NewPageTemplate()
	data.Title = "Settings"
	data.RequestOrigin = "Settings"
	data.Authenticated = authenticated
	data.User = username

	err := templates.Render(w, "settings-page", data)
	if err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}
