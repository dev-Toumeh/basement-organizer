//go:build prod

package env

import (
	"os"
)

const configFile string = "config.conf"

var homeDir string = os.Getenv("HOME")

var defaultProdConfigPreset Configuration = Configuration{
	env:              env_prod,
	alwaysAuthorized: false, // temporary until auth is done
	defaultTableSize: 15,
	infoLogsEnabled:  true,
	debugLogsEnabled: false,
	errorLogsEnabled: true,
	useMemoryDB:      false,
	dbPath:           homeDir + "/.local/share/basement-organizer/internal/database/sqlite-database.db",
	staticPath:       homeDir + "/.local/share/basement-organizer/internal/static",
	templatePath:     homeDir + "/.local/share/basement-organizer/internal",
}

// Copy of preset production config.
func DefaultProductionConfig() Configuration {
	return defaultProdConfig
}

var defaultProdConfig = defaultProdConfigPreset
var configInstance *Configuration = &defaultProdConfig
