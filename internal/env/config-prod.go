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
	debugLogsEnabled: true,
	errorLogsEnabled: true,
	useMemoryDB:      false,
	dbPath:           homeDir + "/.local/share/basement-organizer/internal/database/sqlite-database.db",
	staticPath:       homeDir + "/.local/share/basement-organizer/internal/static",
	templatePath:     homeDir + "/.local/share/basement-organizer/internal",
}

// Copy of preset production config.
func DefaultConfigPreset() Configuration {
	return defaultProdConfigPreset
}

var defaultConfig = defaultProdConfigPreset
var configInstance *Configuration = &defaultConfig
