package shelves

import (
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// retrieve the Page Number Value from the Request and filter, sanitize it
func pageNumber(r *http.Request) int {
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageNumber = 1
	}

	return pageNumber
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

// filledShelfRows returns ListRows of Shelves with empty entries filled up to match limit.
// count - The total number of records found from the search query.
func filledShelfRows(db ShelfDB, searchString string, limit int, pageNr int, count int) ([]common.ListRow, error) {
	shelvesMaps := make([]common.ListRow, limit)
	shelves, err := db.ShelfListRows(searchString, limit, pageNr)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	for i, box := range shelves {
		shelvesMaps[i] = box
	}

	// If count is less than limit, add empty maps to reach the limit
	for i := count; i < limit; i++ {
		shelvesMaps[i] = common.ListRow{}
	}
	return shelvesMaps, nil
}
