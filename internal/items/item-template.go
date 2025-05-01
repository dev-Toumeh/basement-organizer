package items

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"net/http"
)

// Render Item Root page where you can search the available Items
func PageTemplate(db ItemDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := getTemplateData(r, db, w)
		data.SetEnvDevelopment(env.Development())
		data.SetPlaceHolder(true)
		data.SetRequestOrigin("Items")
		server.MustRender(w, r, "item-page-template", data.TypeMap)
	}
}

// Render create Item Template with default values
func CreateTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item := newItem()
		renderItemTemplate(r, w, item.Map(), common.CreateMode)
	}
}

// Render Update Item Template
func UpdateTemplate(db ItemDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := server.ValidID(w, r, "the requested Item doesn't exist")
		if id.IsNil() {
			return
		}
		item, err := db.ItemById(id)
		if err != nil {
			server.WriteInternalServerError("can't query items please comeback later", err, w, r)
			return
		}
		renderItemTemplate(r, w, item.Map(), common.EditMode)
	}
}

// Render Preview Item Template
func PreviewTemplate(db ItemDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := server.ValidID(w, r, "the requested Item doesn't exist")
		if id.IsNil() {
			return
		}
		item, err := db.ItemById(id)
		if err != nil {
			server.WriteInternalServerError("can't query items please comeback later", err, w, r)
			return
		}
		renderItemTemplate(r, w, item.Map(), common.PreviewMode)
	}
}

// Prepare the necessary Data for the items-list-rows
func getTemplateData(r *http.Request, db ItemDatabase, w http.ResponseWriter) common.Data {
	data := common.InitData(r, true)

	count, err := db.ItemListCounter(data.GetSearchInputValue())
	if err != nil {
		server.WriteInternalServerError("error items counter", err, w, r)
		return common.Data{}
	}

	data.SetTitle("Items")
	data.SetSearchInput(true)
	data.SetSearchInputLabel("Search Items")
	data.SetFormHXGet("/items")
	data.SetRowHXGet("/item")
	data.SetShowLimit(env.CurrentConfig().ShowTableSize())
	data.SetCount(count)

	data = common.Pagination2(data)
	var items []common.ListRow
	if count > 0 {
		data.SetListRowTemplateOptions(common.ListRowTemplateOptions{RowHXGet: "item"})
		items, err = filledItemRows(db, data)
		if err != nil {
			server.WriteInternalServerError("can't query items please comeback later", err, w, r)
			return common.Data{}
		}
	}

	data.SetRows(items)
	return data
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
		itemsMaps[i].ListRowTemplateOptions = data.GetListRowTemplateOptions()
	}

	// If count is less than limit, add empty maps to reach the limit
	for i := data.GetCount(); i < limit; i++ {
		itemsMaps[i] = common.ListRow{}
	}
	return itemsMaps, nil
}

func renderItemTemplate(r *http.Request, w http.ResponseWriter, values map[string]any, typeMode common.Mode) {
	data := common.InitData(r, false)
	data.SetEnvDevelopment(env.Development())
	data.SetTypeMode(typeMode)
	data.SetDetailesData(values)
	data.SetOrigin("Item")
	err := templates.Render(w, "item-template", data.TypeMap)
	if err != nil {
		logg.Warningf("An Error accrue while fetching item Extra Info", err)
	}
}
