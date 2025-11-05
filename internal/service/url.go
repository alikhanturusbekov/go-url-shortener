package service

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
	appError "github.com/alikhanturusbekov/go-url-shortener/pkg/error"
)

type URLService struct {
	repo    repository.URLRepository
	baseURL string
}

func NewURLService(repo repository.URLRepository, baseURL string) *URLService {
	return &URLService{
		repo:    repo,
		baseURL: baseURL,
	}
}

// ShortenURL Хэширует url и возвращает первые 7 символов
// Есть вероятность, что будут одинаковые 7 символов у разных url, но так как проект маленький - закрыл глаза)
func (s *URLService) ShortenURL(url string) (string, *appError.HTTPError) {
	validatedURL, err := s.validateURL(url)
	if err != nil {
		return "", appError.NewHTTPError(http.StatusBadRequest, "Invalid URL was provided", err)
	}

	hash := sha1.Sum([]byte(validatedURL))
	urlPath := base64.URLEncoding.EncodeToString(hash[:])[:7]

	err = s.repo.Save(urlPath, url)
	if err != nil {
		return "", appError.NewHTTPError(http.StatusInternalServerError, "Failed to save URL", err)
	}

	shortURL := fmt.Sprintf("%s/%s", s.baseURL, urlPath)

	return shortURL, nil
}

func (s *URLService) ResolveShortURL(shortURL string) (string, *appError.HTTPError) {
	originalURL, isFound := s.repo.GetByShort(shortURL)

	if isFound {
		return originalURL, nil
	}

	return "", appError.NewHTTPError(
		http.StatusNotFound,
		"Could not resolve provided URL",
		errors.New("url not found"),
	)
}

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
