package items

import (
	"basement/main/internal/common"
)

type SearchItemData struct {
	Query          string
	TotalCount     int
	Records        []common.ListRow
	PaginationData []PaginationData
}

type PaginationData struct {
	Offset     int
	PageNumber int
}

type SearchInputTemplate struct {
	SearchInputLabel    string
	SearchInputHxPost   string
	SearchInputHxTarget string
	SearchInputValue    string
}

func NewSearchInputTemplate() *SearchInputTemplate {
	return &SearchInputTemplate{SearchInputLabel: "Search",
		SearchInputHxPost:   "/api/v1/implement-me",
		SearchInputHxTarget: "#item-list-body",
		SearchInputValue:    "",
	}
}

func NewSearchItemInputTemplate() *SearchInputTemplate {
	return &SearchInputTemplate{SearchInputLabel: "Search items", SearchInputHxPost: "/api/v1/search/item", SearchInputHxTarget: "#item-list-body"}
}

func (tmpl *SearchInputTemplate) Map() map[string]any {
	return map[string]any{
		"SearchInputLabel":    tmpl.SearchInputLabel,
		"SearchInputHxPost":   tmpl.SearchInputHxPost,
		"SearchInputHxTarget": tmpl.SearchInputHxTarget,
		"SearchInputValue":    tmpl.SearchInputValue,
	}
}

// update the item based on ID
