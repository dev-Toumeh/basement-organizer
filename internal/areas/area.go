package areas

import (
	"basement/main/internal/common"
	"basement/main/internal/server"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type AreaDatabase interface {
	CreateArea(newArea Area) (uuid.UUID, error)
	UpdateArea(area Area, ignorePicture bool) error
	DeleteArea(id uuid.UUID) error
	AreaById(id uuid.UUID) (Area, error)
	AreaIDs() ([]uuid.UUID, error)
	AreaListRows(query string, limit int, page int) ([]common.ListRow, error)
	AreaListRowByID(id uuid.UUID) (common.ListRow, error)
	AreaListCounter(searchString string) (count int, err error)
	BoxListCounter(searchQuery string) (count int, err error)
	ShelfListCounter(searchQuery string) (count int, err error)
	BoxListRows(searchQuery string, limit int, page int) ([]common.ListRow, error)
	ShelfListRows(searchQuery string, limit int, page int) (shelfRows []common.ListRow, err error)
	InnerListRowsFrom2(belongsToTable string, belongsToTableID uuid.UUID, listRowsTable string) ([]common.ListRow, error)
}

type Area struct {
	common.BasicInfo
}

func NewArea() Area {
	b := common.NewBasicInfoWithLabel("Area")
	return Area{BasicInfo: b}
}

func (area Area) Map() map[string]any {
	m := area.BasicInfo.Map()
	return m
}

func areaFromPostFormValue(id uuid.UUID, r *http.Request) (area Area, ignorePicture bool) {
	ignorePicture = server.ParseIgnorePicture(r)
	area.BasicInfo = common.BasicInfoFromPostFormValue(id, r, false)
	return area, ignorePicture
}
