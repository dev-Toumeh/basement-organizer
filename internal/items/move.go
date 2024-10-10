package items

import (
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

// update the item based on ID
func MoveItemHandler(db ItemDatabase) func(w http.ResponseWriter, r *http.Request) {
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
