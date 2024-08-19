package routes

import (
	"fmt"
	"io"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/database"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
)

// Temporary workaround for simplicity.
// No need to pass reference around in functions for database access.
var pdb *database.DB

func RegisterRoutes(db *database.DB) {
	pdb = db
	staticRoutes()
	authRoutes(db)
	apiRoutes(db)
	experimentalRoutes(db)
}

// MustRender will only render valid templates or throw http.StatusInternalServerError.
func MustRender(w http.ResponseWriter, r *http.Request, name string, data any) {
	err := templates.SafeRender(w, name, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logg.Debug(http.StatusText(http.StatusInternalServerError))
		return
	}
}

func authRoutes(db auth.AuthDatabase) {
	http.HandleFunc("/login", auth.LoginHandler(db))
	http.HandleFunc("/login-form", auth.LoginForm)
	http.HandleFunc("/register", auth.RegisterHandler(db))
	http.HandleFunc("/register-form", func(w http.ResponseWriter, r *http.Request) {
		MustRender(w, r, templates.TEMPLATE_REGISTER_FORM, nil)
	})
	http.HandleFunc("/logout", auth.LogoutHandler)
}

func apiRoutes(db *database.DB) {
	http.HandleFunc("/item", items.ReadItemHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEM_CONTAINER, data)
	}))
	http.HandleFunc("/items", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEMS_CONTAINER, data)
	}))

	http.HandleFunc("/api/v1/create/item", items.CreateItemHandler(db))
	http.HandleFunc("/api/v1/read/item/{id}", items.ReadItemHandler(db, func(w io.Writer, data any) {
		fmt.Fprint(w, data)
	}))
	http.HandleFunc("/api/v1/update/item", items.UpdateItemHandler(db))
	http.HandleFunc("/api/v1/delete/item", items.DeleteItemHandler(db))
	http.HandleFunc("/api/v1/read/items", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		fmt.Fprint(w, data)
	}))
}

var testStyle = templates.DEBUG_STYLE

func SwitchDebugStyle(w http.ResponseWriter, r *http.Request) {
	if testStyle {
		templates.InitTemplates()
		templates.RedefineFromOtherTemplateDefinition("style", templates.InternalTemplate(), "style-debug", templates.InternalTemplate())
		templates.Render(w, templates.TEMPLATE_STYLE, nil)
	} else {
		templates.InitTemplates()
		templates.RedefineTemplateDefinition(templates.InternalTemplate(), "style", "<style></style>")
		templates.Render(w, templates.TEMPLATE_STYLE, nil)
	}
	testStyle = !testStyle
}

func staticRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/static"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/personal-page", PersonalPage)
}

func experimentalRoutes(db *database.DB) {
	http.HandleFunc("/sample-page", SamplePage)
	http.HandleFunc("/switch-debug-style", SwitchDebugStyle)
	http.HandleFunc("/snackbar-success", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderSuccessSnackbar(w, "success")
	})
	http.HandleFunc("/snackbar-warning", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderWarningSnackbar(w, "warning")
	})
	http.HandleFunc("/box", BoxRequestHandler)
}

func BoxRequestHandler(w http.ResponseWriter, r *http.Request) {
	b := items.Box{Label: "asdfasdf"}
	b2 := items.Box{Label: "box 2"}
	err := b2.MoveTo(&b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		logg.Err(err)
	}
	ids, _ := pdb.ItemIDs()
	item, _ := pdb.Item(ids[0])
	item.Picture = ""

	b.Items = []*database.Item{&item}
	// data, _ := json.Marshal(b)
	data, _ := b.MarshalJSON()
	fmt.Fprintf(w, "%s", data)
}
