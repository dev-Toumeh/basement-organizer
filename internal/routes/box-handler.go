package routes

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

// boxDatabaseSuccess never returns errors.
type mockBoxDB struct{}

func (db *mockBoxDB) CreateBox() (string, error) {
	return "fa2e3db6-fcf8-49c6-ac9c-54ce5855bf0b", nil
}

func (db *mockBoxDB) Box(id string) (items.Box, error) {
	return items.Box{}, nil
}
func registerBoxRoutes(db BoxDatabase) {
	http.HandleFunc("/api/v2/box", BoxHandler(FprintWriteFunc, db))
	http.HandleFunc("/api/v2/box/{id}", BoxHandler(FprintWriteFunc, db))
	http.HandleFunc("/box", BoxHandler(func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_BOX, data)
	}, db))
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

func BoxHandler(rw items.ResponseWriter, db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			const errorMsg = "Can't get box"

			id := validID(w, r, errorMsg)
			if id == "" {
				return
			}

			_, err := db.Box(id)
			if err != nil {
				logg.Debugf("Can't find box with id: '%v'. %v", id, err)
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, errorMsg, id)
				return
			}
			// ids, _ := db.ItemIDs()
			// for _, id := range ids {
			// 	item, _ := db.Item(id)
			// 	b.Items = append(b.Items, &item)
			// }
			// b.Description = fmt.Sprintf("This box has %v items", len(ids))
			// rw(w, b)

			// @TODO: Implement
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprint(w, "Method:'", r.Method, "' not implemented")
			break
		case http.MethodPost:
			_, err := db.CreateBox()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, fmt.Errorf("Can't create new box. %w", err))
				return
			}
			// @TODO: Implement
			w.WriteHeader(http.StatusNotImplemented)
			break
		case http.MethodDelete:
			id := validID(w, r, "Can't delete box")
			if id == "" {
				return
			}
			// @TODO: Implement
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprint(w, "Method:'", r.Method, "' not implemented")
			break
		case http.MethodPut:
			id := validID(w, r, "Can't update box")
			if id == "" {
				return
			}
			// @TODO: Implement
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprint(w, "Method:'", r.Method, "' not implemented")
			break
		default:
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		}
	}
}

// validID returns valid id string or if errors occurs
// writes correct response header status code with errorMessage and returns empty string.
func validID(w http.ResponseWriter, r *http.Request, errorMessage string) string {
	id := r.FormValue("id")
	logg.Debugf("Query param id: '%v'.", id)
	if id == "" {
		id = r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			logg.Debug("Empty id")
			fmt.Fprintf(w, `%s ID="%v"`, errorMessage, id)
			return ""
		}
		logg.Debugf("path value id: '%v'.", id)
	}

	_, err := uuid.FromString(id)
	if err != nil {
		logg.Debugf("Wrong id: '%v'. %v", id, err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `%s ID="%v"`, errorMessage, id)
		return ""
	}
	return id
}

