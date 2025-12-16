package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type URLDatabaseRepository struct {
	db *sql.DB
}

func NewURLDatabaseRepository(db *sql.DB) *URLDatabaseRepository {
	return &URLDatabaseRepository{db: db}
}

func (r *URLDatabaseRepository) Save(ctx context.Context, urlPair *model.URLPair) error {
	query := `
        INSERT INTO url_pairs (uid, short, long)
        VALUES ($1, $2, $3)
    `
	_, err := r.db.ExecContext(ctx, query, urlPair.ID, urlPair.Short, urlPair.Long)

	if err != nil {
		var pgErr *pgconn.PgError

		if ok := errors.As(err, &pgErr); !ok {
			return err
		}

		if pgErr.Code == pgerrcode.UniqueViolation {
			return ErrorOnConflict
		}
	}

	return err
}

func (r *URLDatabaseRepository) GetByShort(ctx context.Context, short string) (*model.URLPair, bool) {
	var result model.URLPair

	query := `
        SELECT short, long
        FROM url_pairs
        WHERE short = $1;
    `

	err := r.db.QueryRowContext(ctx, query, short).Scan(&result.Short, &result.Long)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, false
	}

	return &result, err == nil
}

func (r *URLDatabaseRepository) SaveMany(ctx context.Context, urlPairs []*model.URLPair) error {
	tx, err := r.db.Begin()
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO url_pairs (uid, short, long) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, urlPair := range urlPairs {
		_, fail := stmt.Exec(urlPair.ID, urlPair.Short, urlPair.Long)
		if fail != nil {
			return fail
		}
	}

	return tx.Commit()
}

func (r *URLDatabaseRepository) Close() error {
	return r.db.Close()
}
