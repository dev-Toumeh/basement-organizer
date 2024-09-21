package auth

import (
	"basement/main/internal/logg"
	"basement/main/internal/server"
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

	username, _ := UserSessionData(r)
	server.RedirectWithSuccessNotification(w, "/", fmt.Sprintf("Good bye %s", username))

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "logged out")
}
