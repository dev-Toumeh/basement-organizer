package main

import (
	"basement/main/internal/database"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/routes"
	"basement/main/internal/templates"
	"net/http"
)

func main() {
	// Move this only if necessary
	env.SetDevelopment()
	if env.Development() {
		logg.EnableDebugLogger()
		logg.EnableInfoLogger()
	}

	db := &database.DB{}

	db.Connect()
	defer db.Sql.Close()

	routes.RegisterRoutes(db)
	templates.InitTemplates()

	http.ListenAndServe("localhost:8000", nil)
}
