package routes

import (
	"fmt"
	"io"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/database"
	"basement/main/internal/items"
	"basement/main/internal/templates"
)

const (
	STATIC                      string = "/static/"
	ITEMS_FILE_PATH             string = "internal/auth/items.json"
	USERS_FILE_PATH             string = "internal/auth/users2.json"
	API_V1_READ_ITEM            string = "/api/v1/read/item/{id}"
	PERSONAL_PAGE_ROUTE         string = "/personal-page"
	PERSONAL_PAGE_TEMPLATE_PATH string = "internal/templates/personal-page.html"
)

func RegisterRoutes(db *database.JsonDB) {
	http.Handle(STATIC, http.StripPrefix("/static/", http.FileServer(http.Dir("internal/static"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc(PERSONAL_PAGE_ROUTE, PersonalPage)
	http.HandleFunc("/sample-page", SamplePage)
	http.HandleFunc("/item", items.ReadItemHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.ITEM_CONTAINER, data)
	}))
	http.HandleFunc("/items", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		templates.Render(w, templates.ITEMS_CONTAINER, data)
	}))
	http.HandleFunc("/switch-debug-style", SwitchDebugStyle)
	http.HandleFunc("/login-form", auth.LoginForm)

	authRoutes(db)
	apiRoutes(db)
}

func authRoutes(db *database.JsonDB) {
	http.HandleFunc("/login", auth.LoginHandler(db))
	http.HandleFunc("/register", auth.RegisterHandler(db))
	http.HandleFunc("/logout", auth.LogoutHandler)
}

func apiRoutes(db *database.JsonDB) {
	http.HandleFunc("/api/v1/create/item", items.CreateItemHandler(db))
	http.HandleFunc(API_V1_READ_ITEM, items.ReadItemHandler(db, func(w io.Writer, data any) {
		fmt.Fprint(w, data)
	}))
	http.HandleFunc("/api/v1/update/item/id", UpdateItem)
	http.HandleFunc("/api/v1/delete/item", DeleteItem)
	http.HandleFunc("/api/v1/read/items", items.ReadItemsHandler(db, func(w io.Writer, data any) {
		fmt.Fprint(w, data)
	}))
}

var testStyle = templates.DEBUG_STYLE

func SwitchDebugStyle(w http.ResponseWriter, r *http.Request) {
	if testStyle {
		templates.InitTemplates()
		templates.RedefineFromOtherTemplateDefinition("style", templates.InternalTemplate(), "style-debug", templates.InternalTemplate())
		templates.Render(w, templates.STYLE, nil)
	} else {
		templates.InitTemplates()
		templates.RedefineTemplateDefinition(templates.InternalTemplate(), "style", "<style></style>")
		templates.Render(w, templates.STYLE, nil)
	}
	testStyle = !testStyle
}
