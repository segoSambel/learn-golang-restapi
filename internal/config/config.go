package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
	URL string `env:"DATABASE_URL,required,notEmpty"`
}

type HTTPConfig struct {
	Host                   string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port                   int    `env:"HTTP_PORT" envDefault:"8080"`
	ShutdownTimeoutSeconds int    `env:"HTTP_SHUTDOWN_TIMEOUT_SECONDS" envDefault:"10"`
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	return env.ParseAs[Config]()
}
