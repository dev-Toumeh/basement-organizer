package auth

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"fmt"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, COOKIE_NAME)
	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		logg.Debugf("LogloutHandler - ok: %v authenticated: %v", ok, authenticated)
		w.WriteHeader(http.StatusBadRequest)
		templates.RenderErrorNotification(w, "logout failed")
		return
	}
	session.Values["authenticated"] = false
	session.Save(r, w)
	logg.Info("LogoutHandler logged out")

	w.Header().Add("HX-Location", "/")
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "logged out")
}
