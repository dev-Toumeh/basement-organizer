package boxes

import (
	"basement/main/internal/common"
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

// BoxHandler handles read, create, update and delete for single box.
func BoxHandler(db BoxDatabase) http.HandlerFunc {
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

func createBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	box := NewBox()
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
	// @TODO: Find a better solution?
	// This is done because box is missing OuterBox field after it's parsed.
	box, err = db.BoxById(id)
	if err != nil {
		server.WriteInternalServerError("can't get box after update succeeded, should not happen!", err, w, r)
		return
	}
	if server.WantsTemplateData(r) {
		boxTemplate := BoxTemplateData{Box: &box, Edit: false}
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

// BoxesHandler handles read and delete for multiple boxes.
func BoxesHandler(db BoxDatabase) http.HandlerFunc {
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
				a := BoxListTemplateData{Boxes: boxs}
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

func deleteBoxes(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
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
		ListPage(db).ServeHTTP(w, r)
		for _, id := range toDelete {
			templates.RenderSuccessNotification(w, "Box deleted: "+id.String())
		}
		return
	}
	fmt.Fprint(w, nil)
	w.WriteHeader(http.StatusOK)
}

// boxFromPostFormValue returns items.Box without references to inner boxes, outer box and items.
func boxFromPostFormValue(id uuid.UUID, r *http.Request) Box {
	box := Box{}
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

func renderBoxTemplate(box *Box, w http.ResponseWriter, r *http.Request) {
	editParam := r.FormValue("edit")
	edit := false
	if editParam == "true" {
		edit = true
	}
	b := BoxTemplateData{Box: box, Edit: edit}
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
