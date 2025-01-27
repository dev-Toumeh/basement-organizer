package routes

import (
	"net/http"

	"basement/main/internal/areas"
	"basement/main/internal/boxes"
	"basement/main/internal/common"
	"basement/main/internal/database"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/shelves"
	"basement/main/internal/templates"

	"github.com/gofrs/uuid/v5"
)

// Fetch the Object (Boxes/Shelves/Areas) List
// this list will determined where to add the Object
func AddTo(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := common.InitData(r)
		thing := r.PathValue("thing")

		var err error
		var count int
		actionName := "Add to"
		rowActionType := "move"
		post := "/element" + "/" + thing

		switch thing {
		case "box":
			data.SetRowHXGet("/box")
			count, err = db.BoxListCounter("")
			if err != nil {
				server.WriteInternalServerError("no box list counter", err, w, r)
				return
			}
			data.SetCount(count)

			var rows []common.ListRow
			if count > 0 {
				rowOptions := common.ListRowTemplateOptions{
					RowHXGet:              "/box",
					RowAction:             true,
					RowActionType:         rowActionType,
					RowActionHXTarget:     "#" + thing + "-target",
					RowActionName:         actionName,
					RowActionHXPostWithID: post,
				}
				rows, err = common.FilledRows(db.BoxListRows, data.GetSearchInputValue(), data.GetLimit(), data.GetPageNumber(), count, rowOptions)
				if err != nil {
					server.WriteInternalServerError("cant query "+thing+" please comeback later", err, w, r)
				}
			}
			data.SetRows(rows)

		case "shelf":
			data.SetRowHXGet("/shelves")
			count, err = db.ShelfListCounter("")
			if err != nil {
				server.WriteInternalServerError("no shelf list counter", err, w, r)
				return
			}
			data.SetCount(count)

			var rows []common.ListRow
			if count > 0 {
				rowOptions := common.ListRowTemplateOptions{
					RowHXGet:              "/shelves",
					RowAction:             true,
					RowActionType:         rowActionType,
					RowActionHXTarget:     "#" + thing + "-target",
					RowActionName:         actionName,
					RowActionHXPostWithID: post,
				}
				rows, err = common.FilledRows(db.ShelfListRows, data.GetSearchInputValue(), data.GetLimit(), data.GetPageNumber(), count, rowOptions)
				if err != nil {
					server.WriteInternalServerError("cant query "+thing+" please comeback later", err, w, r)
				}
			}
			data.SetRows(rows)
			break

		case "area":
			data.SetRowHXGet("/area")
			count, err = db.AreaListCounter("")
			if err != nil {
				server.WriteInternalServerError("no area list counter", err, w, r)
				return
			}
			data.SetCount(count)

			var rows []common.ListRow
			if count > 0 {
				rowOptions := common.ListRowTemplateOptions{
					RowHXGet:              "/area",
					RowAction:             true,
					RowActionType:         rowActionType,
					RowActionHXTarget:     "#" + thing + "-target",
					RowActionName:         actionName,
					RowActionHXPostWithID: post,
				}
				rows, err = common.FilledRows(db.AreaListRows, data.GetSearchInputValue(), data.GetLimit(), data.GetPageNumber(), count, rowOptions)
				if err != nil {
					server.WriteInternalServerError("cant query "+thing+" please comeback later", err, w, r)
				}
			}
			data.SetRows(rows)
			break

		default:
			server.WriteInternalServerError("can't move box to \""+thing+"\"", logg.NewError("can't move box to \""+thing+"\""), w, r)
			return
		}

		data = common.Pagination2(data)
		data.SetShowLimit(env.Config().ShowTableSize())

		data.SetFormHXPost(post)
		data.SetFormID(thing + "-list")
		data.SetFormHXTarget("#" + thing + "-list")

		data.SetSearchInput(true)
		data.SetSearchInputLabel("Search " + thing)

		data.SetRowAction(true)
		data.SetRowActionName(actionName)

		server.MustRender(w, r, templates.TEMPLATE_LIST, data.TypeMap)
	}
}

// return the Chosen Object "item" / "box" / "shelf" / "area"
func Element(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thing := r.PathValue("thing")
		thingIDFromRequest := r.PathValue("thingid")

		var err error
		var otherthingLabel string
		var otherthingID string
		var otherthingHref string

		switch thing {
		case "box":
			var box boxes.Box
			box, err = db.BoxById(uuid.FromStringOrNil(thingIDFromRequest))
			if err == nil {
				otherthingLabel = box.Label
				otherthingID = box.ID.String()
				otherthingHref = "/box/" + otherthingID
			}
			break
		case "shelf":
			var shelf *shelves.Shelf
			shelf, err = db.Shelf(uuid.FromStringOrNil(thingIDFromRequest))

			if err == nil {
				otherthingLabel = shelf.Label
				otherthingID = shelf.ID.String()
				otherthingHref = "/shelf/" + otherthingID
			}
			break
		case "area":
			var area areas.Area
			area, err = db.AreaById(uuid.FromStringOrNil(thingIDFromRequest))

			if err == nil {
				otherthingLabel = area.Label
				otherthingID = area.ID.String()
				otherthingHref = "/area/" + otherthingID
			}
			break
		default:
			server.WriteInternalServerError("can't move box to \""+thing+"\"", logg.NewError("can't move box to \""+thing+"\""), w, r)
			return
		}

		if err != nil {
			server.WriteNotFoundError("can't find "+thing+" "+thingIDFromRequest, err, w, r)
		}

		inputElements := pickerInputElements(thing, otherthingID, otherthingHref, otherthingLabel)
		server.TriggerSuccessNotification(w, otherthingLabel+`" was added to the "`+thing+`"`)
		server.WriteFprint(w, inputElements)
		server.WriteFprint(w, `<div id="place-holder" hx-swap-oob="true"></div>`)
	}
}

func pickerInputElements(thing string, otherthingID string, aHref string, otherthingLabel string) string {
	label := `<label for="` + thing + `_id">Is inside of ` + common.ToUpper(thing) + `</label>`
	inputID := `<input  type="text"  name="` + thing + `_id" value="` + otherthingID + `" hidden>`
	a := `<a href="` + aHref + `" class="clickable" hx-boost="true" style="">` + otherthingLabel + `</a>`
	button := `<button id="move-btn" hx-target="#place-holder" hx-post="/addto/` + thing + `" hx-push-url="false"> Add to another ` + common.ToUpper(thing) + ` </button>`
	return label + inputID + a + button
}
