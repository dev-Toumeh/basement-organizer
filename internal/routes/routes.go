package routes

import (
	"fmt"
	"io"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/database"
	"basement/main/internal/items"
	"basement/main/internal/server"
	"basement/main/internal/templates"
)

func RegisterRoutes(db *database.DB) {
	staticRoutes()
	authRoutes(db)
	itemsRoutes(db)
	shelvesRoutes()
	registerBoxRoutes(db)
	navigationRoutes()
	experimentalRoutes()
}

func authRoutes(db auth.AuthDatabase) {
	http.HandleFunc("/login", auth.LoginHandler(db))
	http.HandleFunc("/login-form", auth.LoginForm)
	http.HandleFunc("/register", auth.RegisterHandler(db))
	http.HandleFunc("/register-form", func(w http.ResponseWriter, r *http.Request) {
		server.MustRender(w, r, templates.TEMPLATE_REGISTER_FORM, nil)
	})
	http.HandleFunc("/update", auth.UpdateHandler(db))
	http.HandleFunc("/logout", auth.LogoutHandler)
}

func itemsRoutes(db items.ItemDatabase) {
	http.HandleFunc("/api/v1/implement-me", server.ImplementMeHandler)
	http.HandleFunc("/items", itemsPage)
	http.HandleFunc("/template/item-form", itemTemp)
	http.HandleFunc("/template/item-search", searchItemTemp)
	http.HandleFunc("/template/item-dummy", func(w http.ResponseWriter, r *http.Request) {
		db.InsertDummyItems()
		templates.RenderSuccessNotification(w, "dummy items has been added")
	})
	http.HandleFunc("/items-pagination", items.ItemPaginationHandler(db))

	http.HandleFunc("/delete-item", items.DeleteItemHandler(db))
	http.HandleFunc("/move-item", moveItem)
	http.HandleFunc("/item", items.ReadItemHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEM_CONTAINER, data)
	}))
	http.HandleFunc("/items-ids", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEMS_CONTAINER, data)
	}))

	http.HandleFunc("/api/v1/create/item", items.CreateItemHandler(db))
	http.HandleFunc("/api/v1/read/item/{id}", items.ReadItemHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEM_CONTAINER, data)
	}))
	http.HandleFunc("/api/v1/search/item", items.SearchItemHandler(db))
	http.HandleFunc("/api/v1/update/item", items.UpdateItemHandler(db))
	http.HandleFunc("/api/v1/move/item", items.MoveItemHandler(db))
	http.HandleFunc("/api/v1/delete/item", items.DeleteItemHandler(db))
	http.HandleFunc("/api/v1/read/items", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		fmt.Fprint(w, data)
	}))
}

func shelvesRoutes() {
	// @TODO: Implement shelves page.
	http.HandleFunc("/shelves", server.ImplementMeHandler)
}

var testStyle = templates.DEBUG_STYLE

func SwitchDebugStyle(w http.ResponseWriter, r *http.Request) {
	if testStyle {
		templates.InitTemplates("")
		templates.RedefineFromOtherTemplateDefinition("style", templates.InternalTemplate(), "style-debug", templates.InternalTemplate())
		templates.Render(w, templates.TEMPLATE_STYLE, nil)
	} else {
		templates.InitTemplates("")
		templates.RedefineTemplateDefinition(templates.InternalTemplate(), "style", "<style></style>")
		templates.Render(w, templates.TEMPLATE_STYLE, nil)
	}
	testStyle = !testStyle
}

func staticRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/static"))))
	http.HandleFunc("/", AuthPage)
}

func experimentalRoutes() {
	http.HandleFunc("/switch-debug-style", SwitchDebugStyle)
	http.HandleFunc("/notification-success", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderSuccessNotification(w, "success")
	})
	http.HandleFunc("/notification-warning", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderWarningNotification(w, "warning")
	})
}

func navigationRoutes() {

	http.HandleFunc("/settings", SettingsPage)
	http.HandleFunc("/sample-page", SamplePage)
	http.HandleFunc("/personal-page", PersonalPage)
}
