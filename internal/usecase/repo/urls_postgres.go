package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	service "service-url-shortener/internal/errors"
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

// Last -.
func (r *UrlsRepo) Last(ctx context.Context) (value int, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`SELECT value
             FROM last
             WHERE id=true`,
	).Scan(&value)
	if err != nil {
		err = fmt.Errorf("UrlsRepo - Last - r.Pool.QueryRow.Scan: %w", err)
		return
	}

	return
}

// GetID -.
func (r *UrlsRepo) GetID(ctx context.Context, url string) (int, error) {
	var id int

	err := r.Pool.QueryRow(
		ctx,
		`SELECT id
             FROM urls
             WHERE original = $1`,
		url,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, service.ErrNotFoundURL
		}
		return 0, fmt.Errorf("UrlsRepo - GetID - r.Pool.QueryRow.Scan: %w", err)
	}

	return id, nil
}

// GetURL -.
func (r *UrlsRepo) GetURL(ctx context.Context, id int) (original string, liveTill time.Time, err error) {
	err = r.Pool.QueryRow(
		ctx,
		`SELECT original, live_till
             FROM urls
             WHERE id = $1 AND live_till >= $2`,
		id,
		time.Now().UTC(),
	).Scan(&original, &liveTill)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = service.ErrNotFoundURL
			return
		}

		err = fmt.Errorf("UrlsRepo - GetURL - r.Pool.QueryRow.Scan: %w", err)
		return
	}

	return
}

// Activate -.
func (r *UrlsRepo) Activate(ctx context.Context, id int, ttl time.Duration) (err error) {
	now := time.Now().UTC()

	pg, err := r.Pool.Exec(
		ctx,
		`UPDATE urls
             SET activated = $1, live_till = $2
             WHERE id = $3;`,
		now,
		now.Add(ttl),
		id,
	)
	if err != nil {
		err = fmt.Errorf("UrlsRepo - Activate - r.Pool.Exec: %w", err)
		return
	}

	if pg.RowsAffected() <= 0 {
		err = fmt.Errorf("UrlsRepo - Activate - pg.RowsAffected: %s",
			"rows not affected by sql execution")
		return
	}

	return
}

// Rewrite -.
func (r *UrlsRepo) Rewrite(ctx context.Context, url string, ttl time.Duration) (id int, err error) {
	now := time.Now().UTC()

	err = r.Pool.QueryRow(
		ctx,
		`UPDATE urls
			SET activated = $1, live_till = $2, original = $3
            WHERE activated =
            	(
            		SELECT MIN(activated)
            		FROM urls
            	)
            RETURNING id;`,
		now,
		now.Add(ttl),
		url,
	).Scan(&id)
	if err != nil {
		err = fmt.Errorf("UrlsRepo - Rewrite - r.Pool.QueryRow.Scan: %w", err)
		return
	}

	return
}

// Create -.
func (r *UrlsRepo) Create(ctx context.Context, url string, ttl time.Duration) (int, error) {
	transaction, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("UrlsRepo - Create - r.Pool.BeginTx: %w", err)
	}
	defer transaction.Rollback(ctx)

	var id int
	now := time.Now().UTC()

	err = transaction.QueryRow(
		ctx,
		`INSERT INTO urls (original, live_till, activated)
             VALUES($1, $2, $3)
             RETURNING id;`,
		url,
		now.Add(ttl),
		now,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("UrlsRepo - Create - transaction.QueryRow: %w", err)
	}

	pg, err := transaction.Exec(
		ctx,
		`UPDATE last
             SET value = $1
             WHERE id = true;`,
		id,
	)
	if err != nil {
		return 0, fmt.Errorf("UrlsRepo - Create - transaction.Exec: %w", err)
	}

	if pg.RowsAffected() <= 0 {
		return 0, fmt.Errorf("UrlsRepo - Create - pg.RowsAffected: %s",
			"rows not affected by sql execution")
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("UrlsRepo - Create - transaction.Commit: %w", err)
	}

	return id, nil
}
