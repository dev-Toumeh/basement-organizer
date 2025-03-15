//go:build prod

package env

const configFile string = "config.conf"

var defaultProdConfigPreset Configuration = Configuration{
	env:              env_prod,
	alwaysAuthorized: true, // temporary until auth is done
	defaultTableSize: 15,
	infoLogsEnabled:  true,
	debugLogsEnabled: false,
	errorLogsEnabled: true,
	useMemoryDB:      false,
	dbPath:           "./internal/database/sqlite-database-prod-v1.db",
	staticPath:       "/opt/basement-organizer/internal/static/",
	templatePath:     "/opt/basement-organizer/internal/",
}

// Copy of preset production config.
func DefaultProductionConfig() Configuration {
	return defaultProdConfig
}

var defaultProdConfig = defaultProdConfigPreset
var configInstance *Configuration = &defaultProdConfig
