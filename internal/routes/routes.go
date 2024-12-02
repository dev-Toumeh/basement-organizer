package routes

import (
	"fmt"
	"io"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/boxes"
	"basement/main/internal/database"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/shelves"
	"basement/main/internal/templates"
)

func Handle(route string, handler http.HandlerFunc) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		msg := ""
		msg = fmt.Sprintf(`%s "%s" http://%s%s%s`, r.Method, route, r.URL.Scheme, r.Host, r.URL)
		colorMsg := fmt.Sprintf("%s%s%s", logg.Yellow, msg, logg.Reset)
		logg.Debug(colorMsg)
		if r.Method == http.MethodPost {
			// @TODO: Fix. Breaks some post requests because r.ParseForm is empty after this.
			// r.ParseForm()
			// colorMsg := fmt.Sprintf("%sPostFormValue: %v%s", logg.Yellow, r.PostForm, logg.Reset)
			// logg.Debug(colorMsg)
		}
		handler.ServeHTTP(w, r)
	})
}

func RegisterRoutes(db *database.DB) {
	staticRoutes()
	authRoutes(db)
	itemsRoutes(db)
	itemsRoutes2(db)
	registerBoxRoutes(db)
	shelvesRoutes(db)
	navigationRoutes()
	experimentalRoutes(db)
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

	// Handle("/api/v1/create/item", items.CreateItemHandler(db))
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

func itemsRoutes2(db items.ItemDatabase) {
	Handle("/items", items.PageTemplate(db))
	Handle("/items/create", items.CreateTemplate())

	// API's
	http.Handle("/api/v1/create/item", items.ItemHandler(db))
}

func registerBoxRoutes(db *database.DB) {
	// Box templates
	Handle("/box", boxes.BoxHandler(db))
	Handle("/box/{id}", boxes.DetailsPage(db))
	Handle("/box/{id}/moveto/box", boxes.BoxMovePicker("box", db))
	Handle("/box/{id}/moveto/box/{value}", boxes.BoxMovePickerConfirm("box", db))
	Handle("/box/{id}/moveto/shelf", boxes.BoxMovePicker("shelf", db))
	Handle("/box/{id}/moveto/shelf/{value}", boxes.BoxMovePickerConfirm("shelf", db))

	// Box api
	Handle("/api/v1/box", boxes.BoxHandler(db))
	Handle("/api/v1/box/{id}", boxes.BoxHandler(db))
	Handle("/api/v1/box/{id}/move/{toid}", boxes.MoveBoxAPI(db))

	// Boxes templates
	Handle("/boxes", boxes.BoxesHandler(db))
	Handle("/boxes/move", boxes.ListPageMoveToBoxPicker(db))
	Handle("/boxes/moveto/box/{id}", boxes.ListPageMoveToBoxPickerConfirm(db))

	// Boxes api
	Handle("/api/v1/boxes", boxes.BoxesHandler(db))
	Handle("/api/v1/boxes/moveto/box/{id}", boxes.MoveBoxesToBoxAPI(db))
}

func shelvesRoutes(db shelves.ShelfDB) {
	//Template
	Handle("/shelves", shelves.PageTemplate(db))
	Handle("/shelves/create", shelves.CreateTemplate())
	Handle("/shelves/{id}", shelves.DetailsTemplate(db))
	Handle("/shelves/add-list", shelves.AddListTemplate(db))
	Handle("/shelves/add-input/{id}", shelves.AddInputTemplate(db))

	// Api
	http.HandleFunc("/api/v1/create/shelf", shelves.ShelfHandler(db))
	Handle("/api/v1/delete/shelf", shelves.ShelfHandler(db))
	Handle("/api/v1/update/shelf", shelves.ShelfHandler(db))
	Handle("/api/v1/delete/shelves", shelves.DeleteShelves(db))
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

func experimentalRoutes(db *database.DB) {
	Handle("/switch-debug-style", SwitchDebugStyle)
	Handle("/notification-success", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderSuccessNotification(w, "success")
	})
	Handle("/notification-warning", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderWarningNotification(w, "warning")
	})
	Handle("/templates/list", handleSampleListTemplate(db))
	Handle("/samples/return-selected-row-as-input/{id}", handleReturnSelectedInput(db))
	Handle("/samples/notification/{id}", handleReturnSelectedInputAsNotification(db))
}

func navigationRoutes() {

	Handle("/settings", SettingsPage)
	Handle("/sample-page", SamplePage)
	Handle("/personal-page", PersonalPage)
}
