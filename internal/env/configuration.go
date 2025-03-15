package env

import (
	"basement/main/internal/logg"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// Don't check for setter and getter methods for these fields in Configuration struct.
const ignore_config_field_checks string = "fields,methods,env,fieldValues,fieldMethods"

// Don't check if these methods were called in applyConfig().
const ignore_config_method_checks string = "Fields,Methods,Init,FieldValues"

// Every field needs to have a setter and getter method except ignored fields.
// Example field: defaultTableSize needs to implement SetDefaultTablesize() and DefaultTableSize().
type Configuration struct {
	fields           []string                 // not part of user configuration
	methods          []string                 // not part of user configuration
	fieldValues      map[string]fieldMetaData // not part of user configuration
	env              environment              // not part of user configuration
	alwaysAuthorized bool
	defaultTableSize int
	showTableSize    bool
	infoLogsEnabled  bool
	debugLogsEnabled bool
	errorLogsEnabled bool
	useMemoryDB      bool
	dbPath           string
	staticPath       string
	templatePath     string
}

// Init returns false if some Get or Set methods are missing from struct.
func (config *Configuration) Init() error {
	logg.Info("init config")
	configRT := reflect.TypeOf(*config)
	configStructName := configRT.Name()
	if configRT.Kind() != reflect.Struct {
		panic("invalid Configuration struct \"" + configStructName + "\"")
	}

	// Populating methods field.
	populateMethods(config, ignore_config_method_checks)

	missingMethods := populateFieldValues(config, ignore_config_field_checks)
	if len(missingMethods) != 0 {
		mm := make([]string, len(missingMethods))
		for _, e := range missingMethods {
			// logg.Errf("missing %d %s", i, e)
			mm = append(mm, e.Error())
		}
		methods := strings.Join(mm, "")
		return logg.NewError("missing methods " + methods)
	}
	return nil
}

// populateMethods.
//
//	excludeMethodNames: Comma separated string like: "MethodA,MethodB,MethodZ".
func populateMethods(config *Configuration, excludeMethodNames string) {
	configPtrRT := reflect.TypeOf(config)
	excludeMethods := strings.Split(excludeMethodNames, ",")
	for i := 0; i < configPtrRT.NumMethod(); i++ {
		methodName := configPtrRT.Method(i).Name
		if slices.Contains(excludeMethods, methodName) {
			continue
		}
		config.methods = append(config.methods, methodName)
	}
}

func populateFieldValues(config *Configuration, excludeMethodNames string) (missingMethods []error) {
	configRT := reflect.TypeOf(*config)
	configStructName := configRT.Name()

	configRV := reflect.ValueOf(*config)
	excludeFields := strings.Split(ignore_config_field_checks, ",")
	config.fieldValues = make(map[string]fieldMetaData)
	for i := 0; i < configRV.NumField(); i++ {
		fieldName := configRT.Field(i).Name
		logg.Debugf("field: \"%s\"=\"%v\"", fieldName, configRV.Field(i))
		if slices.Contains(excludeFields, fieldName) {
			continue
		}

		expectedPublicFieldSetter := "Set" + strings.ToUpper(fieldName[:1]) + fieldName[1:]
		expectedPublicFieldGetter := strings.ToUpper(fieldName[:1]) + fieldName[1:]

		var err error
		if !slices.Contains(config.methods, expectedPublicFieldSetter) {
			err = logg.NewError("Missing setter method \"" + expectedPublicFieldSetter + "\" for field \"" + fieldName + "\" in \"type " + configStructName + " struct {...}\". Implement this method.")
			missingMethods = append(missingMethods, err)
		}
		if !slices.Contains(config.methods, expectedPublicFieldGetter) {
			err = logg.NewError("Missing getter method \"" + expectedPublicFieldGetter + "\" for field \"" + fieldName + "\" in \"type " + configStructName + " struct {...}\". Implement this method.")
			missingMethods = append(missingMethods, err)
		}

		config.fields = append(config.fields, fieldName)
		field := fieldMetaData{
			Value:  fmt.Sprintf("%v", configRV.Field(i)),
			Setter: expectedPublicFieldSetter,
			Getter: expectedPublicFieldGetter,
			Kind:   configRV.Field(i).Kind(),
		}
		logg.Debugf("Init() fieldname: %s, fieldsetter: %s kind: %s", fieldName, field.Setter, field.Kind.String())
		config.fieldValues[fieldName] = field
	}
	return missingMethods
}

func (config *Configuration) Fields() []string {
	return config.fields
}

func (config *Configuration) Methods() []string {
	return config.methods
}

func (c Configuration) Description() string {
	return fmt.Sprintf("environment config: isProduction=%t, isDevelopment=%t, defaultTableSize=%d, showTableSize=%v", Production(), Development(), configInstance.defaultTableSize, configInstance.showTableSize)
}

// SetDevelopment sets environment setting to development.
//
// This setting is for making development tasks easier by reducing
// checks, validation strictness, adding more logging information etc.
func (c *Configuration) SetDevelopment() *Configuration {
	c.env = env_dev
	logg.InfoForceOutput(4, "environment is development")
	return c
}

// Development returns true if development environment is set.
// Default is false.
func Development() bool {
	return configInstance.env == env_dev
}

// SetDbPath sets the path to the configured database.
func (c *Configuration) SetDbPath(path string) *Configuration {
	if c.useMemoryDB && path != ":memory:" {
		logg.Fatalf("Trying to set db path to \"%s\". Can't set DB path other than \":memory:\" and use in memory db at the same time", path)
	}
	if path == "" {
		logg.Fatal("Can't leave DB path empty.")
	}
	c.dbPath = path
	loadLog("set dbpath to "+path, 1)
	return c
}

// DbPath returns the path to the configured database.
func (c *Configuration) DbPath() string {
	return c.dbPath
}

// SetShowTableSize sets if the table size option should be shown in the list template.
func (c *Configuration) SetShowTableSize(show bool) *Configuration {
	c.showTableSize = show
	loadLog(fmt.Sprintf("set showTablesize to %v", show), 2)
	return c
}

// ShowTableSize returns true if the table size option should be shown in the list template.
func (c *Configuration) ShowTableSize() bool {
	return c.showTableSize
}

// SetDefaultTableSize sets the default amount of shown rows in a table on a page with a list.
func (c *Configuration) SetDefaultTableSize(size int) *Configuration {
	if size < 1 {
		logg.Fatalf("[SetDefaultTableSize] size must be above 1 but is %d", size)
	}
	c.defaultTableSize = size
	loadLog(fmt.Sprintf("set tablesize to %d", size), 2)
	return c
}

// DefaultTableSize returns the default amount of shown rows in a table on a page with a list.
func (c *Configuration) DefaultTableSize() int {
	return configInstance.defaultTableSize
}

// SetUseMemoryDB sets if DB should use memory instead of files.
func (c *Configuration) SetUseMemoryDB(useMemory bool) *Configuration {
	c.useMemoryDB = useMemory
	loadLog(fmt.Sprintf("set use memory db to %v", c.useMemoryDB), 2)
	return c
}

// UseMemoryDB returns true if DB should use memory instead of file.
func (c *Configuration) UseMemoryDB() bool {
	return configInstance.useMemoryDB
}

// SetProduction sets environment setting to production.
//
// This is setting is for production use.
// Enables all checks, increases strictness, removes debug logging, etc.
func (c *Configuration) SetProduction() *Configuration {
	c.env = env_prod
	if configLoaded {
		logg.InfoForceOutput(4, "environment is production")
	}
	return c
}

// Production returns true if production environment is set.
// Default is true.
func Production() bool {
	return configInstance.env == env_prod
}

// SetTest sets environment setting to test.
//
// This is setting is for experimental testing.
func (c *Configuration) SetTest() *Configuration {
	c.env = env_test
	loadLog("environment set to TEST", 1)
	return c
}

func (c *Configuration) SetInfoLogsEnabled(enabled bool) *Configuration {
	c.infoLogsEnabled = enabled
	if enabled {
		logg.EnableInfoLoggerS()
	} else {
		logg.DisableInfoLoggerS()
	}
	loadLog(fmt.Sprintf("set info logs enabled = %v ", enabled), 2)
	return c
}

func (c *Configuration) InfoLogsEnabled() bool {
	return logg.InfoLoggerEnabled()
}

func (c *Configuration) SetDebugLogsEnabled(enabled bool) *Configuration {
	c.debugLogsEnabled = enabled
	if enabled {
		logg.EnableDebugLoggerS()
	} else {
		logg.DisableDebugLoggerS()
	}
	loadLog(fmt.Sprintf("set debug logs enabled = %v ", enabled), 2)
	return c
}

func (c *Configuration) DebugLogsEnabled() bool {
	return logg.DebugLoggerEnabled()
}

func (c *Configuration) SetErrorLogsEnabled(enabled bool) *Configuration {
	c.errorLogsEnabled = enabled
	if enabled {
		logg.EnableErrorLoggerS()
	} else {
		logg.DisableErrorLoggerS()
	}
	loadLog(fmt.Sprintf("set error logs enabled = %v ", enabled), 2)
	return c
}

func (c *Configuration) ErrorLogsEnabled() bool {
	return logg.ErrorLoggerEnabled()
}

// Convenience logging function for loading an option.
func loadLog(s string, l int) {
	if configLoaded {
		if logg.DebugLoggerEnabled() {
			logg.Alog(logg.DebugLogger(), l+1, s)
		}
	}
}

// SetTemplatePath sets the path to the configured template directory.
func (c *Configuration) SetTemplatePath(path string) *Configuration {
	if path == "" {
		logg.Fatal("Can't set TemplatePath to \"\".")
	}
	c.templatePath = path
	loadLog("set TemplatePath to "+path, 1)
	return c
}

// SetStaticPath sets the path to the configured static files directory.
func (c *Configuration) SetStaticPath(path string) *Configuration {
	if path == "" {
		logg.Fatal("Can't set StaticPath to \"\".")
	}
	c.staticPath = path
	loadLog("set StaticPath to "+path, 1)
	return c
}

// TemplatePath returns the path to the configured template directory.
func (c *Configuration) TemplatePath() string {
	return c.templatePath
}

// StaticPath returns the path to the configured static files directory.
func (c *Configuration) StaticPath() string {
	return c.staticPath
}

// if SetAlwaysAuthorized is set to true, all authorization checks will be skipped.
func (c *Configuration) SetAlwaysAuthorized(state bool) *Configuration {
	c.alwaysAuthorized = state
	loadLog(fmt.Sprintf("set alwaysAuthorized to %t", state), 2)
	if state && Production() {
		logg.Warning("alwaysAuthorized=true, all auth checks will be skipped")
	}
	return c
}

// if AlwaysAuthorized is true, all authorization checks will be skipped.
func (c *Configuration) AlwaysAuthorized() bool {
	return c.alwaysAuthorized
}
