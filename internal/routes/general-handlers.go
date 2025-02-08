package routes

import (
	"fmt"
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

		data.SetFormHXPost("/addto/" + thing)
		data.SetFormID(thing + "-list")
		data.SetFormHXTarget("#" + thing + "-list")

		data.SetSearchInput(true)
		data.SetSearchInputLabel("Search " + thing)

		data.SetRowAction(true)
		data.SetRowActionName(actionName)

		server.MustRender(w, r, templates.TEMPLATE_LIST, data.TypeMap)
	}
}

// // Handles the HTTP request, extracts parameters,
func Element(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thing := r.PathValue("thing")
		thingID := uuid.FromStringOrNil(r.PathValue("thingid"))
		if thingID == uuid.Nil {
			server.WriteNotFoundError("Invalid or missing thing ID", fmt.Errorf("empty thing ID"), w, r)
			return
		}

		data, err := fetchObjects(db, thing, thingID)
		if err != nil {
			server.WriteNotFoundError(fmt.Sprintf("No matching objects found for %s", thing), err, w, r)
			return
		}

		var (
			rendered string
			boxObj   boxes.Box
			shelfObj *shelves.Shelf
			areaObj  areas.Area
		)

		// Default to empty values if not found
		boxFound, shelfFound, areaFound := false, false, false
		for _, m := range data {
			if val, ok := m["box"]; ok {
				boxObj = val.(boxes.Box)
				boxFound = true
			}
			if val, ok := m["shelf"]; ok {
				shelfObj = val.(*shelves.Shelf)
				shelfFound = true
			}
			if val, ok := m["area"]; ok {
				areaObj = val.(areas.Area)
				areaFound = true
			}
		}

		// Always render three inputs
		if boxFound {
			rendered += renderPicker("box", boxObj.ID, boxObj.Label)
		} else {
			rendered += renderEmptyPicker("box")
		}
		if shelfFound {
			rendered += renderPicker("shelf", shelfObj.ID, shelfObj.Label)
		} else {
			rendered += renderEmptyPicker("shelf")
		}
		if areaFound {
			rendered += renderPicker("area", areaObj.ID, areaObj.Label)
		} else {
			rendered += renderEmptyPicker("area")
		}

		server.TriggerSuccessNotification(w, "Objects were added successfully.")
		server.WriteFprint(w, rendered)
		server.WriteFprint(w, `<div id="place-holder" hx-swap-oob="true"></div>`)
	}
}

// Fetch the requested object and any related objects
func fetchObjects(db *database.DB, thing string, thingID uuid.UUID) ([]map[string]interface{}, error) {
	var obj interface{}
	var err error
	response := []map[string]interface{}{}

	switch thing {
	case "box":
		obj, err = db.BoxById(thingID)
	case "shelf":
		obj, err = db.Shelf(thingID)
	case "area":
		obj, err = db.AreaById(thingID)
	default:
		return nil, fmt.Errorf("unknown type: %s", thing)
	}

	if err != nil {
		return nil, err
	}

	if thing == "box" {
		box := obj.(boxes.Box)
		response = append(response, map[string]interface{}{"box": box})

		if box.ShelfID != uuid.Nil {
			if shelf, err := db.Shelf(box.ShelfID); err == nil {
				response = append(response, map[string]interface{}{"shelf": shelf})

				if shelf.AreaID != uuid.Nil {
					if area, err := db.AreaById(shelf.AreaID); err == nil {
						response = append(response, map[string]interface{}{"area": area})
					}
				}
			}
		}
	} else if thing == "shelf" {
		shelf := obj.(*shelves.Shelf)
		response = append(response, map[string]interface{}{"shelf": shelf})

		if shelf.AreaID != uuid.Nil {
			if area, err := db.AreaById(shelf.AreaID); err == nil {
				response = append(response, map[string]interface{}{"area": area})
			}
		}
	} else if thing == "area" {
		response = append(response, map[string]interface{}{"area": obj})
	}

	return response, nil
}

func renderPicker(thing string, id uuid.UUID, label string) string {
	targetID := thing + "-target"
	htmlLabel := `<label for="` + thing + `_id">Is inside of ` + common.ToUpper(thing) + `</label>`
	hidden := `<input type="hidden" name="` + thing + `_id" value="` + id.String() + `">`
	a := `<a href="/` + thing + `/` + id.String() + `" class="clickable" hx-boost="true">` + label + `</a>`
	button := `<button hx-post="/addto/` + thing + `" hx-target="#place-holder" hx-swap="innerHTML"
                hx-push-url="false" type="button">Add to another ` + common.ToUpper(thing) + `</button>`

	return `<div id="` + targetID + `" hx-swap-oob="true">` + htmlLabel + hidden + a + button + `</div>`
}

func renderEmptyPicker(thing string) string {
	targetID := thing + "-target"
	htmlLabel := `<label for="` + thing + `_id">Is inside of ` + common.ToUpper(thing) + `</label>`
	emptySpan := `<span id='outerbox-link'>None</span>`
	button := `<button hx-post="/addto/` + thing + `" hx-target="#place-holder" hx-swap="innerHTML"
                hx-push-url="false" type="button">Add to ` + common.ToUpper(thing) + `</button>`

	return `<div id="` + targetID + `" hx-swap-oob="true">` + htmlLabel + emptySpan + button + `</div>`
}
