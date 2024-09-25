package routes

import (
	"fmt"
	"net/http"

	"basement/main/internal/auth"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/templates"
)

func itemsPage(w http.ResponseWriter, r *http.Request) {
	authenticated, _ := auth.Authenticated(r)
	user, _ := auth.UserSessionData(r)
	data := templates.PageTemplate{
		Title:         "Personal",
		Authenticated: authenticated,
		User:          user,
	}

	err := templates.Render(w, "items-page", data.Map())
	if err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

// this function will generate a new from
func itemTemp(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "create-item-form", ""); err != nil {
		fmt.Fprintln(w, "failed")
		return
	}
}

// generate Search Item Template, in case of get request
func searchItemTemp(w http.ResponseWriter, r *http.Request) {
	err := templates.Render(w, "search-item-form", items.NewSearchItemInputTemplate().Map())
	if err != nil {
		logg.Debug(err)
		templates.RenderErrorNotification(w, "something wrong happened")
	}
}

func moveItem(w http.ResponseWriter, r *http.Request) {
	err := templates.Render(w, "item-move-to", "")
	if err != nil {
		logg.Debug(err)
		templates.RenderErrorNotification(w, "something wrong happened")
	}
}
