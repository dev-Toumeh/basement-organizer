package main

import (
	"basement/main/internal/database"
	"basement/main/internal/routes"
	"log"
	"net/http"
)

func main() {
	var db *database.JsonDB
	var err error

	db, err = database.CreateJsonDB()
	if err != nil {
		log.Fatalf("Can't create DB, shutting server down")
	}
	routes.RegisterRoutes(db)

	http.ListenAndServe("localhost:8000", nil)
}
