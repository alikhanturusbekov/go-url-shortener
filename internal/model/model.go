// Package model defines application data structures
package model

import (
	"github.com/google/uuid"
)

// Request represents a shorten URL request
type Request struct {
	URL string `json:"url"`
}

// Response represents a shorten URL response
type Response struct {
	Result string `json:"result"`
}

// URLPair represents a stored URL entity
type URLPair struct {
	ID        string `json:"uid"`
	Short     string `json:"short"`
	Long      string `json:"long"`
	UserID    string `json:"user_id"`
	IsDeleted bool   `json:"is_deleted"`
}

// BatchShortenURLRequest represents a single batch shorten request item
type BatchShortenURLRequest struct {
	CorrelationID *string `json:"correlation_id"`
	OriginalURL   string  `json:"original_url"`
}

// BatchShortenURLResponse represents a single batch shorten response item
type BatchShortenURLResponse struct {
	CorrelationID *string `json:"correlation_id"`
	ShortURL      string  `json:"short_url"`
}

// URLPairsResponse represents a user URL listing response
type URLPairsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// DeleteURLTask represents a background deletion task
type DeleteURLTask struct {
	UserID string `json:"user_id"`
	Short  string `json:"short"`
}

// NewURLPair creates a new URLPair instance
// If id is not provided, it is generated UUID
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
