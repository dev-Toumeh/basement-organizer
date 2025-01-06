package server

import (
	"basement/main/internal/logg"
	"basement/main/internal/templates"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
func WriteBadRequestError(message string, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, message)
	if err == nil {
		logg.Alog(logg.ErrorLogger(), 3, `%s "%s": %s`, r.Method, r.URL.String(), message)
	} else {
		logg.Alog(logg.ErrorLogger(), 3, `%s "%s": %s: %s`, r.Method, r.URL.String(), message, err.Error())
	}
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
	logg.Debugf("Form param id: '%v'.", id)

	if id == "" {
		id = r.PathValue("id")
		logg.Debugf("Path value id: '%v'.", id)

		if id == "" {
			id = r.URL.Query().Get("id")
			if id == "" {
				w.WriteHeader(http.StatusNotFound)
				logg.Debug("Empty id")
				fmt.Fprintf(w, `%s ID="%v"`, errorMessage, id)
				return uuid.Nil
			}
			logg.Debugf("Query param id: '%v'.", id)
		}
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

// wantsTemplateData checks if current request requires template data.
// Helps deciding how to write the data.
func WantsTemplateData(r *http.Request) bool {
	return !strings.Contains(r.URL.Path, "/api/")
}

// parseIDsFromFormWithKey parses r.Form by searching all HTML input elements that start with `key` name and returns a list of valid uuid.UUIDs
//
// `r.ParseForm()` must be called before using this function!
//
// Example:
//
//	// search for all ID values that start with "delete:" key
//	// like "delete:f47ac10b-58cc-0372-8567-0e02b2c3d479"
//	toDeleteIDs := parseIDsFromFormWithKey(r.Form, "delete")
func ParseIDsFromFormWithKey(form url.Values, key string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)
	// logg.Debug("FormValues length: " + strconv.Itoa(len(form)))
	// i := 0
	for k := range form {
		// logg.Debug("FormValue[" + strconv.Itoa(i) + "]: \"" + k + "\"")
		if strings.Contains(k, key+":") {
			// logg.Debug("\"" + k + "\"" + " contains " + "\"" + key + "\"")
			idStr := strings.Split(k, fmt.Sprintf("%s:", key))
			if len(idStr) != 2 {
				return nil, logg.NewError(fmt.Sprintf("Wrong delete key value pair: '%v'", k))
			}
			// logg.Debug("clean value \"" + idStr[1] + "\"")
			id, err := uuid.FromString(idStr[1])
			if err != nil {
				logg.Err(err)
				return nil, logg.WrapErr(err)
			}
			ids = append(ids, id)
		}
		// i++
	}
	// logg.Debugf("clean ids: %v", ids)
	return ids, nil
}

func DeleteThingsFromList(w http.ResponseWriter, r *http.Request, deleteFunc func(id uuid.UUID) error, listPageHandler http.Handler) {
	r.ParseForm()

	parseIDs, _ := ParseIDsFromFormWithKey(r.Form, "delete")
	r.FormValue("delete")

	logg.Debug(len(parseIDs))

	ids := make([]uuid.UUID, len(parseIDs))
	notifications := Notifications{}
	if len(parseIDs) != 0 {
		for i, v := range parseIDs {
			logg.Debug("deleting " + v.String())
			ids[i] = v
			err := deleteFunc(v)
			if err != nil {
				notifications.AddError(`can't delete: "` + v.String() + `"`)
				logg.Err(err)
			} else {
				notifications.AddSuccess(`delete: "` + v.String() + `"`)
			}
		}
	} else {
		notifications.AddWarning("Nothing selected to delete")
	}

	TriggerNotifications(w, notifications)
	listPageHandler.ServeHTTP(w, r)
}

func MoveThingToThing(w http.ResponseWriter, r *http.Request, moveFunc func(id1 uuid.UUID, id2 uuid.UUID) error) Notifications {
	r.ParseForm()
	moveToBoxID := ValidID(w, r, "can't move things invalid id")
	if moveToBoxID == uuid.Nil {
		return Notifications{}
	}

	parseIDs := r.PostForm["id-to-be-moved"]
	ids := make([]uuid.UUID, len(parseIDs))

	logg.Debug(len(parseIDs))

	notifications := Notifications{}
	for i, v := range parseIDs {
		logg.Debug(v)
		id := uuid.FromStringOrNil(v)
		ids[i] = id
		err := moveFunc(id, moveToBoxID)
		if err != nil {
			notifications.AddError("can't move \"" + ids[i].String() + "\" to \"" + moveToBoxID.String() + "\"")
			logg.Err(err)
		} else {
			notifications.AddSuccess("moved \"" + ids[i].String() + "\" to \"" + moveToBoxID.String() + "\"")
		}
	}
	return notifications
}
