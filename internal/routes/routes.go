package routes

import (
	"fmt"
	"io"
	"net/http"

	"basement/main/internal/areas"
	"basement/main/internal/auth"
	"basement/main/internal/boxes"
	"basement/main/internal/common"
	"basement/main/internal/database"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/shelves"
	"basement/main/internal/templates"

	"github.com/gofrs/uuid/v5"
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
	common.RegisterDBInstance(db)
	staticRoutes()
	navigationRoutes()
	authRoutes(db)
	itemsRoutes(db)
	itemsRoutes2(db)
	boxesRoutes(db)
	shelvesRoutes(db)
	areaRoutes(db)
	experimentalRoutes(db)

	Handle("/addto/{thing}", AddTo(db))
	Handle("/element/{thing}/{thingid}", Element(db))
}

func staticRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/opt/basement-organizer/internal/static"))))
	Handle("/", common.Handle404NotFoundPage)
	Handle("/auth", AuthPage)
}

func navigationRoutes() {
	Handle("/settings", SettingsPage)
	Handle("/sample-page", SamplePage)
	Handle("/personal-page", PersonalPage)
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
	Handle("/template/item-dummy", func(w http.ResponseWriter, r *http.Request) {
		db.InsertSampleItems()
		templates.RenderSuccessNotification(w, "dummy items has been added")
	})

	Handle("/delete-item", items.DeleteItemHandler(db))
	Handle("/items-ids", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEMS_CONTAINER, data)
	}))

	// Handle("/api/v1/create/item", items.CreateItemHandler(db))
	Handle("/api/v1/read/item/{id}", items.ReadItemHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.TEMPLATE_ITEM_CONTAINER, data)
	}))
	Handle("/api/v1/move/item", items.MoveItemHandler(db))
	Handle("/api/v1/delete/item", items.DeleteItemHandler(db))
	Handle("/api/v1/read/items", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		fmt.Fprint(w, data)
	}))
}

func itemsRoutes2(db items.ItemDatabase) {
	Handle("/items", items.ItemsHandler(db))
	Handle("/item/{id}", items.DetailsTemplate(db))
	Handle("/items/create", items.CreateTemplate())

	// Move multiple items from list.
	Handle("/items/moveto/{thing}", common.ListPageMovePicker(common.THING_ITEM, db))
	Handle("/items/moveto/box/{id}", func(w http.ResponseWriter, r *http.Request) {
		common.ListPageMovePickerConfirm(db.MoveItemToBox, "/items").ServeHTTP(w, r)
	})
	Handle("/items/moveto/shelf/{id}", func(w http.ResponseWriter, r *http.Request) {
		common.ListPageMovePickerConfirm(db.MoveItemToShelf, "/items").ServeHTTP(w, r)
	})
	Handle("/items/moveto/area/{id}", func(w http.ResponseWriter, r *http.Request) {
		common.ListPageMovePickerConfirm(db.MoveItemToArea, "/items").ServeHTTP(w, r)
	})

	// API's
	http.Handle("/api/v1/create/item", items.ItemHandler(db))
	http.Handle("/api/v1/update/item", items.ItemHandler(db))
	http.Handle("/api/v1/delete/item/{id}", items.ItemHandler(db))
}

