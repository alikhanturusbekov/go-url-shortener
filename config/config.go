package config

import "flag"

type Config struct {
	Address string
	BaseURL string
}

func NewConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.Address, "a", ":8080", "HTTP server start address")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "The base URL of shortened url")

	flag.Parse()

	return config
}
