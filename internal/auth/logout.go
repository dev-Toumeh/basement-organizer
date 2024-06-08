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
		fmt.Fprint(w, "logout failed")
	}
	session.Values["authenticated"] = false
	session.Save(r, w)
	log.Println("LogoutHandler logged out")
	fmt.Fprint(w, "logged out")
}
