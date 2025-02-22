package env

import (
	"basement/main/internal/logg"
	"fmt"
)

type Environment int

const (
	DEVELOPMENT = iota
	PRODUCTION
	TEST
)

type Configuration struct {
	env              Environment
	defaultTableSize int
	showTableSize    bool
	infoLogsEnabled  bool
	debugLogsEnabled bool
	errorLogsEnabled bool
	useMemoryDB      bool
	dbPath           string
	TempPath         string
}

var defaultDevConfig Configuration = Configuration{
	env:              DEVELOPMENT,
	defaultTableSize: 10,
	infoLogsEnabled:  true,
	debugLogsEnabled: true,
	errorLogsEnabled: true,
	useMemoryDB:      false,
	dbPath:           "./internal/database/sqlite-database.db",
	TempPath:         "./internal",
}

var defaultProdConfig Configuration = Configuration{
	env:              PRODUCTION,
	defaultTableSize: 15,
	infoLogsEnabled:  true,
	debugLogsEnabled: false,
	errorLogsEnabled: true,
	useMemoryDB:      false,
	dbPath:           "/opt/basement-organizer/internal/database/sqlite-database-prod-v1.db",
	TempPath:         "/opt/basement-organizer/internal",
}

var config *Configuration = &defaultDevConfig

func DefaultProductionConfig() Configuration {
	return defaultProdConfig
}

func DefaultDevelopmentConfig() Configuration {
	return defaultDevConfig
}

func Config() *Configuration {
	return config
}

var loadConfig = false

func LoadConfig(c Configuration) {
	if loadConfig {
		logg.InfoForceOutput(4, "reload config")
	}
	loadConfig = true
	config.SetInfoLogsEnabled(c.infoLogsEnabled)
	config.SetDebugLogsEnabled(c.debugLogsEnabled)
	config.SetErrorLogsEnabled(c.errorLogsEnabled)
	switch c.env {
	case DEVELOPMENT:
		config.SetDevelopment()
		break
	case PRODUCTION:
		config.SetProduction()
		break
	case TEST:
		config.SetTest()
		break
	}
	config.SetShowTableSize(c.showTableSize)
	config.SetDefaultTableSize(c.defaultTableSize)
	if c.useMemoryDB {
		config.SetDBPath(":memory:")
		config.SetUseMemoryDB(c.useMemoryDB)
	} else {
		config.SetDBPath(c.dbPath)
	}
	loadLog("config loaded", 2)
}

func (c Configuration) Description() string {
	return fmt.Sprintf("environment config: isProduction=%t, isDevelopment=%t, defaultTableSize=%d, showTableSize=%v", Production(), Development(), config.defaultTableSize, config.showTableSize)
}

// SetDevelopment sets environment setting to development.
//
// This setting is for making development tasks easier by reducing
// checks, validation strictness, adding more logging information etc.
func (c *Configuration) SetDevelopment() *Configuration {
	c.env = DEVELOPMENT
	loadLog("environment is development", 1)
	return c
}

// Development returns true if development environment is set.
// Default is false.
func Development() bool {
	return config.env == DEVELOPMENT
}

// SetDBPath sets the path to the configured database.
func (c *Configuration) SetDBPath(path string) *Configuration {
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

// DBPath returns the path to the configured database.
func (c *Configuration) DBPath() string {
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
func DefaultTableSize() int {
	return config.defaultTableSize
}

// SetUseMemoryDB sets if DB should use memory instead of files.
func (c *Configuration) SetUseMemoryDB(useMemory bool) *Configuration {
	c.useMemoryDB = useMemory
	loadLog(fmt.Sprintf("set use memory db to %v", c.useMemoryDB), 2)
	return c
}

// UseMemoryDB returns true if DB should use memory instead of file.
func (c *Configuration) UseMemoryDB() bool {
	return config.useMemoryDB
}

// SetProduction sets environment setting to production.
//
// This is setting is for production use.
// Enables all checks, increases strictness, removes debug logging, etc.
func (c *Configuration) SetProduction() *Configuration {
	c.env = PRODUCTION
	if loadConfig {
		logg.InfoForceOutput(4, "environment is production")
	}
	return c
}

// Production returns true if production environment is set.
// Default is true.
func Production() bool {
	return config.env == PRODUCTION
}

// SetTest sets environment setting to test.
//
// This is setting is for experimental testing.
func (c *Configuration) SetTest() *Configuration {
	c.env = TEST
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

func loadLog(s string, l int) {
	if loadConfig {
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
