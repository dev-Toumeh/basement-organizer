package boxes

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"net/http"

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

// CreateHandler
//
//	GET = create new box page
//	POST = submit new box from create box page
func CreateHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			createPage(db).ServeHTTP(w, r)
			break
		case http.MethodPost:
			createBoxFrom(w, r, db)
			break
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Add("Allowed", http.MethodGet)
			w.Header().Add("Allowed", http.MethodPost)
		}
	}
}

// createPage renders a page with initial box details to create a new box. No box is created yet in the backend.
func createPage(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		box := NewBox()
		data := BoxPageTemplateData()
		data.Box = &box

		data.Title = fmt.Sprintf("Box - %s", box.Label)
		data.Authenticated = authenticated
		data.User = user
		data.Create = true

		dataForTemplate := data.Map()

		dataForTemplate["ListRows"] = templates.SliceToSliceMaps(box.Items)

		server.MustRender(w, r, templates.TEMPLATE_BOX_DETAILS_PAGE, dataForTemplate)
	}
}

// createBoxFrom creates a box from an existing one.
func createBoxFrom(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	id := server.ValidID(w, r, "can't create new box")
	if id == uuid.Nil {
		return
	}
	box := boxFromPostFormValue(id, r)
	logg.Debug("create box: ", box)
	id, err := db.CreateBox(&box)
	if err != nil {
		server.WriteNotFoundError("error while creating the box", err, w, r)
		return
	}
	server.RedirectWithSuccessNotification(w, "/boxes", "Created new box: "+box.Label)
}

// BoxesHandler handles read and delete for multiple boxes.
func BoxesHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			if !server.WantsTemplateData(r) {
				boxs, err := db.BoxListRows("", 100, 1)
				if err != nil {
					server.WriteNotFoundError("Can't find boxes", err, w, r)
					return
				}
				server.WriteJSON(w, boxs)
			} else {
				listPage(db).ServeHTTP(w, r)
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

func deleteBoxes(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	r.ParseForm()

	parseIDs, _ := common.ParseIDsFromFormWithKey(r.Form, "delete")
	r.FormValue("delete")

	logg.Debug(len(parseIDs))

	ids := make([]uuid.UUID, len(parseIDs))
	notifications := server.Notifications{}
	for i, v := range parseIDs {
		logg.Debug("deleting " + v.String())
		ids[i] = v
		err := db.DeleteBox(v)
		if err != nil {
			notifications.AddError(`can't delete: "` + v.String() + `"`)
			logg.Err(err)
		} else {
			notifications.AddSuccess(`delete: "` + v.String() + `"`)
		}
	}

	server.TriggerNotifications(w, notifications)
	listPage(db).ServeHTTP(w, r)
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