func boxesRoutes(db *database.DB) {
	boxes.RegisterDBInstance(db)
	// Box templates
	Handle("/box/create", boxes.CreateHandler(db))
	Handle("/box/{id}", boxes.BoxHandler(db))
	Handle("/box/{id}/boxDetailsForm", boxes.RenderBoxDetailsForm(db))
	Handle("/box/{id}/innerBoxes", common.HandleListTemplateInnerThingsData(common.THING_BOX, common.THING_BOX))
	Handle("/box/{id}/innerItems", common.HandleListTemplateInnerThingsData(common.THING_ITEM, common.THING_BOX))

	Handle("/box/{id}/addto/{thing}", boxes.BoxPicker(boxes.PICKER_TYPE_ADDTO, db))
	Handle("/box/{id}/addto/{thing}/{thingid}", boxes.BoxPickerConfirm(boxes.PICKER_TYPE_ADDTO, db))
	Handle("/box/{id}/moveto/{thing}", boxes.BoxPicker(boxes.PICKER_TYPE_MOVE, db))
	Handle("/box/{id}/moveto/{thing}/{thingid}", boxes.BoxPickerConfirm(boxes.PICKER_TYPE_MOVE, db))

	// Box api
	Handle("/api/v1/box", boxes.BoxHandler(db))
	Handle("/api/v1/box/{id}", boxes.BoxHandler(db))
	// Moves a box to another box.
	Handle("/api/v1/box/{id}/move/{toid}", func(w http.ResponseWriter, r *http.Request) {
		id := uuid.FromStringOrNil(r.PathValue("id"))
		moveToBoxID := uuid.FromStringOrNil(r.PathValue("toid"))
		err := db.MoveBoxToBox(id, moveToBoxID)
		if err != nil {
			server.WriteBadRequestError("can't move box", err, w, r)
			logg.Err(err)
		} else {
			w.WriteHeader(200)
		}
	})

	// Boxes templates
	Handle("/boxes", boxes.BoxesHandler(db))
	Handle("/boxes/moveto/{thing}", common.ListPageMovePicker(common.THING_BOX, db))
	Handle("/boxes/moveto/box/{id}", func(w http.ResponseWriter, r *http.Request) {
		common.ListPageMovePickerConfirm(db.MoveBoxToBox, "/boxes").ServeHTTP(w, r)
	})
	Handle("/boxes/moveto/shelf/{id}", func(w http.ResponseWriter, r *http.Request) {
		common.ListPageMovePickerConfirm(db.MoveBoxToShelf, "/boxes").ServeHTTP(w, r)
	})
	Handle("/boxes/moveto/area/{id}", func(w http.ResponseWriter, r *http.Request) {
		common.ListPageMovePickerConfirm(db.MoveBoxToArea, "/boxes").ServeHTTP(w, r)
	})

	// Boxes api
	Handle("/api/v1/boxes", boxes.BoxesHandler(db))
	Handle("/api/v1/boxes/moveto/{thing}/{id}", func(w http.ResponseWriter, r *http.Request) {
		var notifications server.Notifications

		switch r.PathValue("thing") {
		case "box":
			notifications = server.MoveThingToThing(w, r, db.MoveBoxToBox)
			break
		case "shelf":
			notifications = server.MoveThingToThing(w, r, db.MoveBoxToShelf)
			break
		case "area":
			notifications = server.MoveThingToThing(w, r, db.MoveBoxToArea)
			break
		}
		server.TriggerNotifications(w, notifications)
	})
}

func shelvesRoutes(db shelves.ShelfDB) {
	//Template
	Handle("/shelves", shelves.ShelvesHandler(db))
	Handle("/shelf/create", shelves.CreateTemplate())
	Handle("/shelf/{id}", shelves.DetailsTemplate(db))

	Handle("/shelf/{id}/innerItems", common.HandleListTemplateInnerThingsData(common.THING_ITEM, common.THING_SHELF))
	Handle("/shelf/{id}/innerBoxes", common.HandleListTemplateInnerThingsData(common.THING_BOX, common.THING_SHELF))

	// Move multiple items from list.
	Handle("/shelves/moveto/{thing}", common.ListPageMovePicker(common.THING_SHELF, db))
	Handle("/shelves/moveto/area/{id}", func(w http.ResponseWriter, r *http.Request) {
		common.ListPageMovePickerConfirm(db.MoveShelfToArea, "/shelves").ServeHTTP(w, r)
	})

	// Api
	http.HandleFunc("/api/v1/create/shelf", shelves.ShelfHandler(db))
	Handle("/api/v1/delete/shelf", shelves.ShelfHandler(db))
	Handle("/api/v1/update/shelf", shelves.ShelfHandler(db))
	Handle("/api/v1/delete/shelves", shelves.DeleteShelves(db))
}

func areaRoutes(db *database.DB) {
	// Single area
	Handle("/area/{id}", areas.AreaHandler(db))
	Handle("/area", areas.CreateHandler(db))
	Handle("/area/create", areas.CreateHandler(db))
	Handle("/area/{id}/innerItems", common.HandleListTemplateInnerThingsData(common.THING_ITEM, common.THING_AREA))
	Handle("/area/{id}/innerBoxes", common.HandleListTemplateInnerThingsData(common.THING_BOX, common.THING_AREA))
	Handle("/area/{id}/innerShelves", common.HandleListTemplateInnerThingsData(common.THING_SHELF, common.THING_AREA))

	// Multiple areas
	Handle("/areas", areas.AreasHandler(db))

	// API
	Handle("/api/v1/area/{id}", areas.AreaHandler(db))
	Handle("/api/v1/area/create", areas.CreateHandler(db))
	Handle("/api/v1/areas", areas.AreasHandler(db))
}

func experimentalRoutes(db *database.DB) {
	Handle("/switch-debug-style", SwitchDebugStyle)
	Handle("/notification-success", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderSuccessNotification(w, "success")
	})
	Handle("/notification-warning", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderWarningNotification(w, "warning")
	})
	Handle("/notification-error", func(w http.ResponseWriter, r *http.Request) {
		templates.RenderErrorNotification(w, "error")
	})
	Handle("/templates/list", handleSampleListTemplate(db))
	Handle("/samples/return-selected-row-as-input/{id}", handleReturnSelectedInput(db))
	Handle("/samples/notification/{id}", handleReturnSelectedInputAsNotification(db))
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
