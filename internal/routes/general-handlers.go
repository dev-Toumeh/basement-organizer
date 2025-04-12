package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
)

// Handle registers a route that requires authentication.
// If the user is not authenticated, they are redirected to the /auth page.
func Handle(route string, handler http.HandlerFunc) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		if !authenticated {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}

		msg := ""
		msg = fmt.Sprintf(`%s "%s" http://%s%s%s`, r.Method, route, r.URL.Scheme, r.Host, r.URL)
		colorMsg := fmt.Sprintf("%s%s%s", logg.Yellow, msg, logg.Reset)
		logg.Debug(colorMsg)

		if r.Method == http.MethodPost {
			// @TODO: Fix. Breaks some post requests because r.ParseForm is empty after this.
			// r.ParseForm()
			// colorMsg := fmt.Sprintf("%sPostFormValue: %v%s", logg.Yellow, r.PostForm, logg.Reset)
			// logg.Debug(colorMsg)
		}

		handler.ServeHTTP(w, r)
	})
}

// HandlePublic registers a public route that does not require authentication.
// Useful for pages like login or registration.
func HandlePublic(route string, handler http.HandlerFunc) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf(`%s "%s" http://%s%s%s`, r.Method, route, r.URL.Scheme, r.Host, r.URL)
		colorMsg := fmt.Sprintf("%s%s%s", logg.Yellow, msg, logg.Reset)
		logg.Debug(colorMsg)

		handler.ServeHTTP(w, r)
	})
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	if authenticated {
		http.Redirect(w, r, "/items", http.StatusPermanentRedirect)
		return
	}

	data := common.InitData(r, false)
	data.SetTitle("Authentication Page")
	server.MustRender(w, r, templates.TEMPLATE_AUTH_PAGE, data.TypeMap)
}

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
