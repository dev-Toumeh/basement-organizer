package boxes

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

// shows the page where client can choose where to move the box.
// thing = item / box / shelf / area
func BoxPageMove(thing string, db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := common.InitData(r)

		var err error
		var count int

		switch thing {
		case "box":
			count, err = db.BoxListCounter("")
			if err != nil {
				server.WriteInternalServerError("no box list counter", err, w, r)
				return
			}
			data.SetCount(count)

			var boxes []common.ListRow
			if count > 0 {
				boxes, err = common.FilledRows(db.BoxListRows, data.GetSearchInputValue(), data.GetLimit(), data.GetPageNumber(), count)
				if err != nil {
					server.WriteInternalServerError("cant query "+thing+" please comeback later", err, w, r)
				}
			}
			data.SetRows(boxes)

		case "shelf":
			count, err = db.ShelfListCounter("")
			if err != nil {
				server.WriteInternalServerError("no shelf list counter", err, w, r)
				return
			}
			data.SetCount(count)

			var shelves []common.ListRow
			if count > 0 {
				shelves, err = common.FilledRows(db.ShelfListRows, data.GetSearchInputValue(), data.GetLimit(), data.GetPageNumber(), count)
				if err != nil {
					server.WriteInternalServerError("cant query "+thing+" please comeback later", err, w, r)
				}
			}
			data.SetRows(shelves)
			break

		default:
			server.WriteInternalServerError("can't move box to \""+thing+"\"", logg.NewError("can't move box to \""+thing+"\""), w, r)
			return
		}

		data = common.Pagination2(data)
		data.SetShowLimit(env.Config().ShowTableSize())

		errMsgForUser := "Can't move " + thing
		id := server.ValidID(w, r, errMsgForUser)
		if id.IsNil() {
			return
		}
		data.SetFormHXPost("/box/" + id.String() + "/moveto/" + thing)
		data.SetFormID(thing + "-list")
		data.SetFormHXTarget("#place-holder")

		data.SetSearchInput(true)
		data.SetSearchInputLabel("Search " + thing)

		data.SetRowAction(true)
		data.SetRowActionType("move")
		data.SetRowActionHXTarget("#" + thing + "_id")
		data.SetRowActionName("Move here")
		data.SetRowActionHXPostWithID("/box/" + id.String() + "/moveto/" + thing)

		server.MustRender(w, r, templates.TEMPLATE_LIST, data.TypeMap)
	}
}

// boxMoveConfirm handles data after a box move action is clicked from boxPageMove().
// thing = item / box / shelf / area
func BoxMoveConfirm(thing string, db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moveToThingID := r.PathValue("value")
		boxID := uuid.FromStringOrNil(r.PathValue("id"))

		var err1 error
		var err2 error
		var otherThingLabel string
		var otherThingElementID string
		var otherThingHref string
		switch thing {
		case "box":
			err1 = db.MoveBoxToBox(boxID, uuid.FromStringOrNil(moveToThingID))
			var outerbox Box
			outerbox, err2 = db.BoxById(uuid.FromStringOrNil(moveToThingID))
			if err2 == nil {
				otherThingLabel = outerbox.Label
				otherThingElementID = "outerbox-link"
				otherThingHref = "/boxes/" + moveToThingID
			}
			break
		case "shelf":
			err1 = db.MoveBoxToShelf(boxID, uuid.FromStringOrNil(moveToThingID))
			// var outerbox Box
			// outerbox, err2 = db.BoxById(uuid.FromStringOrNil(moveToThingID))
			if err1 == nil {
				otherThingLabel = moveToThingID
				otherThingElementID = "shelf-link"
				otherThingHref = "/shelves/" + moveToThingID
			}
			break
		case "area":
			break
		default:
			server.WriteInternalServerError("can't move box to \""+thing+"\"", logg.NewError("can't move box to \""+thing+"\""), w, r)
			return
		}
		if err1 != nil {
			logg.Err(err1)
			server.WriteBadRequestError(`can't move "`+boxID.String()+`" to "`+moveToThingID+`"`, err1, w, r)
			return
		}
		if err2 != nil {
			server.WriteNotFoundError("can't find "+thing+" "+moveToThingID, err1, w, r)
		}

		inputThingID := `<input hx-swap-oob="true" type="text" id="` + thing + `_id" name="` + thing + `_id" value="` + moveToThingID + `" readonly>`
		aThingLabel := `
			<a id="` + otherThingElementID + `" hx-swap-oob="true" href="` + otherThingHref + `" 
			class="clickable"
			hx-boost="true"
			style="">` + otherThingLabel + `</a>`

		server.TriggerSuccessNotification(w, `moved"`+boxID.String()+`" to "`+moveToThingID+`"`)
		server.WriteFprint(w, inputThingID)
		server.WriteFprint(w, aThingLabel)
		server.WriteFprint(w, `<div id="place-holder" hx-swap-oob="true"></div>`)
	}
}

// BoxesPageMove shows the page where client can choose where to move the selected boxes from the boxes page.
func BoxesPageMove(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.MustRender(w, r, "move-to", nil)
		// server.WriteNotImplementedWarning("Move multiple boxes page", w, r)
	}
}

// MoveBox moves a box to another box. For direct API calls.
func MoveBox(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.FromStringOrNil(r.PathValue("id"))
		moveToBoxID := uuid.FromStringOrNil(r.PathValue("toid"))
		err := db.MoveBoxToBox(id, moveToBoxID)
		if err != nil {
			server.WriteBadRequestError("can't move box", err, w, r)
			logg.Err(err)
		} else {
			w.WriteHeader(200)
		}
	}
}

func MoveBoxesToBoxHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notifications := MoveBoxesToBox(w, r, db)
		params := common.ListPageParams(r)
		server.RedirectWithNotifications(w, "/boxes"+params, notifications)
	}
}

func MoveBoxesToBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) server.Notifications {
	r.ParseForm()
	moveToBoxID := server.ValidID(w, r, "can't move boxes invalid id")
	if moveToBoxID == uuid.Nil {
		return server.Notifications{}
	}

	parseIDs := r.PostForm["id-to-be-moved"]
	ids := make([]uuid.UUID, len(parseIDs))

	logg.Debug(len(parseIDs))

	notifications := server.Notifications{}
	for i, v := range parseIDs {
		logg.Debug(v)
		id := uuid.FromStringOrNil(v)
		ids[i] = id
		err := db.MoveBoxToBox(id, moveToBoxID)
		if err != nil {
			notifications.AddError(fmt.Sprintf(`can't move "%s" to "%s"`, ids[i].String(), moveToBoxID.String()))
			logg.Err(err)
		} else {
			notifications.AddSuccess(fmt.Sprintf(`moved "%s" to "%s"`, ids[i].String(), moveToBoxID.String()))
		}
	}
	return notifications
}

func MoveBoxesToBoxAPI(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notifications := MoveBoxesToBox(w, r, db)
		server.TriggerNotifications(w, notifications)
	}
}

// GetMoveBoxesPage handles list form for moving things.
func GetMoveBoxesPage(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Uses POST so client can send long list
		// of IDs of boxes inside PostForm body
		if r.Method != http.MethodPost {
			w.Header().Add("Allow", http.MethodPost)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
			return
		}

		// Request doesn't come from this move template.
		isRequestFromMovePage := r.FormValue("move") != ""

		var toMove []uuid.UUID
		if isRequestFromMovePage { // IDs are stored as "id-to-be-moved":UUID
			ids := r.PostForm["id-to-be-moved"]
			toMove = make([]uuid.UUID, len(ids))
			for i, id := range ids {
				toMove[i] = uuid.FromStringOrNil(id)
			}
		} else { // IDs are stored as "move:UUID"
			var err error
			toMove, err = common.ParseIDsFromFormWithKey(r.Form, "move")
			if err != nil {
				server.WriteInternalServerError(fmt.Sprintf("can't move boxes %v", toMove), err, w, r)
				return
			}
			if len(toMove) == 0 {
				server.WriteBadRequestError("No box selected to move", nil, w, r)
				return
			}
		}

		searchString := common.SearchString(r)
		pageNr := common.ParsePageNumber(r)
		limit := common.ParseLimit(r)

		additionalData := make([]common.DataInput, len(toMove))
		for i, id := range toMove {
			additionalData[i] = common.DataInput{Key: "id-to-be-moved", Value: id.String()}
		}
		if isRequestFromMovePage {
			// Store values to return to the original page where the move was requested.
			additionalData = append(additionalData,
				common.DataInput{Key: "return:page", Value: r.FormValue("return:page")},
				common.DataInput{Key: "return:limit", Value: r.FormValue("return:limit")},
				common.DataInput{Key: "return:query", Value: r.FormValue("return:query")},
			)
		} else {
			additionalData = append(additionalData,
				common.DataInput{Key: "return:page", Value: pageNr},
				common.DataInput{Key: "return:limit", Value: limit},
				common.DataInput{Key: "return:query", Value: searchString},
			)
		}

		listTmpl := common.ListTemplate{
			FormID:       "list-move",
			FormHXPost:   "/get-boxes-move-page",
			FormHXTarget: "this",
			RowHXGet:     "/boxes",
			ShowLimit:    env.Config().ShowTableSize(),

			RowAction:             true,
			RowActionName:         "Move here",
			RowActionHXPostWithID: "/boxes/moveto/box",
			RowActionHXTarget:     "#list-move",
			AdditionalDataInputs:  additionalData,
			// I added those extra variables (Naseem)
			PlaceHolder:   false,
			RowActionType: "move",
		}

		// search-input template
		// Clear search when move template is requested the first time.
		if !isRequestFromMovePage {
			searchString = ""
		}
		listTmpl.SearchInput = true
		listTmpl.SearchInputLabel = "Search boxes"
		listTmpl.SearchInputValue = searchString

		// pagination
		listTmpl.Pagination = true

		count, err := db.BoxListCounter(searchString)
		if err != nil {
			server.WriteInternalServerError("cant query boxes", err, w, r)
			return
		}

		var page int
		// Show first page when move template is requested the first time.
		if isRequestFromMovePage {
			page = pageNr
		} else {
			page = 1
		}
		data := common.Pagination(map[string]any{}, count, limit, page)

		listTmpl.Limit = limit
		listTmpl.PaginationButtons = data["Pages"].([]common.PaginationButton)

		// box rows
		var boxes []common.ListRow
		// if there are search results
		if count > 0 {
			boxes, err = common.FilledRows(db.BoxListRows, searchString, limit, page, count)
			if err != nil {
				server.WriteInternalServerError("can't query boxes", err, w, r)
				return
			}
		}
		listTmpl.Rows = boxes
		err = listTmpl.Render(w)
		if err != nil {
			server.WriteInternalServerError("can't render move page", err, w, r)
			return
		}
	}
}
