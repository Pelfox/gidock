package internal

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	DatabaseURL string `envconfig:"database_url"`
}

func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load()

	var config AppConfig
	if err := envconfig.Process("gidock", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
