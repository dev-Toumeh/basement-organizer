package common

import (
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type ListTemplate struct {
	FormID       string // Unique ID to distinguish between multiple ListTemplate forms.
	FormHXGet    string // Replaces form. Is triggered by search input and pagination buttons. "/boxes, /shelves, ..."
	FormHXPost   string // Replaces form. Is triggered by search input and pagination buttons. "/boxes, /shelves, ..."
	FormHXTarget string // Changes the element which the response will replace.

	SearchInput      bool   // Show search input
	SearchInputLabel string // Label of input
	SearchInputValue string // Current value for search input to "remember" last input after form is replaced

	Pagination        bool // Show pagination buttons
	CurrentPageNumber int  // Sets "page" input element. Used in requests as query params or POST body value.
	Limit             int  // Sets "limit" input element. How many things will be shown or requested. Used in requests as query params or POST body value.
	ShowLimit         bool // Show limit input element.
	PaginationButtons []PaginationButton
	Pages             []PaginationButton
	NextPage          int
	PrevPage          int

	PlaceHolder        bool   // Responsible of rendering the placeholder for the Additional listTemplate (move, add etc..)
	MoveButtonHXTarget string // Change response target. Default: "#move-to-list"
	Count              int
	Rows               []ListRow
	RowAction          bool   // If an action button should be displayed inside each row instead of checkmarks.
	RowActionType      string // the type of the RequestAction (add, move, preview, ...etc)
	RowActionName      string // button and column header name

	AdditionalDataInputs []DataInput // Used for query params or POST body values for new requests from this template.

	AlternativeView bool // Displays Alternative View button. (unused currently).
	RequestOrigin   string
	HideMoveCol     bool
}

func (tmpl ListTemplate) Map() map[string]any {
	return map[string]any{
		"FormID":               tmpl.FormID,
		"FormHXGet":            tmpl.FormHXGet,
		"FormHXPost":           tmpl.FormHXPost,
		"FormHXTarget":         tmpl.FormHXTarget,
		"SearchInput":          tmpl.SearchInput,
		"SearchInputLabel":     tmpl.SearchInputLabel,
		"SearchInputValue":     tmpl.SearchInputValue,
		"Pagination":           tmpl.Pagination,
		"CurrentPageNumber":    tmpl.CurrentPageNumber,
		"Limit":                tmpl.Limit,
		"ShowLimit":            tmpl.ShowLimit,
		"PaginationButtons":    tmpl.PaginationButtons,
		"MoveButtonHXTarget":   tmpl.MoveButtonHXTarget,
		"Rows":                 tmpl.Rows,
		"RowActionType":        tmpl.RowActionType,
		"RowActionName":        tmpl.RowActionName,
		"PlaceHolder":          tmpl.PlaceHolder,
		"AdditionalDataInputs": tmpl.AdditionalDataInputs,
		"AlternativeView":      tmpl.AlternativeView,
		"RequestOrigin":        tmpl.RequestOrigin,
		"HideMoveCol":          tmpl.HideMoveCol,
	}
}

type PaginationButton struct {
	PageNumber int
	Selected   bool
	Disabled   bool
}

type listTemplate2 struct {
	ListTemplate
	ActionInputs []map[string]string
}

type DataInput struct {
	Key   string
	Value any
}

func (tmpl ListTemplate) Render(w http.ResponseWriter) error {
	return templates.SafeRender(w, templates.TEMPLATE_LIST, tmpl)
}

// FilledRows returns ListRows with empty entries filled up to match limit.
//
// listRowsFunc is a DB function like "db.BoxListRows()" and will be called like this internally:
//
//	rows, err := listRowsFunc(searchString, limit, count)
//
// count - The total number of records found from the search query.
func FilledRows(listRowsFunc func(query string, limit int, page int) ([]ListRow, error), searchString string, limit int, pageNr int, count int, listRowOptions ListRowTemplateOptions) ([]ListRow, error) {
	filledRows := make([]ListRow, limit)

	// Fetch the Records from the Database and pack it into map
	rows, err := listRowsFunc(searchString, limit, pageNr)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	if len(rows) > limit {
		return nil, logg.NewError(fmt.Sprintf("found rows (%d) is greater than limit (%d)", len(rows), limit))
	}

	for i, b := range rows {
		filledRows[i] = b
		filledRows[i].ListRowTemplateOptions = listRowOptions
	}

	logg.Debugf("filled rows: %v", filledRows[0])

	// Fill up empty rows to keep same table size
	if count < limit {
		for i := count; i < limit; i++ {
			filledRows[i] = ListRow{}
		}
	}
	return filledRows, nil
}

func ListPageParams(r *http.Request) string {
	query := "query=" + r.FormValue("return:query")
	limit := "limit=" + r.FormValue("return:limit")
	page := "page=" + r.FormValue("return:page")
	return "?" + query + "&" + limit + "&" + page
}

func PickerInputElements(iID string, iValue string, aID string, aHref string, aLabel string) string {
	input := `<input hx-swap-oob="true" type="text" id="` + iID + `_id" name="` + iID + `_id" value="` + iValue + `" hidden>`
	a := ` <a id="` + aID + `" hx-swap-oob="true" href="` + aHref + `" 
			class="clickable"
			hx-boost="true"
			style="">` + aLabel + `</a>`
	return input + a
}

// Database implements partial thing Database (BoxDatabase, ...) interface.
type Database interface {
	BoxListCounter(searchQuery string) (count int, err error)
	ShelfListCounter(searchQuery string) (count int, err error)
	BoxListRows(searchQuery string, limit int, page int) ([]ListRow, error)
	ShelfListRows(searchQuery string, limit int, page int) (shelfRows []ListRow, err error)
}

func ListPageMovePicker(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Uses POST so client can send long list
		// of IDs of boxes inside PostForm body
		if r.Method != http.MethodPost {
			w.Header().Add("Allow", http.MethodPost)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
			return
		}

		moveTo := r.PathValue("thing")

		// Request doesn't come from this move template.
		isRequestFromMovePage := r.FormValue("move") != ""
		// logg.Debugf("isRequestFromMovePage=%v", isRequestFromMovePage)

		var toMove []uuid.UUID
		if isRequestFromMovePage { // IDs are stored as "id-to-be-moved":UUID
			ids := r.PostForm["id-to-be-moved"]
			toMove = make([]uuid.UUID, len(ids))
			for i, id := range ids {
				toMove[i] = uuid.FromStringOrNil(id)
			}
		} else { // IDs are stored as "move:UUID"
			var err error
			toMove, err = server.ParseIDsFromFormWithKey(r.Form, "move")
			if err != nil {
				server.WriteInternalServerError(fmt.Sprintf("can't move boxes %v", toMove), err, w, r)
				return
			}
			if len(toMove) == 0 {
				server.WriteBadRequestError("No box selected to move", nil, w, r)
				return
			}
		}

		searchString := SearchString(r)
		pageNr := ParsePageNumber(r)
		limit := ParseLimit(r)
		origin := ParseOrigin(r)

		additionalData := make([]DataInput, len(toMove))
		for i, id := range toMove {
			additionalData[i] = DataInput{Key: "id-to-be-moved", Value: id.String()}
		}
		if isRequestFromMovePage {
			// Store values to return to the original page where the move was requested.
			additionalData = append(additionalData,
				DataInput{Key: "return:page", Value: r.FormValue("return:page")},
				DataInput{Key: "return:limit", Value: r.FormValue("return:limit")},
				DataInput{Key: "return:query", Value: r.FormValue("return:query")},
			)
		} else {
			additionalData = append(additionalData,
				DataInput{Key: "return:page", Value: pageNr},
				DataInput{Key: "return:limit", Value: limit},
				DataInput{Key: "return:query", Value: searchString},
			)
		}
		// logg.Debugf("additionalData=%v", additionalData)

		listTmpl := ListTemplate{
			FormID:       "list-move",
			FormHXPost:   "/boxes/moveto/" + moveTo,
			FormHXTarget: "this",
			ShowLimit:    env.Config().ShowTableSize(),

			RowAction:            true,
			RowActionType:        "move",
			AdditionalDataInputs: additionalData,
			// I added those extra variables (Naseem)
			PlaceHolder:   false,
			RequestOrigin: origin,
		}
		logg.Debug("move to: " + moveTo)

		// search-input template
		// Clear search when move template is requested the first time.
		if !isRequestFromMovePage {
			searchString = ""
		}
		listTmpl.SearchInput = true
		listTmpl.SearchInputLabel = "Search " + moveTo
		listTmpl.SearchInputValue = searchString

		// pagination
		listTmpl.Pagination = true

		var page int
		// Show first page when move template is requested the first time.
		if isRequestFromMovePage {
			page = pageNr
		} else {
			page = 1
		}
		listTmpl.Limit = limit

		var count int
		var err error
		var rowHXGet string
		switch moveTo {
		case "box":
			rowHXGet = "/box"
			count, err = db.BoxListCounter(searchString)
			break

		case "shelf":
			rowHXGet = "/shelves"
			count, err = db.ShelfListCounter(searchString)
			break

		case "area":
			server.WriteNotImplementedWarning("area", w, r)
			return
		}

		if err != nil {
			server.WriteInternalServerError("can't query "+moveTo, err, w, r)
			return
		}

		// box rows
		var rows []ListRow
		// if there are search results
		if count > 0 {
			rowOptions := ListRowTemplateOptions{
				RowHXGet:              rowHXGet,
				RowAction:             true,
				RowActionName:         "Move here",
				RowActionHXPostWithID: "/boxes/moveto/" + moveTo,
				RowActionHXTarget:     "#list-move",
				RowActionType:         "move",
			}
			switch moveTo {
			case "box":
				rows, err = FilledRows(db.BoxListRows, searchString, limit, page, count, rowOptions)
				break
			case "shelf":
				rows, err = FilledRows(db.ShelfListRows, searchString, limit, page, count, rowOptions)
				break
			}

			if err != nil {
				server.WriteInternalServerError("can't query "+moveTo, err, w, r)
				return
			}
		}

		data := Pagination(map[string]any{}, count, limit, page)
		listTmpl.PaginationButtons = data["Pages"].([]PaginationButton)
		listTmpl.Rows = rows
		err = listTmpl.Render(w)
		if err != nil {
			server.WriteInternalServerError("can't render move page", err, w, r)
			return
		}
	}
}

func ListPageMovePickerConfirm(DBMoveToThing func(thing1 uuid.UUID, thing2 uuid.UUID) error, redirectURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var notifications server.Notifications
		notifications = server.MoveThingToThing(w, r, DBMoveToThing)
		params := ListPageParams(r)
		server.RedirectWithNotifications(w, redirectURL+params, notifications)
	}
}
