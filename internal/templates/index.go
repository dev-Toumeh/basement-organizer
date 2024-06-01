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
		fmt.Fprintln(w, "failed")
		return
	}

	templateData := IndexTemplate{Title: "index"}

	if err := tmpl.ExecuteTemplate(w, "index", templateData); err != nil {
		log.Printf("index page: %v\n", err)
		fmt.Fprintln(w, "failed")
	}

	if err := tmpl.ExecuteTemplate(w, "body", templateData); err != nil {
		log.Printf("index page: %v\n", err)
		fmt.Fprintln(w, "failed")
	}
}
