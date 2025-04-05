package common

import (
	"basement/main/internal/logg"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
)

// BasicInfo is present in item, box, shelf and area
type BasicInfo struct {
	ID             uuid.UUID
	Label          string
	Description    string
	Picture        string
	PreviewPicture string
	QRCode         string
}

func (b BasicInfo) Map() map[string]any {
	return map[string]interface{}{
		"ID":             b.ID,
		"Label":          b.Label,
		"Description":    b.Description,
		"Picture":        b.Picture,
		"PreviewPicture": b.PreviewPicture,
		"QRCode":         b.QRCode,
	}
}

func NewBasicInfo() BasicInfo {
	return BasicInfo{ID: uuid.Must(uuid.NewV4())}.MakeLabelWithTime("thing")
}

func NewBasicInfoWithLabel(label string) BasicInfo {
	return BasicInfo{ID: uuid.Must(uuid.NewV4())}.MakeLabelWithTime(label)
}

func (b BasicInfo) MakeLabelWithTime(label string) BasicInfo {
	t := time.Now().Format("2006-01-02_15_04_05")
	b.Label = fmt.Sprintf("%s_%s", label, t)
	return b
}

func BasicInfoFromPostFormValue(id uuid.UUID, r *http.Request, ignorePicture bool) BasicInfo {
	info := BasicInfo{}
	info.ID = id
	info.Label = r.PostFormValue("label")
	info.Description = r.PostFormValue("description")
	logg.Debugf("ignorePicture=%v", ignorePicture)
	if !ignorePicture {
		info.Picture = ParsePicture(r)
	}

	info.QRCode = r.PostFormValue("qrcode")
	return info
}
