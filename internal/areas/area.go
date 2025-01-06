package areas

import (
	"basement/main/internal/common"
	"basement/main/internal/templates"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type AreaDatabase interface {
	CreateArea(newArea Area) (uuid.UUID, error)
	UpdateArea(area Area) error
	DeleteArea(id uuid.UUID) error
	AreaById(id uuid.UUID) (Area, error)
	AreaIDs() ([]uuid.UUID, error)
	AreaListRows(query string, limit int, page int) ([]common.ListRow, error)
	AreaListRowByID(id uuid.UUID) (common.ListRow, error)
	AreaListCounter(searchString string) (count int, err error)
}

type Area struct {
	common.BasicInfo
}

func NewArea() Area {
	b := common.NewBasicInfoWithLabel("Area")
	return Area{b}
}

func (area Area) Map() map[string]any {
	m := area.BasicInfo.Map()
	return m
}

type AreaDetailsPageData struct {
	templates.PageTemplate
	Area
	Edit   bool
	Create bool
}

// NewAreaDetailsPageData returns struct needed for "templates.TEMPLATE_area_DETAILS_PAGE" with default values.
func NewAreaDetailsPageData() (data AreaDetailsPageData) {
	data.PageTemplate = templates.NewPageTemplate()
	return data
}

func areaFromPostFormValue(id uuid.UUID, r *http.Request) Area {
	area := Area{}
	area.BasicInfo = common.BasicInfoFromPostFormValue(id, r)
	return area
}
