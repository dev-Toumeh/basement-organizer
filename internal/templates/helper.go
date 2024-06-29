package templates

import (
	"log"
	"net/http"
	"text/template"
)

const (
	PAGE_TEMPLATE                       string = "page"
	PAGE_TEMPLATE_FILENAME              string = PAGE_TEMPLATE + ".html"
	PAGE_TEMPLATE_FILE_WTH_PATH         string = "internal/templates/" + PAGE_TEMPLATE_FILENAME
	REGISTER_TEMPLATE_FILE_WITH_PATH    string = "internal/templates/auth/register.html"
	LOGIN_TEMPLATE_FILE_WITH_PATH       string = "internal/templates/auth/login.html"
	CREATE_ITEM_TEMPLATE_FILE_WITH_PATH string = "internal/templates/items/create-item.html"
)

type PageTemplate struct {
	Title         string
	Authenticated bool
	User          string
}

func ApplyPageTemplate(w http.ResponseWriter, bodyTemplateFile string, data interface{}) error {
	tmpl, err := template.ParseFiles(PAGE_TEMPLATE_FILE_WTH_PATH, bodyTemplateFile)
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
