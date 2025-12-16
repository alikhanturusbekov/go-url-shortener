package repository

import (
	"context"
	"encoding/json"
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"io"
	"os"
	"sync"
)

type URLFileRepository struct {
	filePath string
	data     map[string]*model.URLPair
	mu       sync.RWMutex
}

func NewURLFileRepository(filePath string) (*URLFileRepository, error) {
	repo := &URLFileRepository{
		filePath: filePath,
		data:     make(map[string]*model.URLPair),
	}

	if _, err := os.Stat(filePath); err == nil {
		if err := repo.load(); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

func (r *URLFileRepository) Save(_ context.Context, urlPair *model.URLPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	itemJSON, err := json.Marshal(urlPair)
	if err != nil {
		return err
	}

	err = r.addRecord(file, itemJSON, stat.Size() == 0)
	if err != nil {
		return err
	}

	r.data[urlPair.Short] = urlPair

	return nil
}

func (r *URLFileRepository) GetByShort(_ context.Context, short string) (*model.URLPair, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	urlPair, ok := r.data[short]
	return urlPair, ok
}

func (r *URLFileRepository) SaveMany(_ context.Context, urlPairs []*model.URLPair) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	isFirstRecord := stat.Size() == 0

	for _, urlPair := range urlPairs {
		urlJSON, err := json.Marshal(urlPair)
		if err != nil {
			return err
		}

		if err := r.addRecord(file, urlJSON, isFirstRecord); err != nil {
			return err
		}

		isFirstRecord = false
	}

	return nil
}

func (r *URLFileRepository) load() error {
	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var urlPairs []*model.URLPair
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&urlPairs); err != nil {
		return nil
	}

	for _, urlPair := range urlPairs {
		r.data[urlPair.Short] = urlPair
	}
	return nil
}

func (r *URLFileRepository) addRecord(file *os.File, itemJSON []byte, isFirst bool) error {
	var err error
	var firstAppendItem string

	if isFirst {
		firstAppendItem = "[\n"
	} else {
		_, err = file.Seek(-1, io.SeekEnd)
		if err != nil {
			return err
		}

		firstAppendItem = ",\n"
	}

	_, err = file.Write([]byte(firstAppendItem))
	if err != nil {
		return err
	}
	_, err = file.Write(itemJSON)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte("\n]"))
	if err != nil {
		return err
	}

	return nil
}
