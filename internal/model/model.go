package model

import (
	"github.com/google/uuid"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type URLPair struct {
	ID        string `json:"uid"`
	Short     string `json:"short"`
	Long      string `json:"long"`
	UserID    string `json:"user_id"`
	IsDeleted bool   `json:"is_deleted"`
}

type BatchShortenURLRequest struct {
	CorrelationID *string `json:"correlation_id"`
	OriginalURL   string  `json:"original_url"`
}

type BatchShortenURLResponse struct {
	CorrelationID *string `json:"correlation_id"`
	ShortURL      string  `json:"short_url"`
}

type URLPairsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type DeleteURLTask struct {
	UserID string `json:"user_id"`
	Short  string `json:"short"`
}

func NewURLPair(short, long string, id *string, userID string, isDeleted bool) *URLPair {
	urlPair := &URLPair{
		ID:        uuid.NewString(),
		Short:     short,
		Long:      long,
		UserID:    userID,
		IsDeleted: isDeleted,
	}

	if id != nil {
		urlPair.ID = *id
	}

	return urlPair
}
