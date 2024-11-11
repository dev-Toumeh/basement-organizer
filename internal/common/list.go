package common

import (
	"basement/main/internal/templates"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type ListTemplate struct {
	FormID       string // Unique ID to distinguish between multiple ListTemplate forms.
	FormHXGet    string // Replaces form. Is triggered by search input and pagination buttons. "/boxes, /shelves, ..."
	FormHXPost   string // Replaces form. Is triggered by search input and pagination buttons. "/boxes, /shelves, ..."
	FormHXTarget string // Changes the element which the response will replace.
	RowHXGet     string // hx-get="{{ .RowHXGet }}/{{.ID}}"

	SearchInput      bool   // Show search input
	SearchInputLabel string // Label of input
	SearchInputValue string // Current value for search input to "remember" last input after form is replaced

	Pagination        bool // Show pagination buttons
	CurrentPageNumber int  // Sets "page" input element. Used in requests as query params or POST body value.
	Limit             int  // Sets "limit" input element. How many things will be shown or requested. Used in requests as query params or POST body value.
	ShowLimit         bool // Show limit input element.
	PaginationButtons []PaginationButton

	MoveButtonHXTarget string // Change response target. Default: "#move-to-list"

	Rows                              []ListRow
	RowAction                         bool   // If an action button should be displayed inside each row instead of checkmarks.
	RowActionHXPost                   string // hx-post="{{ .ActionHXPost }}"
	RowActionHXPostWithID             string // hx-post="{{ .ActionHXPost }}/{{.ID}}"
	RowActionHXPostWithIDAsQueryParam string // hx-post="{{ .ActionHXPost }}?id={{.ID}}"
	RowActionName                     string // button and column header name
	RowActionHXTarget                 string // hx-target="{{$RowActionHXTarget}}" in action button

	AdditionalDataInputs []DataInput // Used for query params or POST body values for new requests from this template.
}

func (tmpl ListTemplate) Map() map[string]any {
	return map[string]any{
		"FormID":                            tmpl.FormID,
		"FormHXGet":                         tmpl.FormHXGet,
		"FormHXPost":                        tmpl.FormHXPost,
		"FormHXTarget":                      tmpl.FormHXTarget,
		"RowHXGet":                          tmpl.RowHXGet,
		"SearchInput":                       tmpl.SearchInput,
		"SearchInputLabel":                  tmpl.SearchInputLabel,
		"SearchInputValue":                  tmpl.SearchInputValue,
		"Pagination":                        tmpl.Pagination,
		"CurrentPageNumber":                 tmpl.CurrentPageNumber,
		"Limit":                             tmpl.Limit,
		"ShowLimit":                         tmpl.ShowLimit,
		"PaginationButtons":                 tmpl.PaginationButtons,
		"MoveButtonHXTarget":                tmpl.MoveButtonHXTarget,
		"Rows":                              tmpl.Rows,
		"RowAction":                         tmpl.RowAction,
		"RowActionHXPost":                   tmpl.RowActionHXPost,
		"RowActionHXPostWithID":             tmpl.RowActionHXPostWithID,
		"RowActionHXPostWithIDAsQueryParam": tmpl.RowActionHXPostWithIDAsQueryParam,
		"RowActionName":                     tmpl.RowActionName,
		"RowActionHXTarget":                 tmpl.RowActionHXTarget,
		"AdditionalDataInputs":              tmpl.AdditionalDataInputs,
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

// ListRow is a single row entry used for list templates.
type ListRow struct {
	ID             uuid.UUID
	Label          string
	BoxID          uuid.UUID
	BoxLabel       string
	ShelfID        uuid.UUID
	ShelfLabel     string
	AreaID         uuid.UUID
	AreaLabel      string
	PreviewPicture string
}

func (row *ListRow) Map() map[string]any {
	return map[string]interface{}{
		"ID":             row.ID,
		"Label":          row.Label,
		"BoxID":          row.BoxID,
		"BoxLabel":       row.BoxLabel,
		"ShelfID":        row.ShelfID,
		"ShelfLabel":     row.ShelfLabel,
		"AreaID":         row.AreaID,
		"AreaLabel":      row.AreaLabel,
		"PreviewPicture": row.PreviewPicture,
	}
}
