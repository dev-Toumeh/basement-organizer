package items

import (
	"fmt"
	"net/http"

	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
)

// Handles read, create, update, and delete for a single item.
func ItemHandler(db ItemDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// const errMsgForUser = "Can't find item"
			//
			// id := server.ValidID(w, r, errMsgForUser)
			// if id.IsNil() {
			// 	return
			// }
			//
			// item, err := db.ItemById(id)
			// if err != nil {
			// 	server.WriteNotFoundError(errMsgForUser, err, w, r)
			// 	return
			// }
			break

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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func createItem(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	var responseMessage []string
	newItem, err := item(r)
	if err != nil {
		logg.Err(err)
		templates.RenderErrorNotification(w, "Error while generating the User please comeback later")
	}
	if newItem, err = validateItem(newItem, &responseMessage); err != nil {
		responseGenerator(w, responseMessage, false)
		return
	}
	if err := db.CreateNewItem(newItem); err != nil {
		if err == db.ErrorExist() {
			templates.RenderErrorNotification(w, "the Label is already token please choice another one")
		} else {
			templates.RenderErrorNotification(w, "Unable to add new item due to technical issues. Please try again later.")
		}
	}
	logg.Debug("the Item with id: " + newItem.ID.String() + " was created")
	server.RedirectWithSuccessNotification(w, "/items", "The Item was created successfully")
	return
}

func deleteItem(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	var responseMessage []string
	newItem, err := item(r)
	if err != nil {
		logg.Err(err)
		templates.RenderErrorNotification(w, "Error while generating the User please comeback later")
	}
	if newItem, err = validateItem(newItem, &responseMessage); err != nil {
		responseGenerator(w, responseMessage, false)
		return
	}
	if err := db.CreateNewItem(newItem); err != nil {
		if err == db.ErrorExist() {
			templates.RenderErrorNotification(w, "the Label is already token please choice another one")
		} else {
			templates.RenderErrorNotification(w, "Unable to add new item due to technical issues. Please try again later.")
		}
	}
	logg.Debug("the Item with id: " + newItem.ID.String() + " was created")
	server.RedirectWithSuccessNotification(w, "/items", "The Item was created successfully")
	return
}

func updateItem(w http.ResponseWriter, r *http.Request, db ItemDatabase) {
	var errorMessages []string
	item, err := item(r)
	if err != nil {
		logg.Errf("error while parsing item data: %v", err)
		templates.RenderErrorNotification(w, "Invalid Item data")
		return
	}
	valiedItem, err := validateItem(item, &errorMessages)
	fmt.Print(valiedItem)

	ignorePicture := server.ParseIgnorePicture(r)
	pictureFormat := ""
	if !ignorePicture {
		pictureFormat, _ = common.ParsePictureFormat(r)
		if err != nil {
			logg.Debug("no picture format")
		}
	}

	err = db.UpdateItem(valiedItem, ignorePicture, pictureFormat)
	if err != nil {
		server.WriteNotFoundError("Can't update item. "+logg.CleanLastError(err), err, w, r)
		return
	}
	url := "/item/" + item.ID.String()
	server.RedirectWithSuccessNotification(w, url, "item updated successfully")
}
