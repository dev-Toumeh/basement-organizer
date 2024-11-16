package shelves

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"fmt"
	"net/http"
	"strings"
)

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

func shelfListRowTemplate(r *http.Request, db ShelfDB, w http.ResponseWriter) {
	data := common.Data
	pageNr := common.ParsePageNumber(r)
	limit := common.ParseLimit(r)

	// list template
	data.SetFormHXGet("/shelves")
	data.SetRowHXGet("/shelves")
	data.SetShowLimit(env.Config().ShowTableSize())

	// search-input template
	searchString := common.SearchString(r)
	data.SetSearchInput(true)
	data.SetSearchInputLabel("Search Shelves")
	data.SetSearchInputValue(searchString)

	count, err := db.ShelfListCounter(searchString)
	if err != nil {
		server.WriteInternalServerError("error shelves counter", err, w, r)
	}

	var shelves []common.ListRow

	// pagination
	data.SetPageNumber(pageNr)
	data.SetLimit(limit)
	data.SetCount(count)

	if count > 0 {
		data.SetPagination(true)
		common.Pagination2()
		shelves, err = filledShelfRows(db, searchString, limit, pageNr, count)
		if err != nil {
			server.WriteInternalServerError("cant query shelves please comeback later", err, w, r)
		}
	}

	data.SetRows(shelves)
}
