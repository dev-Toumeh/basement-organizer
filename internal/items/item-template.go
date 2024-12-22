package items

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"maps"
	"net/http"
)

// Render shelf Root page where you can search the available Shelves
func PageTemplate(db ItemDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := getTemplateData(r, db, w)
		data.SetPlaceHolder(true)
		data.SetRequestOrigin("Items")
		server.MustRender(w, r, "items-page", data.TypeMap)
	}
}

// Render create Item Template with default values
func CreateTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		page := templates.NewPageTemplate()
		page.Title = "Add new Item"
		page.Authenticated = authenticated
		page.User = user

		item := newItem()
		data := page.Map()
		maps.Copy(data, item.Map())

		templates.Render(w, ITEM_CREATE_TEMPLATE, data)
	}
}

// Prepare the necessary Data for the items-list-rows
func getTemplateData(r *http.Request, db ItemDatabase, w http.ResponseWriter) common.Data {
	data := common.InitData(r)

	count, err := db.ItemListCounter(data.GetSearchInputValue())
	if err != nil {
		server.WriteInternalServerError("error items counter", err, w, r)
	}

	data.SetTitle("Items")
	data.SetSearchInput(true)
	data.SetSearchInputLabel("Search Items")
	data.SetFormHXGet("/items")
	data.SetRowHXGet("/items")
	data.SetShowLimit(env.Config().ShowTableSize())
	data.SetCount(count)

	data = common.Pagination2(data)
	var items []common.ListRow
	if count > 0 {
		items, err = filledItemRows(db, data)
		if err != nil {
			server.WriteInternalServerError("can't query items please comeback later", err, w, r)
		}
	}

	data.SetRows(items)
	return data
}

// Render Shelf Details Template where you can preview the shelf and update the relevant Data
func DetailsTemplate(db ItemDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errMsgForUser := "the requested item doesn't exist"

		data := common.InitData(r)
		id := server.ValidID(w, r, errMsgForUser)
		if id.IsNil() {
			return
		}

		// item, err := db.ItemById(id)
		// if err != nil {
		// 	server.TriggerErrorNotification(w, errMsgForUser)
		// }

		templates.Render(w, "shelf-details-page", data)
	}
}

// filledItemRows returns ListRows of Items with empty entries filled up to match limit.
// count - The total number of records found from the search query.
func filledItemRows(db ItemDatabase, data common.Data) ([]common.ListRow, error) {
	limit := data.GetLimit()
	itemsMaps := make([]common.ListRow, limit)
	items, err := db.ItemListRows(data.GetSearchInputValue(), limit, data.GetPageNumber())
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	for i, item := range items {
		itemsMaps[i] = item
	}

	// If count is less than limit, add empty maps to reach the limit
	for i := data.GetCount(); i < limit; i++ {
		itemsMaps[i] = common.ListRow{}
	}
	return itemsMaps, nil
}
