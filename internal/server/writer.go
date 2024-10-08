package server

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
)

func WriteFprint(w io.Writer, data any) {
	fmt.Fprint(w, data)
}

func WriteJSON(w io.Writer, data any) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}

// WriteNotFoundError sets not found status code 404, logs error and writes error message to client.
func WriteNotFoundError(message string, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, message)
	logg.Infof("%s: %s", message, err.Error())
}

// WriteNotFoundError sets not found status code 404, logs error and writes error message to client.
func WriteNotImplementedWarning(message string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, message+" not implemented")
	logg.Info(message + " not implemented")
}

// writeNotFoundError sets not found status code 404, logs error and writes error message to client.
func WriteInternalServerError(message string, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, message)
	logg.Alog(logg.ErrorLogger(), 3, "%s: %s", message, err)
}

func ImplementMeHandler(w http.ResponseWriter, r *http.Request) {
	returnMessage := "This is a placeholder API"
	logs := make([]string, 0)
	logs = append(logs, "Headers: ")
	for i, v := range r.Header {
		if !strings.HasPrefix(i, "Hx-") {
			continue
		}
		logs = append(logs, fmt.Sprintf("%s: %v", i, v))
	}
	trigger := r.Header.Get("HX-Trigger")
	if trigger != "" {
		returnMessage = fmt.Sprintf("%s triggered request: %s", trigger, returnMessage)
	}

	logg.Info(fmt.Sprintf("\"%s: %s\". %s called /api/v1/implement-me", r.Header.Get("hx-trigger"), returnMessage, r.Referer()))
	logg.Debug(strings.Join(logs, "\n\t"))
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, returnMessage)
}

// MustRender will only render valid templates or throw http.StatusInternalServerError.
func MustRender(w http.ResponseWriter, r *http.Request, name string, data any) {
	err := templates.SafeRender(w, name, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logg.Alog(logg.ErrorLogger(), 3, "%s: %s", http.StatusText(http.StatusInternalServerError), logg.WrapErr(err))
		return
	}
}

// Render applies data to a defined template and writes result back to the writer.
func RenderWithSuccessNotification(w http.ResponseWriter, r *http.Request, name string, data any, successMessage string) error {
	err := templates.CanRender(name, data)
	if err != nil {
		TriggerSuccessNotification(w, successMessage)
		return logg.Errorf("Template rendering failed %w", err)
	}
	TriggerSuccessNotification(w, successMessage)
	templates.Render(w, name, data)
	return nil
}

// ValidID returns valid uuid from request and handles errors.
// Check for uuid.Nil! If error occurs return will be uuid.Nil.
func ValidID(w http.ResponseWriter, r *http.Request, errorMessage string) uuid.UUID {
	id := r.FormValue("id")
	logg.Debugf("Query param id: '%v'.", id)
	if id == "" {
		id = r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			logg.Debug("Empty id")
			fmt.Fprintf(w, `%s ID="%v"`, errorMessage, id)
			return uuid.Nil
		}
		logg.Debugf("path value id: '%v'.", id)
	}

	newId, err := uuid.FromString(id)
	if err != nil {
		logg.Debugf("Wrong id: '%v'. %v", id, err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `%s ID="%v"`, errorMessage, id)
		return uuid.Nil
	}
	return newId
}
