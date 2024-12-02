package boxes

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"maps"
	"net/http"
)

// ListPage shows a page with a box list.
func ListPage(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		// page template
		page := templates.NewPageTemplate()
		page.Title = "Boxes"
		page.RequestOrigin = "Boxes"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// list template
		listTmpl := common.ListTemplate{
			FormHXGet:   "/boxes",
			RowHXGet:    "/boxes",
			PlaceHolder: true,
			ShowLimit:   env.Config().ShowTableSize(),
		}

		// search-input template
		searchString := common.SearchString(r)
		listTmpl.SearchInput = true
		listTmpl.SearchInputLabel = "Search boxes"
		listTmpl.SearchInputValue = searchString

		count, err := db.BoxListCounter(searchString)
		if err != nil {
			server.WriteInternalServerError("cant query boxes", err, w, r)
			return
		}

		// pagination
		pageNr := common.ParsePageNumber(r)
		limit := common.ParseLimit(r)
		data = common.Pagination(data, count, limit, pageNr)
		listTmpl.Pagination = true
		listTmpl.CurrentPageNumber = data["PageNumber"].(int)
		listTmpl.Limit = limit
		listTmpl.PaginationButtons = data["Pages"].([]common.PaginationButton)

		// box-list-row to fill box-list template
		var boxes []common.ListRow

		// Boxes found
		if count > 0 {
			boxes, err = common.FilledRows(db.BoxListRows, searchString, limit, pageNr, count)
			if err != nil {
				server.WriteInternalServerError("cant query boxes", err, w, r)
				return
			}
		}
		listTmpl.Rows = boxes

		maps.Copy(data, listTmpl.Map())
		server.MustRender(w, r, templates.TEMPLATE_BOXES_LIST_PAGE, data)
	}
}
