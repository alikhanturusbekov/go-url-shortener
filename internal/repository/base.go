package repository

import (
	"errors"

	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
)

var ErrorOnConflict = errors.New("conflict")

type URLRepository interface {
	Save(urlPair *model.URLPair) error
	GetByShort(short string) (*model.URLPair, bool)
	SaveMany(urlPairs []*model.URLPair) error
}
