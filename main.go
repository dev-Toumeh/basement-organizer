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
	_, err := env.LoadConfig()
	if err != nil {
		logg.Err(err)
		logg.Fatal("can't load config " + logg.CleanLastError(err))
	}

	db := &database.DB{}

	db.Connect()
	defer db.Sql.Close()

	routes.RegisterRoutes(db)
	err = templates.InitTemplates(env.CurrentConfig().TemplatePath())
	if err != nil {
		logg.Fatal("Templates failed to initialize", err)
	}

	err = http.ListenAndServe("0.0.0.0:8101", nil)
	if err != nil {
		panic(err)
	}
}
