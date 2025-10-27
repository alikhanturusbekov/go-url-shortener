package service

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"

	"github.com/alikhanturusbekov/go-url-shortener/internal/repository"
)

type UrlService struct {
	repo repository.UrlRepository
}

func NewUrlService(repo repository.UrlRepository) *UrlService {
	return &UrlService{repo: repo}
}

// ShortenUrl Хэширует url и возвращает первые 7 символов
// Есть вероятность, что будут одинаковые 7 символов у разных url, но так как проект маленький - закрыл глаза)
func (s *UrlService) ShortenUrl(url string) (string, error) {
	hash := sha1.Sum([]byte(url))
	shortUrl := base64.URLEncoding.EncodeToString(hash[:])[:7]

	err := s.repo.Save(shortUrl, url)
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}

func (s *UrlService) ResolveShortUrl(shortUrl string) (string, error) {
	url, isFound := s.repo.GetByShort(shortUrl)

	if isFound {
		return url, nil
	}

	return "", errors.New("could not resolve provided url")
}
