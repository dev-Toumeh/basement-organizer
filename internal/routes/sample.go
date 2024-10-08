package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"errors"
	"net/http"
)

func SamplePage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	user, _ := auth.UserSessionData(r)
	data := templates.NewPageTemplate()
	data.Title = "Sample Page"
	data.Authenticated = authenticated
	data.User = user

	logg.Debug("debug log")
	logg.Debugf("debugf %s", "log")
	logg.Info("info log")
	logg.Infof("infof %s", "log")
	err := errors.New("Long error chain start !!!")
	err2 := logg.Errorf("Error 2 %w", err)
	err3 := logg.Errorf("Error 3 %w", err2)
	err4 := logg.Errorf("Error 4 %w", err3)
	logg.Err("Error happened:", err4)
	logg.Errf("Err %s", "fff")

	server.TriggerAllServerNotifications(w)
	// server.MustRender(w, r, templates.TEMPLATE_SAMPLE_PAGE, "sdf")
	server.MustRender(w, r, templates.TEMPLATE_SAMPLE_PAGE, data.Map())
}
