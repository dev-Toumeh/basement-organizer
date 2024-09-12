package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type VirtualItem struct {
	Item_Id        uuid.UUID
	Label          string
	Box_label      string
	Box_id         uuid.UUID
	Shelve_label   string
	Area_label     string
	PreviewPicture string
}

// update the item based on ID
func SearchItemHandler(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			prepereResponse(w, r, db)
		} else {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		}
	}
}

// retrieve the items from the database based on the search query and rendering the results and send HTTP response .
func prepereResponse(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	searchQuery := r.FormValue("query")
	virtualItems, err := db.ItemFuzzyFinder(searchQuery)
	if err != nil {
		logg.Err(err)
	}
	logg.Debug("the search was triggered")
	templates.Render(w, "item-list-units", virtualItems)
	if err != nil {

		logg.Debug(err)
		templates.RenderErrorNotification(w, "something wrong happened")
	}
}
