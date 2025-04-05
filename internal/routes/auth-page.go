package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/templates"
)

// this is the first route inside of the application, so this page will open by Default
func AuthPage(w http.ResponseWriter, r *http.Request) {

	authenticated, _ := auth.Authenticated(r)
	if authenticated {
		http.Redirect(w, r, "/items", http.StatusPermanentRedirect)
	}
	data := templates.NewPageTemplate()
	data.Title = "home"
	data.Authenticated = authenticated

	err := templates.Render(w, "auth-page", data)
	if err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}
