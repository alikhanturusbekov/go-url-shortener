package repository

import (
	"database/sql"
	"errors"
	"github.com/alikhanturusbekov/go-url-shortener/internal/model"
)

type URLDatabaseRepository struct {
	db *sql.DB
}

func NewURLDatabaseRepository(db *sql.DB) *URLDatabaseRepository {
	return &URLDatabaseRepository{db: db}
}

func (r *URLDatabaseRepository) Save(urlPair *model.URLPair) error {
	query := `
        INSERT INTO url_pairs (uid, short, long)
        VALUES ($1, $2, $3)
        ON CONFLICT (short) DO UPDATE SET long = EXCLUDED.long;
    `
	_, err := r.db.Exec(query, urlPair.ID, urlPair.Short, urlPair.Long)
	return err
}

func (r *URLDatabaseRepository) GetByShort(short string) (*model.URLPair, bool) {
	var result model.URLPair

	query := `
        SELECT short, long
        FROM url_pairs
        WHERE short = $1;
    `

	err := r.db.QueryRow(query, short).Scan(&result.Short, &result.Long)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, false
	}

	return &result, err == nil
}

func (r *URLDatabaseRepository) SaveMany(urlPairs []*model.URLPair) error {
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO url_pairs (uid, short, long) VALUES ($1, $2, $3)")
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
