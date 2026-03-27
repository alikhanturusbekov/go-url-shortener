// Package config provides application configuration loading
package config

import (
	"encoding/json"
	"flag"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/logger"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"os"
	"strings"
)

// Config structure of application configuration
type Config struct {
	Address          string `env:"SERVER_ADDRESS" json:"server_address"`
	BaseURL          string `env:"BASE_URL" json:"base_url"`
	LogLevel         string `env:"LOG_LEVEL" json:"log_level"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN      string `env:"DATABASE_DSN" json:"database_dsn"`
	AuthorizationKey string `env:"AUTHORIZATION_KEY" json:"authorization_key"`
	AuditFile        string `env:"AUDIT_FILE" json:"audit_file"`
	AuditURL         string `env:"AUDIT_URL" json:"audit_url"`
	EnableHTTPS      bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	HTTPSCertFile    string `env:"HTTPS_CERT_FILE" json:"https_cert_file"`
	HTTPSKeyFile     string `env:"HTTPS_KEY_FILE" json:"https_key_file"`
}

// NewConfig loads configuration from defaults, environment variables and flags
func NewConfig() (*Config, error) {
	config := Config{
		Address:          ":8080",
		BaseURL:          "http://localhost:8080",
		LogLevel:         "info",
		FileStoragePath:  "",
		DatabaseDSN:      "",
		AuthorizationKey: "secret_auth_key",
		AuditFile:        "",
		AuditURL:         "",
		EnableHTTPS:      false,
		HTTPSCertFile:    "certs/server.crt",
		HTTPSKeyFile:     "certs/server.key",
	}
	var configPath string

	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	flag.StringVar(&config.Address, "a", config.Address, "HTTP server start address")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "The base URL of shortened url")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "The file path for url pairs storage")
	flag.StringVar(&config.DatabaseDSN, "d", config.DatabaseDSN, "Database connection string")
	flag.StringVar(&config.AuthorizationKey, "ak", config.AuthorizationKey, "Authorization Key")
	flag.StringVar(&config.AuditFile, "audit-file", config.AuditFile, "Path to audit log file")
	flag.StringVar(&config.AuditURL, "audit-url", config.AuditURL, "Remote audit server URL")
	flag.BoolVar(&config.EnableHTTPS, "s", config.EnableHTTPS, "Enable HTTPS")
	flag.StringVar(&config.HTTPSCertFile, "https-cert", config.HTTPSCertFile, "Path to TLS certificate")
	flag.StringVar(&config.HTTPSKeyFile, "https-key", config.HTTPSKeyFile, "Path to TLS private key")
	flag.StringVar(&configPath, "c", os.Getenv("CONFIG"), "Path to config file")
	flag.Parse()

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			logger.Log.Warn("failed to open config file", zap.Error(err))
		} else {
			defer func() {
				if err := file.Close(); err != nil {
					logger.Log.Warn("failed to close config file", zap.Error(err))
				}
			}()

			decoder := json.NewDecoder(file)
			if err = decoder.Decode(&config); err != nil {
				logger.Log.Error("error while parsing config file: " + err.Error())
			}
		}
	}

	if config.EnableHTTPS && strings.HasPrefix(config.BaseURL, "http://") {
		config.BaseURL = "https://" + strings.TrimPrefix(config.BaseURL, "http://")
	}

	return &config, nil
}
