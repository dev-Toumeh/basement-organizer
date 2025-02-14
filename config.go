package main

import (
	"basement/main/internal/env"
)

func LoadConfig() *env.Configuration {
	c := env.DefaultDevelopmentConfig()
	env.LoadConfig(c)

	return env.Config()
}
