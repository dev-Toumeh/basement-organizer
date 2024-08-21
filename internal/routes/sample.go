package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/templates"
	"fmt"
	"net/http"
)

func SamplePage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	user, _ := auth.UserSessionData(r)
	data := templates.NewPageTemplate()
	data.Title = "Sample Page"
	data.Authenticated = authenticated
	data.User = user

	triggerAllServerNotifications(w)
	MustRender(w, r, templates.TEMPLATE_SAMPLE_PAGE, data)
}

func triggerAllServerNotifications(w http.ResponseWriter) {
	w.Header().Set("HX-Trigger-After-Swap", `{"ServerNotificationEvents":[{"message":"error", "type":"error" },{"message":"warning", "type":"warning" },{"message":"success", "type":"success" }, {"message":"default", "type":"" }]}`)
}

func TriggerErrorNotification(w http.ResponseWriter, message string) {
	msg := fmt.Sprintf(`{"ServerNotificationEvents":[{"message":"%s", "type":"error" }]}`, message)
	w.Header().Set("HX-Trigger-After-Swap", msg)
}

func TriggerSuccessNotification(w http.ResponseWriter, message string) {
	msg := fmt.Sprintf(`{"ServerNotificationEvents":[{"message":"%s", "type":"success" }]}`, message)
	w.Header().Set("HX-Trigger-After-Swap", msg)
}

func TriggerWarningNotification(w http.ResponseWriter, message string) {
	msg := fmt.Sprintf(`{"ServerNotificationEvents":[{"message":"%s", "type":"warning" }]}`, message)
	w.Header().Set("HX-Trigger-After-Swap", msg)
}
