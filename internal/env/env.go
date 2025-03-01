package env

import (
	"basement/main/internal/logg"
	"fmt"
)

type environment int

const (
	env_dev environment = iota
	env_prod
	env_test
)

const ignoreConfigFieldChecks string = "fields,methods,env"   // Don't check for setter and getter methods for these fields in Configuration struct.
const ignoreConfigMethodChecks string = "Fields,Methods,Init" // Don't check if these methods were called in ApplyConfig().

// Every field needs to have a setter and getter method except ignored fields.
// Example field: defaultTableSize needs to implement SetDefaultTablesize() and DefaultTableSize().
type Configuration struct {
	fields           []string
	methods          []string
	env              environment
	defaultTableSize int
	showTableSize    bool
	infoLogsEnabled  bool
	debugLogsEnabled bool
	errorLogsEnabled bool
	useMemoryDB      bool
	dbPath           string
}

// Copy of preset development config.
func DefaultDevelopmentConfig() Configuration {
	return defaultDevConfig
}

var defaultDevConfig Configuration = Configuration{
	env:              env_dev,
	defaultTableSize: 10,
	infoLogsEnabled:  true,
	debugLogsEnabled: false,
	errorLogsEnabled: true,
	useMemoryDB:      false,
	dbPath:           "./internal/database/sqlite-database.db",
}

// Copy of preset production config.
func DefaultProductionConfig() Configuration {
	return defaultProdConfig
}

var defaultProdConfig Configuration = Configuration{
	env:              env_prod,
	defaultTableSize: 15,
	infoLogsEnabled:  true,
	debugLogsEnabled: false,
	errorLogsEnabled: true,
	useMemoryDB:      false,
	dbPath:           "./internal/database/sqlite-database-prod-v1.db",
}

// CurrentConfig returns the currently applied config instance across the project.
func CurrentConfig() *Configuration {
	return configInstance
}

var configInstance *Configuration = &Configuration{}

// Keeps track if config was already loaded.
var configLoaded = false

// ApplyConfig uses a copy of a config, checks for correctness and applies the options to the currently shared config instance.
func ApplyConfig(c Configuration) {
	missingMethods := configInstance.Init()
	if missingMethods {
		logg.Fatal("config has missing methods")
	}

	// logg.Debugf("configInstance.Methods() %v", configInstance.Methods())
	if configLoaded {
		logg.InfoForceOutput(4, "reload config")
	}
	configLoaded = true

	configInstance.SetInfoLogsEnabled(c.infoLogsEnabled)
	configInstance.SetDebugLogsEnabled(c.debugLogsEnabled)
	configInstance.SetErrorLogsEnabled(c.errorLogsEnabled)
	configInstance.SetShowTableSize(c.showTableSize)
	configInstance.SetDefaultTableSize(c.defaultTableSize)

	if c.useMemoryDB {
		configInstance.SetDbPath(":memory:")
		configInstance.SetUseMemoryDB(c.useMemoryDB)
	} else {
		configInstance.SetDbPath(c.dbPath)
	}

	switch c.env {
	case env_dev:
		configInstance.SetDevelopment()
		break
	case env_prod:
		configInstance.SetProduction()
		break
	case env_test:
		configInstance.SetTest()
		break
	}

	applyConfigFuncName := "ApplyConfig"
	funcs := calledFunctionsFrom("./internal/env/env.go", applyConfigFuncName)
	// logg.Infof("applied methods %v", funcs)
	allApplied := checkIfAllSettingsWereApplied(configInstance, funcs)
	if !allApplied {
		logg.Fatalf("Not all Config settings were applied in \"%s()\"", applyConfigFuncName)
	}
	loadLog("config loaded", 2)
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
	loadLog("environment is development", 1)
	return c
}

// Development returns true if development environment is set.
// Default is false.
func Development() bool {
	return configInstance.env == env_dev
}

// SetDbPath sets the path to the configured database.
func (c *Configuration) SetDbPath(path string) *Configuration {
	if c.useMemoryDB {
		logg.Fatal("Can't set DB path. It is set to use memory")
	}
	if path == "" {
		logg.Fatal("Can't set DB path to \"\".")
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
		logg.Fatalf("size must be above 1 but is %d", size)
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
	loadLog(fmt.Sprintf("set info logs enabled = %v ", enabled), 3)
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
	loadLog(fmt.Sprintf("set debug logs enabled = %v ", enabled), 3)
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
	loadLog(fmt.Sprintf("set error logs enabled = %v ", enabled), 3)
	return c
}

func (c *Configuration) ErrorLogsEnabled() bool {
	return logg.ErrorLoggerEnabled()
}

func loadLog(s string, l int) {
	if configLoaded {
		if logg.DebugLoggerEnabled() {
			logg.Alog(logg.DebugLogger(), l+2, s)
		}
	}
}

// var configurations = make(map[string]string)

// DescribeDevelopmentOverride adds a label and description for a development config change.
//
// Use this to document and track environment-specific settings when custom logic is applied. Returns an error if the label already exists.
//
// Example:
// err := DescribeDevelopmentOverride("SkipPasswordValidations", "Disables password validations (password length, password strength) for registering a user.")
// func DescribeDevelopmentOverride(label string, description string) error {
// 	val, ok := configurations[label]
// 	if ok {
// 		return errors.New(fmt.Sprintf("Environment configuration '%v' already exists. Description:'%v'\n", label, val))
// 	}
// 	configurations[label] = description
// 	return nil
// }

func (config *Configuration) Fields() []string {
	return config.fields
}

func (config *Configuration) Methods() []string {
	return config.methods
}
