package boxes

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"maps"
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

// BoxDetailsPageTemplateData returns struct needed for "templates.TEMPLATE_BOX_DETAILS_PAGE" with default values.
func BoxDetailsPageTemplateData() boxDetailsPageTemplate {
	pageTmpl := templates.NewPageTemplate()
	boxTmpl := boxDetailsPageTemplate{
		PageTemplate: pageTmpl,
	}
	return boxTmpl
}

func (tmpl boxDetailsPageTemplate) Map() map[string]any {
	data := make(map[string]any, 0)
	maps.Copy(data, tmpl.PageTemplate.Map())
	maps.Copy(data, tmpl.Box.Map())
	data["Edit"] = tmpl.Edit
	data["Create"] = tmpl.Create
	return data
}

// DetailsPage shows a page with details of a specific box.
func DetailsPage(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		id := server.ValidID(w, r, "no box")
		if id.IsNil() {
			return
		}
		logg.Debug(id)

		notFound := false
		box, err := db.BoxById(id)
		if err != nil {
			logg.Errf("%s", err)
			notFound = true
		}
		box.ID = id
		logg.Debug("area id" + box.AreaID.String())
		tmpl := BoxDetailsPageTemplateData()
		tmpl.Box = box
		tmpl.Title = fmt.Sprintf("Box - %s", box.Label)
		tmpl.Authenticated = authenticated
		tmpl.User = user
		tmpl.NotFound = notFound

		searchInput := NewSearchInputTemplate()
		searchInput.SearchInputLabel = "Search boxes"
		searchInput.SearchInputHxTarget = "#box-list"
		searchInput.SearchInputHxPost = "/boxes"

		// innerBoxes := common.ListTemplate{
		// 	FormID:        "inner-boxes",
		// 	FormHXGet:     "/boxes",
		// 	Rows:          box.InnerBoxes,
		// 	RequestOrigin: "Boxes",
		// 	FormHXTarget:  "#inner-boxes",
		// 	// PlaceHolder:   true,
		// }
		// innerBoxes.AddRowOptions(common.ListRowTemplateOptions{
		// 	RowHXGet: "/box",
		// })

		var notifications server.Notifications
		tmpl.InnerBoxesList, err = common.ListTemplateInnerThingsFrom(common.THING_BOX, common.THING_BOX, w, r)
		if err != nil {
			notifications.AddError("could not load inner boxes")
		}

		tmpl.InnerItemsList, err = common.ListTemplateInnerThingsFrom(common.THING_ITEM, common.THING_BOX, w, r)
		if err != nil {
			notifications.AddError("could not load inner items")
		}

		// server.TriggerNotifications(w, )

		// tmpl.InnerBoxesList.FormID = "inner-boxes"
		// tmpl.InnerBoxesList.FormHXGet = "/boxes"
		// tmpl.InnerBoxesList.Rows = box.InnerBoxes
		// tmpl.InnerBoxesList.RequestOrigin = "Boxes"
		// tmpl.InnerBoxesList.FormHXTarget = "#inner-boxes"

		// tmpl.InnerBoxesList.AddRowOptions(common.ListRowTemplateOptions{
		// 	RowHXGet: "/box",
		// })
		// if err != nil {
		// 	logg.Err(err)
		// 	return
		// }
		// listTmpl.FormID = "innerboxes"
		// tmpl.InnerBoxesList = innerBoxes

		// innerItems := common.ListTemplate{
		// 	FormID:        "inner-items",
		// 	Rows:          box.Items,
		// 	RequestOrigin: "Items",
		// }
		// innerItems.AddRowOptions(common.ListRowTemplateOptions{
		// 	RowHXGet: "/items",
		// })
		// tmpl.InnerItemsList = innerItems
		if len(notifications.ServerNotificationEvents) > 0 {
			server.TriggerNotifications(w, notifications)
		}
		server.MustRender(w, r, templates.TEMPLATE_BOX_DETAILS_PAGE, tmpl)
	}
}

func RenderBoxDetailsForm(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		box, err := readBoxFromRequest(w, r, db)
		if err != nil {
			server.WriteNotFoundError("box not found", err, w, r)
			return
		}

		editParam := r.FormValue("edit")
		edit := false
		if editParam == "true" {
			edit = true
		}
		b := boxDetailsPageTemplate{Box: box, Edit: edit}
		opts := common.ListRowTemplateOptions{RowHXGet: "box"}
		common.AddRowOptionsToListRows2(b.Box.InnerBoxes, opts)

		server.MustRender(w, r, templates.TEMPLATE_BOX_DETAILS, b.Map())
	}
}

func NewSearchInputTemplate() *SearchInputTemplate {
	return &SearchInputTemplate{SearchInputLabel: "Search",
		SearchInputHxPost:   "/api/v1/implement-me",
		SearchInputHxTarget: "#item-list-body",
		SearchInputValue:    "",
	}
}
