package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/boxes"
	"basement/main/internal/common"
	"basement/main/internal/database"
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type BoxDatabase interface {
	CreateBox(newBox *boxes.Box) (uuid.UUID, error)
	MoveBoxToBox(id1 uuid.UUID, id2 uuid.UUID) error
	UpdateBox(box boxes.Box) error
	DeleteBox(boxId uuid.UUID) error
	BoxById(id uuid.UUID) (boxes.Box, error)
	BoxIDs() ([]string, error) // @TODO: Change string to uuid.UUID
	BoxListRows(query string, limit int, page int) ([]common.ListRow, error)
	BoxListRowByID(id uuid.UUID) (common.ListRow, error)
	BoxListCounter(searchString string) (count int, err error)
}

func registerBoxRoutes(db *database.DB) {
	// Box templates
	// http.HandleFunc("/box", boxHandler(db))
	Handle("/box", boxHandler(db))
	Handle("/box/{id}/move", boxPageMove(db))
	// Box api
	Handle("/api/v1/box", boxHandler(db))
	Handle("/api/v1/box/{id}", boxHandler(db))
	Handle("/api/v1/box/{id}/move", moveBox(db))
	// Boxes templates
	Handle("/boxes", boxesPage(db))
	Handle("/boxes/{id}", boxDetailsPage(db))
	Handle("/boxes/move", boxesPageMove(db))
	Handle("/boxes-list", boxesHandler(db))
	// Boxes api
	Handle("/api/v1/boxes", boxesHandler(db))
	Handle("/api/v1/boxes/moveto/box/{id}", moveBoxes(db))
	Handle("/get-boxes-move-page", getMoveBoxesPage(db))
}

// boxHandler handles read, create, update and delete for single box.
func boxHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			const errMsgForUser = "Can't find box"

			id := server.ValidID(w, r, errMsgForUser)
			if id.IsNil() {
				return
			}

			box, err := db.BoxById(id)
			if err != nil {
				server.WriteNotFoundError(errMsgForUser, err, w, r)
				return
			}

			// Use API data writer
			if !server.WantsTemplateData(r) {
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

// boxesHandler handles read and delete for multiple boxes.
func getMoveBoxesPage(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Uses POST so client can send long list
		// of IDs of boxes inside PostForm body
		if r.Method != http.MethodPost {
			w.Header().Add("Allow", http.MethodPost)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
			return
		}

		boxs, err := db.BoxListRows("", 100, 1)
		if err != nil {
			server.WriteInternalServerError("cant query boxes", logg.Errorf("%w", err), w, r)
			return
		}

		a := boxes.BoxListTemplateData{Boxes: boxs}
		d := a.Map()
		d["Move"] = true
		for i := range d["Boxes"].([]map[string]any) {
			d["Boxes"].([]map[string]any)[i]["Move"] = true

		}

		// parse selected boxes
		r.ParseForm()
		toMove, err := parseIDsFromFormWithKey(r.Form, "move")
		if err != nil {
			server.WriteInternalServerError(fmt.Sprintf("can't move boxes %v", toMove), err, w, r)
			return
		}
		moveData := make([]map[string]string, len(toMove))
		for i, id := range toMove {
			moveData[i] = map[string]string{
				"Key":   "id-to-be-moved",
				"Value": id.String(),
			}
		}
		d["Data"] = moveData
		server.MustRender(w, r, templates.TEMPLATE_BOX_LIST, d)
	}
}

// boxesHandler handles read and delete for multiple boxes.
func boxesHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			if !server.WantsTemplateData(r) {
				boxes, err := db.BoxListRows("", 5, 1)
				if err != nil {
					server.WriteInternalServerError("cant query boxes", logg.Errorf("%w", err), w, r)
					return
				}
				server.WriteJSON(w, boxes)
				return
			}

			boxs, err := db.BoxListRows("", 100, 1)
			if err != nil {
				server.WriteNotFoundError("Can't find boxes", err, w, r)
				return
			}
			if server.WantsTemplateData(r) {
				a := boxes.BoxListTemplateData{Boxes: boxs}
				d := a.Map()
				d["Move"] = true
				for i := range d["Boxes"].([]map[string]any) {
					d["Boxes"].([]map[string]any)[i]["Move"] = true

				}
				server.MustRender(w, r, templates.TEMPLATE_BOX_LIST, d)
			} else {
				server.WriteJSON(w, boxs)
			}
			break

		case http.MethodPost:
			query := r.PostFormValue("query")
			logg.Debugf("search query: %s", query)
			boxes, err := db.BoxListRows(query, 5, 1)
			if err != nil {
				server.WriteInternalServerError("cant query boxes", err, w, r)
				return
			}
			err = renderBoxesListTemplate(w, r, db, boxes, query)
			if err != nil {
				server.WriteInternalServerError("cant render boxlist", err, w, r)
				return
			}

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

