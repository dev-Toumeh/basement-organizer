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
	data := templates.PageTemplate{
		Title:         "home",
		Authenticated: authenticated,
	}
	if err := templates.ApplyPageTemplate(w, "internal/templates/home.html", data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}
