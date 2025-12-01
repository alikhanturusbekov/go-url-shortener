package repository

import (
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"sync"
)

type URLInMemoryRepository struct {
	data map[string]model.URLPair
	mu   sync.RWMutex
}

func NewURLInMemoryRepository() *URLInMemoryRepository {
	return &URLInMemoryRepository{data: make(map[string]model.URLPair)}
}

func (r *URLInMemoryRepository) Save(urlPair model.URLPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[urlPair.Short] = urlPair
	return nil
}

func (r *URLInMemoryRepository) GetByShort(short string) (model.URLPair, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	urlPair, ok := r.data[short]
	return urlPair, ok
}
