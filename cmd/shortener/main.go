package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alikhanturusbekov/go-url-shortener/internal/config"
	"github.com/alikhanturusbekov/go-url-shortener/internal/handler"
	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/authorization"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/compress"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/logger"
)

func main() {
	log.Print("Starting the app...")

	if err := run(); err != nil {
		log.Printf("Error while starting the app: %s", err)

		os.Exit(1)
	}
}

func run() error {
	appConfig := config.NewConfig()

	if err := logger.Initialize(appConfig.LogLevel); err != nil {
		return err
	}

	database, err := sql.Open("pgx", appConfig.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	if appConfig.DatabaseDSN != "" {
		if err := applyMigrations(database, "migrations"); err != nil {
			log.Fatalf("failed to apply migrations: %v", err)
		}
	}

	urlRepo, cleanUp, err := setupRepository(appConfig)
	if err != nil {
		return err
	}
	defer cleanUp()

	urlService := service.NewURLService(urlRepo, appConfig.BaseURL)
	urlHandler := handler.NewURLHandler(urlService, database)

	r := chi.NewRouter()

	r.Use(logger.RequestLogger())
	r.Use(compress.GzipCompressor())
	r.Use(authorization.AuthMiddleware([]byte(appConfig.AuthorizationKey)))

	r.Get("/ping", urlHandler.Ping)
	r.Get(`/{id}`, urlHandler.ResolveURL)
	r.With(middleware.AllowContentType("text/plain")).
		Post(`/`, urlHandler.ShortenURLAsText)
	r.With(middleware.AllowContentType("application/json")).
		Post(`/api/shorten`, urlHandler.ShortenURLAsJSON)
	r.With(middleware.AllowContentType("application/json")).
		Post(`/api/shorten/batch`, urlHandler.BatchShortenURL)

	r.Get(`/api/user/urls`, urlHandler.GetUserURLs)
	r.With(middleware.AllowContentType("application/json")).
		Delete(`/api/user/urls`, urlHandler.DeleteUserURLs)

	logger.Log.Info("running server...", zap.String("address", appConfig.Address))
	return http.ListenAndServe(appConfig.Address, r)
}

func setupRepository(config *config.Config) (repository.URLRepository, func(), error) {
	if config.DatabaseDSN != "" {
		logger.Log.Info("Using the database for storage...")

		database, err := sql.Open("pgx", config.DatabaseDSN)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}

		databaseRepo := repository.NewURLDatabaseRepository(database)
		cleanup := func() { databaseRepo.Close() }

		return databaseRepo, cleanup, nil
	}

	if config.FileStoragePath != "" {
		logger.Log.Info("Using the file system for storage...")

		fileRepo, err := repository.NewURLFileRepository(config.FileStoragePath)

		return fileRepo, func() {}, err
	}

	logger.Log.Info("Using the in-memory repository...")

	return repository.NewURLInMemoryRepository(), func() {}, nil
}

func applyMigrations(db *sql.DB, dir string) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY
        );
    `)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		return err
	}

	for _, file := range files {
		version := filepath.Base(file)

		var exists string
		err := db.QueryRow("SELECT version FROM schema_migrations WHERE version=$1", version).Scan(&exists)
		if err == nil {
			continue
		} else if err != sql.ErrNoRows {
			return fmt.Errorf("failed to check migration %s: %w", version, err)
		}

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}

		_, err = db.Exec("INSERT INTO schema_migrations(version) VALUES ($1)", version)
		if err != nil {
			return fmt.Errorf("failed to record applied migration %s: %w", version, err)
		}

		logger.Log.Info("Applied migration: " + version)
	}

	return nil
}
