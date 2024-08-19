package routes

import (
	"basement/main/internal/database"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"fmt"
	"net/http"
)

func registerBoxRoute() {
	http.HandleFunc("/box", BoxHandler)
}

func BoxHandler(w http.ResponseWriter, r *http.Request) {
	b := items.Box{Label: "asdfasdf"}
	b2 := items.Box{Label: "box 2"}
	err := b2.MoveTo(&b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		logg.Err(err)
	}
	db := r.Context().Value("db").(*database.DB)
	ids, _ := db.ItemIDs()
	item, _ := db.Item(ids[0])
	item.Picture = ""

	b.Items = []*database.Item{&item}
	// data, _ := json.Marshal(b)
	data, _ := b.MarshalJSON()
	fmt.Fprintf(w, "%s", data)
}
