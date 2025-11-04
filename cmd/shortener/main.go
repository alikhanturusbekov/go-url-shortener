package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"

	"github.com/alikhanturusbekov/go-url-shortener/internal/config"
	"github.com/alikhanturusbekov/go-url-shortener/internal/handler"
	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
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

	urlRepo := repository.NewURLMapRepository()
	urlService := service.NewURLService(urlRepo, appConfig.BaseURL)
	urlHandler := handler.NewURLHandler(urlService)

	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("text/plain"))

	r.Post(`/`, urlHandler.ShortenURL)
	r.Get(`/{id}`, urlHandler.ResolveURL)

	return http.ListenAndServe(appConfig.Address, r)
}
