package repository

import (
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
)

type URLRepository interface {
	Save(urlPair *model.URLPair) error
	GetByShort(short string) (*model.URLPair, bool)
	SaveMany(urlPairs []*model.URLPair) error
}
