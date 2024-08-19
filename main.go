package main

import (
	"basement/main/internal/database"
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/routes"
	"basement/main/internal/templates"
	"context"
	"net"
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

	// Accessable context in every http.Request.
	// Access db with r.Context().Value("db").(*database.DB) inside an http.Handler.
	DBctx := context.WithValue(context.Background(), "db", db)
	server := http.Server{
		Addr:        "localhost:8000",
		BaseContext: func(_ net.Listener) context.Context { return DBctx },
	}

	server.ListenAndServe()
}
