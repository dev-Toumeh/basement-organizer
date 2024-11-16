package items

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"maps"
	"net/http"
)

// Render shelf Root page where you can search the available Shelves
func PageTemplate(db ItemDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Initialize page template
		user, _ := auth.UserSessionData(r)
		authenticated, _ := auth.Authenticated(r)

		page := templates.NewPageTemplate()
		page.Title = "items"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// list template
		listTmpl := common.ListTemplate{
			FormHXGet: "/items",
			RowHXGet:  "/items",
			ShowLimit: env.Config().ShowTableSize(),
		}

		// search-input template
		searchString := common.SearchString(r)
		listTmpl.SearchInput = true
		listTmpl.SearchInputLabel = "Search items"
		listTmpl.SearchInputValue = searchString

		count := 0

		// box-list-row to fill box-list template
		var items []common.ListRow

		// pagination
		pageNr := common.ParsePageNumber(r)
		limit := common.ParseLimit(r)
		data = common.Pagination(data, count, limit, pageNr)
		listTmpl.Pagination = true
		listTmpl.CurrentPageNumber = data["PageNumber"].(int)
		listTmpl.Limit = limit
		listTmpl.PaginationButtons = data["Pages"].([]common.PaginationButton)

		listTmpl.Rows = items

		maps.Copy(data, listTmpl.Map())
		server.MustRender(w, r, ITEM_PAGE_TEMPLATE, data)
	}
}

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
