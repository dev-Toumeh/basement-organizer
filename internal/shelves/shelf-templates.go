package shelves

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"

	"maps"
	"net/http"
)

// Render shelf Root page where you can search the available Shelves
func PageTemplate(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Initialize page template
		user, _ := auth.UserSessionData(r)
		authenticated, _ := auth.Authenticated(r)

		page := templates.NewPageTemplate()
		page.Title = "Shelves"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// list template
		listTmpl := common.ListTemplate{
			FormHXGet: "/shelves",
			RowHXGet:  "/shelves",
			ShowLimit: env.Config().ShowTableSize(),
		}

		// search-input template
		searchString := common.SearchString(r)
		listTmpl.SearchInput = true
		listTmpl.SearchInputLabel = "Search Shelves"
		listTmpl.SearchInputValue = searchString

		count, err := db.ShelfListCounter(searchString)
		if err != nil {
			server.WriteInternalServerError("error shelves counter", err, w, r)
			return
		}

		// box-list-row to fill box-list template
		var shelves []common.ListRow

		// pagination
		pageNr := common.ParsePageNumber(r)
		limit := common.ParseLimit(r)
		data = common.Pagination(data, count, limit, pageNr)
		listTmpl.Pagination = true
		listTmpl.CurrentPageNumber = data["PageNumber"].(int)
		listTmpl.Limit = limit
		listTmpl.PaginationButtons = data["Pages"].([]common.PaginationButton)

		if count > 0 {
			shelves, err = filledShelfRows(db, searchString, limit, pageNr, count)
			if err != nil {
				server.WriteInternalServerError("cant query shelves please comeback later", err, w, r)
				return
			}
		}
		fmt.Print(len(shelves))
		listTmpl.Rows = shelves

		maps.Copy(data, listTmpl.Map())
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
