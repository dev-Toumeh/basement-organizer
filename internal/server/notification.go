package server

import (
	"basement/main/internal/logg"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type event struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Duration int    `json:"duration"`
}

func (e event) toJSONString() string {
	data, err := json.Marshal(e)
	if err != nil {
		logg.WrapErr(err)
		return "error"
	}
	return string(data)
}

type Notifications struct {
	ServerNotificationEvents []event
}

func (e *Notifications) AddInfo(msg string) {
	e.Add(msg, "", 2000)
}

func (e *Notifications) AddSuccess(msg string) {
	e.Add(msg, "success", 4000)
}

func (e *Notifications) AddWarning(msg string) {
	e.Add(msg, "warning", 5000)
}

func (e *Notifications) AddError(msg string) {
	e.Add(msg, "error", 10000)
}

// Add adds a server notification that is triggered on the client.
//
// `notificationType` can be "" (info), "success", "warning", "error".
//
// `duration` in milliseconds: 2000 = 2s
func (e *Notifications) Add(msg string, notificationType string, duration int) {
	if duration < 500 {
		duration = 500
		logg.Infof("duration for server notification is below 500ms (%d), setting to 500", duration)
	}
	if notificationType != "" && notificationType != "success" && notificationType != "warning" && notificationType != "error" {
		logg.Infof(`notificationType "%s" is not valid. Setting default type=""`, notificationType)
		notificationType = ""
	}
	e.ServerNotificationEvents = append(e.ServerNotificationEvents, event{Type: notificationType, Message: msg, Duration: duration})
}

func (e Notifications) Messages() []string {
	messages := make([]string, len(e.ServerNotificationEvents))
	for i, m := range e.ServerNotificationEvents {
		messages[i] = m.Message
	}
	return messages
}

func (e Notifications) toJSONString() string {
	data, err := json.Marshal(e)
	if err != nil {
		logg.WrapErr(err)
		return "error"
	}
	return fmt.Sprint(string(data))
}

// TriggerNotifications will trigger all notifications using "HX-Trigger-After-Swap" to trigger client side event to create notifications.
func TriggerNotifications(w http.ResponseWriter, e Notifications) {
	msgReturn := e.toJSONString()
	logg.Debugf("TriggerNotificationEvents: %s", msgReturn)
	w.Header().Set("HX-Trigger-After-Swap", msgReturn)
}

// Triggers all notifications for testing purposes.
func TriggerAllServerNotifications(w http.ResponseWriter) {
	n := Notifications{}
	n.AddInfo("info")
	n.AddSuccess("success")
	n.AddWarning("warning")
	n.AddError("error")
	w.Header().Set("HX-Trigger-After-Swap", n.toJSONString())
}

// TriggerErrorNotification uses "HX-Trigger-After-Swap" to trigger client side event to create notification.
// Must be called before any writing to w happens.
func TriggerErrorNotification(w http.ResponseWriter, message string) {
	n := Notifications{}
	n.AddError(message)
	w.Header().Set("HX-Trigger-After-Swap", n.toJSONString())
}

// TriggerSuccessNotification uses "HX-Trigger-After-Swap" to trigger client side event to create notification.
// Must be called before any writing to w happens.
func TriggerSuccessNotification(w http.ResponseWriter, message string) {
	n := Notifications{}
	n.AddSuccess(message)
	w.Header().Set("HX-Trigger-After-Swap", n.toJSONString())
}

// TriggerWarningNotification uses "HX-Trigger-After-Swap" to trigger client side event to create notification.
// Must be called before any writing to w happens.
func TriggerWarningNotification(w http.ResponseWriter, message string) {
	n := Notifications{}
	n.AddWarning(message)
	w.Header().Set("HX-Trigger-After-Swap", n.toJSONString())
}

func TriggerWarningNotifications(w http.ResponseWriter, messages []string) {
	eves := make([]event, len(messages))
	serverEvents := Notifications{ServerNotificationEvents: eves}
	for i := range eves {
		serverEvents.ServerNotificationEvents[i].Message = messages[i]
		serverEvents.ServerNotificationEvents[i].Type = "warning"
	}
	msgReturn := serverEvents.toJSONString()
	logg.Infof("warning %s", msgReturn)
	w.Header().Set("HX-Trigger-After-Swap", msgReturn)
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
	n := Notifications{}
	n.Add(message, notificationType, duration)
	notification := fmt.Sprintf(`{ "path":"%s", "headers":{ "ServerNotificationEvents":[%s] } }`, path, n.ServerNotificationEvents[0].toJSONString())
	logg.Debug(notification)
	w.Header().Set("HX-Location", notification)
}

// RedirectWithNotifications redirects client to another page and shows a notification after.
//
// It uses "HX-Redirect" to change page and triggers javascript notification after swap.
func RedirectWithNotifications(w http.ResponseWriter, path string, notifications Notifications) {
	n := strings.Trim(notifications.toJSONString(), "{}")
	notification := fmt.Sprintf(`{ "path":"%s", "headers":{ %s } }`, path, n)
	logg.Debug(notification)
	w.Header().Set("HX-Location", notification)
}
