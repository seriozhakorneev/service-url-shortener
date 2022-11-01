package repo

import (
	"context"
	"fmt"

	"service-url-shortener/pkg/postgres"
)

const _defaultEntityCap = 64

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
		Insert("postgres").
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

// GetHistory -.
//func (r *ShortenerRepo) GetHistory(ctx context.Context) ([]entity.Translation, error) {
//	sql, _, err := r.Builder.
//		Select("source, destination, original, translation").
//		From("history").
//		ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("ShortenerRepo - GetHistory - r.Builder: %w", err)
//	}
//
//	rows, err := r.Pool.Query(ctx, sql)
//	if err != nil {
//		return nil, fmt.Errorf("ShortenerRepo - GetHistory - r.Pool.Query: %w", err)
//	}
//	defer rows.Close()
//
//	entities := make([]entity.Translation, 0, _defaultEntityCap)
//
//	for rows.Next() {
//		e := entity.Translation{}
//
//		err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
//		if err != nil {
//			return nil, fmt.Errorf("ShortenerRepo - GetHistory - rows.Scan: %w", err)
//		}
//
//		entities = append(entities, e)
//	}
//
//	return entities, nil
//}
