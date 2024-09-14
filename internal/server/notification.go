package server

import (
	"fmt"
	"net/http"
)

func TriggerAllServerNotifications(w http.ResponseWriter) {
	w.Header().Set("HX-Trigger-After-Swap", `{"ServerNotificationEvents":[{"message":"error", "type":"error" },{"message":"warning", "type":"warning" },{"message":"success", "type":"success" }, {"message":"default", "type":"" }]}`)
}

// TriggerErrorNotification uses "HX-Trigger-After-Swap" to trigger client side event to create notification.
func TriggerErrorNotification(w http.ResponseWriter, message string) {
	msg := fmt.Sprintf(`{"ServerNotificationEvents":[{"message":"%s", "type":"error" }]}`, message)
	w.Header().Set("HX-Trigger-After-Swap", msg)
}

// TriggerSuccessNotification uses "HX-Trigger-After-Swap" to trigger client side event to create notification.
func TriggerSuccessNotification(w http.ResponseWriter, message string) {
	msg := fmt.Sprintf(`{"ServerNotificationEvents":[{"message":"%s", "type":"success" }]}`, message)
	w.Header().Set("HX-Trigger-After-Swap", msg)
}

// TriggerWarningNotification uses "HX-Trigger-After-Swap" to trigger client side event to create notification.
func TriggerWarningNotification(w http.ResponseWriter, message string) {
	msg := fmt.Sprintf(`{"ServerNotificationEvents":[{"message":"%s", "type":"warning" }]}`, message)
	w.Header().Set("HX-Trigger-After-Swap", msg)
}

// RedirectWithSuccessNotification redirects client to another page and shows a notification after for 2 seconds.
//
// It uses "HX-Redirect" to change page and triggers javascript success notification after swap.
func RedirectWithSuccessNotification(w http.ResponseWriter, path string, message string) {
	redirectWithNotification(w, path, message, "success", 2000)
}

// RedirectWithWarningNotification redirects client to another page and shows a notification after for 2 seconds.
//
// It uses "HX-Redirect" to change page and triggers javascript warning notification after swap.
func RedirectWithWarningNotification(w http.ResponseWriter, path string, message string) {
	redirectWithNotification(w, path, message, "warning", 2000)
}

// RedirectWithErrorNotification redirects client to another page and shows a notification after for 2 seconds.
//
// It uses "HX-Redirect" to change page and triggers javascript error notification after swap.
func RedirectWithErrorNotification(w http.ResponseWriter, path string, message string) {
	redirectWithNotification(w, path, message, "error", 2000)
}

// RedirectWithInfoNotification redirects client to another page and shows a notification after for 2 seconds.
//
// It uses "HX-Redirect" to change page and triggers javascript info notification after swap.
func RedirectWithInfoNotification(w http.ResponseWriter, path string, message string) {
	redirectWithNotification(w, path, message, "info", 2000)
}

// RedirectWithNotification redirects client to another page and shows a notification after.
//
// It uses "HX-Redirect" to change page and triggers javascript notification after swap.
//
// Possible values for notificationType: "success", "warning", "error" and "" for default.
//
// duration in milliseconds.
func redirectWithNotification(w http.ResponseWriter, path string, message string, notificationType string, duration int) {
	notification := fmt.Sprintf(`{ "path":"%s", "headers":{ "notification":{ "message":"%s", "type":"%s", "duration":"%v" } } }`, path, message, notificationType, duration)
	w.Header().Set("HX-Location", notification)
}
