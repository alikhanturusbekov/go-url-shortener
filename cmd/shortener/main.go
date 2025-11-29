package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	"github.com/alikhanturusbekov/go-url-shortener/internal/config"
	"github.com/alikhanturusbekov/go-url-shortener/internal/handler"
	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
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

	urlRepo, err := repository.NewURLFileRepository(appConfig.FileStoragePath)
	if err != nil {
		return err
	}
	urlService := service.NewURLService(urlRepo, appConfig.BaseURL)
	urlHandler := handler.NewURLHandler(urlService, database)

	r := chi.NewRouter()

	r.Use(logger.RequestLogger())
	r.Use(compress.GzipCompressor())

	r.Get("/ping", urlHandler.Ping)
	r.Get(`/{id}`, urlHandler.ResolveURL)
	r.With(middleware.AllowContentType("text/plain")).
		Post(`/`, urlHandler.ShortenURLAsText)
	r.With(middleware.AllowContentType("application/json")).
		Post(`/api/shorten`, urlHandler.ShortenURLAsJSON)

	logger.Log.Info("running server...", zap.String("address", appConfig.Address))
	return http.ListenAndServe(appConfig.Address, r)
}
