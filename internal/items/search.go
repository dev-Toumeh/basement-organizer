package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"fmt"
	"net/http"
)

// update the item based on ID
func SearchItemHandler(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			prepereResponse(w, r, db)
		} else if r.Method == http.MethodGet {
			renderSearchForm(w)
		}
	}
}

// retrieve the items from the database based on the search query and rendering the results and send HTTP response .
func prepereResponse(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	searchQuery := r.FormValue("query")
	items, err := db.SearchItemsByLabel(searchQuery)
	if err != nil {
		logg.Err(err)
	}
	searchResult := renderSearchResult(items)
	fmt.Fprintf(w, searchResult)
}

// generate Search Item Template, in case of get request
func renderSearchForm(w http.ResponseWriter) {
	err := templates.Render(w, "search-item-form", "")
	if err != nil {
		logg.Debug(err)
		templates.RenderErrorSnackbar(w, "something wrong happened")
	}
}

// The function will take relevant data and dynamically generate a series of HTML form tags.
func renderSearchResult(items []struct {
	Id    string
	Label string
}) string {

	var result string
	for _, item := range items {
		result += fmt.Sprintf(`
   <form hx-get="/item" id="form-item" hx-target="#main-container">
			         <input type="hidden" name="id" value="%s">
              <button type="submit">%s</button>
          </form>`, item.Id, item.Label)
	}
	return result
}
