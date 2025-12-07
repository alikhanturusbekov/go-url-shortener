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
	ID    string `json:"uid"`
	Short string `json:"short"`
	Long  string `json:"long"`
}

type BatchShortenURLRequest struct {
	CorrelationID *string `json:"correlation_id"`
	OriginalURL   string  `json:"original_url"`
}

type BatchShortenURLResponse struct {
	CorrelationID *string `json:"correlation_id"`
	ShortURL      string  `json:"short_url"`
}

func NewURLPair(short, long string, id *string) *URLPair {
	urlPair := &URLPair{
		ID:    uuid.NewString(),
		Short: short,
		Long:  long,
	}

	if id != nil {
		urlPair.ID = *id
	}

	return urlPair
}
