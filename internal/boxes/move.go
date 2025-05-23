package boxes

import (
	"basement/main/internal/common"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

const (
	PICKER_TYPE_ADDTO = iota
	PICKER_TYPE_MOVE
)

// shows the page where client can choose where to move the box.
//
// thing = "item" / "box" / "shelf" / "area"
//
// pickerType:
//
//	PICKER_TYPE_ADDTO: Current box is not created yet. No move operation will be execuded.
//	PICKER_TYPE_MOVE: Current box is already created. Will use execute a move operation.
func BoxPicker(pickerType int, db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := common.InitData(r, true)
		thing := r.PathValue("thing")
		errMsgForUser := "Can't move " + thing
		id := server.ValidID(w, r, errMsgForUser)
		if id.IsNil() {
			return
		}

		var post string
		var actionName string
		if pickerType == PICKER_TYPE_ADDTO {
			post = "/box/" + id.String() + "/addto/" + thing
			actionName = "Add to"
		} else if pickerType == PICKER_TYPE_MOVE {
			post = "/box/" + id.String() + "/moveto/" + thing
			actionName = "Move here"
		}
		rowActionType := "move"

		var err error
		var count int

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
					RowActionHXTarget:     "#" + thing + "_id",
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
					RowActionHXTarget:     "#" + thing + "_id",
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
					RowActionHXTarget:     "#" + thing + "_id",
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
		data.SetShowLimit(env.CurrentConfig().ShowTableSize())

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

// boxMoveConfirm handles data after a box move action is clicked from boxPageMove().
//
// thing = "item" / "box" / "shelf" / "area"
//
// pickerType:
//
//	PICKER_TYPE_ADDTO: Current box is not created yet. No move operation will be execuded.
//	PICKER_TYPE_MOVE: Current box is already created. Will use execute a move operation.
func BoxPickerConfirm(pickerType int, db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thing := r.PathValue("thing")
		moveToThingID := r.PathValue("thingid")
		boxID := uuid.FromStringOrNil(r.PathValue("id"))

		var err1 error
		var err2 error
		var otherThingLabel string
		var otherThingElementID string
		var otherThingHref string
		switch thing {
		case "box":
			var outerbox Box

			// no move necessary
			if pickerType == PICKER_TYPE_ADDTO {
				err1 = nil
			} else if pickerType == PICKER_TYPE_MOVE {
				err1 = db.MoveBoxToBox(boxID, uuid.FromStringOrNil(moveToThingID))
			}

			outerbox, err2 = db.BoxById(uuid.FromStringOrNil(moveToThingID))
			if err2 == nil {
				otherThingLabel = outerbox.Label
				otherThingElementID = "outerbox-link"
				otherThingHref = "/box/" + moveToThingID
			}
			break
		case "shelf":
			// no move necessary
			if pickerType == PICKER_TYPE_ADDTO {
				err1 = nil
			} else if pickerType == PICKER_TYPE_MOVE {
				err1 = db.MoveBoxToShelf(boxID, uuid.FromStringOrNil(moveToThingID))
			}

			if err1 == nil {
				otherThingLabel = moveToThingID
				otherThingElementID = "shelf-link"
				otherThingHref = "/shelves/" + moveToThingID
			}
			break
		case "area":
			if pickerType == PICKER_TYPE_ADDTO {
				err1 = nil
			} else if pickerType == PICKER_TYPE_MOVE {
				err1 = db.MoveBoxToArea(boxID, uuid.FromStringOrNil(moveToThingID))
			}
			if err1 == nil {
				otherThingLabel = moveToThingID
				otherThingElementID = "area-link"
				otherThingHref = "/area/" + moveToThingID
			}
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
			server.WriteNotFoundError("can't find "+thing+" "+moveToThingID, err2, w, r)
		}

		inputElements := common.PickerInputElements(thing, moveToThingID, otherThingElementID, otherThingHref, otherThingLabel)
		server.TriggerSuccessNotification(w, `moved"`+boxID.String()+`" to "`+moveToThingID+`"`)
		server.WriteFprint(w, inputElements)
		server.WriteFprint(w, `<div id="place-holder" hx-swap-oob="true"></div>`)
	}
}
