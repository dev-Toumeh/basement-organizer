package shelves

import (
	"net/http"

	"basement/main/internal/server"
)

// Handles read, create, update, and delete for multiple shelves.
func ShelvesHandler(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			PageTemplate(db).ServeHTTP(w, r)
			return
		case http.MethodDelete:
			server.DeleteThingsFromList(w, r, db.DeleteShelf2, PageTemplate(db))
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
