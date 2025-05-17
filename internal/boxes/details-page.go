package boxes

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"net/http"
)

type boxDetailsPageTemplate struct {
	templates.PageTemplate
	Box
	Edit           bool
	Create         bool
	InnerItemsList common.ListTemplate
	InnerBoxesList common.ListTemplate
}

type SearchInputTemplate struct {
	SearchInputLabel    string
	SearchInputHxPost   string
	SearchInputHxTarget string
	SearchInputValue    string
}

// DetailsPage shows a page with details of a specific box.
func DetailsPage(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		notFound := false
		id := server.ValidID(w, r, "no box")
		if id.IsNil() {
			return
		}
		logg.Debug(id)

		box, err := db.BoxById(id)
		if err != nil {
			logg.Errf("%s", err)
			notFound = true
		}
		logg.Debug("area id" + box.AreaID.String())

		data := common.InitData(r, false)
		data.SetDetailesData(box.Map())
		data.SetNotFound(notFound)

		// Set innerBoxes and innerItems
		var notifications server.Notifications
		err = data.SetInnerBoxesList(w, r)
		if err != nil {
			notifications.AddError("could not load inner boxes")
		}
		err = data.SetInnerItemsList(w, r)
		if err != nil {
			notifications.AddError("could not load inner items")
		}
		if len(notifications.ServerNotificationEvents) > 0 {
			server.TriggerNotifications(w, notifications)
		}
		renderBoxTemplate(w, r, data.TypeMap, common.PreviewMode)
	}
}

// RenderBoxDetailsForm render DetailsPage for editing
func RenderBoxDetailsForm(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notFound := false
		id := server.ValidID(w, r, "no box")
		if id.IsNil() {
			return
		}
		logg.Debug(id)

		box, err := db.BoxById(id)
		if err != nil {
			logg.Errf("%s", err)
			notFound = true
		}
		logg.Debug("area id" + box.AreaID.String())

		data := common.InitData(r, false)
		data.SetDetailesData(box.Map())
		data.SetNotFound(notFound)

		// Set innerBoxes and innerItems
		var notifications server.Notifications
		err = data.SetInnerBoxesList(w, r)
		if err != nil {
			notifications.AddError("could not load inner boxes")
		}
		err = data.SetInnerItemsList(w, r)
		if err != nil {
			notifications.AddError("could not load inner items")
		}
		if len(notifications.ServerNotificationEvents) > 0 {
			server.TriggerNotifications(w, notifications)
		}
		renderBoxTemplate(w, r, data.TypeMap, common.EditMode)
	}
}

func renderBoxTemplate(w http.ResponseWriter, r *http.Request, values map[string]any, typeMode common.Mode) {
	data := common.InitData(r, false)
	data.SetEnvDevelopment(env.Development())
	data.SetDetailesData(values)
	data.SetTypeMode(typeMode)
	data.SetOrigin("box")
	data.SetTitle(common.ToUpper(string(typeMode)) + " Box")
	err := templates.Render(w, templates.TEMPLATE_BOX_DETAILS_PAGE, data.TypeMap)
	if err != nil {
		logg.Warningf("An Error accrue while fetching item Extra Info", err)
	}
}
