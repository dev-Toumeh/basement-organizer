package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/templates"
)

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	username, _ := auth.UserSessionData(r)
	data := templates.NewPageTemplate()
	data.Title = "Setting"
	data.Authenticated = authenticated
	data.User = username

	err := templates.Render(w, "settings", data)
	if err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}
