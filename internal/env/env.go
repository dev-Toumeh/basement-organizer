package env

import (
	"basement/main/internal/logg"
	"os"
	"runtime"
	"strings"
)

type environment int

const (
	env_dev environment = iota
	env_prod
	env_test
)

func LoadConfig() *Configuration {
	c := configInstance

	err := c.Init()
	if err != nil {
		logg.Fatal("LoadConfig failed" + err.Error())
	}

	// use defaults
	_, err = os.Stat(configFile)
	if err != nil {
		applyConfig(*c)
		CreateFileFromConfiguration(configFile, c)
		return c
	}

	// override options from config file
	applyConfigFileOptions(configFile, c)

	applyConfig(*c)

	return CurrentConfig()
}

// CurrentConfig returns the currently applied config instance across the project.
func CurrentConfig() *Configuration {
	return configInstance
}

// Keeps track if config was already loaded.
var configLoaded = false

// applyConfig uses a copy of a config, checks for correctness and applies the options to the currently shared config instance.
func applyConfig(c Configuration) {
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
	// configInstance.SetConfigFile(c.configFile)

	if c.useMemoryDB {
		configInstance.SetDbPath(":memory:")
		configInstance.SetUseMemoryDB(c.useMemoryDB)
	} else {
		configInstance.SetDbPath(c.dbPath)
	}
	configInstance.SetTemplatePath(c.templatePath)
	configInstance.SetStaticPath(c.staticPath)

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
	configInstance.SetAlwaysAuthorized(c.alwaysAuthorized)

	thisFile, thisFunc := thisFuncAndFileName()
	validateInternals(thisFile, thisFunc)
	loadLog("config loaded", 2)
}

func CreateFileFromConfiguration(path string, config *Configuration) {
	logg.Info("creating config file \"" + path + "\"")
	lines := make([]string, len(config.fieldValues)+1)
	lines[0] = "# Default config values, uncomment and change to override"
	i := 1
	for k, v := range config.fieldValues {
		lines[i] = "#" + k + "=" + string(v.Value)
		logg.Debug(lines[i])
		i++
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic("can't create config file \"" + path + "\": " + err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		panic("can't write config file to disk \"" + path + "\": " + err.Error())
	}
}

// applyConfigFileOptions only applies options that are valid.
// Will exit program on error with error message.
func applyConfigFileOptions(configFile string, c *Configuration) {
	parsed := parseConfigFile(configFile, c)
	errs := applyParsedOptions(c, parsed)
	if len(errs) != 0 {
		errorMessages := ""
		for i, e := range errs {
			nl := ""
			if i != 0 {
				nl = "\n"
			}
			errorMessages += nl + logg.CleanLastError(e)
		}
		logg.Fatalf("config parser for \"%s\" has encountered errors\n%s", configFile, errorMessages)
	}

	errs = validateOptions(c)
	if len(errs) > 0 {
		for _, e := range errs {
			logg.Err(e)
		}
		logg.Fatalf("config check failed for \"%s\"", configFile)
	}
}

// returns full path file name and the function name of the caller.
func thisFuncAndFileName() (fileName string, funcName string) {
	pc, fileName, _, _ := runtime.Caller(1)
	fullFuncName := runtime.FuncForPC(pc).Name()
	funSplit := strings.Split(fullFuncName, "/")
	shortFuncName := funSplit[len(funSplit)-1]
	noPackageShortFuncName := strings.Split(shortFuncName, ".")[1]
	return fileName, noPackageShortFuncName
}
