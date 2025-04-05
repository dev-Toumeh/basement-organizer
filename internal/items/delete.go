package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type Response struct {
	Successful bool `json:"successful"`
}

// delete Item based on Id
func DeleteItemHandler(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			ids, err := itemIDS(r)
			if err != nil {
				logg.Errf("error while checking the Ids: %v", err)
				http.Error(w, "Invalid delete Item request", http.StatusBadRequest)
				return
			}

			if err := db.DeleteItems(ids); err != nil {
				templates.RenderErrorNotification(w, "we was not able to delete the item please comeback later")
				return
			}

			w.Header().Set("HX-Trigger-After-On-Load", "handleDeleteRows")
			w.WriteHeader(http.StatusOK)
			templates.RenderSuccessNotification(w, "items was deleted successfully")
			break

		case http.MethodDelete:
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

func itemIDS(r *http.Request) ([]uuid.UUID, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading request body: %v", err)
	}

	queryValues := strings.Split(string(body), "&")
	var selectedItemIDs []uuid.UUID

	for _, pair := range queryValues {
		parts := strings.Split(pair, "=")
		key, value := parts[0], parts[1]
		if value == "on" {
			id, err := uuid.FromString(key)
			if err != nil {
				return nil, fmt.Errorf("Error while converting to uuid during DeleteItems: %v", err)
			}
			selectedItemIDs = append(selectedItemIDs, id)
		}
	}
	return selectedItemIDs, nil
}
