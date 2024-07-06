package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/templates"
)

// this is the first route inside of the application, so this page will open by Default
func HomePage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	data := templates.NewPageTemplate()
	data.Title = "home"
	data.Authenticated = authenticated

	err := templates.Render(w, "home-page", data)
	if err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}
