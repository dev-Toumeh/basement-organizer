package shelves

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"

	"maps"
	"net/http"
)

// Render shelf Root page where you can search the available Shelves
func PageTemplate(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchString := common.SearchString(r)
		pageNumber := pageNumber(r)
		limit := env.DefaultTableSize()

		// Initialize page template
		user, _ := auth.UserSessionData(r)
		authenticated, _ := auth.Authenticated(r)

		page := templates.NewPageTemplate()
		page.Title = "Shelves"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// search-input template
		searchInput := items.NewSearchInputTemplate()
		searchInput.SearchInputLabel = "Search Shelves"
		searchInput.SearchInputValue = searchString
		maps.Copy(data, searchInput.Map())

		count, err := db.ShelfListCounter(searchString)
		if err != nil {
			server.WriteInternalServerError("error shelves counter", err, w, r)
			return
		}

		shelves, err := db.ShelfListRows(searchString, limit, pageNumber)
		if err != nil {
			server.WriteInternalServerError("cant query Shelves", err, w, r)
			return
		}

		// pagination
		data = common.Pagination(data, count, limit, pageNumber)

		shelvesMaps := make([]map[string]any, limit)
		for i, box := range shelves {
			mappedBox := box.Map()
			if mappedBox == nil {
				shelvesMaps[i] = map[string]any{}
			} else {
				shelvesMaps[i] = mappedBox
			}
		}
		// If count is less than limit, add empty maps to reach the limit
		for i := count; i < limit; i++ {
			shelvesMaps[i] = map[string]any{}
		}
		maps.Copy(data, map[string]any{"Shelves": shelvesMaps})

		server.MustRender(w, r, "shelves-page", data)
	}
}

// Render create Shelf Template with defaults Values
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

// Render Shelf Details Template where you can preview the shelf and update the relevant Data
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

// Render the SearchTemplate, you can use it with difference Args
// "/shelves/search?type=add" will return SearchTemplate with add functionality
// "/shelves/search?type=move" will return SearchTemplate with move functionality
// "/shelves/search?type=search" will return SearchTemplate with search functionality
// func SearchTemplate(db ShelfDB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		rowType, err := typeFromRequest(r)
// 		if err != nil {
// 			server.WriteInternalServerError("false request", err, w, r)
// 			return
// 		}
//
// 		shelvesMaps := shelfRowListData(w, r, db, rowType)
//
// 		data := searchInput.Map()
// 		data[rowType] = true
// 		data["Shelves"] = shelvesMaps
// 		data["Pagination"] = pagination
//
// 		server.MustRender(w, r, "shelf-list", data)
// 	}
// }

// render html element from type input Field with the desired Id and label
// the response element should be placed inside of the create/update of Item/Box forms
func InputTemplate(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		errMsgForUser := "can't add, please comeback later"
		id := server.ValidID(w, r, errMsgForUser)

		shelf, err := db.Shelf(id)
		if err != nil {
			server.WriteInternalServerError("false request", err, w, r)
			return
		}

		inputHTML := fmt.Sprintf(`
      <label for="shelf_id">Shelf</label>
      <input type="text" name="shelf_id" value="%s" hidden>
      <input type="text" name="label" value="%s" disabled>
      <button
       hx-get="/shelves/search?type=add"
       hx-target="#tg"
       hx-push-url="false">
       Add to another Shelf
      </button>
      `,
			shelf.ID.String(), shelf.Label)

		w.Header().Set("Content-Type", "text/html")

		w.Write([]byte(inputHTML))
	}
}
