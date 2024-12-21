package areas

import (
	"basement/main/internal/common"

	"github.com/gofrs/uuid/v5"
)

type AreaDatabase interface {
	CreateArea(newArea Area) (uuid.UUID, error)
	UpdateArea(box Area) error
	DeleteArea(boxId uuid.UUID) error
	AreaById(id uuid.UUID) (Area, error)
	AreaIDs() ([]uuid.UUID, error)
	AreaListRows(query string, limit int, page int) ([]common.ListRow, error)
	AreaListRowByID(id uuid.UUID) (common.ListRow, error)
	AreaListCounter(searchString string) (count int, err error)
}

type Area struct {
	common.BasicInfo
}

func (area Area) Map() map[string]any {
	m := area.BasicInfo.Map()
	return m
}
