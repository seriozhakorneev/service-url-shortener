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

// Count -.
func (r *ShortenerRepo) Count(ctx context.Context) (value int, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`SELECT value FROM count WHERE id=true`,
	).Scan(&value)
	if err != nil {
		err = fmt.Errorf("ShortenerRepo - Count - r.Pool.QueryRow.Scan: %w", err)
		return
	}

	return
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
			"rows not affected by sql execution")
		return
	}

	return
}

// Rewrite -.
func (r *ShortenerRepo) Rewrite(ctx context.Context, URL string) (id int, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`
			UPDATE urls
			SET touched = $1, original = $2
            WHERE touched =
            	(
            		SELECT MIN(touched)
            		FROM urls
            	)
            RETURNING id;`,
		time.Now(), URL,
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("ShortenerRepo - Rewrite - r.Pool.QueryRow.Scan: %w", err)
		return
	}
	return
}

// Create -.
func (r *ShortenerRepo) Create(ctx context.Context, URL string) (id int, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`INSERT INTO urls (original,touched) VALUES($1, $2) RETURNING id;`,
		URL, time.Now(),
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("ShortenerRepo - Create - r.Pool.QueryRow.Scan: %w", err)
		return
	}

	err = r.updCount(ctx)
	if err != nil {
		err = fmt.Errorf("ShortenerRepo - Create - %w", err)
		return
	}

	return
}

func (r *ShortenerRepo) updCount(ctx context.Context) (err error) {
	pg, err := r.Pool.Exec(
		ctx,
		`UPDATE count SET value = value + 1 WHERE id = true;`,
	)
	if err != nil {
		err = fmt.Errorf("updCount - r.Pool.Exec: %w", err)
		return
	}

	if pg.RowsAffected() <= 0 {
		err = fmt.Errorf("updCount - pg.RowsAffected: %s",
			"rows not affected by sql execution")
		return
	}

	return
}
