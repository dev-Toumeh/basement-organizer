package common

import (
	"basement/main/internal/auth"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"net/http"
)

func Handle404NotFoundPage(w http.ResponseWriter, r *http.Request) {
	msg := "\"" + r.URL.Path + "\" page doesn't exist"

	// Render full page
	if server.WantsTemplateData(r) && r.Referer() == "" {
		w.WriteHeader(http.StatusNotFound)
		logg.Infof("%s: %s", msg, logg.NewError(msg))

		tmpl := templates.NewPageTemplate()
		tmpl.Title = "Page not found"
		ok, _ := auth.Authenticated(r)
		tmpl.Authenticated = ok
		tmpl.PageText = msg
		templates.Render(w, templates.TEMPLATE_NOT_FOUND_PAGE, tmpl)
	} else {
		server.WriteNotFoundError(msg, logg.NewError(msg), w, r)
	}
	return
}
