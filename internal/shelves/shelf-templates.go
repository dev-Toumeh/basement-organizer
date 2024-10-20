package shelves

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/server"
	"basement/main/internal/templates"

	"maps"
	"math"
	"net/http"
	"strconv"
)

func ShelvesPage(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		// Initialize page template
		page := templates.NewPageTemplate()
		page.Title = "Shelves"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// search-input template
		query := ""
		searchInput := items.NewSearchInputTemplate()
		searchInput.SearchInputLabel = "Search boxes"
		searchInput.SearchInputValue = query

		maps.Copy(page.Map(), searchInput.Map())
		var err error
		var shelves []*items.ListRow
		var count, offset int
		limit := 100

		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			offset = 0
		}

		if len(query) == 0 {
			shelves, count, err = db.ShelfListRowsPaginated(limit, offset)
			if count > limit {
				createPagination()
			}
		} else {
			err = handleShelfSearchRequest(query, db)
		}
		if err != nil {
			server.WriteInternalServerError("cant query Shelves", err, w, r)
			return
		}

		// Map shelves to template data
		shelvesMaps := make(map[int]any, count)
		for i, shelf := range shelves {
			shelvesMaps[i] = shelf.Map()
		}
		data["Shelves"] = shelvesMaps
		server.MustRender(w, r, "shelves-page", data)
	}
}

func createPagination() {
	panic("unimplemented")
}

// generate create new Shelf Template with defaults Values
func CreateTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		page := templates.NewPageTemplate()
		page.Title = "Add new Shelf"
		page.Authenticated = authenticated
		page.User = user

		shelf := newShelf()
		data := page.Map()
		maps.Copy(data, shelf.Map())

		templates.Render(w, "shelf-create-page", data)
	}
}

// Generate Shelf Details Page where you can preview the shelf and update the relevant Data
func DetailsTemplate(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errMsgForUser := "the requested Shelf doesn't exist"

		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		id := server.ValidID(w, r, errMsgForUser)
		if id.IsNil() {
			return
		}

		shelf, err := db.Shelf(id)
		if err != nil {
			server.TriggerErrorNotification(w, errMsgForUser)
		}

		page := templates.NewPageTemplate()
		page.Title = "Shelf Details"
		page.Authenticated = authenticated
		page.User = user

		maps := []map[string]any{
			page.Map(),
			shelf.Map(),
			{"Edit": common.CheckEditMode(r)},
		}

		data := common.MergeMaps(maps)

		templates.Render(w, "shelf-details-page", data)
	}
}

func handleShelfSearchRequest(query string, db ShelfDB) (err error) {
	// shelves, count, err := db.ShelfSearchListRowsPaginated(1, 20, query)
	// shelvesMaps := make([]map[string]any, count)
	return err
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
