package repository

import (
	"context"
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"sync"
)

type URLInMemoryRepository struct {
	data map[string]*model.URLPair
	mu   sync.RWMutex
}

func NewURLInMemoryRepository() *URLInMemoryRepository {
	return &URLInMemoryRepository{data: make(map[string]*model.URLPair)}
}

func (r *URLInMemoryRepository) Save(_ context.Context, urlPair *model.URLPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[urlPair.Short] = urlPair
	return nil
}

func (r *URLInMemoryRepository) GetByShort(_ context.Context, short string) (*model.URLPair, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	urlPair, ok := r.data[short]
	return urlPair, ok
}

func (r *URLInMemoryRepository) SaveMany(_ context.Context, urlPairs []*model.URLPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, urlPair := range urlPairs {
		r.data[urlPair.Short] = urlPair
	}

	return nil
}
