package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/templates"
)

// PersonalPage renders full HTML page.
func PersonalPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	data := templates.PageTemplate{
		Title:         "Personal",
		Authenticated: authenticated,
		User:          auth.Username(r),
	}

	err := templates.Render(w, templates.TEMPLATE_PERSONAL_PAGE, data)
	if err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}
