package routes

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
)

func registerBoxRoutes(db BoxDatabase) {
	http.HandleFunc("/api/v2/box", BoxHandler(WriteJSON, db))
	http.HandleFunc("/api/v2/box/{id}", BoxHandler(WriteJSON, db))
	http.HandleFunc("/box", BoxHandler(WriteBoxTemplate, db))
}

func WriteBoxTemplate(w io.Writer, data any) {
	templates.Render(w, templates.TEMPLATE_BOX, data)
}

func WriteFprint(w io.Writer, data any) {
	fmt.Fprint(w, data)
}

func WriteJSON(w io.Writer, data any) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

type BoxDatabase interface {
	// CreateBox returns id of box if successful, otherwise error.
	CreateBox() (string, error)
	Box(id string) (items.Box, error)
	BoxByField(field string, value string) (*items.Box, error)
	MoveBox(id1 uuid.UUID, id2 uuid.UUID) error
	BoxIDs() ([]string, error)
	BoxExist(field string, value string) bool
	CreateNewBox(newBox *items.Box) (uuid.UUID, error)
	ErrorExist() error
	UpdateBox(box *items.Box) error
	DeleteBox(id string) error
}
}

func BoxHandler(writeData items.DataWriteFunc, db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			const errorMsg = "Can't get box"

			id := validID(w, r, errorMsg)
			if id == "" {
				return
			}

			box, err := db.Box(id)
			if err != nil {
				logg.Debugf("Can't find box with id: '%v'. %v", id, err)
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, errorMsg, id)
				return
			}

			// Use API data writer
			if !wantsTemplateData(r) {
				writeData(w, box)
				return
			}

			// Template writer
			editParam := r.FormValue("edit")
			edit := false
			if editParam == "true" {
				edit = true
			}
			b := struct {
				items.Box
				Edit bool
			}{box, edit}

			writeData(w, b)
			break

		case http.MethodPost:
			id, err := db.CreateBox()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, fmt.Errorf("Can't create new box. %w", err))
				return
			}
			if wantsTemplateData(r) {
				templates.Render(w, "box-list-item", struct{ Id string }{Id: id})
			} else {
				writeData(w, id)
			}
			break

		case http.MethodDelete:
			id := validID(w, r, "Can't delete box")
			if id == "" {
				return
			}
			// @TODO: Implement
			// w.WriteHeader(http.StatusNotImplemented)
			// fmt.Fprint(w, "Method:'", r.Method, "' not implemented")
			break

		case http.MethodPut:
			id := validID(w, r, "Can't update box")
			if id == "" {
				return
			}

			label := r.FormValue("label")
			// picture := r.FormValue("picture")
			// qrcode := r.FormValue("qrcode")
			description := r.FormValue("description")

			if wantsTemplateData(r) {
				b := items.BoxTemplateData{}
				b.Id = uuid.Must(uuid.FromString(id))
				b.Label = label
				// b.Picture = picture
				// b.QRcode = qrcode
				b.Description = description
				logg.Debug(b)
				writeData(w, b)
			} else {
				// b := items.Box{}
				// b.Label = label
				// b.Picture = picture
				// b.QRcode = qrcode
				// b.Description = description
				// writeData(w, b)
				w.WriteHeader(http.StatusOK)
			}

			// items := r.FormValue("items")
			// innerboxes := r.FormValue("innerboxes")
			// outerbox := r.FormValue("outerbox")
			// @TODO: Implement
			// w.WriteHeader(http.StatusNotImplemented)
			// fmt.Fprint(w, "Method:'", r.Method, "' not implemented")
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
// validID returns valid id string.
// Check for empty string! If error occurs return will be "".
// But the error is already handled.
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

// wantsTemplateData checks if current request requires template data.
func wantsTemplateData(r *http.Request) bool {
	return !strings.Contains(r.URL.Path, "/api/")
}
