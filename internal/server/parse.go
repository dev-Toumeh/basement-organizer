package server

import (
	"basement/main/internal/logg"
	"net/http"
)

func ParseIgnorePicture(r *http.Request) (ignorePicture bool) {
	updatePicture := r.PostFormValue("updatepicture")
	logg.Debug("updatepicture=" + updatePicture)
	if updatePicture != "on" {
		return true
	} else {
		return false
	}
}
