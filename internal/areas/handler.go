package areas

import (
	"basement/main/internal/server"
	"fmt"
	"net/http"
)

func AreaHandler(db AreaDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.WriteNotImplementedWarning("todo", w, r)
			break

		case http.MethodPost:
			server.WriteNotImplementedWarning("todo", w, r)
			break

		case http.MethodDelete:
			server.WriteNotImplementedWarning("todo", w, r)
			return

		case http.MethodPut:
			server.WriteNotImplementedWarning("todo", w, r)
			break

		default:
			// Other methods are not allowed.
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		}
	}
}

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
			server.WriteNotImplementedWarning("todo", w, r)
			return

		case http.MethodPut:
			server.WriteNotImplementedWarning("todo", w, r)
			break

		default:
			// Other methods are not allowed.
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		}
	}
}
