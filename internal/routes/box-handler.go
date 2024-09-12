package routes

import (
	"basement/main/internal/auth"
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
	http.HandleFunc("/boxes", boxesPage)
	http.HandleFunc("/boxes-list", BoxesHandler(WriteJSON, db))
}

func boxesPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	user, _ := auth.UserSessionData(r)
	data := templates.NewPageTemplate()
	data.Title = "Boxes"
	data.Authenticated = authenticated
	data.User = user

	triggerAllServerNotifications(w)
	MustRender(w, r, templates.TEMPLATE_BOXES_PAGE, data)
}

func WriteBoxTemplate(w io.Writer, data any) {
	boxTemplate, ok := data.(items.BoxTemplateData)
	errMessageForUser := "Something went wrong"
	if !ok {
		ww, ok := w.(http.ResponseWriter)
		if !ok {
			logg.Fatal("Can't write box template, writer is not http.ResponseWriter")
		}
		ww.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(ww, errMessageForUser)
		logg.Info(errMessageForUser)
	}
	err := templates.Render(w, templates.TEMPLATE_BOX, boxTemplate)
	if err != nil {
		ww, ok := w.(http.ResponseWriter)
		if !ok {
			logg.Fatal("Can't write box template, writer is not http.ResponseWriter")
		}
		ww.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(ww, errMessageForUser)
		logg.Info(fmt.Errorf("Can't render TEMPLATE_BOX. %w", err))
		return
	}
}

func WriteFprint(w io.Writer, data any) {
	fmt.Fprint(w, data)
}

func WriteJSON(w io.Writer, data any) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

type BoxDatabase interface {
	CreateBox(newBox *items.Box) (uuid.UUID, error)
	MoveBox(id1 uuid.UUID, id2 uuid.UUID) error
	UpdateBox(box items.Box) error
	DeleteBox(boxId uuid.UUID) error
	BoxById(id uuid.UUID) (items.Box, error)
	BoxIDs() ([]string, error)                // @TODO: Change string to uuid.UUID
	BoxExist(field string, value string) bool // @TODO: Change string to uuid.UUID
	ErrorExist() error
}

