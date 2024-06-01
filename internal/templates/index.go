package templates

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type IndexTemplate struct {
	Title string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("internal/templates/index.html")
	if err != nil {
		log.Printf("Index Page: %v\n", err)
		fmt.Fprintln(w, "failed", err)
		return
	}

	data := IndexTemplate{"index page"}
	terr := tmpl.Execute(w, data)
	if terr != nil {
		log.Printf("index page: %v\n", terr)
		fmt.Fprintln(w, "failed")
	}
}
