package usecase

import (
	"context"
	"errors"
	"fmt"
	urlParse "net/url"
	"strings"
	"time"

	internal "service-url-shortener/internal/errors"
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
	url string,
	ttl time.Duration,
) (short string, err error) {
	// setting cache short:original
	// on usecase success
	defer func() {
		if err == nil {
			err = uc.cache.Set(short, url, ttl)
			if err != nil {
				err = fmt.Errorf("%w: %s", internal.ErrCaching, err)
			}
		}
	}()

	// check in storage for already existed URL and return it
	short, err = uc.exist(ctx, url, ttl)
	if err != nil {
		// return if any error, but internal.ErrNotFoundURL
		// if url not found we continue
		if !errors.Is(err, internal.ErrNotFoundURL) {
			return "", fmt.Errorf("ShortenerUseCase - Shorten - %w", err)
		}
	} else {
		return uc.blank + short, nil
	}

	// TODO: Should be transaction probably
	// ---------------------------------------------
	// count of urls in storage
	count, err := uc.repo.Count(ctx)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.repo.Count: %w", err)
	}

	var id int
	// creating new URL or rewriting the oldest URL
	if count < uc.digitiser.Max() {
		id, err = uc.repo.Create(ctx, url, ttl)
		if err != nil {
			return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.repo.Create: %w", err)
		}
	} else {
		id, err = uc.repo.Rewrite(ctx, url, ttl)
		if err != nil {
			return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.repo.Rewrite: %w", err)
		}
	}

	short, err = uc.digitiser.String(id)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.digitiser.String: %w", err)
	}
	// ---------------------------------------------
	return uc.blank + short, nil
}

// exist - checks the repository for an already existing URL,
// if found, returns it short representation in URL.
func (uc *ShortenerUseCase) exist(ctx context.Context, url string, ttl time.Duration) (string, error) {
	id, err := uc.repo.GetID(ctx, url)
	if err != nil {
		if errors.Is(err, internal.ErrNotFoundURL) {
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

// Lengthen - returns the URL associated with the given short URL
func (uc *ShortenerUseCase) Lengthen(ctx context.Context, shortURL string) (string, error) {
	// cache error
	// TODO cErr ошибки тест,не работает дефер
	var cErr error
	// check if in cache and return
	// continue if not
	value, err := uc.cache.Get(shortURL)
	if err != nil {
		cErr = err
	} else if value != nil {
		return *value, nil
	}

	defer func() {
		if err == nil && cErr != nil {
			err = fmt.Errorf("%w: %s", internal.ErrCaching, cErr)
		}
	}()

	short, err := uc.parseShort(shortURL)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.parseShort: %w", err)
	}

	id, err := uc.digitiser.Digit(short)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.digitiser.Digit: %w", err)
	}

	original, liveTill, err := uc.repo.GetURL(ctx, id)
	if err != nil {
		if errors.Is(err, internal.ErrNotFoundURL) {
			return "", err
		}
		return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.repo.GetURL: %w", err)
	}

	// add to cash after getting from db
	err = uc.cache.Set(shortURL, original, liveTill.Sub(time.Now().UTC()))
	if err != nil {
		cErr = fmt.Errorf("%w: %s", internal.ErrCaching, cErr)
	}

	cErr = fmt.Errorf("%w: %s", internal.ErrCaching, cErr)

	return original, nil
}

// parseShort - returns the short identifier after / from short url
func (uc *ShortenerUseCase) parseShort(url string) (string, error) {
	u, err := urlParse.Parse(url)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - parseShort: %w", err)
	}

	short := strings.TrimLeft(u.EscapedPath(), "/")

	if len(short) > uc.digitiser.Length() {
		return "", internal.ErrImpossibleShortURL
	}

	return short, nil
}
