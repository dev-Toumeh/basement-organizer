package routes

import (
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

type BoxDatabase interface {
	// CreateBox returns id of box if successful, otherwise error.
	CreateBox() (string, error)
	Box(id string) (items.Box, error)
	// // BoxIDs returns IDs of all boxes.
	// BoxIDs() ([]string, error)
	// // MoveBox moves box with id1 into box with id2.
	// MoveBox(id1 string, id2 string) error
}

func BoxHandler(rw items.ResponseWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			db := r.Context().Value("db").(BoxDatabase)
			id := r.FormValue("id")
			_, err := db.Box(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, fmt.Errorf("Can't find box with id: %v. %w", id, err))
				return
			}
			// ids, _ := db.ItemIDs()
			// for _, id := range ids {
			// 	item, _ := db.Item(id)
			// 	b.Items = append(b.Items, &item)
			// }
			// b.Description = fmt.Sprintf("This box has %v items", len(ids))
			// rw(w, b)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Method:'", r.Method, "' not implemented")
			break
		case http.MethodPost:
			db := r.Context().Value("db").(BoxDatabase)
			_, err := db.CreateBox()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, fmt.Errorf("Can't create new box. %w", err))
				return
			}
			w.WriteHeader(http.StatusOK)
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
