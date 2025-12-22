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
	ID     string `json:"uid"`
	Short  string `json:"short"`
	Long   string `json:"long"`
	UserId string `json:"user_id"`
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

func NewURLPair(short, long string, id *string, userId string) *URLPair {
	urlPair := &URLPair{
		ID:     uuid.NewString(),
		Short:  short,
		Long:   long,
		UserId: userId,
	}

	if id != nil {
		urlPair.ID = *id
	}

	return urlPair
}
