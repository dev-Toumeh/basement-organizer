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

// Prepare the necessary Data for the Shelf-list-rows
func getTemolateData(r *http.Request, db ShelfDB, w http.ResponseWriter) common.Data {
	data := common.InitData()
	data.SetBaseData(r)

	count, err := db.ShelfListCounter(data.GetSearchInputValue())
	if err != nil {
		server.WriteInternalServerError("error shelves counter", err, w, r)
	}

	data.SetTitle("Shelves")
	data.SetSearchInput(true)
	data.SetSearchInputLabel("Search Shelves")
	data.SetFormHXGet("/shelves")
	data.SetRowHXGet("/shelves")
	data.SetShowLimit(env.Config().ShowTableSize())
	data.SetCount(count)

	data = common.Pagination2(data)
	var shelves []common.ListRow
	if count > 0 {
		shelves, err = filledShelfRows(db, data)
		if err != nil {
			server.WriteInternalServerError("cant query shelves please comeback later", err, w, r)
		}
	}

	data.SetRows(shelves)
	return data
}

// filledShelfRows returns ListRows of Shelves with empty entries filled up to match limit.
// count - The total number of records found from the search query.
func filledShelfRows(db ShelfDB, data common.Data) ([]common.ListRow, error) {
	limit := data.GetLimit()
	shelvesMaps := make([]common.ListRow, limit)
	shelves, err := db.ShelfListRows(data.GetSearchInputValue(), limit, data.GetPageNumber())
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	for i, box := range shelves {
		shelvesMaps[i] = box
	}

	// If count is less than limit, add empty maps to reach the limit
	for i := data.GetCount(); i < limit; i++ {
		shelvesMaps[i] = common.ListRow{}
	}
	return shelvesMaps, nil
}
