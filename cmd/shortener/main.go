package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/alikhanturusbekov/go-url-shortener/config"
	"github.com/alikhanturusbekov/go-url-shortener/internal/handler"
	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
)

func main() {
	if err := run(); err != nil {
		panic(err)
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
