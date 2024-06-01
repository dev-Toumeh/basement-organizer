package templates

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type RootPageTemplate struct {
	Title string
}

const ROOT_PAGE_TEMPLATE string = "page"
const ROOT_PAGE_TEMPLATE_FILENAME string = ROOT_PAGE_TEMPLATE + ".html"
const ROOT_PAGE_TEMPLATE_FILE string = "internal/templates/" + ROOT_PAGE_TEMPLATE_FILENAME

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(ROOT_PAGE_TEMPLATE_FILE, "internal/templates/home.html")
	if err != nil {
		log.Printf("%v: %v\n", ROOT_PAGE_TEMPLATE, err)
		fmt.Fprintln(w, "failed")
		return
	}

	templateData := RootPageTemplate{Title: "home"}

	if err := tmpl.ExecuteTemplate(w, ROOT_PAGE_TEMPLATE, templateData); err != nil {
		log.Printf("%v: %v\n", ROOT_PAGE_TEMPLATE, err)
		fmt.Fprintln(w, "failed")
	}

}
