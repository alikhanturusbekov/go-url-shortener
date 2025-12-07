package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
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

	urlPath, err := s.generateShortURLPath(validatedURL)
	if err != nil {
		return "", appError.NewHTTPError(http.StatusInternalServerError, "Failed to generate short URL", err)
	}

	if _, isFound := s.repo.GetByShort(urlPath); !isFound {
		err = s.repo.Save(model.NewURLPair(urlPath, validatedURL, nil))
		if err != nil {
			return "", appError.NewHTTPError(http.StatusInternalServerError, "Failed to save URL", err)
		}
	}

	shortURL := fmt.Sprintf("%s/%s", s.baseURL, urlPath)

	return shortURL, nil
}

func (s *URLService) ResolveShortURL(shortURL string) (string, *appError.HTTPError) {
	urlPair, isFound := s.repo.GetByShort(shortURL)

	if isFound {
		return urlPair.Long, nil
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

func (s *URLService) generateShortURLPath(originalURL string) (string, error) {
	urlPath := s.hashURL(originalURL)

	for {
		urlPair, isFound := s.repo.GetByShort(urlPath)

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

func (s *URLService) addSalt(url string) (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return url + ":" + base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *URLService) hashURL(url string) string {
	hash := sha1.Sum([]byte(url))
	return base64.URLEncoding.EncodeToString(hash[:])[:7]
}
