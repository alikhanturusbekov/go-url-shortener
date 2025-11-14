package repository

import "sync"

type URLRepository interface {
	Save(short string, long string) error
	GetByShort(short string) (string, bool)
}

type URLMapRepository struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewURLMapRepository() *URLMapRepository {
	return &URLMapRepository{
		data: make(map[string]string),
		mu:   sync.RWMutex{},
	}
}

func (r *URLMapRepository) Save(short string, long string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[short] = long
	return nil
}

func (r *URLMapRepository) GetByShort(short string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	val, ok := r.data[short]
	return val, ok
}
