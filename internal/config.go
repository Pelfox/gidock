package internal

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// AppConfig holds the global application configuration.
type AppConfig struct {
	// DatabaseURL is the connection string (DSN) for the PostgreSQL database.
	DatabaseURL string `envconfig:"database_url"`
}

// LoadConfig loads the application configuration from environment variables.
// It first attempts to load a .env file (if present) using godotenv. Then, it
// processes environment variables with the prefix "GIDOCK_" and maps them to
// the AppConfig.
func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load()

	var config AppConfig
	if err := envconfig.Process("gidock", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
