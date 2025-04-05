package templates

import "basement/main/internal/env"

// If set to true, will show DebugStyle button with "SwitchDebugStyle()"
const DEBUG_STYLE bool = false

type PageTemplate struct {
	Title          string
	Authenticated  bool
	User           string
	Debug          bool
	NotFound       bool
	RequestOrigin  string
	EnvDevelopment bool
	PageText       string
}

func (tmpl PageTemplate) Map() map[string]any {
	return map[string]any{
		"Title":          tmpl.Title,
		"Authenticated":  tmpl.Authenticated,
		"User":           tmpl.User,
		"Debug":          tmpl.Debug,
		"NotFound":       tmpl.NotFound,
		"RequestOrigin":  tmpl.RequestOrigin,
		"EnvDevelopment": tmpl.EnvDevelopment,
		"PageText":       tmpl.PageText,
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
		RequestOrigin:  "",
		EnvDevelopment: env.Development(),
	}
}

// ApplyPageTemplate generates a complete page from the "page.html" template.
// func ApplyPageTemplate(w http.ResponseWriter, bodyTemplateFile string, data interface{}) error {
// 	tmpl, err := template.ParseFiles(TEMPLATE_PAGE_PATH, bodyTemplateFile)
// 	if err != nil {
// 		log.Printf("%v, %v: %v\n", TEMPLATE_PAGE, bodyTemplateFile, err)
// 		return err
// 	}
//
// 	if err := tmpl.ExecuteTemplate(w, TEMPLATE_PAGE, data); err != nil {
// 		log.Printf("%v: %v\n", TEMPLATE_PAGE, err)
// 		return err
// 	}
// 	return nil
// }
