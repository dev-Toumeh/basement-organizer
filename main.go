package main

import (
	"basement/main/internal/database"
	"basement/main/internal/logg"
	"basement/main/internal/routes"
	"basement/main/internal/templates"
	"net/http"
)

func main() {
	// Execute before the first run:
	//
	//	go run ./tools/setup/setup.go
	//
	// This will create config.go where LoadConfig() is defined.
	env := LoadConfig()

	db := &database.DB{}

	db.Connect()
	defer db.Sql.Close()

	routes.RegisterRoutes(db)
	err := templates.InitTemplates(env.TemplatePath())
	if err != nil {
		logg.Fatal("Templates failed to initialize", err)
	}

	err = http.ListenAndServe("0.0.0.0:8101", nil)
	if err != nil {
		panic(err)
	}
}
