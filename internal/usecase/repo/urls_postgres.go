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

// UrlsRepo -.
type UrlsRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *UrlsRepo {
	return &UrlsRepo{pg}
}

// Count -.
func (r *UrlsRepo) Count(ctx context.Context) (value int, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`SELECT value FROM count WHERE id=true`,
	).Scan(&value)
	if err != nil {
		err = fmt.Errorf("UrlsRepo - Count - r.Pool.QueryRow.Scan: %w", err)
		return
	}

	return
}

// GetID -.
func (r *UrlsRepo) GetID(ctx context.Context, URL string) (int, error) {
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
		return 0, fmt.Errorf("UrlsRepo - GetID - r.Pool.QueryRow.Scan: %w", err)
	}

	return id, nil
}

// GetURL -.
func (r *UrlsRepo) GetURL(ctx context.Context, id int) (original string, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`SELECT original FROM urls WHERE id = $1`,
		id,
	).Scan(&original)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = internal.ErrNotFoundURL
			return
		}
		err = fmt.Errorf("UrlsRepo - GetURL - r.Pool.QueryRow.Scan: %w", err)
		return
	}

	return
}

// Touch -.
func (r *UrlsRepo) Touch(ctx context.Context, id int) (err error) {
	pg, err := r.Pool.Exec(
		ctx,
		`UPDATE urls SET touched = $1 WHERE id = $2;`,
		time.Now(), id,
	)
	if err != nil {
		err = fmt.Errorf("UrlsRepo - Touch - r.Pool.Exec: %w", err)
		return
	}

	if pg.RowsAffected() <= 0 {
		err = fmt.Errorf("UrlsRepo - Touch - pg.RowsAffected: %s",
			"rows not affected by sql execution")
		return
	}

	return
}

// Rewrite -.
func (r *UrlsRepo) Rewrite(ctx context.Context, URL string) (id int, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`UPDATE urls
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
		err = fmt.Errorf("UrlsRepo - Rewrite - r.Pool.QueryRow.Scan: %w", err)
		return
	}
	return
}

// Create -.
func (r *UrlsRepo) Create(ctx context.Context, URL string) (id int, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`INSERT INTO urls (original,touched) VALUES($1, $2) RETURNING id;`,
		URL, time.Now(),
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("UrlsRepo - Create - r.Pool.QueryRow.Scan: %w", err)
		return
	}

	err = r.updCount(ctx)
	if err != nil {
		err = fmt.Errorf("UrlsRepo - Create - %w", err)
		return
	}

	return
}

func (r *UrlsRepo) updCount(ctx context.Context) (err error) {
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