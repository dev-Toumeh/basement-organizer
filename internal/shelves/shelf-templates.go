package shelves

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"regexp"
	"strings"

	"maps"
	"net/http"
	"strconv"
)

var offsetData *[]int

func ShelvesPage(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var shelves []*items.ListRow

		limit := env.DefaultTableSize()
		query := queryFromRequest(r)
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		// Initialize page template
		page := templates.NewPageTemplate()
		page.Title = "Shelves"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// search-input template
		searchInput := items.NewSearchInputTemplate()
		searchInput.SearchInputLabel = "Search Shelves"
		searchInput.SearchInputValue = query
		maps.Copy(data, searchInput.Map())

		count, err := db.ShelfCounter(query)
		if err != nil {
			server.WriteInternalServerError("error shelves counter", err, w, r)
			return
		}

		offsetRequest := offsetFromRequest(r)
		pagination := paginate(count, limit, query)
		shelves, err = db.SearchShelves(limit, offsetRequest, count, query)
		if err != nil {
			server.WriteInternalServerError("cant query Shelves", err, w, r)
			return
		}

		// Map shelves to template data
		shelvesMaps := mapShelvesToTemplateData(shelves, count, limit)

		data["Shelves"] = shelvesMaps
		data["Pagination"] = pagination
		data["Limit"] = limit
		server.MustRender(w, r, "shelves-page", data)
	}
}

// Generate create new Shelf Template with defaults Values
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

// check the count of available shelves and generate the pagination
// if the number of the available Shelves was less than the number of allowed rows per page it will return nil
// if the number was bigger it will generate the Pagination Data
func paginate(count, limit int, query string) []map[string]any {
	if count <= limit {
		return nil
	}

	totalPages := (count + limit - 1) / limit
	paginationData := make([]map[string]any, totalPages)

	for i := 0; i < totalPages; i++ {
		pageData := make(map[string]any)
		pageData["Query"] = query
		pageData["PageNumber"] = i + 1
		pageData["Offset"] = i * limit
		paginationData[i] = pageData
	}

	return paginationData
}

// return query Value from the Request and filter, sanitize it
func queryFromRequest(r *http.Request) string {
	searchTrimmed := strings.TrimSpace(r.FormValue("query"))
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(searchTrimmed, "*")
}

func offsetFromRequest(r *http.Request) int {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	return offset
}

func mapShelvesToTemplateData(shelves []*items.ListRow, count, limit int) map[int]any {
	shelvesMaps := make(map[int]any, count)

	// Original mapping logic
	for i, shelf := range shelves {
		if shelf == nil {
			shelvesMaps[i] = map[string]any{}
		} else {
			shelvesMaps[i] = shelf.Map()
		}
	}

	// If count is less than limit, add empty maps to reach the limit
	for i := count; i < limit; i++ {
		shelvesMaps[i] = map[string]any{}
	}

	return shelvesMaps
}
