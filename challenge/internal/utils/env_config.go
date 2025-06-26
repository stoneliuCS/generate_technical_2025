package utils

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type EnvConfig struct {
	// Required environment variables.
	DB_HOST     string `env:"DB_HOST, required"`
	DB_PORT     string `env:"DB_PORT, required"`
	DB_USER     string `env:"DB_USER, required"`
	DB_PASSWORD string `env:"DB_PASSWORD, required"`
	DB_NAME     string `env:"DB_NAME, required"`

	// Optional, default environment variables.
	PORT      int    `env:"PORT, default=8081"`
	LOG_LEVEL string `env:"LOG_LEVEL, default=INFO"`
}

// Loads the environment variables as an EnvConfig
func LoadEnv() EnvConfig {
	var config EnvConfig
	envFun := func() error { return envconfig.Process(context.Background(), &config) }
	FatalCallErrorSupplier(envFun)
	return config
}
