package repository

import "sync"

type UrlRepository interface {
	Save(short string, long string) error
	GetByShort(short string) (string, bool)
}

type UrlMapRepository struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewUrlMapRepository() *UrlMapRepository {
	return &UrlMapRepository{
		data: make(map[string]string),
		mu:   sync.RWMutex{},
	}
}

func (r *UrlMapRepository) Save(short string, long string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[short] = long
	return nil
}

func (r *UrlMapRepository) GetByShort(short string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	val, ok := r.data[short]
	return val, ok
}
