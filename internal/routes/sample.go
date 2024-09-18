package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"net/http"
)

func SamplePage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	user, _ := auth.UserSessionData(r)
	data := templates.NewPageTemplate()
	data.Title = "Sample Page"
	data.Authenticated = authenticated
	data.User = user

	server.TriggerAllServerNotifications(w)
	server.MustRender(w, r, templates.TEMPLATE_SAMPLE_PAGE, data.Map())
}
