package routes

import (
	"basement/main/internal/database"
	"basement/main/internal/items"
	"basement/main/internal/templates"
	"fmt"
	"io"
	"net/http"
)

func registerBoxRoutes() {
	http.HandleFunc("/api/v2/box", BoxHandler(FprintWriteFunc))
	http.HandleFunc("/api/v2/box/{id}", BoxHandler(FprintWriteFunc))
	http.HandleFunc("/box", BoxHandler(func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_BOX, data)
	}))
}

func FprintWriteFunc(w io.Writer, data any) { fmt.Fprint(w, data) }

func BoxHandler(rw items.ResponseWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			b := items.NewBox()
			db := r.Context().Value("db").(*database.DB)
			ids, _ := db.ItemIDs()
			for _, id := range ids {
				item, _ := db.Item(id)
				b.Items = append(b.Items, &item)
			}
			b.Description = fmt.Sprintf("This box has %v items", len(ids))
			rw(w, b)
			break
		case http.MethodPost:
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "Method:'", r.Method, "' not implemented")
			break
		case http.MethodDelete:
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "Method:'", r.Method, "' not implemented")
			break
		case http.MethodPut:
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "Method:'", r.Method, "' not implemented")
			break
		default:
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		}
	}
}
