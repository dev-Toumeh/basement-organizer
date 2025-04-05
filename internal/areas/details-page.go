package areas

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"net/http"
)

type AreaDetailsPageData struct {
	templates.PageTemplate
	Area
	InnerItemsList   common.ListTemplate
	InnerBoxesList   common.ListTemplate
	InnerShelvesList common.ListTemplate
	Edit             bool
	Create           bool
}

// NewAreaDetailsPageData returns struct needed for "templates.TEMPLATE_area_DETAILS_PAGE" with default values.
func NewAreaDetailsPageData() (data AreaDetailsPageData) {
	data.PageTemplate = templates.NewPageTemplate()
	return data
}

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

		data.InnerItemsList, err = common.ListTemplateInnerThingsFrom(common.THING_ITEM, common.THING_AREA, w, r)
		data.InnerBoxesList, err = common.ListTemplateInnerThingsFrom(common.THING_BOX, common.THING_AREA, w, r)
		data.InnerShelvesList, err = common.ListTemplateInnerThingsFrom(common.THING_SHELF, common.THING_AREA, w, r)
		logg.Debugf("inner boxes %v", data.InnerBoxesList.Rows)

		// {{ template "list" .InnerBoxesList }}

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
