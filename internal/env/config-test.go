//go:build test

package env

const configFile string = "config-test.conf"

var defaultTestConfigPreset Configuration = Configuration{
	env:              env_prod,
	alwaysAuthorized: false,
	defaultTableSize: 15,
	infoLogsEnabled:  false,
	debugLogsEnabled: false,
	errorLogsEnabled: false,
	useMemoryDB:      false,
	dbPath:           ":memory:",
	staticPath:       "./internal/static",
	templatePath:     "./internal",
}

// Copy of preset test config.
func DefaultConfigPreset() Configuration {
	return defaultTestConfigPreset
}

var defaultConfig = defaultTestConfigPreset
var configInstance *Configuration = &defaultConfig
