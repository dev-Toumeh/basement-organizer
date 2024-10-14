package shelves

import (
	"basement/main/internal/auth"
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/server"
	"math"

	"basement/main/internal/templates"
	_ "maps"
	"net/http"
	"strconv"
)

func ShelvesPage(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		// page template
		page := templates.NewPageTemplate()
		page.Title = "Shelves"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		server.MustRender(w, r, "shelves-page", data)

		//
		// // search-input template
		// query := r.FormValue("query")
		// searchInput := items.NewSearchInputTemplate()
		// searchInput.SearchInputLabel = "Search boxes"
		// searchInput.SearchInputValue = query
		//
		// maps.Copy(page.Map(), searchInput.Map())
		//
		// // @TODO: Implement move page
		// var err error
		// var count int
		// var shelves []*items.ListRow
		// usedSearch := false
		// urlQuery := r.URL.Query()
		//
		// pageNr, limit, err := searchPaginationData(r)
		//
		// err = nil
		// if urlQuery.Has("query") && query != "" {
		// 	shelves, count, err = db.ShelfSearchListRowsPaginated(pageNr, limit, query)
		// 	usedSearch = true
		// } else {
		// 	shelves, count, err = db.ShelfSearchListRowsPaginated(pageNr, 10, query)
		// }
		// if err != nil {
		// 	server.WriteInternalServerError("cant query boxes", err, w, r)
		// 	return
		// }
		//
		// pagination(pageNr, count, limit, usedSearch, shelves)
	}
}

// retrieve the data from the request to create the pagination
// return the number of pages and the amount of rows for every page
func searchPaginationData(r *http.Request) (int, int, error) {
	pageNr, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || pageNr < 1 {
		pageNr = 1
	}
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil || limit == 0 {
		limit = env.DefaultTableSize()
	}
	return pageNr, limit, err
}

func pagination(pageNr int, count int, limit int, usedSearch bool, shelves []*items.ListRow) {
	totalPages := 1
	totalPagesF := float64(count) / float64(limit)
	totalPagesCeil := math.Ceil(float64(totalPagesF))
	totalPages = int(totalPagesCeil)
	if totalPages < 1 {
		totalPages = 1
	}

	currentPage := 1
	if pageNr > totalPages {
		currentPage = totalPages
	} else {
		currentPage = pageNr
	}

	if totalPages == 0 {
		totalPages = 1
	}

	nextPage := currentPage + 1
	if nextPage < 1 {
		nextPage = 1
	}
	if nextPage > totalPages {
		nextPage = totalPages
	}

	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}
	if prevPage > totalPages {
		prevPage = totalPages
	}

	if usedSearch {
		fromOffset := (currentPage - 1) * limit
		toOffset := currentPage * limit
		if toOffset > count {
			toOffset = count
		}
		if toOffset < 0 {
			toOffset = 0
		}
		if fromOffset < 0 {
			fromOffset = 0
		}
		shelves = shelves[fromOffset:toOffset]
	}
}
