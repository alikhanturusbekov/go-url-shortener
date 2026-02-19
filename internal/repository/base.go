// Package repository defines storage abstractions for URL entities.
package repository

import (
	"context"
	"errors"

	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
)

// ErrorOnConflict is returned when a save operation conflicts with existing data
var ErrorOnConflict = errors.New("conflict")

// URLRepository defines persistence methods for URL pairs
type URLRepository interface {
	// Save stores a single URL pair
	Save(ctx context.Context, urlPair *model.URLPair) error

	// GetByShort retrieves a URL pair by its short URL
	GetByShort(ctx context.Context, short string) (*model.URLPair, bool)

	// SaveMany stores multiple URL pairs
	SaveMany(ctx context.Context, urlPairs []*model.URLPair) error

	// GetAllByUserID returns all URL pairs for a user
	GetAllByUserID(ctx context.Context, userID string) ([]*model.URLPair, error)

	// DeleteByShorts marks URL pairs as deleted for a user
	DeleteByShorts(ctx context.Context, userID string, shorts []string) error
}
