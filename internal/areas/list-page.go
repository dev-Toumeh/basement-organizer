package areas

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"maps"
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
			FormHXGet:     "/areas",
			PlaceHolder:   true,
			ShowLimit:     env.CurrentConfig().ShowTableSize(),
			HideMoveCol:   true,
			RequestOrigin: common.ParseOrigin(r),
		}

		// search-input template
		searchString := common.SearchString(r)
		listTmpl.SearchInput = true
		listTmpl.SearchInputLabel = "Search areas"
		listTmpl.SearchInputValue = searchString

		count, err := db.AreaListCounter(searchString)
		if err != nil {
			server.WriteInternalServerError("cant query areas", err, w, r)
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
		var rows []common.ListRow

		// rows found
		if count > 0 {
			rowTemplateOptions := common.ListRowTemplateOptions{
				HideMoveCol: true,
				RowHXGet:    "/area",
			}
			rows, err = common.FilledRows(db.AreaListRows, searchString, limit, pageNr, count, rowTemplateOptions)
			if err != nil {
				server.WriteInternalServerError("cant query areas", err, w, r)
				return
			}
		}
		listTmpl.Rows = rows

		maps.Copy(data, listTmpl.Map())
		server.MustRender(w, r, templates.TEMPLATE_AREAS_LIST_PAGE, data)
	}
}
