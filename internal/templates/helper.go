package templates

import (
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"maps"
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

const errorNotification notificationType = "error"
const SuccessNotification notificationType = "success"
const warningNotification notificationType = "warning"

type notificationType string

type oobNotificationData struct {
	Message        string
	NotificationId int
	Duration       int
	Type           notificationType
}

type PageTemplate struct {
	Title          string
	Authenticated  bool
	User           string
	Debug          bool
	NotFound       bool
	EnvDevelopment bool
}

func (tmpl PageTemplate) Map() map[string]any {
	return map[string]any{
		"Title":          tmpl.Title,
		"Authenticated":  tmpl.Authenticated,
		"User":           tmpl.User,
		"Debug":          tmpl.Debug,
		"NotFound":       tmpl.NotFound,
		"EnvDevelopment": tmpl.EnvDevelopment,
	}
}

// NewPageTemplate returns default data struct for page templates.
func NewPageTemplate() PageTemplate {
	return PageTemplate{
		Title:          "Default Page",
		Authenticated:  false,
		User:           "Default User",
		Debug:          DEBUG_STYLE,
		NotFound:       false,
		EnvDevelopment: env.Development(),
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
func InitTemplates(directory string) error {
	var err error
	dir := ""
	if directory == "" {
		dir = "internal/templates"
	} else {
		dir = directory
	}
	internalTemplate, _, err = ParseDirectory(dir)
	if err != nil {
		return logg.Errorf("Init Templates failed %w", err)
	}

	log.Println("Templates initialized")
	return nil
}

// templateInlineMap is used inside a template directly.
type templateInlineMap struct {
	Map map[string]string
}

// newMap defines a template function "map" for inline definition of data to be passed to other templates.
//
//	Example usage:
//	// Defines a map[string]string = {"key":"value", "key2":"value2", ...}
//	{{ $inlineMap := map "key" "value" "key2" "value2" ... }}
//
//	// Pass variable to a template "another-template" with .Map()
//	{{ template "another-template" $inlineMap.Map() }}
func newMap(values ...string) (*templateInlineMap, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("missing keys or values")
	}

	m := make(map[string]string, 0)
	for i := 0; i < len(values); i += 2 {
		k := values[i]
		v := values[i+1]
		m[k] = v
	}
	return &templateInlineMap{m}, nil
}

func (v *templateInlineMap) Set(key string, value string) string {
	v.Map[key] = value
	return ""
}

// Recursively parse all files in directory, including sub-directories.
func ParseDirectory(dirpath string) (*template.Template, []string, error) {
	internalTemplate = template.New("main")
	internalTemplate.Funcs(template.FuncMap{"map": newMap})
	paths, err := allFilePathsInDirectory(dirpath)
	if err != nil {
		return nil, nil, err
	}
	tmpl, err := internalTemplate.ParseFiles(paths...)
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
		return logg.Errorf("Can't render template %w", err)
	}
	return nil
}

// Mapable is used for rendering templates.
// It is more convenient than using only structs.
// You can check if keys are available without gettings errors during runtime.
//
// For example a struct which implements this interface.
//
//	Mystruct {.FieldX string}
//
// Mystruct doesn't have a "FieldY".
// A check check inside a template `{{ if .FieldY }}` will usually return an error.
// But if this struct implements this interface, passing the Map() to the template will work as expected.
type Mapable interface {
	Map() map[string]any
}

type Renderable interface {
	Render() error
}

// RenderMaps builds multiple templates and renders it.
func RenderMaps(w io.Writer, baseTemplateName string, templates []Mapable) error {
	data := make(map[string]any, 0)
	for _, tmpl := range templates {
		maps.Copy(data, tmpl.Map())
		// logg.Debug(tmpl.Map())
	}

	// logg.Debug(data)
	err := SafeRender(w, baseTemplateName, data)
	if err != nil {
		return logg.Errorf("Can't render maps %w", err)
	}
	return nil
}

// SafeRender will write to the writer "w" only if there are no errors executing the template.
func SafeRender(w io.Writer, name string, data any) error {
	wbs := bytes.NewBufferString("")
	err := internalTemplate.ExecuteTemplate(wbs, name, data)
	if err != nil {
		// return logg.Errorf("Can't execute template \"%s\" with data \"%v\"", err)
		return logg.Errorf("Can't execute template \"%s\" with data \"%v\".\n\t%w", name, data, err)
	}
	w.Write(wbs.Bytes())
	return nil
}

// CanRender will return nil if rendering is ok.
func CanRender(name string, data any) error {
	wbs := bytes.NewBufferString("")
	err := internalTemplate.ExecuteTemplate(wbs, name, data)
	if err != nil {
		return logg.Errorf("Can't execute template %w", err)
	}
	return nil
}

// RenderErrorNotification shows a brief error notification.
func RenderErrorNotification(w io.Writer, message string) error {
	data := newOobNotificationData(message, errorNotification)
	err := internalTemplate.ExecuteTemplate(w, TEMPLATE_OOB_NOTIFICATION, data)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return err
	}
	return nil
}

// RenderWarningNotification shows a brief warning notification.
func RenderWarningNotification(w io.Writer, message string) error {
	data := newOobNotificationData(message, warningNotification)
	logg.Debug(data.Message, data.Type, data.Duration, data.NotificationId)
	err := internalTemplate.ExecuteTemplate(w, TEMPLATE_OOB_NOTIFICATION, data)
	if err != nil {
		logg.Fatal(err)
		return err
	}
	return nil
}

// RenderSuccessNotification shows a brief success notification.
func RenderSuccessNotification(w io.Writer, message string) error {
	data := newOobNotificationData(message, SuccessNotification)
	err := internalTemplate.ExecuteTemplate(w, TEMPLATE_OOB_NOTIFICATION, data)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return err
	}
	return nil
}

// newOobNotificationData creates newOobNotificationData for error notification with default values.
//
// Defaut Duration = 2000 // 2 seconds
func newOobNotificationData(message string, messageType notificationType) oobNotificationData {
	return oobNotificationData{
		Message:        message,
		NotificationId: rand.Int(),
		Duration:       3000,
		Type:           messageType,
	}
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
