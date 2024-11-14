package shelves

import (
	"fmt"
	"net/http"

	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"

	"github.com/gofrs/uuid/v5"
)

// handles read, create, update and delete for single shelf.
func ShelfHandler(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			const errMsgForUser = "Can't find box"

			id := server.ValidID(w, r, errMsgForUser)
			if id.IsNil() {
				return
			}

			shelf, err := db.Shelf(id)
			if err != nil {
				server.WriteNotFoundError(errMsgForUser, err, w, r)
				return
			}

			// Use API data writer
			if !server.WantsTemplateData(r) {
				server.WriteJSON(w, shelf)
				return
			}

		case http.MethodPost:
			createShelf(w, r, db)
			break

		case http.MethodDelete:
			deleteShelf(w, r, db)
			return

		case http.MethodPut:
			updateShelf(w, r, db)
			break

		default:
			// Other methods are not allowed.
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		}
	}
}

func createShelf(w http.ResponseWriter, r *http.Request, db ShelfDB) {
	shelf, err := shelf(r)
	if err != nil {
		logg.Errf("error while parsing the shelf request data: %v", err)
		templates.RenderErrorNotification(w, "Invalid shelf data")
		return
	}
	// @Todo validate shelf request data
	err = db.CreateShelf(shelf)
	if err != nil {
		templates.RenderErrorNotification(w, "Error while creating a new shelf, please try again later")
		return
	}
	server.RedirectWithSuccessNotification(w, "/shelves", "The Shelf was created successfully")
}

func readShelf(w http.ResponseWriter, r *http.Request, db ShelfDB) {
	// Extract the shelf ID from the request
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		logg.Errf("Shelf ID is missing in the request")
		templates.RenderErrorNotification(w, "Shelf was not Found")
		return
	}
	id, err := uuid.FromString(idStr)
	if err != nil {
		logg.Errf("Invalid shelf ID: %v", err)
		templates.RenderErrorNotification(w, "Invalid shelf ID")
		return
	}
	_, err = db.Shelf(id)
	if err != nil {
		logg.Errf("Shelf not found: %v", err)
		templates.RenderErrorNotification(w, "Shelf not found")
		return
	}
	// Render the shelf data (assuming a template exists)
	templates.Render(w, "", "")
}

func updateShelf(w http.ResponseWriter, r *http.Request, db ShelfDB) {
	shelf, err := shelf(r)
	if err != nil {
		logg.Errf("error while parsing shelf data: %v", err)
		templates.RenderErrorNotification(w, "Invalid shelf data")
		return
	}
	// @Todo validate shelf request data
	err = db.UpdateShelf(shelf)
	if err != nil {
		templates.RenderErrorNotification(w, "Error while updating the shelf, please try again later")
		return
	}
	url := fmt.Sprintf("/shelves/update?id=%s", shelf.ID.String())
	server.RedirectWithSuccessNotification(w, url, "Shelf updated successfully")
}

// delete single shelf
func deleteShelf(w http.ResponseWriter, r *http.Request, db ShelfDB) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		logg.Errf("Shelf ID is missing in the request")
		templates.RenderErrorNotification(w, "Shelf ID is required")
		return
	}
	id, err := uuid.FromString(idStr)
	if err != nil {
		logg.Errf("Invalid shelf ID: %v", err)
		templates.RenderErrorNotification(w, "Invalid shelf ID")
		return
	}
	label, err := db.DeleteShelf(id)
	if err != nil {
		if err == db.ErrorNotEmpty() {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "the Shelf"+label+"is not Empty")
			return
		}
		fmt.Printf("error while deleting the shelf: %v", err)
		templates.RenderErrorNotification(w, "Error deleting the shelf, please try again later")
		return
	}
	server.RedirectWithSuccessNotification(w, "/shelves", "Shelf deleted successfully")
}

// delete multiple Shelves
func DeleteShelves(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		errMsgForUser := "Can't delete the Shelves please try again later"
		r.ParseForm()
		toDelete, err := common.ParseIDsFromFormWithKey(r.Form, "delete")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, errMsgForUser)
			logg.Err(err)
		}

		notifications := server.Notifications{}

		for _, deleteId := range toDelete {
			label, err := db.DeleteShelf(deleteId)
			if err != nil {
				if err == db.ErrorNotEmpty() {
					logg.Debug("the Shelf with the label: " + label + " and id:" + deleteId.String() +
						"could not be deleted as it is not empty \n")
					notifications.AddWarning("the Shelf: " + label +
						"could not be deleted as it is not empty")
					continue
				}
				fmt.Printf("an error accrue while deleting the Shelf with id: %s : %v",
					deleteId.String(), err)
				templates.RenderErrorNotification(w, errMsgForUser)
			}
			logg.Debug("the Shelf with the label: " + label + " and id:" +
				deleteId.String() + " was deleted \n")
			notifications.AddSuccess("the Shelf " + label + "was deleted")
		}
		server.RedirectWithNotifications(w, "/shelves", notifications)
	}
}
