package repository

import (
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
)

type URLRepository interface {
	Save(url model.URLPair) error
	GetByShort(short string) (model.URLPair, bool)
}
