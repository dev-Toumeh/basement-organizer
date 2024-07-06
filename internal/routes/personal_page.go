package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/database"
	"basement/main/internal/templates"
	"fmt"
	"net/http"
)

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

func ItemHandler(db *database.JsonDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := r.PathValue("id")
			// get "Water Bottle" "123e4567-e89b-12d3-a456-426614174002"
			data := db.Items[id]
			templates.Render(w, "item-container", data)
			return
		}
		w.Header().Add("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		return
	}
}
