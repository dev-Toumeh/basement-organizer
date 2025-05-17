package routes

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"basement/main/internal/env"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
)

func SettingsPageConfiguration(w http.ResponseWriter, r *http.Request) {
	data := settingsPageData(r).Map()
	data["Configuration"] = env.CurrentConfig()
	server.MustRender(w, r, templates.TEMPLATE_SETTINGS_PAGE, data)
}

func SettingsPageConfigurationDefaultValue(w http.ResponseWriter, r *http.Request) {
	opt := r.PathValue("option")
	if opt == "" {
		server.WriteBadRequestError("option empty", logg.NewError("option empty"), w, r)
		return
	}

	value, err := env.CurrentConfig().DefaultValue(opt)
	if err != nil {
		server.WriteBadRequestError(logg.CleanLastError(err), err, w, r)
		return
	}
	inp := fmt.Sprintf("<input value=\"%s\">", value)
	server.WriteFprint(w, inp)
	return
}

// Serves only the configuration form.
func SettingsConfigurationForm(w http.ResponseWriter, r *http.Request) {
	data := env.CurrentConfig()
	server.MustRender(w, r, templates.TEMPLATE_CONFIGURATION, data)
}

func SettingsConfigurationUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Add("Allow", http.MethodPut)
		w.WriteHeader(http.StatusMethodNotAllowed)
		logg.Err("Not allowed")
		return
	}

	errs := updateConfiguration(w, r)
	if len(errs) > 0 {
		e := server.Notifications{}
		for _, v := range errs {
			e.AddError(logg.CleanLastError(v))
		}
		server.WriteBadRequestError(strings.Join(e.Messages(), "\n\n"), nil, w, r)
		return
	}
	server.TriggerSuccessNotification(w, "Configuration updated and reloaded")
	return
}

func updateConfiguration(w http.ResponseWriter, r *http.Request) (errs []error) {
	parsed, err := parseConfigOptions(r)
	if err != nil {
		errs = append(errs, logg.WrapErr(err))
		return errs
	}

	c := env.CurrentConfig()
	errs = env.ApplyParsedConfigOptions(parsed, c)
	if len(errs) > 0 {
		logg.Warningf("applying parsed config options produced \"%d\" errors", len(errs))
		return errs
	}
	err = env.WriteCurrentConfigToFile()
	if err != nil {
		errs = append(errs, logg.WrapErr(err))
	}
	return errs
}

func parseConfigOptions(r *http.Request) (map[string]string, error) {
	fields := env.CurrentConfig().FieldValues()
	parsed := make(map[string]string, len(fields))
	for f, val := range fields {
		formValue := r.FormValue(f)
		newValue := ""
		switch formValue {
		case "":
			// checkbox not checked
			if val.Kind == reflect.Bool {
				newValue = "false"
			}
			break
		// checkbox checked
		case "on":
			newValue = "true"
			break
		default:
			newValue = formValue
		}

		logg.Debugf("parsed[%s] = %s", f, newValue)
		parsed[f] = newValue
	}

	return parsed, nil
}
