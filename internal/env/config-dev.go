//go:build dev

package env

const configFile string = "config-dev.conf"

var defaultDevConfigPreset Configuration = Configuration{
	env:              env_dev,
	alwaysAuthorized: true,
	defaultTableSize: 10,
	infoLogsEnabled:  true,
	debugLogsEnabled: true,
	errorLogsEnabled: true,
	useMemoryDB:      false,
	dbPath:           "./internal/database/sqlite-database.db",
	staticPath:       "./internal/static",
	templatePath:     "./internal",
}

// Copy of preset development config.
func DefaultConfigPreset() Configuration {
	return defaultDevConfigPreset
}

var defaultConfig = defaultDevConfigPreset
var configInstance *Configuration = &defaultConfig
