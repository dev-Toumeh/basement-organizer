package shelves

import (
	"net/http"

	"basement/main/internal/logg"
	"basement/main/internal/templates"

	"github.com/gofrs/uuid/v5"
)

func createShelf(w http.ResponseWriter, r *http.Request, db ShelfDB) {
	shelf, err := shelf(r)
	if err != nil {
		logg.Errf("error while parsing the shelf request data: %v", err)
		templates.RenderErrorNotification(w, "Invalid shelf data")
		return
	}
	// Generate a new UUID for the shelf
	shelf.Id, err = uuid.NewV4()
	if err != nil {
		logg.Errf("error generating UUID for the shelf: %v", err)
		templates.RenderErrorNotification(w, "Error generating shelf ID")
		return
	}
	// @Todo validate shelf request data
	err = db.CreateShelf(shelf)
	if err != nil {
		templates.RenderErrorNotification(w, "Error while creating a new shelf, please try again later")
		return
	}
	templates.RenderSuccessNotification(w, "The new shelf was created successfully")
}

func readShelf(w http.ResponseWriter, r *http.Request, db ShelfDB) {
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
		templates.RenderErrorNotification(w, "Error updating the shelf, please try again later")
		return
	}
	templates.RenderSuccessNotification(w, "Shelf updated successfully")
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
