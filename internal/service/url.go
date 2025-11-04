package service

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
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
func (s *URLService) ShortenURL(url string) (string, error) {
	validatedURL, isValid := validateURL(url)
	if !isValid {
		return "", errors.New("invalid URL")
	}

	hash := sha1.Sum([]byte(validatedURL))
	urlPath := base64.URLEncoding.EncodeToString(hash[:])[:7]

	err := s.repo.Save(urlPath, url)
	if err != nil {
		return "", err
	}

	shortURL := fmt.Sprintf("%s/%s", s.baseURL, urlPath)

	return shortURL, nil
}

func (s *URLService) ResolveShortURL(shortURL string) (string, error) {
	url, isFound := s.repo.GetByShort(shortURL)

	if isFound {
		return url, nil
	}

	return "", errors.New("could not resolve provided url")
}

func validateURL(originalURL string) (string, bool) {
	trimmedURL := strings.TrimSpace(originalURL)
	if trimmedURL == "" {
		return "", false
	}

	resultURL, err := url.ParseRequestURI(trimmedURL)
	if err != nil {
		return "", false
	}

	if resultURL.Scheme == "" || resultURL.Host == "" {
		return "", false
	}

	return resultURL.String(), true
}