// boxesPage shows a page with a box list.
func boxesPage(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		// page template
		page := templates.NewPageTemplate()
		page.Title = "Boxes"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// list template
		listTmpl := common.ListTemplate{
			FormHXGet: "/boxes",
			RowHXGet:  "/boxes",
		}

		searchString := common.SearchString(r)

		// search-input template
		listTmpl.SearchInput = true
		listTmpl.SearchInputLabel = "Search boxes"
		listTmpl.SearchInputValue = searchString

		// @TODO: Implement move page
		// move := false
		urlQuery := r.URL.Query()
		logg.Debugf("has query: %v", urlQuery.Has("query"))

		// pagination
		listTmpl.Pagination = true
		pageNr, err := strconv.Atoi(r.FormValue("page"))
		if err != nil || pageNr < 1 {
			pageNr = 1
		}
		limit, err := strconv.Atoi(r.FormValue("limit"))
		if err != nil {
			limit = env.DefaultTableSize()
		}
		if limit == 0 {
			limit = env.DefaultTableSize()
		}

		count, err := db.BoxListCounter(searchString)
		if err != nil {
			server.WriteInternalServerError("cant query boxes", err, w, r)
			return
		}

		data = common.Pagination(data, count, limit, pageNr)
		listTmpl.CurrentPageNumber = pageNr
		listTmpl.Limit = limit
		listTmpl.PaginationButtons = data["Pages"].([]common.PaginationButton)

		// box-list-row to fill box-list template
		boxes := make([]common.ListRow, limit)

		// Boxes found
		if count > 0 {
			// Fetch the Records from the Database and pack it into map
			rows, err := db.BoxListRows(searchString, limit, pageNr)
			if err != nil {
				server.WriteInternalServerError("cant query boxes", err, w, r)
				return
			}

			for i, b := range rows {
				boxes[i] = b
			}
			// Fill up empty rows to keep same table size
			if count < limit {
				for i := count; i < limit; i++ {
					boxes[i] = common.ListRow{}
				}
			}
		}
		listTmpl.Rows = boxes

		maps.Copy(data, listTmpl.Map())
		server.MustRender(w, r, templates.TEMPLATE_BOXES_PAGE, data)
	}
}

// boxDetailsPage shows a page with details of a specific box.
func boxDetailsPage(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		id := server.ValidID(w, r, "no box")
		if id.IsNil() {
			return
		}
		logg.Debug(id)

		notFound := false
		box, err := db.BoxById(id)
		if err != nil {
			logg.Errf("%s", err)
			notFound = true
		}
		box.ID = id
		data := boxes.BoxPageTemplateData()
		data.Box = &box

		data.Title = fmt.Sprintf("Box - %s", box.Label)
		data.Authenticated = authenticated
		data.User = user
		data.NotFound = notFound

		searchInput := items.NewSearchInputTemplate()
		searchInput.SearchInputLabel = "Search boxes"
		searchInput.SearchInputHxTarget = "#box-list"
		searchInput.SearchInputHxPost = "/boxes"

		dataForTemplate := data.Map()
		maps.Copy(dataForTemplate, searchInput.Map())

		dataForTemplate["ListRows"] = templates.SliceToSliceMaps(box.Items)

		server.MustRender(w, r, templates.TEMPLATE_BOX_DETAILS_PAGE, dataForTemplate)
	}
}

func createBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	box := boxes.NewBox()
	logg.Debug("create box: ", box)
	id, err := db.CreateBox(&box)
	if err != nil {
		server.WriteNotFoundError("error while creating the box", err, w, r)
		return
	}
	if server.WantsTemplateData(r) {
		box, err := db.BoxListRowByID(id)
		logg.Debug(box)
		if err != nil {
			server.WriteNotFoundError("error while fetching the box based on Id", err, w, r)
			return
		}
		server.MustRender(w, r, templates.TEMPLATE_BOX_LIST_ROW, box.Map())
	} else {
		server.WriteJSON(w, id)
	}
}

func updateBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	errMsgForUser := "Can't update box."
	id := server.ValidID(w, r, errMsgForUser)
	if id.IsNil() {
		return
	}

	box := boxFromPostFormValue(id, r)
	err := db.UpdateBox(box)
	if err != nil {
		server.WriteNotFoundError(errMsgForUser, err, w, r)
		return
	}
	if server.WantsTemplateData(r) {
		boxTemplate := boxes.BoxTemplateData{Box: &box, Edit: false}
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
	toDelete, err := parseIDsFromFormWithKey(r.Form, "delete")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, errMsgForUser)
		logg.Err(err)
	}
	deleteErrorIds := []string{}
	errOccurred := false
	for _, deleteId := range toDelete {
		err = nil
		err = db.DeleteBox(deleteId)
		if err != nil {
			errOccurred = true
			deleteErrorIds = append(deleteErrorIds, deleteId.String())
			logg.Errorf("%v: %v. %w", errMsgForUser, deleteId, err)
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

	if server.WantsTemplateData(r) {
		boxesPage(db).ServeHTTP(w, r)
		for _, id := range toDelete {
			templates.RenderSuccessNotification(w, "Box deleted: "+id.String())
		}
		return
	}
	fmt.Fprint(w, nil)
	w.WriteHeader(http.StatusOK)
}

// parseIDsFromFormWithKey parses r.Form by searching for all keys that start with `key` name and returns a list of valid uuid.UUIDs
//
// `r.ParseForm()` must be called before using this function!
//
// Example:
//
//	// search for all ID values that start with "delete:" key
//	// like "delete:f47ac10b-58cc-0372-8567-0e02b2c3d479"
//	toDeleteIDs := parseIDsFromFormWithKey(r.Form, "delete")
func parseIDsFromFormWithKey(form url.Values, key string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)
	for k := range form {
		// logg.Debugf("k: %v, v:%v", k, v)
		if strings.Contains(k, fmt.Sprintf("%s:", key)) {
			idStr := strings.Split(k, fmt.Sprintf("%s:", key))
			if len(idStr) != 2 {
				return nil, logg.NewError(fmt.Sprintf("Wrong delete key value pair: '%v'", k))
			}
			id, err := uuid.FromString(idStr[1])
			if err != nil {
				return nil, logg.WrapErr(err)
			}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// deleteBox deletes a single box.
func deleteBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	errMsgForUser := "Can't delete box"
	id := server.ValidID(w, r, errMsgForUser)
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
// boxPageMove shows the page where client can choose where to move the box.
func boxPageMove(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move single box page", w, r)
	}
}

// boxesPageMove shows the page where client can choose where to move the selected boxes from the boxes page.
func boxesPageMove(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.MustRender(w, r, "move-to", nil)
		// server.WriteNotImplementedWarning("Move multiple boxes page", w, r)
	}
}

// @TODO: Implement moveBox.
func moveBox(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move single box", w, r)
	}
}

func moveBoxes(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		moveToBoxID := server.ValidID(w, r, "can't move boxes invalid id")
		if moveToBoxID == uuid.Nil {
			return
		}

		parseIDs := r.PostForm["id-to-be-moved"]
		ids := make([]uuid.UUID, len(parseIDs))

		logg.Debug(len(parseIDs))

		notifications := server.Notifications{}
		for i, v := range parseIDs {
			logg.Debug(v)
			id := uuid.FromStringOrNil(v)
			ids[i] = id
			err := db.MoveBoxToBox(id, moveToBoxID)
			if err != nil {
				notifications.AddError(fmt.Sprintf(`can't move "%s" to "%s"`, ids[i].String(), moveToBoxID.String()))
				logg.Err(err)
			} else {
				notifications.AddSuccess(fmt.Sprintf(`moved "%s" to "%s"`, ids[i].String(), moveToBoxID.String()))
			}
		}
		server.TriggerNotifications(w, notifications)
	}
}

// boxFromPostFormValue returns items.Box without references to inner boxes, outer box and items.
func boxFromPostFormValue(id uuid.UUID, r *http.Request) boxes.Box {
	box := boxes.Box{}
	box.ID = id
	box.Label = r.PostFormValue("label")
	box.Description = r.PostFormValue("description")
	box.Picture = common.ParsePicture(r)
	box.QRCode = r.PostFormValue("qrcode")
	return box
}

func renderBoxTemplate(box *boxes.Box, w http.ResponseWriter, r *http.Request) {
	editParam := r.FormValue("edit")
	edit := false
	if editParam == "true" {
		edit = true
	}
	b := boxes.BoxTemplateData{Box: box, Edit: edit}
	server.MustRender(w, r, templates.TEMPLATE_BOX_DETAILS, b.Map())
}

func renderBoxesListTemplate(w http.ResponseWriter, r *http.Request, db BoxDatabase, boxes []common.ListRow, query string) error {
	searchInput := items.NewSearchInputTemplate()
	searchInput.SearchInputLabel = "Search boxes"
	searchInput.SearchInputHxTarget = "#box-list"
	searchInput.SearchInputHxPost = "/boxes-list"
	searchInput.SearchInputValue = query
	logg.Debugf("searchInput %v", searchInput.Map())
	boxesMaps := make([]map[string]any, len(boxes))
	for i := range boxes {
		boxesMaps[i] = boxes[i].Map()
	}
	data := map[string]any{"Boxes": boxesMaps}
	maps.Copy(data, searchInput.Map())
	logg.Debug("renderBoxesListTemplate: Boxes=", len(data["Boxes"].([]map[string]any)))

	err := templates.SafeRender(w, templates.TEMPLATE_BOX_LIST, data)
	if err != nil {
		return logg.Errorf("%w", err)
	}
	return nil
}

// // sampleBoxDB to test request handler.
// type sampleBoxDB struct {
// 	Boxes map[string]*boxes.Box
// }
//
// func newSampleBoxDB() *sampleBoxDB {
// 	db := sampleBoxDB{Boxes: make(map[string]*boxes.Box, 100)}
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

