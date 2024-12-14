package areas

import (
	"basement/main/internal/common"
)

type Area struct {
	common.BasicInfo
}

func (area Area) Map() map[string]any {
	m := area.BasicInfo.Map()
	return m
}
