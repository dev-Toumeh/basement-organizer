package areas

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

func AreaHandler(db AreaDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if !server.WantsTemplateData(r) {
				APIArea(db).ServeHTTP(w, r)
				return
			}
			DetailsPage(db).ServeHTTP(w, r)
			break

		case http.MethodPost:
			createArea(w, r, db)
			break

		case http.MethodPut:
			updateArea(w, r, db)
			break

		case http.MethodDelete:
			deleteArea(w, r, db)
			return

		default:
			// Other methods are not allowed.
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			server.WriteFprint(w, "Method:'"+r.Method+"' not allowed")
		}
	}
}

func createArea(w http.ResponseWriter, r *http.Request, db AreaDatabase) {
	area := NewArea()
	logg.Debug("create area: ", area)
	id, err := db.CreateArea(area)
	if err != nil {
		server.WriteNotFoundError("error while creating the area", err, w, r)
		return
	}
	if server.WantsTemplateData(r) {
		area, err := db.AreaListRowByID(id)
		logg.Debug(area)
		if err != nil {
			server.WriteNotFoundError("error while fetching the area based on Id", err, w, r)
			return
		}

		area.RowHXGet = "/area"
		area.HideMoveCol = true
		server.MustRender(w, r, templates.TEMPLATE_LIST_ROW, area)
	} else {
		server.WriteJSON(w, id)
	}
}

func updateArea(w http.ResponseWriter, r *http.Request, db AreaDatabase) {
	errMsgForUser := "Can't update area."

	area, validator, err := ValidateArea(w, r)
	if err != nil {
		if err == validator.Err() {
			logg.Warning("validation error while updating the Area: %v", validator.Messages.Map())
			templates.Render(w, "area-details", validator.AreaFormData(true))
		} else {
			logg.Debugf("error happened while updating the Area: %v", err)
			server.TriggerSingleErrorNotification(w, "Error while updating the Area, please come back later")
		}
		return
	}
	ignorePicture := !server.ParseIgnorePicture(r)
	pictureFormat := ""
	if !ignorePicture {
		var err error
		pictureFormat, err = common.ParsePictureFormat(r)
		if err != nil {
			logg.Debug("no picture format")
		}
	}

	err = db.UpdateArea(area, ignorePicture, pictureFormat)
	logg.Debugf("this is the area: %v", area.Map())
	if err != nil {
		server.WriteNotFoundError(errMsgForUser+" "+logg.CleanLastError(err), err, w, r)
		return
	}

	// @TODO: Find a better solution. Picture is not included in request if ignorePicture is true and will be missing in response.
	area, err = db.AreaById(area.ID)
	if err != nil {
		server.WriteNotFoundError("no area found with id: "+area.ID.String(), err, w, r)
		return
	}
	if server.WantsTemplateData(r) {
		areaTemplate := AreaDetailsPageData{Area: area, Edit: false}
		err := server.RenderWithSuccessNotification(w, r, "area-details", areaTemplate, "Updated area: "+areaTemplate.Label)
		if err != nil {
			server.WriteInternalServerError(errMsgForUser, err, w, r)
			return
		}
	} else {
		server.WriteJSON(w, area)
	}
	logg.Debug("Updated Area: ", area)
}

// deleteArea deletes a single area.
func deleteArea(w http.ResponseWriter, r *http.Request, db AreaDatabase) {
	errMsgForUser := "Can't delete area"
	id := server.ValidID(w, r, errMsgForUser)
	if id.IsNil() {
		return
	}
	err := db.DeleteArea(id)
	if err != nil {
		server.WriteNotFoundError(errMsgForUser, err, w, r)
		return
	}
	server.RedirectWithSuccessNotification(w, "/areas", id.String()+" deleted")
}

// CreateHandler
//
//	GET = create new area page
//	POST = submit new area from create area page
func CreateHandler(db AreaDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			createPage(db).ServeHTTP(w, r)
			break
		case http.MethodPost:
			idStr := r.PostFormValue("id")
			if idStr == "" {
				createArea(w, r, db)
			} else {
				id := server.ValidID(w, r, "can't create new area")
				if id == uuid.Nil {
					return
				}
				createAreaWithID(w, r, db, id)
			}
			break
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Add("Allowed", http.MethodGet)
			w.Header().Add("Allowed", http.MethodPost)
		}
	}
}

// createPage renders a page with initial area details to create a new area. No area is created yet in the backend.
func createPage(db AreaDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		data := NewAreaDetailsPageData()
		area := NewArea()
		data.Area = area

		data.Title = "area - " + area.Label
		data.Authenticated = authenticated
		data.User = user
		data.Create = true
		data.RequestOrigin = "Areas"
		data.DescriptionError = ""
		data.LabelError = ""

		server.MustRender(w, r, "area-details-page", data)
	}
}

// createAreaWithID creates an area from an existing one.
func createAreaWithID(w http.ResponseWriter, r *http.Request, db AreaDatabase, id uuid.UUID) {

	area, validator, err := ValidateArea(w, r)
	if err != nil {
		if err == validator.Err() {
			logg.Warning("validation error while creating the Area: %v", validator.Messages.Map())
			templates.Render(w, "area-details", validator.AreaFormData(false))
		} else {
			logg.Debugf("error happened while creating the Area: %v", err)
			server.TriggerSingleErrorNotification(w, "Error while creating the Area, please come back later")
		}
		return
	}

	logg.Debug("create area: ", area)
	id, err = db.CreateArea(area)
	if err != nil {
		server.WriteNotFoundError("error while creating the area", err, w, r)
		return
	}
	server.RedirectWithSuccessNotification(w, "/areas", "Created new area: "+area.Label)
}
