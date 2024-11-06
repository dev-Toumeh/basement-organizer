package common

import (
	"basement/main/internal/templates"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type ListTemplate struct {
	FormHXGet   string // /boxes, /shelves
	RowHXGet    string // hx-get="{{ .RowHXGet }}/{{.ID}}"
	SearchInput any
	PaginationButtons []PaginationButton

	Rows                              []ListRow
	RowAction                         bool   // If an action button should be displayed inside each row instead of checkmarks.
	RowActionHXPost                   string // hx-post="{{ .ActionHXPost }}"
	RowActionHXPostWithID             string // hx-post="{{ .ActionHXPost }}/{{.ID}}"
	RowActionHXPostWithIDAsQueryParam string // hx-post="{{ .ActionHXPost }}?id={{.ID}}"
	RowActionName                     string // button and column header name
	RowActionHXTarget                 string // hx-target="{{$RowActionHXTarget}}" in action button

	AdditionalDataInputs      []DataInput
	AdditionalDataInputValues []string
	AdditionalDataInputName   string
	ReturnSelectedRowInput    bool
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
	Value string
}

func (tmpl ListTemplate) Render(w http.ResponseWriter) error {

	// t2 := listTemplate2{ListTemplate: tmpl}
	//
	// t2.ActionInputs = make([]map[string]string, len(t2.AdditionalDataInputValues))
	// for i, v := range t2.AdditionalDataInputValues {
	// 	t2.ActionInputs[i] = map[string]string{
	// 		"Key":   t2.AdditionalDataInputName,
	// 		"Value": v,
	// 	}
	// }

	// logg.Debug(tmpl.Rows)
	// d := ListTemplate{HXGet: "/boxes"}
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
