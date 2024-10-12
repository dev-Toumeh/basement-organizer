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
	Handle("/login", auth.LoginHandler(db))
	Handle("/login-form", auth.LoginForm)
	Handle("/register", auth.RegisterHandler(db))
	Handle("/register-form", func(w http.ResponseWriter, r *http.Request) {
		server.MustRender(w, r, templates.TEMPLATE_REGISTER_FORM, nil)
	})
	Handle("/update", auth.UpdateHandler(db))
	Handle("/logout", auth.LogoutHandler)
}

func itemsRoutes(db items.ItemDatabase) {
	Handle("/api/v1/implement-me", server.ImplementMeHandler)
	Handle("/items", itemsPage)
	Handle("/template/item-form", itemTemp)
	Handle("/template/item-search", searchItemTemp)
	Handle("/template/item-dummy", func(w http.ResponseWriter, r *http.Request) {
		db.InsertSampleItems()
		templates.RenderSuccessNotification(w, "dummy items has been added")
	})
	Handle("/items-pagination", items.ItemPaginationHandler(db))

	Handle("/delete-item", items.DeleteItemHandler(db))
	Handle("/move-item", moveItem)
	Handle("/item", items.ReadItemHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEM_CONTAINER, data)
	}))
	Handle("/items-ids", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEMS_CONTAINER, data)
	}))

	Handle("/api/v1/create/item", items.CreateItemHandler(db))
	Handle("/api/v1/read/item/{id}", items.ReadItemHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEM_CONTAINER, data)
	}))
	Handle("/api/v1/search/item", items.SearchItemHandler(db))
	Handle("/api/v1/update/item", items.UpdateItemHandler(db))
	Handle("/api/v1/move/item", items.MoveItemHandler(db))
	Handle("/api/v1/delete/item", items.DeleteItemHandler(db))
	Handle("/api/v1/read/items", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		fmt.Fprint(w, data)
	}))
}

func shelvesRoutes() {
	// @TODO: Implement shelves page.
	Handle("/shelves", server.ImplementMeHandler)
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
	Handle("/", AuthPage)
}

func experimentalRoutes() {
	Handle("/switch-debug-style", SwitchDebugStyle)
	Handle("/notification-success", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderSuccessNotification(w, "success")
	})
	Handle("/notification-warning", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderWarningNotification(w, "warning")
	})
}

func navigationRoutes() {

	Handle("/settings", SettingsPage)
	Handle("/sample-page", SamplePage)
	Handle("/personal-page", PersonalPage)
}
