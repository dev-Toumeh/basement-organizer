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
	"strings"

	"github.com/gofrs/uuid/v5"
)

func registerBoxRoutes(db *database.DB) {
	// Box templates
	Handle("/box", boxHandler(db))
	Handle("/box/{id}/moveto/box", boxes.BoxPageMove("box", db))
	Handle("/box/{id}/moveto/box/{value}", boxes.BoxMoveConfirm("box", db))
	Handle("/box/{id}/moveto/shelf", boxes.BoxPageMove("shelf", db))
	Handle("/box/{id}/moveto/shelf/{value}", boxes.BoxMoveConfirm("shelf", db))

	// Box api
	Handle("/api/v1/box", boxHandler(db))
	Handle("/api/v1/box/{id}", boxHandler(db))
	Handle("/api/v1/box/{id}/move/{toid}", boxes.MoveBox(db))

	// Boxes templates
	Handle("/boxes", boxesPage(db))
	Handle("/boxes/{id}", boxDetailsPage(db))
	Handle("/boxes/move", boxes.BoxesPageMove(db))
	Handle("/boxes/moveto/box/{id}", boxes.MoveBoxesToBoxHandler(db))
	Handle("/boxes-list", boxesHandler(db))

	// Boxes api
	Handle("/api/v1/boxes", boxesHandler(db))
	Handle("/api/v1/boxes/moveto/box/{id}", boxes.MoveBoxesToBoxAPI(db))
	Handle("/get-boxes-move-page", boxes.GetMoveBoxesPage(db))
}

// boxHandler handles read, create, update and delete for single box.
func boxHandler(db boxes.BoxDatabase) http.HandlerFunc {
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
func boxesHandler(db boxes.BoxDatabase) http.HandlerFunc {
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
func boxesPage(db boxes.BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		// page template
		page := templates.NewPageTemplate()
		page.Title = "Boxes"
		page.RequestOrigin = "Boxes"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// list template
		listTmpl := common.ListTemplate{
			FormHXGet:   "/boxes",
			RowHXGet:    "/boxes",
			PlaceHolder: true,
			ShowLimit:   env.Config().ShowTableSize(),
		}

		// search-input template
		searchString := common.SearchString(r)
		listTmpl.SearchInput = true
		listTmpl.SearchInputLabel = "Search boxes"
		listTmpl.SearchInputValue = searchString

		count, err := db.BoxListCounter(searchString)
		if err != nil {
			server.WriteInternalServerError("cant query boxes", err, w, r)
			return
		}

		// pagination
		pageNr := common.ParsePageNumber(r)
		limit := common.ParseLimit(r)
		data = common.Pagination(data, count, limit, pageNr)
		listTmpl.Pagination = true
		listTmpl.CurrentPageNumber = data["PageNumber"].(int)
		listTmpl.Limit = limit
		listTmpl.PaginationButtons = data["Pages"].([]common.PaginationButton)

		// box-list-row to fill box-list template
		var boxes []common.ListRow

		// Boxes found
		if count > 0 {
			boxes, err = common.FilledRows(db.BoxListRows, searchString, limit, pageNr, count)
			if err != nil {
				server.WriteInternalServerError("cant query boxes", err, w, r)
				return
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

func createBox(w http.ResponseWriter, r *http.Request, db boxes.BoxDatabase) {
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

func updateBox(w http.ResponseWriter, r *http.Request, db boxes.BoxDatabase) {
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
	// @TODO: Find a better solution?
	// This is done because box is missing OuterBox field after it's parsed.
	box, err = db.BoxById(id)
	if err != nil {
		server.WriteInternalServerError("can't get box after update succeeded, should not happen!", err, w, r)
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

func deleteBoxes(w http.ResponseWriter, r *http.Request, db boxes.BoxDatabase) {
	errMsgForUser := "Can't delete boxes"
	r.ParseForm()
	toDelete, err := common.ParseIDsFromFormWithKey(r.Form, "delete")
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

// deleteBox deletes a single box.
func deleteBox(w http.ResponseWriter, r *http.Request, db boxes.BoxDatabase) {
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

// boxFromPostFormValue returns items.Box without references to inner boxes, outer box and items.
func boxFromPostFormValue(id uuid.UUID, r *http.Request) boxes.Box {
	box := boxes.Box{}
	box.ID = id
	box.Label = r.PostFormValue("label")
	box.Description = r.PostFormValue("description")
	box.Picture = common.ParsePicture(r)
	box.QRCode = r.PostFormValue("qrcode")
	box.OuterBoxID = uuid.FromStringOrNil(r.PostFormValue("box_id"))
	box.ShelfID = uuid.FromStringOrNil(r.PostFormValue("shelf_id"))
	box.AreaID = uuid.FromStringOrNil(r.PostFormValue("area_id"))
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

func renderBoxesListTemplate(w http.ResponseWriter, r *http.Request, db boxes.BoxDatabase, boxes []common.ListRow, query string) error {
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
