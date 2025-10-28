package service

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"

	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
)

type URLService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

// ShortenURL Хэширует url и возвращает первые 7 символов
// Есть вероятность, что будут одинаковые 7 символов у разных url, но так как проект маленький - закрыл глаза)
func (s *URLService) ShortenURL(url string) (string, error) {
	hash := sha1.Sum([]byte(url))
	shortURL := base64.URLEncoding.EncodeToString(hash[:])[:7]

	err := s.repo.Save(shortURL, url)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *URLService) ResolveShortURL(shortURL string) (string, error) {
	url, isFound := s.repo.GetByShort(shortURL)

	if isFound {
		return url, nil
	}

	return "", errors.New("could not resolve provided url")
}
