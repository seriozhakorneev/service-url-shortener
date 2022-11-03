package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	internal "service-url-shortener/internal/errors"
	"service-url-shortener/pkg/postgres"
)

// ShortenerRepo -.
type ShortenerRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *ShortenerRepo {
	return &ShortenerRepo{pg}
}

// GetID -.
func (r *ShortenerRepo) GetID(ctx context.Context, URL string) (int, error) {
	var id int
	err := r.Pool.QueryRow(
		ctx,
		`SELECT id FROM urls WHERE original = $1`,
		URL,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, internal.ErrNotFoundURL
		}
		return 0, fmt.Errorf("ShortenerRepo - GetID - r.Pool.QueryRow.Scan: %w", err)
	}

	return id, nil
}

// Touch -.
func (r *ShortenerRepo) Touch(ctx context.Context, id int) (err error) {
	pg, err := r.Pool.Exec(
		ctx,
		`UPDATE urls SET touched = $1 WHERE id = $2;`,
		time.Now(), id,
	)
	if err != nil {
		err = fmt.Errorf("ShortenerRepo - Touch - r.Pool.Exec: %w", err)
		return
	}

	if pg.RowsAffected() <= 0 {
		err = fmt.Errorf("ShortenerRepo - Touch - pg.RowsAffected: %s",
			"rows not affected by execution")
		return
	}

	return
}

// Insert -.
func (r *ShortenerRepo) Insert(ctx context.Context, URL string) (id int, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`INSERT INTO urls (original,touched) VALUES($1, $2) RETURNING id;`,
		URL, time.Now(),
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("ShortenerRepo - Insert - r.Pool.QueryRow.Scan: %w", err)
		return
	}

	return
}
