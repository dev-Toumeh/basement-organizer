package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/items"
	"basement/main/internal/templates"
)

const (
	STATIC                      string = "/static/"
	ITEMS_FILE_PATH             string = "internal/auth/items.json"
	USERS_FILE_PATH             string = "internal/auth/users2.json"
	API_V1_READ_ITEM            string = "/api/v1/read/item/id"
	PERSONAL_PAGE_ROUTE         string = "/personal-page"
	PERSONAL_PAGE_TEMPLATE_PATH string = "internal/templates/personal-page.html"
)

func RegisterRoutes(db *auth.JsonDB) {
	http.Handle(STATIC, http.StripPrefix("/static/", http.FileServer(http.Dir("internal/static"))))
	http.HandleFunc("/", HomePage)
	http.HandleFunc(PERSONAL_PAGE_ROUTE, PersonalPage)

	authRoutes(db)
	apiRoutes(db)
}

func authRoutes(db *auth.JsonDB) {
	http.HandleFunc("/login", db.LoginHandler)
	http.HandleFunc("/register", auth.RegisterHandler(db))
	http.HandleFunc("/logout", auth.LogoutHandler)
}

func apiRoutes(db *auth.JsonDB) {
	http.HandleFunc("/api/v1/create/item", items.CreateItemHandler(db))
	http.HandleFunc("/api/v1/read/items", ReadItems)
	http.HandleFunc(API_V1_READ_ITEM, ReadItem)
	http.HandleFunc("/api/v1/update/item/id", UpdateItem)
	http.HandleFunc("/api/v1/delete/item", DeleteItem)
}

func PersonalPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	data := templates.PageTemplate{
		Title:         "Personal",
		Authenticated: authenticated,
		User:          auth.Username(r),
	}

	if err := templates.ApplyPageTemplate(w, PERSONAL_PAGE_TEMPLATE_PATH, data); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}
