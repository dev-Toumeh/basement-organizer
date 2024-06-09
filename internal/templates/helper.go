package templates

import (
	"log"
	"net/http"
	"text/template"
)

const (
	PAGE_TEMPLATE          string = "page"
	PAGE_TEMPLATE_FILENAME string = PAGE_TEMPLATE + ".html"
	PAGE_TEMPLATE_FILE     string = "internal/templates/" + PAGE_TEMPLATE_FILENAME
)

type PageTemplate struct {
	Title         string
	Authenticated bool
	User          string
}

func ApplyPageTemplate(w http.ResponseWriter, bodyTemplateFile string, data interface{}) error {
	tmpl, err := template.ParseFiles(PAGE_TEMPLATE_FILE, bodyTemplateFile)
	if err != nil {
		log.Printf("%v, %v: %v\n", PAGE_TEMPLATE, bodyTemplateFile, err)
		return err
	}

	if err := tmpl.ExecuteTemplate(w, PAGE_TEMPLATE, data); err != nil {
		log.Printf("%v: %v\n", PAGE_TEMPLATE, err)
		return err
	}
	return nil
}
