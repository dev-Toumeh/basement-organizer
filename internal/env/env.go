package env

import (
	"basement/main/internal/logg"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

type environment int

const (
	env_dev environment = iota
	env_prod
	env_test
)

func LoadConfig() (*Configuration, error) {
	c := configInstance

	err := c.Init()
	if err != nil {
		logg.Fatal("LoadConfig failed" + err.Error())
	}

	// use defaults
	_, err = os.Stat(configFile)
	if err != nil {
		applyConfig(*c)
		err = createFileFromConfiguration(configFile, c)
		if err != nil {
			return c, logg.WrapErr(err)
		}
		return c, nil
	}

	// override options from config file
	errs := applyConfigFileOptions(configFile, c)
	if len(errs) != 0 {
		return c, logg.WrapErr(errors.Join(errs...))
	}

	applyConfig(*c)

	return CurrentConfig(), nil
}

var backupConfig *Configuration = &defaultConfig

func LoadDefault() (*Configuration, error) {
	logg.Info("load default config")
	defaultConfig = DefaultConfigPreset()
	configInstance = &defaultConfig
	return LoadConfig()
}

// CurrentConfig returns the currently applied config instance across the project.
func CurrentConfig() *Configuration {
	return configInstance
}

// Keeps track if config was already loaded.
var configLoaded = false

// ApplyParsedConfigOptions applies settings from parsed to configuration with validation.
// On error will roll back to previous working config.
func ApplyParsedConfigOptions(parsed map[string]string, c *Configuration) []error {
	backup := *CurrentConfig()
	logg.Debug("applying parsed config options")
	errs := applyParsedOptions(parsed, c)
	if len(errs) > 0 {
		for _, e := range errs {
			logg.Err(e)
		}
		applyConfig(backup)
		c = CurrentConfig()
		logg.Warning("config rolled back")
		return errs
	}

	errs = ValidateOptions(c)
	if len(errs) > 0 {
		for _, e := range errs {
			logg.Err(e)
		}

		applyConfig(backup)
		c = CurrentConfig()
		logg.Warning("config rolled back")
		return errs
	}
	c.Init()

	if configLoaded {
		logg.InfoForceOutput(4, "config reloaded")
	}
	configLoaded = true
	return nil
}

// applyConfig uses a copy of a config, checks for correctness and applies the options to the currently shared config instance.
func applyConfig(c Configuration) {
	defer configInstance.Init()
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

func WriteCurrentConfigToFile() error {
	return createFileFromConfiguration(configFile, CurrentConfig())
}

func createFileFromConfiguration(path string, config *Configuration) error {
	logg.Info("creating config file \"" + path + "\"")
	lines := make([]string, len(config.fieldValues)+1)
	lines[0] = "# Runtime configuration values"
	i := 1
	data := config.FieldValues()
	for k, v := range data {
		lines[i] = k + "=" + string(v.Value)
		i++
	}
	sort.Strings(lines)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		return logg.WrapErr(fmt.Errorf("can't create config file \""+path+"\". %w", err))
	}

	out := strings.Join(lines, "\n")
	logg.Debugf("write to config: \"%s\"", out)
	_, err = file.WriteString(out)
	if err != nil {
		return logg.WrapErr(fmt.Errorf("can't write config file to disk \""+path+"\". %w", err))
	}
	return nil
}

// applyConfigFileOptions only applies options that are valid.
func applyConfigFileOptions(configFile string, c *Configuration) []error {
	parsed, errs := parseConfigFile(configFile, c)
	if len(errs) != 0 {
		errs[len(errs)-1] = logg.WrapErr(errs[len(errs)-1])
		return errs
	}
	errs = ApplyParsedConfigOptions(parsed, c)
	return errs
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
