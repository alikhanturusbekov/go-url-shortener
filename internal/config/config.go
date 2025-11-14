package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	Address string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
}

func NewConfig() *Config {
	config := Config{
		Address: ":8080",
		BaseURL: "http://localhost:8080",
	}

	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&config.Address, "a", config.Address, "HTTP server start address")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "The base URL of shortened url")
	flag.Parse()

	return &config
}
