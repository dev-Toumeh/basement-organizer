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

// Copy of preset production config.
func DefaultTestConfig() Configuration {
	return defaultTestConfig
}

var defaultTestConfig = defaultTestConfigPreset
var configInstance *Configuration = &defaultTestConfig
