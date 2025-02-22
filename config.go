package main

import (
	"basement/main/internal/env"
)

func LoadConfig() *env.Configuration {
	c := env.DefaultProductionConfig()
	env.LoadConfig(c)

	return env.Config()
}
