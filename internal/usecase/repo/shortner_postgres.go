package repo

import (
	"context"
	"fmt"

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

// Test -.
func (r *ShortenerRepo) Test(ctx context.Context) error {
	sql, args, err := r.Builder.
		Insert("history").
		Columns("source").
		Values("testvalue").
		ToSql()
	if err != nil {
		return fmt.Errorf("ShortenerRepo - Test - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ShortenerRepo - Test - r.Pool.Exec: %w", err)
	}

	return nil
}
