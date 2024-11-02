package shelves

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/server"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// generate the necessary Data to render the shelf-list template
func shelfRowListData(w http.ResponseWriter, r *http.Request, db ShelfDB, typeRequest string) (map[any]any, []map[string]any, *items.SearchInputTemplate) {

	offsetRequest := offsetFromRequest(r)
	queryRequest := queryFromRequest(r)

	count, err := db.ShelfCounter(queryRequest)
	limit := env.DefaultTableSize()
	if err != nil {
		server.WriteInternalServerError("error shelves counter", err, w, r)
		return nil, nil, nil
	}

	// search-input template
	searchInput := items.NewSearchInputTemplate()
	searchInput.SearchInputLabel = "Search Shelves"
	searchInput.SearchInputValue = queryRequest

	pagination := paginate(count, limit, queryRequest, typeRequest)

	shelves, err := db.SearchShelves(limit, offsetRequest, count, queryRequest)
	if err != nil {
		server.WriteInternalServerError("cant query Shelves", err, w, r)
		return nil, nil, nil
	}

	// Map shelves to template data
	shelvesMaps := mapShelvesToTemplateData(shelves, count, limit, typeRequest)
	return shelvesMaps, pagination, searchInput
}

// check the count of available shelves and generate the pagination
// if the number of the available Shelves was less than the number of allowed rows per page it will return nil
// if the number was bigger it will generate the Pagination Data
func paginate(count, limit int, query string, typeRequest string) []map[string]any {
	if count <= limit {
		return nil
	}
	if typeRequest != "" {
		typeRequest = strings.ToLower(typeRequest[:1]) + typeRequest[1:]
	}

	totalPages := (count + limit - 1) / limit
	paginationData := make([]map[string]any, totalPages)

	for i := 0; i < totalPages; i++ {
		pageData := make(map[string]any)
		pageData["Query"] = query
		pageData["Type"] = typeRequest
		pageData["PageNumber"] = i + 1
		pageData["Offset"] = i * limit
		paginationData[i] = pageData
	}

	return paginationData
}

// transform the default listRow structs array into listRow maps array and
// include the type of rows what will help while rendering the search Template
func mapShelvesToTemplateData(shelves []*common.ListRow, count, limit int, rowType string) map[any]any {
	shelvesMaps := make(map[any]any, count)

	for i, shelf := range shelves {
		if shelf == nil {
			shelvesMaps[i] = map[string]any{}
		} else {
			shelfMap := shelf.Map()
			shelfMap[rowType] = true
			shelvesMaps[i] = shelfMap
		}
	}

	// If count is less than limit, add empty maps to reach the limit
	for i := count; i < limit; i++ {
		shelvesMaps[i] = map[string]any{}
	}

	return shelvesMaps
}

// retrieve the query Value from the Request and filter, sanitize it
func queryFromRequest(r *http.Request) string {
	searchTrimmed := strings.TrimSpace(r.FormValue("query"))
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(searchTrimmed, "*")
}

// retrieve the Offset Value from the Request and filter, sanitize it
func offsetFromRequest(r *http.Request) int {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	return offset
}

// retrieve the desired type from the Request (add, move or search) and validate it.
func typeFromRequest(r *http.Request) (string, error) {
	t := strings.TrimSpace(r.URL.Query().Get("type"))

	if t == "" {
		return "", nil
	}

	// Validate that 'type' is one of the allowed values
	allowedTypes := map[string]bool{"add": true, "move": true, "search": true}
	if _, ok := allowedTypes[t]; !ok {
		return "", fmt.Errorf("unexpected type: %s, while preparing the search Template", t)
	}

	return strings.ToUpper(t[:1]) + t[1:], nil
}
