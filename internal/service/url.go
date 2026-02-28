package service

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	"github.com/alikhanturusbekov/go-url-shortener/internal/worker"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/audit"
	appError "github.com/alikhanturusbekov/go-url-shortener/pkg/error"
)

// URLService provides URL shortening business logic
type URLService struct {
	repo            repository.URLRepository
	baseURL         string
	deleteURLWorker *worker.DeleteURLWorker
	audit           audit.Publisher
}

// NewURLService creates a new URLService instance
func NewURLService(
	repo repository.URLRepository,
	baseURL string,
	deleteURLWorker *worker.DeleteURLWorker,
	auditPublisher audit.Publisher,
) *URLService {
	return &URLService{
		repo:            repo,
		baseURL:         baseURL,
		deleteURLWorker: deleteURLWorker,
		audit:           auditPublisher,
	}
}

// ShortenURL validates and shortens a URL
func (s *URLService) ShortenURL(url string, userID string) (string, *appError.HTTPError) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	validatedURL, err := s.validateURL(url)
	if err != nil {
		return "", appError.NewHTTPError(http.StatusBadRequest, "Invalid URL was provided", err)
	}

	urlPath, err := s.generateShortURLPath(ctx, validatedURL)
	if err != nil {
		return "", appError.NewHTTPError(http.StatusInternalServerError, "Failed to generate short URL", err)
	}

	_, isFound := s.repo.GetByShort(ctx, urlPath)
	if !isFound {
		err = s.repo.Save(ctx, model.NewURLPair(urlPath, validatedURL, nil, userID, false))
		if err != nil && !errors.Is(err, repository.ErrorOnConflict) {
			return "", appError.NewHTTPError(http.StatusInternalServerError, "Failed to save URL", err)
		}
	}

	shortURL := fmt.Sprintf("%s/%s", s.baseURL, urlPath)

	if isFound || (err != nil && errors.Is(err, repository.ErrorOnConflict)) {
		return shortURL, appError.NewHTTPError(http.StatusConflict, "", nil)
	}

	s.audit.Notify(audit.Event{
		TS:     time.Now().Unix(),
		Action: "shorten",
		UserID: userID,
		URL:    validatedURL,
	})

	return shortURL, nil
}

// BatchShortenURL shortens multiple URLs in a single request.
func (s *URLService) BatchShortenURL(items []model.BatchShortenURLRequest, userID string) ([]*model.BatchShortenURLResponse, *appError.HTTPError) {
	results := make([]*model.BatchShortenURLResponse, 0, len(items))
	urlPairs := make([]*model.URLPair, 0, len(items))

	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	for _, item := range items {
		validatedURL, err := s.validateURL(item.OriginalURL)
		if err != nil {
			return nil, appError.NewHTTPError(http.StatusBadRequest, "Invalid URL was provided", err)
		}

		urlPath, err := s.generateShortURLPath(ctx, validatedURL)
		if err != nil {
			return nil, appError.NewHTTPError(http.StatusInternalServerError, "Failed to generate short URL", err)
		}

		urlPairs = append(urlPairs, model.NewURLPair(urlPath, item.OriginalURL, item.CorrelationID, userID, false))
		results = append(results, &model.BatchShortenURLResponse{
			CorrelationID: item.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", s.baseURL, urlPath),
		})
	}

	if err := s.repo.SaveMany(ctx, urlPairs); err != nil {
		return nil, appError.NewHTTPError(http.StatusInternalServerError, "Failed to batch save URL pairs", err)
	}

	return results, nil
}

// ResolveShortURL resolves a short code to the original URL.
func (s *URLService) ResolveShortURL(shortURL string) (string, *appError.HTTPError) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	urlPair, isFound := s.repo.GetByShort(ctx, shortURL)

	if !isFound {
		return "", appError.NewHTTPError(
			http.StatusNotFound,
			"Could not resolve provided URL",
			errors.New("url not found"),
		)
	}

	if urlPair.IsDeleted {
		return "", appError.NewHTTPError(
			http.StatusGone,
			"URL is no longer active",
			errors.New("URL has been deleted by the user"),
		)
	}

	s.audit.Notify(audit.Event{
		TS:     time.Now().Unix(),
		Action: "follow",
		UserID: urlPair.UserID,
		URL:    urlPair.Long,
	})

	return urlPair.Long, nil
}

// GetUserURLs returns all URLs created by a user
func (s *URLService) GetUserURLs(userID string) ([]*model.URLPairsResponse, *appError.HTTPError) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	urlPairs, err := s.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, appError.NewHTTPError(http.StatusInternalServerError, "Failed to get user URLs", err)
	}

	results := make([]*model.URLPairsResponse, 0, len(urlPairs))

	for _, urlPair := range urlPairs {
		results = append(results, &model.URLPairsResponse{
			OriginalURL: urlPair.Long,
			ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, urlPair.Short),
		})
	}

	return results, nil
}

// DeleteUserURLs enqueues URL deletion tasks for the user
func (s *URLService) DeleteUserURLs(userID string, shorts []string) *appError.HTTPError {
	for _, short := range shorts {
		s.deleteURLWorker.Enqueue(model.DeleteURLTask{
			UserID: userID,
			Short:  short,
		})
	}

	return nil
}

// validateURL validates and normalizes a URL string
func (s *URLService) validateURL(originalURL string) (string, error) {
	trimmedURL := strings.TrimSpace(originalURL)
	if trimmedURL == "" {
		return "", errors.New("couldn't parse URL")
	}

	resultURL, err := url.ParseRequestURI(trimmedURL)
	if err != nil {
		return "", err
	}

	if resultURL.Scheme == "" || resultURL.Host == "" {
		return "", errors.New("url's scheme or host is empty")
	}

	return resultURL.String(), nil
}

// generateShortURLPath generates a unique short path
func (s *URLService) generateShortURLPath(ctx context.Context, originalURL string) (string, error) {
	urlPath := s.hashURL(originalURL)

	for {
		urlPair, isFound := s.repo.GetByShort(ctx, urlPath)

		if !isFound || urlPair.Long == originalURL {
			return urlPath, nil
		}

		salted, err := s.addSalt(originalURL)
		if err != nil {
			return "", err
		}

		originalURL = salted
		urlPath = s.hashURL(originalURL)
	}
}

// addSalt appends random data to avoid hash collisions
func (s *URLService) addSalt(url string) (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return url + ":" + base64.RawURLEncoding.EncodeToString(b), nil
}

// hashURL generates a short hash for a URL
func (s *URLService) hashURL(url string) string {
	hash := sha1.Sum([]byte(url))
	return base64.URLEncoding.EncodeToString(hash[:])[:7]
}
