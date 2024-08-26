package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

// DataWriteFunc should implement a function to write a template response or normal response.
//
// Example:
//
//	func(w io.Writer, data any) {
//		// templates
//		templates.Render(w, "items-container", data)
//		// Fprint
//		fmt.Fprint(w, data)
//	})
type DataWriteFunc func(w io.Writer, data any)

// ReadItemHandler returns a single item.
//
// Accepts "/item?id=" and "/item/id"
func ReadItemHandler(db ItemDatabase, responseWriter DataWriteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			logg.Info("access: ", r.URL)

			id := r.FormValue("id")
			if id == "" {
				id = r.PathValue("id")
			}

			data, err := db.ItemByField("id", id)
			if err != nil && err != db.ErrorExist() {
				w.WriteHeader(http.StatusInternalServerError)
				templates.RenderErrorSnackbar(w, err.Error())
			}
			if data.Id.IsNil() {
				logg.Debug("item not found: ", id)
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, `Item with id "%s" not found`, id)
				return
			}
			logg.Debugf("item: %T=%v, %T=%v", data.Id, data.Id, data.Label, data.Label)
			responseWriter(w, data)
			return
		}
		w.Header().Add("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		return
	}
}

// ReadItemsHandler returns a list of items or list of item IDs.
//
// Accepts "/items" to return all items with all information.
//
//	id := uuid.must(uuid.fromstring(r.FormValue("query")),)
//
// Accepts "/items?query=id" to only return item IDs.
func ReadItemsHandler(db ItemDatabase, responseWriter DataWriteFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := r.FormValue("query")
			switch id {
			// return all item IDs
			case "id":
				ids, err := db.ItemIDs()
				if err != nil {
					fmt.Fprintln(w, "something wrong happened please comeback later")
				}
				responseWriter(w, ids)
			// return all items
			default:
				ids, err := db.ItemIDs()
				if err != nil {
					http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
					return
				}
				tmpl, err := template.New("ids").Parse("{{range .}}<div>{{.}}</div>{{ end}}")
				if err != nil {
					http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
					return
				}
				err = tmpl.Execute(w, ids)
				if err != nil {
					http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
					return
				}
			}
			return
		}
		w.Header().Add("Allowed", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
