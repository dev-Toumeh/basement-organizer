//go:build prod

package env

const configFile string = "config.conf"

var defaultProdConfigPreset Configuration = Configuration{
	env:              env_prod,
	defaultTableSize: 15,
	infoLogsEnabled:  true,
	debugLogsEnabled: false,
	errorLogsEnabled: true,
	useMemoryDB:      false,
	dbPath:           "./internal/database/sqlite-database-prod-v1.db",
}

// Copy of preset production config.
func DefaultProductionConfig() Configuration {
	return defaultProdConfig
}

var defaultProdConfig = defaultProdConfigPreset
var configInstance *Configuration = &defaultProdConfig
