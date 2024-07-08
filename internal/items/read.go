package items

import (
	"basement/main/internal/database"
	"fmt"
	"io"
	"net/http"
)

// ResponseWriter should implement a function to write a template response or normal response.
//
// Example:
//
//	func(w io.Writer, data any) {
//		// templates
//		templates.Render(w, "items-container", data)
//		// Fprint
//		fmt.Fprint(w, data)
//	})
type ResponseWriter func(w io.Writer, data any)

// ReadItemHandler returns a single item.
//
// Accepts "/item?id=" and "/item/id"
func ReadItemHandler(db *database.JsonDB, responseWriter ResponseWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := r.FormValue("id")
			if id == "" {
				id = r.PathValue("id")
			}
			data := db.Items[id]
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
// Accepts "/items?query=id" to only return item IDs.
func ReadItemsHandler(db *database.JsonDB, responseWriter ResponseWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := r.FormValue("query")
			switch id {
			// return all item IDs
			case "id":
				responseWriter(w, Keys(db.Items))
			// return all items
			default:
				responseWriter(w, db.Items)
			}
			return
		}
		w.Header().Add("Allowed", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
//
// # In Go 1.21 a part of the maps package has been moved into the standard library, but not maps.Keys
//
// https://stackoverflow.com/a/69889828
// https://cs.opensource.google/go/x/exp/+/39d4317d:maps/maps.go;l=10
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
