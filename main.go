package main

import (
	"basement/main/internal/database"
	"basement/main/internal/routes"
	"basement/main/internal/templates"
	"log"
	"net/http"
)

func main() {

	db := &database.DB{}

	err := db.Connect()
	if err != nil {
		log.Fatalf("Can't create DB, shutting server down")
	}
	defer db.Sql.Close()
	//  db.PrintItemRecords()

	routes.RegisterRoutes(db)
	templates.InitTemplates()

	http.ListenAndServe("localhost:8000", nil)
}
