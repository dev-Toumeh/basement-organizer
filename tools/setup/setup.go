package main

import (
	"basement/main/internal/logg"
	"go/format"
	"os"
)

// config.go should be in root directory.
// Function LoadConfig() from main.go is defined in inside config.go.
const configFile = "config.go"
const configFileContent = `
package main

import (
	"basement/main/internal/env"
)

func LoadConfig() *env.Configuration {
	c := env.DefaultDevelopmentConfig()
	env.LoadConfig(&c)

	return env.Config()
}`

func main() {
	var file *os.File
	_, err := os.Stat(configFile)
	// config.go exists
	if err == nil {
		logg.InfoForceOutput(2, configFile+" exists.")
		return
	}

	file, err = os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	newFileContentFormatted, err := format.Source([]byte(configFileContent))
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(string(newFileContentFormatted))
	if err != nil {
		panic(err)
	}
	logg.InfoForceOutput(2, configFile+" created. Restart main.go to apply configuration.")
}
