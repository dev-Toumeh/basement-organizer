package boxes

import (
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"net/http"
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

// createPage renders a Box Details Page with initial details, No box is created yet in the Backend.
func createPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		box := NewBox()
		renderBoxTemplate(w, r, box.Map(), common.CreateMode)

	}
}

// createBox generates a Box filled with random data and stores it in db.
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

// createBoxFrom parses the submitted form values, builds a Box, and inserts it into db.
func createBoxFrom(w http.ResponseWriter, r *http.Request, db BoxDatabase) {

	box, validator, err := ValidateBox(w, r)
	if err != nil {
		if err == validator.Err() {
			logg.Debugf("validation error while creating the Box: %v", validator.Messages.Map())
			renderBoxTemplate(w, r, validator.BoxFormData(), common.CreateMode)
		} else {
			logg.Debugf("error happened while creating the Box: %v", err)
			server.TriggerSingleErrorNotification(w, "Error while creating the Box please comeback later")
		}
		return
	}

	logg.Debug("create box: ", box)
	_, err = db.CreateBox(&box)
	server.RedirectWithSuccessNotification(w, "/boxes", "Created new box: "+box.Label)
}

// updateBox update existing Box Record in the Database
func updateBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	const errMsgForUser = "Can't update box."

	// ── Validate input ────────────────────────────────────────────────
	box, validator, err := ValidateBox(w, r)
	if err != nil {
		if err == validator.Err() {
			logg.Debugf("validation error while updating the Box: %v", validator.Messages.Map())
			renderBoxTemplate(w, r, validator.BoxFormData(), common.EditMode)
		} else {
			logg.Debugf("error happened while updating the Box: %v", err)
			server.WriteNotFoundError("error while creating the box", err, w, r)
		}
		return
	}

	// ── Picture handling ─────────────────────────────────────────────
	ignorePicture, pictureFormat := false, ""
	if _, _, fh := r.FormFile("picture"); fh != nil {
		pictureFormat, err = common.ParsePictureFormat(r)
		if err != nil {
			logg.Debug("no picture format")
			ignorePicture = true
		}
	} else {
		ignorePicture = true
	}

	// ── Update in DB ────────────────────────────────────────────────
	if err = db.UpdateBox(box, ignorePicture, pictureFormat); err != nil {
		server.WriteNotFoundError(errMsgForUser+" "+logg.CleanLastError(err), err, w, r)
		return
	}

	logg.Debug("Updated Box: ", box)
	server.RedirectWithSuccessNotification(w, "/box/"+box.ID.String()+"", "Updated box: "+box.Label)
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
