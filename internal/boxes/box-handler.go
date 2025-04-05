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

var boxDB BoxDatabase

// RegisterDBInstance sets db instance for internal package usage.
// Is used for public functions that depend on the DB without the need to pass the instance as a parameter.
func RegisterDBInstance(db BoxDatabase) {
	boxDB = db
	logg.Debug("boxDB in boxes package registered")
}

// BoxHandler handles read, create, update and delete for single box.
func BoxHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			const errMsgForUser = "Can't find box"

			box, err := readBoxFromRequest(w, r, db)
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
			DetailsPage(db).ServeHTTP(w, r)
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

func readBoxFromRequest(w http.ResponseWriter, r *http.Request, db BoxDatabase) (box Box, err error) {
	const errMsgForUser = "Can't find box"

	id := server.ValidID(w, r, errMsgForUser)
	if id.IsNil() {
		return box, err
	}

	box, err = db.BoxById(id)
	if err != nil {
		return box, logg.WrapErr(err)
		// server.WriteNotFoundError(errMsgForUser, err, w, r)
	}
	return box, err
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
		box.RowHXGet = "/box"
		server.MustRender(w, r, templates.TEMPLATE_LIST_ROW, box)
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

	var err error
	box, ignorePicture := boxFromPostFormValue(id, r)
	pictureFormat := ""
	if !ignorePicture {
		pictureFormat, err = common.ParsePictureFormat(r)
		if err != nil {
			logg.Debug("no picture format")
		}
	}
	err = db.UpdateBox(box, ignorePicture, pictureFormat)

	if err != nil {
		server.WriteNotFoundError("Can't update box. "+logg.CleanLastError(err), err, w, r)
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
		boxTemplate := boxDetailsPageTemplate{Box: box, Edit: false}
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
			createPage().ServeHTTP(w, r)
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
func createPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		box := NewBox()
		data := BoxDetailsPageTemplateData()
		data.Edit = true
		data.Box = box

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
	box, _ := boxFromPostFormValue(id, r)
	logg.Debug("create box: ", box)
	id, err := db.CreateBox(&box)
	if err != nil {
		server.WriteNotFoundError("error while creating the box", err, w, r)
		return
	}
	server.RedirectWithSuccessNotification(w, "/boxes", "Created new box: "+box.Label)
}

// boxFromPostFormValue returns items.Box without references to inner boxes, outer box and items.
func boxFromPostFormValue(id uuid.UUID, r *http.Request) (box Box, ignorePicture bool) {
	ignorePicture = server.ParseIgnorePicture(r)
	box.BasicInfo = common.BasicInfoFromPostFormValue(id, r, ignorePicture)
	box.OuterBoxID = uuid.FromStringOrNil(r.PostFormValue("box_id"))
	box.ShelfID = uuid.FromStringOrNil(r.PostFormValue("shelf_id"))
	box.AreaID = uuid.FromStringOrNil(r.PostFormValue("area_id"))
	return box, ignorePicture
}
