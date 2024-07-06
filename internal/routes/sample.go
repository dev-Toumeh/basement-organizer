package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/templates"
	"net/http"
)

func SamplePage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	data := templates.NewPageTemplate()
	data.Title = "Sample Page"
	data.Authenticated = authenticated
	data.User = auth.Username(r)

	templates.RenderPage(w, "sample-page", data)
}
