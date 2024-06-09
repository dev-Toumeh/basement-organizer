package routes

import (
	"basement/main/internal/auth"
	"net/http"
)

func RegisterRoutes(db auth.AuthDatabase) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/static"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login", db.LoginHandler)
	http.HandleFunc("/register", db.RegisterHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)
}
