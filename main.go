package main

import (
	"basement/main/internal/auth"
	"basement/main/internal/routes"
	"log"
	"net/http"
)

func main() {
	var db *auth.JsonDB
	var err error

	db, err = auth.CreateJsonDB()
	if err != nil {
		log.Fatalf("Can't create DB, shutting server down")
	}
	routes.RegisterRoutes(db)

	http.ListenAndServe("localhost:8000", nil)
}
