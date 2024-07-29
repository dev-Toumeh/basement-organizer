package auth

import (
	"fmt"
	"log"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, COOKIE_NAME)
	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		log.Printf("LogloutHandler - ok: %v authenticated: %v", ok, authenticated)
		w.WriteHeader(http.StatusBadRequest)
		templates.RenderErrorSnackbar(w, "logout failed")
		return
	}
	session.Values["authenticated"] = false
	session.Save(r, w)
	log.Println("LogoutHandler logged out")

	w.Header().Add("HX-Location", "/")
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "logged out")
}
