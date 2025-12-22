package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	Address          string `env:"SERVER_ADDRESS"`
	BaseURL          string `env:"BASE_URL"`
	LogLevel         string `env:"LOG_LEVEL"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN      string `env:"DATABASE_DSN"`
	AuthorizationKey string `env:"AUTHORIZATION_KEY"`
}

func NewConfig() *Config {
	config := Config{
		Address:          ":8080",
		BaseURL:          "http://localhost:8080",
		LogLevel:         "info",
		FileStoragePath:  "",
		DatabaseDSN:      "",
		AuthorizationKey: "secret_auth_key",
	}

	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&config.Address, "a", config.Address, "HTTP server start address")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "The base URL of shortened url")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "The file path for url pairs storage")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "Database connection string")
	flag.StringVar(&config.AuthorizationKey, "ak", config.AuthorizationKey, "Authorization Key")
	flag.Parse()

	return &config
}
