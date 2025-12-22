package repository

import (
	"context"
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"slices"
	"sync"
)

type URLInMemoryRepository struct {
	data []*model.URLPair
	mu   sync.RWMutex
}

func NewURLInMemoryRepository() *URLInMemoryRepository {
	return &URLInMemoryRepository{data: make([]*model.URLPair, 0)}
}

func (r *URLInMemoryRepository) Save(_ context.Context, urlPair *model.URLPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = append(r.data, urlPair)
	return nil
}

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

func (r *URLInMemoryRepository) SaveMany(_ context.Context, urlPairs []*model.URLPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, urlPair := range urlPairs {
		r.data = append(r.data, urlPair)
	}

	return nil
}

func (r *URLInMemoryRepository) DeleteByIDs(_ context.Context, userID string, ids []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, urlPair := range r.data {
		if !urlPair.IsDeleted && urlPair.UserID == userID && slices.Contains(ids, urlPair.ID) {
			urlPair.IsDeleted = true
		}
	}

	return nil
}

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
