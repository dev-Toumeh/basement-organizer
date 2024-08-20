package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

// delete Item based on Id
func DeleteItemHandler(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {

			id, err := uuid.FromString(r.FormValue(ID))
			if err != nil {
				logg.Debug("the id is not valid")
				logg.Err(err)
				http.Error(w, "Invalid ID", http.StatusBadRequest)
			} else {
				ctx := context.TODO()
				if err := db.DeleteItem(ctx, id); err != nil {
					templates.RenderErrorSnackbar(w, "we was not able to delete the item please comeback later")
					return
				}
				templates.RenderSuccessSnackbar(w, "deleted Successfully")
				w.WriteHeader(http.StatusAccepted)
			}
		} else {
			logg.Debug("Invalid Request")
			http.Error(w, "Invalid Request", http.StatusBadRequest)
		}
	}
}
