package usecase

import (
	"context"
	"errors"
	"fmt"
	urlParse "net/url"
	"strings"
	"time"

	service "service-url-shortener/internal/errors"
)

// ShortenerUseCase -.
type ShortenerUseCase struct {
	repo      UrlsRepo
	cache     UrlsCache
	digitiser Digitiser
	blank     string
}

// NewShortener -.
func NewShortener(r UrlsRepo, c UrlsCache, d Digitiser, b string) *ShortenerUseCase {
	return &ShortenerUseCase{
		repo:      r,
		cache:     c,
		digitiser: d,
		blank:     b,
	}
}

// Shorten - shortens the URL, returns short URL.
func (uc *ShortenerUseCase) Shorten(
	ctx context.Context,
	original string,
	ttl time.Duration,
) (shortURL string, err error) {
	// setting cache on usecase success
	defer func() {
		if err == nil {
			err = uc.cache.Set(shortURL, original, ttl)
			if err != nil {
				err = fmt.Errorf("%w: %s", service.ErrCaching, err)
			}
		}
	}()

	// check in storage for already existed URL and return it
	shortURL, err = uc.exist(ctx, original, ttl)
	if err != nil {
		// return if any error, but service.ErrNotFoundURL
		// if url not found we continue
		if !errors.Is(err, service.ErrNotFoundURL) {
			return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.exist: %w", err)
		}
	} else {
		return uc.blank + shortURL, nil
	}

	shortURL, err = uc.create(ctx, original, ttl)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.create: %w", err)
	}

	return uc.blank + shortURL, nil
}

// exist - checks the repository for an already existing URL,
// if found, returns it short representation.
func (uc *ShortenerUseCase) exist(
	ctx context.Context,
	original string,
	ttl time.Duration,
) (string, error) {
	id, err := uc.repo.GetID(ctx, original)
	if err != nil {
		if errors.Is(err, service.ErrNotFoundURL) {
			return "", err
		}

		return "", fmt.Errorf("exist - uc.repo.GetID: %w", err)
	}

	short, err := uc.digitiser.String(id)
	if err != nil {
		return "", fmt.Errorf("exist - uc.digitiser.String: %w", err)
	}

	err = uc.repo.Activate(ctx, id, ttl)
	if err != nil {
		return "", fmt.Errorf("exist - uc.repo.Activate: %w", err)
	}

	return short, nil
}

// create - creating new entry of url
// or rewrite old if storage gets its limit
func (uc *ShortenerUseCase) create(
	ctx context.Context,
	original string,
	ttl time.Duration,
) (string, error) {
	// last url int-id in storage
	last, err := uc.repo.Last(ctx)
	if err != nil {
		return "", fmt.Errorf("create - uc.repo.Last: %w", err)
	}

	var id int
	// creating new entry or rewriting the oldest
	if last < uc.digitiser.Max() {
		id, err = uc.repo.Create(ctx, original, ttl)
		if err != nil {
			return "", fmt.Errorf("create - uc.repo.Create: %w", err)
		}
	} else {
		id, err = uc.repo.Rewrite(ctx, original, ttl)
		if err != nil {
			return "", fmt.Errorf("create - uc.repo.Rewrite: %w", err)
		}
	}

	short, err := uc.digitiser.String(id)
	if err != nil {
		return "", fmt.Errorf("create - uc.digitiser.String: %w", err)
	}

	return short, nil
}

// Lengthen - returns the URL associated with the given short URL
func (uc *ShortenerUseCase) Lengthen(
	ctx context.Context,
	shortURL string,
) (original string, err error) {
	var cacheErr error

	// return from cache if exist
	value, err := uc.cache.Get(shortURL)
	if err != nil {
		cacheErr = fmt.Errorf("%w: %s", service.ErrCaching, err)
	} else if value != nil {
		return *value, nil
	}

	var ttl time.Duration
	// setting cache on usecase success
	// trying to handle caching error for log
	// if no other error happens
	defer func() {
		if err == nil {
			err = uc.cache.Set(shortURL, original, ttl)
			if err != nil {
				err = fmt.Errorf("%w: %s", service.ErrCaching, err)
			}
		}

		if err == nil && cacheErr != nil {
			err = cacheErr
		}
	}()

	short, err := uc.parseShort(shortURL)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.parseShort: %w", err)
	}

	var liveTill time.Time

	original, liveTill, err = uc.get(ctx, short)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.get: %w", err)
	}

	ttl = liveTill.Sub(time.Now().UTC())

	return original, nil
}

// parseShort - returns the short identifier after / from short url
func (uc *ShortenerUseCase) parseShort(url string) (string, error) {
	u, err := urlParse.Parse(url)
	if err != nil {
		return "", fmt.Errorf("parseShort - urlParse.Parse: %w", err)
	}

	short := strings.TrimLeft(u.EscapedPath(), "/")

	if len(short) > uc.digitiser.Length() {
		return "", service.ErrImpossibleShortURL
	}

	return short, nil
}

// get - returns original url and his deactivation time, or error not found
func (uc *ShortenerUseCase) get(ctx context.Context, short string) (string, time.Time, error) {
	id, err := uc.digitiser.Digit(short)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("get - uc.digitiser.Digit: %w", err)
	}

	original, liveTill, err := uc.repo.GetURL(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFoundURL) {
			return "", time.Time{}, err
		}
		return "", time.Time{}, fmt.Errorf("get - uc.repo.GetURL: %w", err)
	}

	return original, liveTill, nil
}
