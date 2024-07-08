package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/templates"
	"fmt"
	"net/http"
)

// PersonalPage renders full HTML page.
func PersonalPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	data := templates.PageTemplate{
		Title:         "Personal",
		Authenticated: authenticated,
		User:          auth.Username(r),
	}

	if err := templates.ApplyPageTemplate(w, PERSONAL_PAGE_TEMPLATE_PATH, data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}
