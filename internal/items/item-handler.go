package items

import (
	"net/http"

	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"

	"github.com/gofrs/uuid/v5"
)

// Handles read, create, update, and delete for multiple items.
func ItemsHandler(db ItemDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			PageTemplate(db).ServeHTTP(w, r)
			return
		case http.MethodDelete:
			server.DeleteThingsFromList(w, r, db.DeleteItem, PageTemplate(db))
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}

// Handles read, create, update, and delete for a single item.
func ItemHandler(db ItemDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createItem(w, r, db)
			break

		case http.MethodDelete:
			deleteItem(w, r, db)
			return

		case http.MethodPut:
			updateItem(w, r, db)
			break

		default:
			logg.Debug("Invalid Request")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Add("Allowed", http.MethodPost)
			w.Header().Add("Allowed", http.MethodDelete)
			break
		}
	}
}

func createItem(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	validator, err := ValidateItem(r, w)
	if err != nil {
		if err == validator.Err() {
			renderItemTemplate(r, w, validator.ItemFormData(), common.CreateMode)
		} else {
			logg.Err(err)
			server.TriggerSingleErrorNotification(w, "Error while generating the Item please comeback later")
		}
		return
	}

	item := ToItem(validator.Item)

	if err := db.CreateNewItem(item); err != nil {
		if err == db.ErrorExist() {
			logg.Debugf("the Label is already token please choice another one", err)
			templates.RenderErrorNotification(w, "the Label is already token please choice another one")
		} else {
			logg.Debugf("error while creating new Item: %v", err)
			templates.RenderErrorNotification(w, "Unable to add new item due to technical issues. Please try again later.")
		}
	}
	logg.Debug("the Item with id: " + item.ID.String() + " was created")
	server.RedirectWithSuccessNotification(w, "/items", "The Item was created successfully")
	return
}

func updateItem(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	validator, err := ValidateItem(r, w)
	if err != nil {
		if err == validator.Err() {
			logg.Debugf("validation error while updating the Item: %v", err)
			renderItemTemplate(r, w, validator.ItemFormData(), common.EditMode)
		} else {
			logg.Debugf("error happened while updating the Item: %v", err)
			server.TriggerSingleErrorNotification(w, "Error while generating the Item please comeback later")
		}
		return
	}

	item := ToItem(validator.Item)
	ignorePicture := server.ParseIgnorePicture(r)
	pictureFormat := ""
	if !ignorePicture {
		pictureFormat, err = common.ParsePictureFormat(r)
		if err != nil {
			logg.Debug("no picture format")
		}
	}

	err = db.UpdateItem(item, ignorePicture, pictureFormat)
	if err != nil {
		server.WriteNotFoundError("Can't update item. "+logg.CleanLastError(err), err, w, r)
		return
	}
	url := "/item/" + item.ID.String()
	server.RedirectWithSuccessNotification(w, url, "item updated successfully")
}

func deleteItem(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	w.Header().Add("Allowed", http.MethodGet)
	id := server.ValidID(w, r, "invalid ID")
	if id == uuid.Nil {
		return
	}

	if err := db.DeleteItem(id); err != nil {
		server.WriteBadRequestError(logg.CleanLastError(err), err, w, r)
		return
	}
	server.RedirectWithSuccessNotification(w, "/items", "item deleted "+id.String())
}

func MoveItem(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		// server.WriteNotImplementedWarning("Move item", w, r)
		errMsgForUser := "Can't move item"

		id := server.ValidID(w, r, errMsgForUser)
		if id.IsNil() {
			return
		}
		id2, err := uuid.FromString(r.PostFormValue("id2"))
		if err != nil {
			err = logg.Errorf("%s %w", errMsgForUser, err)
			server.WriteInternalServerError(errMsgForUser, err, w, r)
			return
		}
		err = db.MoveItemToBox(id, id2)
		if err != nil {
			err = logg.Errorf("%s %w", errMsgForUser, err)
			server.WriteInternalServerError(errMsgForUser, err, w, r)
			return
		}
		logg.Infof("move '%s' to '%s'", id, id2)
		return
	}
}
