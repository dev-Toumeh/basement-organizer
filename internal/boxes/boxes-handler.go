package boxes

import (
	"basement/main/internal/server"
	"net/http"
)

// BoxesHandler handles read and delete for multiple boxes.
func BoxesHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			if !server.WantsTemplateData(r) {
				boxs, err := db.BoxListRows("", 100, 1)
				if err != nil {
					server.WriteNotFoundError("Can't find boxes", err, w, r)
					return
				}
				server.WriteJSON(w, boxs)
			} else {
				listPage(db).ServeHTTP(w, r)
			}
			break

		case http.MethodPut:
			server.WriteNotImplementedWarning("Multiple boxes edit?", w, r)
			break

		case http.MethodDelete:
			deleteBoxes(w, r, db)
			break

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			break
		}
	}
}

func deleteBoxes(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	server.DeleteThingsFromList(w, r, db.DeleteBox, listPage(db))
}
