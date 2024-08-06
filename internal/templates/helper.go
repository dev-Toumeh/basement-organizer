package templates

import (
	"basement/main/internal/logg"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

const (
	// If set to true, will show DebugStyle button with "SwitchDebugStyle()"
	DEBUG_STYLE bool = false
	// TEMPLATE_CONSTANTS_PATH points to auto generated constants.go file
	TEMPLATE_CONSTANTS_PATH string = "internal/templates/constants.go"
	// TEMPLATE_DIR defines directory for all HTML templates
	TEMPLATE_DIR string = "internal/templates/"
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
	tmpl, err := template.ParseFiles(TEMPLATE_PAGE_PATH, bodyTemplateFile)
	if err != nil {
		log.Printf("%v, %v: %v\n", TEMPLATE_PAGE, bodyTemplateFile, err)
		return err
	}

	if err := tmpl.ExecuteTemplate(w, TEMPLATE_PAGE, data); err != nil {
		log.Printf("%v: %v\n", TEMPLATE_PAGE, err)
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

// SafeRender will write to the writer "w" only if there are no errors executing the template.
func SafeRender(w io.Writer, name string, data any) error {
	wbs := bytes.NewBufferString("")
	err := internalTemplate.ExecuteTemplate(wbs, name, data)
	if err != nil {
		logg.Err(err)
		return err
	}
	w.Write(wbs.Bytes())
	return nil
}

// RenderErrorSnackbar shows a brief error notification.
func RenderErrorSnackbar(w io.Writer, message string) error {
	data := NewErrorOobSnackbarData(message)
	err := internalTemplate.ExecuteTemplate(w, TEMPLATE_OOB_SNACKBAR, data)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return err
	}
	return nil
}

type oobSnackbarData struct {
	Message    string
	SnackbarId int
	Duration   int
	Type       snackbarType
}

// NewErrorOobSnackbarData creates oobSnackbarData for error snackbar with default values.
//
// Defaut Duration = 2000 // 2 seconds
func NewErrorOobSnackbarData(message string) oobSnackbarData {
	return oobSnackbarData{
		Message:    message,
		SnackbarId: rand.Int(),
		Duration:   2000,
		Type:       errorSnackbar,
	}
}

type snackbarType string

const errorSnackbar snackbarType = "error"

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
