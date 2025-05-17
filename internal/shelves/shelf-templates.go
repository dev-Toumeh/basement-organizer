package shelves

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/server"
	"basement/main/internal/templates"

	"maps"
	"net/http"
)

// Render shelf Root page where you can search the available Shelves
func PageTemplate(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := getTemplateData(r, db, w)
		data.SetPlaceHolder(true)
		data.SetEnvDevelopment(env.Development())
		data.SetRequestOrigin("Shelves")
		data.TypeMap["HideAreaLabel"] = false
		server.MustRender(w, r, "shelves-page", data.TypeMap)
	}
}

// Render create Shelf Template with defaults Values
func CreateTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		page := templates.NewPageTemplate()
		page.Title = "Add new Shelf"
		page.Authenticated = authenticated
		page.User = user

		shelf := newShelf()
		data := page.Map()
		maps.Copy(data, shelf.Map())

		templates.Render(w, "shelf-create-page", data)
	}
}

// Render Shelf Details Template where you can preview the shelf and update the relevant Data
func DetailsTemplate(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errMsgForUser := "the requested Shelf doesn't exist"

		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		id := server.ValidID(w, r, errMsgForUser)
		if id.IsNil() {
			return
		}

		shelf, err := db.Shelf(id)
		if err != nil {
			server.TriggerErrorNotification(w, errMsgForUser)
		}

		page := templates.NewPageTemplate()
		page.Title = "Shelf Details"
		page.Authenticated = authenticated
		page.User = user

		var notifications server.Notifications
		shelf.InnerBoxesList, err = common.ListTemplateInnerThingsFrom(common.THING_BOX, common.THING_SHELF, w, r)
		if err != nil {
			notifications.AddError("could not load inner boxes")
		}

		shelf.InnerItemsList, err = common.ListTemplateInnerThingsFrom(common.THING_ITEM, common.THING_SHELF, w, r)
		if err != nil {
			notifications.AddError("could not load inner items")
		}

		if len(notifications.ServerNotificationEvents) > 0 {
			server.TriggerNotifications(w, notifications)
		}

		maps := []map[string]any{
			page.Map(),
			shelf.Map(),
			{"Edit": common.CheckEditMode(r)},
		}

		data := common.MergeMaps(maps)

		templates.Render(w, "shelf-details-page", data)
	}
}
