package repository

import (
	"context"
	"encoding/json"
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"io"
	"os"
	"slices"
	"sync"
)

type URLFileRepository struct {
	filePath string
	data     []*model.URLPair
	mu       sync.RWMutex
}

func NewURLFileRepository(filePath string) (*URLFileRepository, error) {
	repo := &URLFileRepository{
		filePath: filePath,
		data:     make([]*model.URLPair, 0),
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

	r.data = append(r.data, urlPair)

	return nil
}

func (r *URLFileRepository) GetByShort(_ context.Context, short string) (*model.URLPair, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, urlPair := range r.data {
		if urlPair.Short == short {
			return urlPair, true
		}
	}

	return nil, false
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

func (r *URLFileRepository) GetAllByUserID(_ context.Context, userID string) ([]*model.URLPair, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.URLPair

	for _, urlPair := range r.data {
		if urlPair.UserID == userID && !urlPair.IsDeleted {
			result = append(result, urlPair)
		}
	}

	return result, nil
}

func (r *URLFileRepository) DeleteByIDs(ctx context.Context, userID string, ids []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, urlPair := range r.data {
		if urlPair.UserID == userID && !urlPair.IsDeleted && slices.Contains(ids, urlPair.ID) {
			urlPair.IsDeleted = true
		}
	}

	err := r.syncFile(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *URLFileRepository) syncFile(ctx context.Context) error {
	err := os.Truncate(r.filePath, 0)
	if err != nil {
		return err
	}

	err = r.SaveMany(ctx, r.data)
	if err != nil {
		return err
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

	r.data = append(r.data, urlPairs...)

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
