package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/database"
	"basement/main/internal/env"
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

func handleSampleListTemplate(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// rows := []common.ListRow{{ID: database.BOX_VALID_UUID_1, Label: "THING 1"}}
		boxes, _ := db.BoxFuzzyFinder("", 3, 1)
		// boxes, _ := db.BoxFuzzyFinder(uuid.FromStringOrNil("17973d34-1942-4a15-bcba-80ddca1b29fc"))

		tmpl := common.ListTemplate{
			// FormHXGet: "/items",
			// RowHXGet:  "/api/v1/read/item",
			FormID:    "my-custom-list-template",
			RowHXGet:  "/api/v1/box",
			Rows:      boxes,
			RowAction: true,
			ShowLimit: env.Config().ShowTableSize(),
			// DataInputName:   "id-to-be-moved",
			AdditionalDataInputs: []common.DataInput{
				{Key: "return-hidden-input", Value: "false"},
				{Key: "id-to-be-moved", Value: "1f73d774-8bd5-4246-940f-ef9abd1c480e"},
			},
			// AdditionalDataInputValues: []string{"1f73d774-8bd5-4246-940f-ef9abd1c480e"},
			RowActionName: "Create notification with this id",
			// RowActionHXPost:   "/api/v1/boxes/moveto/box",
			// RowActionHXPost:   "/api/v1/boxes?query=b",
			// RowActionHXPost:   "/api/v1/implement-me",
			// RowActionHXPostWithID: "/samples/return-selected-row-as-input",
			RowActionHXPostWithID: "/samples/notification",
			// RowActionHXPostWithIDAsQueryParam: "/samples/return-selected-row-as-input",
			RowActionHXTarget: "#mytarget",
		}

		err := tmpl.Render(w)
		if err != nil {
			logg.Err(err)
			server.TriggerErrorNotification(w, err.Error())
			return
		}
	}
}

func handleReturnSelectedInput(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			id := r.PathValue("id")
			hidden := r.PostFormValue("hidden")
			server.MustRender(w, r, templates.TEMPLATE_SELECTED_INPUT, map[string]string{"Name": "select", "Value": id, "Hidden": hidden})
		} else {
			return
		}
	}
}

func handleReturnSelectedInputAsNotification(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			id := r.PathValue("id")
			hidden := r.PostFormValue("hidden")
			s := server.Notifications{}
			s.Add("id="+id+"\nhidden="+hidden, "", 1000)
			server.TriggerNotifications(w, s)
			server.MustRender(w, r, templates.TEMPLATE_SELECTED_INPUT, map[string]string{"Name": "select", "Value": id, "Hidden": hidden})
		} else {
			w.Header().Add("Allowed", http.MethodPost)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
