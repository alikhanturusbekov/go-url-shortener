package repository

import (
	"context"
	"errors"

	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
)

var ErrorOnConflict = errors.New("conflict")

type URLRepository interface {
	Save(ctx context.Context, urlPair *model.URLPair) error
	GetByShort(ctx context.Context, short string) (*model.URLPair, bool)
	SaveMany(ctx context.Context, urlPairs []*model.URLPair) error
	GetAllByUserID(ctx context.Context, userID string) ([]*model.URLPair, error)
	DeleteByIDs(ctx context.Context, userID string, ids []string) error
}
