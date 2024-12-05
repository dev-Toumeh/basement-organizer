package areas

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"net/http"
)

// listPage shows a page with a area list.
func listPage(db AreaDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		// page template
		page := templates.NewPageTemplate()
		page.Title = "Areas"
		page.RequestOrigin = "Areas"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// list template
		listTmpl := common.ListTemplate{
			FormHXGet:   "/areas",
			RowHXGet:    "/area",
			PlaceHolder: true,
			ShowLimit:   env.Config().ShowTableSize(),
		}

		// search-input template
		searchString := common.SearchString(r)
		listTmpl.SearchInput = true
		listTmpl.SearchInputLabel = "Search areas"
		listTmpl.SearchInputValue = searchString

		server.MustRender(w, r, "areas-list-page", data)
	}
}
