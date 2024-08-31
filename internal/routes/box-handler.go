package routes

import (
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"encoding/json"
	"errors"
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
			const errMsgForUser = "Can't find box"

			id := validID(w, r, errMsgForUser)
			if id == "" {
				return
			}

			box, err := db.Box(id)
			if err != nil {
				writeNotFoundError(errMsgForUser, err, w, r)
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
			errMsgForUser := "Can't create new box"
			id, err := db.CreateBox()
			if err != nil {
				writeNotFoundError(errMsgForUser, err, w, r)
				return
			}
			if wantsTemplateData(r) {
				templates.Render(w, "box-list-item", struct{ Id string }{Id: id})
			} else {
				writeData(w, id)
			}
			break

		case http.MethodDelete:
			errMsgForUser := "Can't delete box"
			id := validID(w, r, errMsgForUser)
			if id == "" {
				return
			}
			err := db.DeleteBox(id)
			if err != nil {
				writeNotFoundError(errMsgForUser, err, w, r)
				return
			}
			break

		case http.MethodPut:
			errMsgForUser := "Can't update box."
			id := validID(w, r, errMsgForUser)
			if id == "" {
				return
			}

			box := boxFromPostFormValue(id, r)
			err := db.UpdateBox(&box)
			if err != nil {
				writeNotFoundError(errMsgForUser, err, w, r)
				return
			}
			if wantsTemplateData(r) {
				boxTemplate := items.BoxTemplateData{&box, false}
				writeData(w, boxTemplate)
			} else {
				writeData(w, box)
			}
			logg.Debug("Updated Box: ", box)
			break

		default:
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		}
	}
}

// writeNotFoundError sets not found status code 404, logs error and writes error message to client.
func writeNotFoundError(message string, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, message)
	logg.Info(fmt.Errorf("%s: %w", message, err))
}

// boxFromPostFormValue returns items.Box without references to inner boxes, outer box and items.
func boxFromPostFormValue(id string, r *http.Request) items.Box {
	box := items.Box{}
	box.Id = uuid.Must(uuid.FromString(id))
	box.Label = r.PostFormValue("label")
	box.Description = r.PostFormValue("description")
	// box.Picture = r.PostFormValue("picture")
	// box.QRcode = r.PostFormValue("qrcode")
	return box
}

// validID returns valid id string.
// Check for empty string! If error occurs return will be "".
// But the error is already handled.
func validID(w http.ResponseWriter, r *http.Request, errorMessage string) string {
	id := r.FormValue("id")
	logg.Debugf("Query param id: '%v'.", id)
	if id == "" {
		id = r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
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
