package shelves

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"

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

// render the Shelf List template with Add option
func AddListTemplate(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.Header().Add("Allow", http.MethodPost)
			w.WriteHeader(http.StatusMethodNotAllowed)
			logg.Debug(w, "Method:'", r.Method, "' not allowed")
			fmt.Fprint(w, "something went wrong please try again later")
			return
		}

		data := getTemplateData(r, db, w)
		data.SetFormHXTarget("this")
		data.SetRowHXGet("shelves/add-input")
		data.SetFormID("list-add")
		data.SetShowLimit(env.Config().ShowTableSize())
		data.SetPlaceHolder(false)

		data.SetRowAction(true)
		data.SetRowActionName("Add to")
		data.SetRowActionHXPostWithID("/shelves/add-input")
		data.SetRowActionHXTarget("#shelf-target")

		server.MustRender(w, r, "list", data.TypeMap)
	}
}

// render an html element from type input Field with the desired Shelf Id and label
func AddInputTemplate(db ShelfDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		errMsgForUser := "can't add, please comeback later"
		id := server.ValidID(w, r, errMsgForUser)

		shelf, err := db.Shelf(id)
		if err != nil {
			server.WriteInternalServerError("false request", err, w, r)
			return
		}

		inputHTML := fmt.Sprintf(`
      <label for="shelf_id">Put inside of Shelf</label></br>
      <input type="text" name="shelf_id" value="%s" hidden>
      <input type="text" name="label" value="%s" disabled>
      <button
       hx-target="#place-holder"
       hx-post="/shelves/add-list"
       hx-push-url="false">
       Add to another Shelf
      </button>
      <div id="place-holder" hx-swap-oob="true"></div>
      `,
			shelf.ID.String(), shelf.Label)

		w.Header().Set("Content-Type", "text/html")

		w.Write([]byte(inputHTML))
	}
}
