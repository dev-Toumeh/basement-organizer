package shelves

import (
	"fmt"
	"net/http"

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

			// Template writer
			renderShelfTemplate(shelf, w, r)
			break

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

func deleteShelf(w http.ResponseWriter, r *http.Request, db ShelfDB) {
	// Extract the shelf ID from the request
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
	err = db.DeleteShelf(id)
	if err != nil {
		templates.RenderErrorNotification(w, "Error deleting the shelf, please try again later")
		return
	}
	templates.RenderSuccessNotification(w, "Shelf deleted successfully")
}

func renderShelfTemplate(box *Shelf, w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
