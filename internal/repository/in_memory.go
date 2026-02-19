package repository

import (
	"context"
	"slices"
	"sync"

	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
)

// URLInMemoryRepository implements URLRepository using in-memory storage
type URLInMemoryRepository struct {
	data []*model.URLPair
	mu   sync.RWMutex
}

// NewURLInMemoryRepository creates a new URLInMemoryRepository instance
func NewURLInMemoryRepository() *URLInMemoryRepository {
	return &URLInMemoryRepository{data: make([]*model.URLPair, 0)}
}

// Save stores a single URL pair
func (r *URLInMemoryRepository) Save(_ context.Context, urlPair *model.URLPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = append(r.data, urlPair)
	return nil
}

// GetByShort retrieves a URL pair by its short URL
func (r *URLInMemoryRepository) GetByShort(_ context.Context, short string) (*model.URLPair, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, urlPair := range r.data {
		if urlPair.Short == short {
			return urlPair, true
		}
	}

	return nil, false
}

// SaveMany stores multiple URL pairs
func (r *URLInMemoryRepository) SaveMany(_ context.Context, urlPairs []*model.URLPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = append(r.data, urlPairs...)
	return nil
}

// DeleteByShorts marks URL pairs as deleted for a user
func (r *URLInMemoryRepository) DeleteByShorts(_ context.Context, userID string, shorts []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, urlPair := range r.data {
		if !urlPair.IsDeleted && urlPair.UserID == userID && slices.Contains(shorts, urlPair.Short) {
			urlPair.IsDeleted = true
		}
	}

	return nil
}

// GetAllByUserID returns all URL pairs for a user
func (r *URLInMemoryRepository) GetAllByUserID(_ context.Context, userID string) ([]*model.URLPair, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.URLPair

	for _, urlPair := range r.data {
		if urlPair.UserID == userID {
			result = append(result, urlPair)
		}
	}

	return result, nil
}
