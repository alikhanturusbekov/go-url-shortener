package main

import (
	"net/http"

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
	urlRepo := repository.NewURLMapRepository()
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, urlHandler.ShortenUrl)
	mux.HandleFunc(`/{id}`, urlHandler.ResolveUrl)

	return http.ListenAndServe(`:8080`, mux)
}
