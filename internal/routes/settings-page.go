package routes

import (
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/server"
	"basement/main/internal/templates"
)

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	data := settingsPageData(r)
	server.MustRender(w, r, "settings-page", data.Map())
}

func settingsPageData(r *http.Request) templates.PageTemplate {
	authenticated, _ := auth.Authenticated(r)
	username, _ := auth.UserSessionData(r)
	data := templates.NewPageTemplate()
	data.Title = "Settings"
	data.RequestOrigin = "Settings"
	data.Authenticated = authenticated
	data.User = username

	return data
}
