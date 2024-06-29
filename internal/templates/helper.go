package templates

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

// ApplyPageTemplate generates a complete page from the "page.html" template.
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

var internalTemplate *template.Template

// InitTemplates loads all templates from "internal/templates" directory.
func InitTemplates() {
	var err error
	internalTemplate, err = parseDirectory("internal/templates")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Templates initialized")
}

// Recursively parse all files in directory, including sub-directories.
func parseDirectory(dirpath string) (*template.Template, error) {
	paths, err := allFilePathsInDirectory(dirpath)
	if err != nil {
		return nil, err
	}
	return template.ParseFiles(paths...)
}

// Recursively get all file paths in directory, including sub-directories.
func allFilePathsInDirectory(dirpath string) ([]string, error) {
	var paths []string
	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

// Render applies data to a defined template and writes result back to the writer.
func Render(wr io.Writer, name string, data any) error {
	err := internalTemplate.ExecuteTemplate(wr, name, data)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(wr, err)
		return err
	}
	return nil
}
