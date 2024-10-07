package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"net/http"
	"strconv"
)

type SearchItemData struct {
	Query          string
	TotalCount     int
	Records        []ListRow
	PaginationData []PaginationData
}

type PaginationData struct {
	Offset     int
	PageNumber int
}

type SearchInputTemplate struct {
	SearchInputLabel    string
	SearchInputHxPost   string
	SearchInputHxTarget string
	SearchInputValue    string
}

func NewSearchInputTemplate() *SearchInputTemplate {
	return &SearchInputTemplate{SearchInputLabel: "Search",
		SearchInputHxPost:   "/api/v1/implement-me",
		SearchInputHxTarget: "#item-list-body",
		SearchInputValue:    "",
	}
}

func NewSearchItemInputTemplate() *SearchInputTemplate {
	return &SearchInputTemplate{SearchInputLabel: "Search items", SearchInputHxPost: "/api/v1/search/item", SearchInputHxTarget: "#item-list-body"}
}

func (tmpl *SearchInputTemplate) Map() map[string]any {
	return map[string]any{
		"SearchInputLabel":    tmpl.SearchInputLabel,
		"SearchInputHxPost":   tmpl.SearchInputHxPost,
		"SearchInputHxTarget": tmpl.SearchInputHxTarget,
		"SearchInputValue":    tmpl.SearchInputValue,
	}
}

// update the item based on ID
func SearchItemHandler(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			prepereResponse(w, r, db)
		} else {
			templates.RenderErrorNotification(w, "Method Not Allowed")
		}
	}
}

// retrieve the items from the database based on the search query and rendering the results and send HTTP response .
func prepereResponse(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	var SearchItemData SearchItemData
	var err error

	SearchItemData.Query = r.FormValue("query")
	SearchItemData.TotalCount, err = db.NumOfItemRecords(SearchItemData.Query)
	if err != nil {
		logg.Errf("error accrue while checking the number of items %v", err)
		templates.RenderErrorNotification(w, "something wend wrong please comeback later")
	}

	if SearchItemData.TotalCount <= 10 {
		virtualItems, err := db.ItemFuzzyFinder(SearchItemData.Query)
		if err != nil {
			logg.Err(err)
		}
		SearchItemData.PaginationData = []PaginationData{}
		SearchItemData.Records = virtualItems
		err = templates.Render(w, "item-list-units", SearchItemData)
		if err != nil {
			logg.Debug(err)
			templates.RenderErrorNotification(w, "something wrong happened")
		}

	} else {
		SearchItemData.PaginationData = generatePaginationData(SearchItemData.TotalCount, 10)
		SearchItemData.Records, err = db.ItemFuzzyFinderWithPagination(SearchItemData.Query, 10, SearchItemData.PaginationData[0].Offset)
		err = templates.Render(w, "item-list-units-pagination", SearchItemData)
		if err != nil {
			logg.Debug(err)
			templates.RenderErrorNotification(w, "something wrong happened")
		}
	}
}

func ItemPaginationHandler(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			var SearchItemData SearchItemData
			var err error

			queryValues := r.URL.Query()
			query := queryValues.Get("query")
			offsetStr := queryValues.Get("offset")
			logg.Debugf("query %s, offset %s \n", query, offsetStr)
			offset, err := strconv.Atoi(offsetStr)
			if err != nil {
				templates.RenderErrorNotification(w, "Invalid search request")
				return
			}

			SearchItemData.Records, err = db.ItemFuzzyFinderWithPagination(query, 10, offset)
			if err != nil {
				logg.Debug(err)
				templates.RenderErrorNotification(w, "Something wrong happened")
				return
			}

			err = templates.Render(w, "item-list-units-pagination", SearchItemData)
			if err != nil {
				logg.Debug(err)
				templates.RenderErrorNotification(w, "Something wrong happened")
			}
		} else {
			templates.RenderErrorNotification(w, "Method Not Allowed")
		}
	}
}

func generatePaginationData(totalRecords, pageSize int) []PaginationData {
	totalPages := (totalRecords + pageSize - 1) / pageSize

	var pagination []PaginationData
	for i := 0; i < totalPages; i++ {
		pagination = append(pagination, PaginationData{
			Offset:     i * pageSize,
			PageNumber: i + 1,
		})
	}

	return pagination
}
