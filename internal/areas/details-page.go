package areas

import (
	"basement/main/internal/auth"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"fmt"
	"net/http"
)

// DetailsPage shows a page with details of a specific area.
func DetailsPage(db AreaDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		id := server.ValidID(w, r, "no area")
		if id.IsNil() {
			return
		}
		logg.Debug(id)

		notFound := false
		area, err := db.AreaById(id)
		if err != nil {
			logg.Errf("%s", err)
			notFound = true
		}

		area.ID = id
		data := NewAreaDetailsPageData()
		data.RequestOrigin = "Areas"
		data.Area = area

		data.Title = fmt.Sprintf("Area - %s", area.Label)
		data.Authenticated = authenticated
		data.User = user
		data.NotFound = notFound

		editParam := r.FormValue("edit")
		if editParam == "true" {
			data.Edit = true
		}

		server.MustRender(w, r, "area-details-page", data)
	}
}

func APIArea(db AreaDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := server.ValidID(w, r, "no area")
		if id.IsNil() {
			return
		}
		logg.Debug(id)
		area, err := db.AreaById(id)
		if err != nil {
			server.WriteNotFoundError("", err, w, r)
			return
		}
		server.WriteJSON(w, area)
	}
}
