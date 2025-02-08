package areas

import (
	"basement/main/internal/server"
	"net/http"
)

func AreasHandler(db AreaDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listPage(db).ServeHTTP(w, r)
			break

		case http.MethodPost:
			server.WriteNotImplementedWarning("todo", w, r)
			break

		case http.MethodDelete:
			deleteAreas(w, r, db)
			return

		case http.MethodPut:
			server.WriteNotImplementedWarning("todo", w, r)
			break

		default:
			// Other methods are not allowed.
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			server.WriteFprint(w, "Method:'"+r.Method+"' not allowed")
		}
	}
}

func deleteAreas(w http.ResponseWriter, r *http.Request, db AreaDatabase) {
	server.DeleteThingsFromList(w, r, db.DeleteArea, listPage(db))
}
