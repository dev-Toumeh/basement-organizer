package items

import (
	"basement/main/internal/auth"
	"basement/main/internal/templates"
	"maps"
	"net/http"
)

// Render create Item Template with default values
func CreateTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		page := templates.NewPageTemplate()
		page.Title = "Add new Item"
		page.Authenticated = authenticated
		page.User = user

		item := newItem()
		data := page.Map()
		maps.Copy(data, item.Map())

		templates.Render(w, ITEM_CREATE_TEMPLATE, data)
	}
}
