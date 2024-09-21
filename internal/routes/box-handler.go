package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"maps"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type BoxDatabase interface {
	CreateBox(newBox *items.Box) (uuid.UUID, error)
	MoveBox(id1 uuid.UUID, id2 uuid.UUID) error
	UpdateBox(box items.Box) error
	DeleteBox(boxId uuid.UUID) error
	BoxById(id uuid.UUID) (items.Box, error)
	BoxIDs() ([]string, error) // @TODO: Change string to uuid.UUID
}

func registerBoxRoutes(db BoxDatabase) {
	// Box templates
	http.HandleFunc("/box", boxHandler(db))
	http.HandleFunc("/box/{id}/move", boxPageMove(db))
	// Box api
	http.HandleFunc("/api/v1/box", boxHandler(db))
	http.HandleFunc("/api/v1/box/{id}", boxHandler(db))
	http.HandleFunc("/api/v1/box/{id}/move", moveBox(db))
	// Boxes templates
	http.HandleFunc("/boxes", boxesPage)
	http.HandleFunc("/boxes/{id}", boxDetailsPage(db))
	http.HandleFunc("/boxes/move", boxesPageMove(db))
	http.HandleFunc("/boxes-list", boxesHandler(db))
	// Boxes api
	http.HandleFunc("/api/v1/boxes", boxesHandler(db))
	http.HandleFunc("/api/v1/boxes/move", moveBoxes(db))
}

func boxHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			const errMsgForUser = "Can't find box"

			id := ValidID(w, r, errMsgForUser)
			if id.IsNil() {
				return
			}

			box, err := db.BoxById(id)
			if err != nil {
				server.WriteNotFoundError(errMsgForUser, err, w, r)
				return
			}

			// Use API data writer
			if !wantsTemplateData(r) {
				server.WriteJSON(w, box)
				return
			}

			// Template writer
			renderBoxTemplate(&box, w, r)
			break

		case http.MethodPost:
			createBox(w, r, db)
			break

		case http.MethodDelete:
			deleteBox(w, r, db)
			return

		case http.MethodPut:
			updateBox(w, r, db)
			break

		default:
			// Other methods are not allowed.
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		}
	}
}

func boxesHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ids, err := db.BoxIDs()
			if err != nil {
				server.WriteNotFoundError("Can't find boxes", err, w, r)
			}

			if wantsTemplateData(r) {
				renderBoxesListTemplate(w, r, db, ids)
			} else {
				server.WriteJSON(w, ids)
			}
			break

		case http.MethodPut:
			server.WriteNotImplementedWarning("Multiple boxes edit?", w, r)
			break

		case http.MethodDelete:
			deleteBoxes(w, r, db)
			break

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			break
		}
	}
}

func boxesPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	user, _ := auth.UserSessionData(r)
	data := templates.NewPageTemplate()
	data.Title = "Boxes"
	data.Authenticated = authenticated
	data.User = user

	server.MustRender(w, r, templates.TEMPLATE_BOXES_PAGE, data.Map())
}

func boxDetailsPage(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		id := ValidID(w, r, "no box")
		if id.IsNil() {
			return
		}
		logg.Debug(id)

		notFound := false
		box, err := db.BoxById(id)
		if err != nil {
			notFound = true
		}
		box.Id = id
		data := items.BoxPageTemplateData()
		data.Box = &box

		data.Title = fmt.Sprintf("Box - %s", box.Label)
		data.Authenticated = authenticated
		data.User = user
		data.NotFound = notFound
		nd := data.Map()
		maps.Copy(nd, map[string]any{"Boxes": &box.InnerBoxes})
		searchInput := items.NewSearchInputTemplate()
		searchInput.SearchInputLabel = "Search boxes"
		searchInput.SearchInputHxTarget = "#box-list-body"
		searchInput.SearchInputHxPost = "/api/v1/implement-me"
		maps.Copy(nd, searchInput.Map())

		server.MustRender(w, r, templates.TEMPLATE_BOX_DETAILS_PAGE, nd)
	}
}

func createBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	box := items.NewBox()
	logg.Debug("create box: ", box)
	id, err := db.CreateBox(&box)
	if err != nil {
		server.WriteNotFoundError("error while creating the box", err, w, r)
		return
	}
	if wantsTemplateData(r) {
		box, err := db.BoxById(id)
		logg.Debug(box)
		if err != nil {
			server.WriteNotFoundError("error while fetching the box based on Id", err, w, r)
			return
		}
		templates.Render(w, templates.TEMPLATE_BOX_LIST_ITEM, box)
	} else {
		server.WriteJSON(w, id)
	}
}

func updateBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	errMsgForUser := "Can't update box."
	id := ValidID(w, r, errMsgForUser)
	if id.IsNil() {
		return
	}

	box := boxFromPostFormValue(id, r)
	err := db.UpdateBox(box)
	if err != nil {
		server.WriteNotFoundError(errMsgForUser, err, w, r)
		return
	}
	if wantsTemplateData(r) {
		boxTemplate := items.BoxTemplateData{Box: &box, Edit: false}
		err := server.RenderWithSuccessNotification(w, r, templates.TEMPLATE_BOX_DETAILS, boxTemplate.Map(), fmt.Sprintf("Updated box: %v", boxTemplate.Label))
		if err != nil {
			server.WriteInternalServerError(errMsgForUser, err, w, r)
			return
		}
	} else {
		server.WriteJSON(w, box)
	}
	logg.Debug("Updated Box: ", box)
}

func deleteBoxes(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
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
				logg.Errorf(fmt.Sprintf("%s: Malformed uuid \"%s\"", errMsgForUser, k), err)
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
			logg.Errorf(fmt.Sprintf("%v: %v", errMsgForUser, deleteId), err)
		} else {
			logg.Debug("Box deleted: ", deleteId)
		}
	}
	if errOccurred {
		errIds := strings.Join(deleteErrorIds, ",")
		server.TriggerErrorNotification(w, errMsgForUser+errIds)
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
	fmt.Fprint(w, nil)
	w.WriteHeader(http.StatusOK)
}

// deleteBox deletes a single box.
func deleteBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	errMsgForUser := "Can't delete box"
	id := ValidID(w, r, errMsgForUser)
	if id.IsNil() {
		return
	}
	err := db.DeleteBox(id)
	if err != nil {
		server.WriteNotFoundError(errMsgForUser, err, w, r)
		return
	}
	server.RedirectWithSuccessNotification(w, "/boxes", fmt.Sprintf("%s deleted", id))
}

// @TODO: Implement move box.
func boxPageMove(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move single box page", w, r)
		// err := db.MoveBox(uuid.FromStringOrNil("5cca42c2-5f1b-45e7-b2d2-175a0ff99b61"), uuid.FromStringOrNil("a88a1ebd-0551-4008-bdda-9677d375c7eb"))

		// if err != nil {
		// 	writeNotFoundError(errMsgForUser, err, w, r)
		// }
		// w.WriteHeader(http.StatusNotImplemented)
	}
}

func boxesPageMove(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move multiple boxes page", w, r)
	}
}

// @TODO: Implement moveBox.
func moveBox(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move single box", w, r)
	}
}

// @TODO: Implement move boxes.
func moveBoxes(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move multiple boxes", w, r)
		// w.WriteHeader(http.StatusNotImplemented)
	}
}

// boxFromPostFormValue returns items.Box without references to inner boxes, outer box and items.
func boxFromPostFormValue(id uuid.UUID, r *http.Request) items.Box {
	box := items.Box{}
	box.Id = id
	box.Label = r.PostFormValue("label")
	box.Description = r.PostFormValue("description")
	box.Picture = items.ParsePicture(r)
	// box.QRcode = r.PostFormValue("qrcode")
	return box
}

// ValidID returns valid uuid from request and handles errors.
// Check for uuid.Nil! If error occurs return will be uuid.Nil.
func ValidID(w http.ResponseWriter, r *http.Request, errorMessage string) uuid.UUID {
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

func renderBoxTemplate(box *items.Box, w http.ResponseWriter, r *http.Request) {
	editParam := r.FormValue("edit")
	edit := false
	if editParam == "true" {
		edit = true
	}
	b := items.BoxTemplateData{Box: box, Edit: edit}
	server.MustRender(w, r, templates.TEMPLATE_BOX_DETAILS, b.Map())
}

func renderBoxesListTemplate(w http.ResponseWriter, r *http.Request, db BoxDatabase, ids []string) {
	var boxes []*items.Box
	for _, id := range ids {
		box, _ := db.BoxById(uuid.Must(uuid.FromString(id)))
		boxes = append(boxes, &box)
		// items.RenderBoxListItem(w, &box)
	}
	// items.RenderBoxList(w, boxes)
	searchInput := items.NewSearchInputTemplate()
	// logg.Debugf("searchInput %v", searchInput)
	logg.Debugf("searchInput %v", searchInput.Map())
	searchInput.SearchInputLabel = "Search boxes"
	searchInput.SearchInputHxTarget = "#box-list-body"
	searchInput.SearchInputHxPost = "/api/v1/implement-me"
	maps := []templates.Mapable{searchInput, items.BoxListTemplateData{Boxes: boxes}}
	templates.RenderMaps(w, templates.TEMPLATE_BOX_LIST, maps)
	return
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
// 		return logg.Errorf("UpdateBox(): %w", err)
// 	}
// 	db.Boxes[oldBox.Id.String()] = &box
// 	return nil
// }
//
// func (db *sampleBoxDB) DeleteBox(id uuid.UUID) error {
// 	_, err := db.BoxById(id)
// 	if err != nil {
// 		return logg.Errorf("DeleteBox(): %w", err)
// 	}
// 	delete(db.Boxes, id.String())
// 	return nil
// }
