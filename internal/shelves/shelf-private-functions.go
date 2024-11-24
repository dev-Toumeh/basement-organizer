package shelves

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"net/http"
)

// Prepare the necessary Data for the Shelf-list-rows
func getTemplateData(r *http.Request, db ShelfDB, w http.ResponseWriter) common.Data {
	data := common.InitData(r)

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