func BoxesHandler(writeData items.DataWriteFunc, db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ids, err := db.BoxIDs()
			if err != nil {
				writeNotFoundError("Can't find boxes", err, w, r)
			}
			if wantsTemplateData(r) {
				for _, id := range ids {
					box, _ := db.BoxById(uuid.Must(uuid.FromString(id)))
					templates.Render(w, templates.TEMPLATE_BOX_LIST_ITEM, box)
				}
				return
			}
			writeData(w, ids)

		case http.MethodPut:
			// @TODO: Implement move boxes.
			w.WriteHeader(http.StatusNotImplemented)
			return

		case http.MethodDelete:
			errMsgForUser := "Can't delete boxes"
			r.ParseForm()
			toDelete := make([]uuid.UUID, 0)
			for k, v := range r.Form {
				logg.Debugf("k: %v, v:%v", k, v)
				if strings.Contains(k, "delete:") {
					ids := strings.Split(k, "delete:")
					if len(ids) != 2 {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprint(w, errMsgForUser)
						logg.Debugf("Wrong delete key value pair: '%v'\n\t%s", k, errMsgForUser)
						return
					}
					id, err := uuid.FromString(ids[1])
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprint(w, errMsgForUser)
						logg.Debug(fmt.Errorf("%v\n\tMalformed uuid: %v\n\t%w", errMsgForUser, k, err))
						return
					}
					toDelete = append(toDelete, id)
				}
			}
			deleteErrorIds := []string{}
			var err error
			errOccurred := false
			for _, deleteId := range toDelete {
				err = nil
				err = db.DeleteBox(deleteId)
				if err != nil {
					errOccurred = true
					deleteErrorIds = append(deleteErrorIds, deleteId.String())
					logg.Err(fmt.Errorf("%v: %v\n\t%w", errMsgForUser, deleteId, err))
				} else {
					logg.Debug("Box deleted: ", deleteId)
				}
			}
			if errOccurred {
				errIds := strings.Join(deleteErrorIds, ",")
				TriggerErrorNotification(w, errMsgForUser+errIds)
				// @TODO: Update partial table, even if error happens.
				return
			}

			if wantsTemplateData(r) {
				newids, _ := db.BoxIDs()
				for _, id := range newids {
					box, _ := db.BoxById(uuid.Must(uuid.FromString(id)))
					templates.Render(w, templates.TEMPLATE_BOX_LIST_ITEM, box)
				}
				for _, id := range toDelete {
					templates.RenderSuccessNotification(w, "Box deleted: "+id.String())
				}
				return
			}
			writeData(w, nil)
			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func BoxHandler(writeData items.DataWriteFunc, db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			const errMsgForUser = "Can't find box"

			id := validID(w, r, errMsgForUser)
			if id.IsNil() {
				return
			}

			box, err := db.BoxById(id)
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
			b := items.BoxTemplateData{Box: &box, Edit: edit}

			WriteBoxTemplate(w, b)
			break

		case http.MethodPost:
			box := items.NewBox()
			id, err := db.CreateBox(&box)
			if err != nil {
				writeNotFoundError("error while creating the box", err, w, r)
				return
			}
			if wantsTemplateData(r) {
				box, err := db.BoxById(id)
				logg.Debug(box)
				if err != nil {
					writeNotFoundError("error while fetching the box based on Id", err, w, r)
					return
				}
				templates.Render(w, templates.TEMPLATE_BOX_LIST_ITEM, box)
			} else {
				writeData(w, id)
			}
			break

		case http.MethodDelete:
			errMsgForUser := "Can't delete box"
			id := validID(w, r, errMsgForUser)
			if id.IsNil() {
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
			if id.IsNil() {
				return
			}

			box := boxFromPostFormValue(id, r)
			err := db.UpdateBox(box)
			if err != nil {
				writeNotFoundError(errMsgForUser, err, w, r)
				return
			}
			if wantsTemplateData(r) {
				boxTemplate := items.BoxTemplateData{Box: &box, Edit: false}
				RenderWithSuccessNotification(w, templates.TEMPLATE_BOX, boxTemplate, fmt.Sprintf("Updated box: %v", boxTemplate.Label))
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

// Render applies data to a defined template and writes result back to the writer.
func RenderWithSuccessNotification(w http.ResponseWriter, name string, data any, successMessage string) error {
	TriggerSuccessNotification(w, successMessage)
	err := templates.Render(w, name, data)
	if err != nil {
		fmt.Fprintln(w, err)
		return fmt.Errorf("RenderWithSuccessNotification() error: %w", err)
	}
	return nil
}

// writeNotFoundError sets not found status code 404, logs error and writes error message to client.
func writeNotFoundError(message string, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, message)
	logg.Info(fmt.Errorf("%s:\n\t%w", message, err))
}

// boxFromPostFormValue returns items.Box without references to inner boxes, outer box and items.
func boxFromPostFormValue(id uuid.UUID, r *http.Request) items.Box {
	box := items.Box{}
	box.Id = id
	box.Label = r.PostFormValue("label")
	box.Description = r.PostFormValue("description")
	// box.Picture = r.PostFormValue("picture")
	// box.QRcode = r.PostFormValue("qrcode")
	return box
}

// validID returns valid uuid from request and handles errors.
// Check for uuid.Nil! If error occurs return will be uuid.Nil.
func validID(w http.ResponseWriter, r *http.Request, errorMessage string) uuid.UUID {
	id := r.FormValue("id")
	logg.Debugf("Query param id: '%v'.", id)
	if id == "" {
		id = r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			logg.Debug("Empty id")
			fmt.Fprintf(w, `%s ID="%v"`, errorMessage, id)
			return uuid.Nil
		}
		logg.Debugf("path value id: '%v'.", id)
	}

	newId, err := uuid.FromString(id)
	if err != nil {
		logg.Debugf("Wrong id: '%v'. %v", id, err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `%s ID="%v"`, errorMessage, id)
		return uuid.Nil
	}
	return newId
}

// wantsTemplateData checks if current request requires template data.
func wantsTemplateData(r *http.Request) bool {
	return !strings.Contains(r.URL.Path, "/api/")
}

// // sampleBoxDB to test request handler.
// type sampleBoxDB struct {
// 	Boxes map[string]*items.Box
// }
//
// func newSampleBoxDB() *sampleBoxDB {
// 	db := sampleBoxDB{Boxes: make(map[string]*items.Box, 100)}
// 	for i := range 10 {
// 		box := items.NewBox()
// 		box.Label = fmt.Sprintf("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa %d", i)
// 		db.CreateBox(&box)
// 	}
// 	return &db
// }
//
// // func (db *sampleBoxDB) CreateBox() (uuid.UUID, error) {
// // 	box := items.NewBox()
// // 	db.Boxes[box.Id.String()] = &box
// // 	logg.Debug(db.Boxes)
// // 	return box.Id, nil
// // }
//
// func (db *sampleBoxDB) CreateBox(box *items.Box) (uuid.UUID, error) {
// 	// box := items.NewBox()
// 	db.Boxes[box.Id.String()] = box
// 	logg.Debug(db.Boxes)
// 	return box.Id, nil
// }
//
// func (db *sampleBoxDB) BoxById(id uuid.UUID) (items.Box, error) {
// 	box, ok := db.Boxes[id.String()]
// 	if !ok {
// 		// logg.Debug("BoxByID: ",db.Boxes)
// 		return items.Box{}, errors.New("ID " + id.String() + " doesn't exist")
// 	}
// 	return *box, nil
// }
//
// // func (db *sampleBoxDB) BoxIDs() ([]uuid.UUID, error) {
// func (db *sampleBoxDB) BoxIDs() ([]string, error) {
// 	// ids := make([]uuid.UUID, len(db.Boxes))
// 	ids := make([]string, len(db.Boxes))
// 	i := 0
// 	for _, v := range db.Boxes {
// 		// ids[i] = v.Id
// 		ids[i] = v.Id.String()
// 		i++
// 	}
// 	slices.Sort(ids)
// 	return ids, nil
// }
//
// func (db *sampleBoxDB) UpdateBox(box items.Box) error {
// 	oldBox, err := db.BoxById(box.Id)
// 	if err != nil {
// 		return fmt.Errorf("UpdateBox(): %w", err)
// 	}
// 	db.Boxes[oldBox.Id.String()] = &box
// 	return nil
// }
//
// func (db *sampleBoxDB) DeleteBox(id uuid.UUID) error {
// 	_, err := db.BoxById(id)
// 	if err != nil {
// 		return fmt.Errorf("DeleteBox(): %w", err)
// 	}
// 	delete(db.Boxes, id.String())
// 	return nil
// }
