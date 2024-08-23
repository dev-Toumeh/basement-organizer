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

func RegisterRoutes(db *database.DB) {
	staticRoutes()
	authRoutes(db)
	apiRoutes(db)
	experimentalRoutes()
	registerBoxRoutes(&mockBoxDB{})
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
	http.HandleFunc("/update", auth.UpdateHandler(db))
	http.HandleFunc("/logout", auth.LogoutHandler)
}

func apiRoutes(db items.ItemDatabase) {
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

func experimentalRoutes() {
	http.HandleFunc("/sample-page", SamplePage)
	http.HandleFunc("/switch-debug-style", SwitchDebugStyle)
	http.HandleFunc("/snackbar-success", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderSuccessSnackbar(w, "success")
	})
	http.HandleFunc("/snackbar-warning", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderWarningSnackbar(w, "warning")
	})
}
