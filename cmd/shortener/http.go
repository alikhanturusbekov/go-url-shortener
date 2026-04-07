package main

import (
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/alikhanturusbekov/go-url-shortener/internal/certs"
	"github.com/alikhanturusbekov/go-url-shortener/internal/config"
	"github.com/alikhanturusbekov/go-url-shortener/internal/handler"
	appmiddleware "github.com/alikhanturusbekov/go-url-shortener/internal/middleware"
	"github.com/alikhanturusbekov/go-url-shortener/internal/service"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/authorization"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/compress"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/logger"
)

// setupHTTPServer prepares everything to start HTTP/HTTPS server
func setupHTTPServer(appConfig *config.Config, urlService *service.URLService, database *sql.DB) *http.Server {
	urlHandler := handler.NewURLHandler(urlService, database)

	r := chi.NewRouter()

	r.Mount("/debug", chimiddleware.Profiler())

	r.Group(func(r chi.Router) {
		r.Use(logger.RequestLogger())
		r.Use(compress.GzipCompressor())
		r.Use(authorization.AuthMiddleware([]byte(appConfig.AuthorizationKey)))

		r.Get("/ping", urlHandler.Ping)
		r.Get(`/{id}`, urlHandler.ResolveURL)
		r.With(chimiddleware.AllowContentType("text/plain")).
			Post(`/`, urlHandler.ShortenURLAsText)
		r.With(chimiddleware.AllowContentType("application/json")).
			Post(`/api/shorten`, urlHandler.ShortenURLAsJSON)
		r.With(chimiddleware.AllowContentType("application/json")).
			Post(`/api/shorten/batch`, urlHandler.BatchShortenURL)

		r.Get(`/api/user/urls`, urlHandler.GetUserURLs)
		r.With(chimiddleware.AllowContentType("application/json")).
			Delete(`/api/user/urls`, urlHandler.DeleteUserURLs)

		r.Route("/api/internal", func(r chi.Router) {
			r.Use(appmiddleware.TrustedSubnet(appConfig.TrustedSubnet))
			r.Get("/stats", urlHandler.GetStats)
		})
	})

	return &http.Server{
		Addr:              appConfig.Address,
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}

// startHTTPServer starts HTTP/HTTPS server in a goroutine
func startHTTPServer(appConfig *config.Config, srv *http.Server) <-chan error {
	serverErr := make(chan error, 1)

	go func() {
		var err error

		if appConfig.EnableHTTPS {
			if err = certs.EnsureCertificates(appConfig.HTTPSCertFile, appConfig.HTTPSKeyFile); err != nil {
				serverErr <- fmt.Errorf("ensure certificates: %w", err)
				return
			}

			logger.Log.Info("starting HTTPS server",
				zap.String("address", appConfig.Address),
				zap.String("cert_file", appConfig.HTTPSCertFile),
				zap.String("key_file", appConfig.HTTPSKeyFile),
			)

			err = srv.ListenAndServeTLS(
				appConfig.HTTPSCertFile,
				appConfig.HTTPSKeyFile,
			)
		} else {
			logger.Log.Info("running server...", zap.String("address", appConfig.Address))
			err = srv.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}

		serverErr <- nil
	}()

	return serverErr
}
