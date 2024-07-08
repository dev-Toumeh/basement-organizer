package templates

import (
	"bytes"
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
	DEBUG_STYLE                         bool   = false // When true, will show DebugStyle button with "SwitchDebugStyle()"
	TEMPLATE_CONSTANTS_PATH             string = "internal/templates/constants.go"
	TEMPLATE_DIR                        string = "internal/templates/"
)

type PageTemplate struct {
	Title         string
	Authenticated bool
	User          string
	Debug         bool
}

func NewPageTemplate() PageTemplate {
	return PageTemplate{
		Title:         "Default Page",
		Authenticated: false,
		User:          "Default User",
		Debug:         DEBUG_STYLE,
	}
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

func InternalTemplate() *template.Template {
	return internalTemplate
}

// InitTemplates loads all templates from "internal/templates" directory.
func InitTemplates() {
	var err error
	internalTemplate, _, err = ParseDirectory("internal/templates")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Templates initialized")
}

// Recursively parse all files in directory, including sub-directories.
func ParseDirectory(dirpath string) (*template.Template, []string, error) {
	paths, err := allFilePathsInDirectory(dirpath)
	if err != nil {
		return nil, nil, err
	}
	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		return nil, nil, err
	}
	return tmpl, paths, nil
}

// Recursively get all file paths in directory, including sub-directories.
func allFilePathsInDirectory(dirpath string) ([]string, error) {
	var paths []string
	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(info.Name()) == ".html") {
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
func Render(w io.Writer, name string, data any) error {
	err := internalTemplate.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return err
	}
	return nil
}

/*
RedefineTemplateDefinition redefines "targetDefinitionName" template in "tmpl" by using the new definition from "definitionTemplate".

(experimental function)
*/
func RedefineTemplateDefinition(tmpl *template.Template, targetDefinitionName string, definitionTemplate string) {
	// log.Println()
	// tmpl = template.Must(tmpl.New(.Clone())
	newdef := fmt.Sprintf("{{ define `%s`}}%s{{end}}", targetDefinitionName, definitionTemplate)
	template.Must(tmpl.Parse(newdef))
}

/*
RedefineFromOtherTemplateDefinition redefines "targetDefinitionName" template in "targetTmpl" by using the new definition from "sourceDefinitionName"

(experimental function)

Example:

	template1 (source): {{ define "a" }}source definition{{end}}

	template2 (target): {{ define "b" }}to redefine source definition "a"{{end}}

	RedefineFromOtherTemplateDefinition("a", template1, "b", template2)

	remplate2 will become: {{ define "a" }}to redefine source definition "a"{{end}}
*/
func RedefineFromOtherTemplateDefinition(targetDefinitionName string, sourceTmpl *template.Template, sourceDefinitionName string, targetTmpl *template.Template) {
	log.Printf("oldTemplate.Name():%s, oldTemplateDefinitionName:%s, newTemplate.Name():%s, newTemplateDefintionName:%s)", sourceTmpl.Name(), targetDefinitionName, targetTmpl.Name(), sourceDefinitionName)

	var err error

	// newTemplateContent := bytes.NewBufferString("{{define `tyle`}}aaaa{{end}}")
	newTemplateContent := bytes.NewBufferString("")
	tmpTmpl := template.Must(sourceTmpl.Clone())
	err = tmpTmpl.ExecuteTemplate(newTemplateContent, sourceDefinitionName, nil)
	if err != nil {
		panic(err)
	}

	newdef := fmt.Sprintf("{{ define `%s`}}%s{{end}}", targetDefinitionName, newTemplateContent)
	log.Println("newcontent\n", newdef)

	template.Must(targetTmpl.Parse(newdef))
}
