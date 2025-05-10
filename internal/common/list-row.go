package common

import (
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/gofrs/uuid/v5"
)

// ListRow is a single row entry used for list templates.
type ListRow struct {
	ID             uuid.UUID
	Label          string
	Description    string
	BoxID          uuid.UUID
	BoxLabel       string
	ShelfID        uuid.UUID
	ShelfLabel     string
	AreaID         uuid.UUID
	AreaLabel      string
	PreviewPicture string

	ListRowTemplateOptions
}

func (row ListRow) Map() map[string]any {
	m := map[string]interface{}{
		"ID":             row.ID,
		"Label":          row.Label,
		"Description":    row.Description,
		"BoxID":          row.BoxID,
		"BoxLabel":       row.BoxLabel,
		"ShelfID":        row.ShelfID,
		"ShelfLabel":     row.ShelfLabel,
		"AreaID":         row.AreaID,
		"AreaLabel":      row.AreaLabel,
		"PreviewPicture": row.PreviewPicture,
	}
	maps.Copy(row.ListRowTemplateOptions.Map(), m)
	return m
}

func (row ListRow) String() string {
	if env.Development() {
		row.PreviewPicture = ShortenPictureForLogs(row.PreviewPicture)
	}

	data, err := json.Marshal(row)
	if err != nil {
		logg.Err("Can't JSON box to string:", err)
		return ""
	}
	s := fmt.Sprintf("%s", data)
	return s
}

func AddRowOptionsToListRows(rows []ListRow, opts ListRowTemplateOptions) []ListRow {
	for i := range rows {
		rows[i].ListRowTemplateOptions = opts
	}
	return rows
}

func AddRowOptionsToListRows2(rows []ListRow, opts ListRowTemplateOptions) {
	for i := range rows {
		rows[i].ListRowTemplateOptions = opts
	}
}

type ListRowTemplateOptions struct {
	RowHXGet                          string // hx-get="{{ .RowHXGet }}/{{.ID}}"
	RowAction                         bool   // If an action button should be displayed inside each row instead of checkmarks.
	RowActionType                     string // the type of the RequestAction (add, move, preview, ...etc)
	RowActionHXPost                   string // hx-post="{{ .ActionHXPost }}"
	RowActionHXPostWithID             string // hx-post="{{ .ActionHXPost }}/{{.ID}}"
	RowActionHXPostWithIDAsQueryParam string // hx-post="{{ .ActionHXPost }}?id={{.ID}}"
	RowActionName                     string // button and column header name
	RowActionHXTarget                 string // hx-target="{{$RowActionHXTarget}}" in action button
	HideMoveCol                       bool
	HideBoxLabel                      bool
	HideShelfLabel                    bool
	HideAreaLabel                     bool
}

func (row ListRowTemplateOptions) Map() map[string]any {
	return map[string]interface{}{
		"RowHXGet":                          row.RowHXGet,
		"RowActionType":                     row.RowActionType,
		"RowAction":                         row.RowAction,
		"RowActionHXPost":                   row.RowActionHXPost,
		"RowActionHXPostWithID":             row.RowActionHXPostWithID,
		"RowActionHXPostWithIDAsQueryParam": row.RowActionHXPostWithIDAsQueryParam,
		"RowActionName":                     row.RowActionName,
		"RowActionHXTarget":                 row.RowActionHXTarget,
		"HideMoveCol":                       row.HideMoveCol,
		"HideBoxLabel":                      row.HideBoxLabel,
		"HideShelfLabel":                    row.HideShelfLabel,
		"HideAreaLabel":                     row.HideAreaLabel,
	}
}
