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
	env.LoadConfig()

	db := &database.DB{}

	db.Connect()
	defer db.Sql.Close()

	routes.RegisterRoutes(db)
	err := templates.InitTemplates("./internal")
	if err != nil {
		logg.Fatal("Templates failed to initialize", err)
	}

	err = http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		panic(err)
	}
}
