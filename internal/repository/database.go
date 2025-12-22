package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alikhanturusbekov/go-url-shortener/internal/model"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

type URLDatabaseRepository struct {
	db *sql.DB
}

func NewURLDatabaseRepository(db *sql.DB) *URLDatabaseRepository {
	return &URLDatabaseRepository{db: db}
}

func (r *URLDatabaseRepository) Save(ctx context.Context, urlPair *model.URLPair) error {
	query := `
        INSERT INTO url_pairs (uid, short, long, user_id, is_deleted)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.ExecContext(ctx, query, urlPair.ID, urlPair.Short, urlPair.Long, urlPair.UserID, urlPair.IsDeleted)

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
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO url_pairs (uid, short, long, user_id, is_deleted) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, urlPair := range urlPairs {
		_, fail := stmt.Exec(urlPair.ID, urlPair.Short, urlPair.Long, urlPair.UserID, urlPair.IsDeleted)
		if fail != nil {
			return fail
		}
	}

	return tx.Commit()
}

func (r *URLDatabaseRepository) GetAllByUserID(ctx context.Context, userID string) ([]*model.URLPair, error) {
	var result []*model.URLPair

	query := `
        SELECT uid, short, long, user_id
        FROM url_pairs
        WHERE user_id = $1 AND is_deleted = false;
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pair model.URLPair
		if rowsErr := rows.Scan(&pair.ID, &pair.Short, &pair.Long, &pair.UserID); rowsErr != nil {
			return nil, rowsErr
		}
		result = append(result, &pair)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, err
	}

	return result, nil
}

func (r *URLDatabaseRepository) DeleteByIDs(ctx context.Context, userID string, ids []string) error {
	query := `
		UPDATE url_pairs
		SET is_deleted = TRUE
		WHERE user_id = $1 AND uid = ANY($2)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		userID,
		pq.Array(ids),
	)

	return err
}

func (r *URLDatabaseRepository) Close() error {
	return r.db.Close()
}
