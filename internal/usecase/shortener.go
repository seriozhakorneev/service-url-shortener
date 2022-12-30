package usecase

import (
	"context"
	"errors"
	"fmt"
	urlParse "net/url"
	"strings"

	internal "service-url-shortener/internal/errors"
)

// ShortenerUseCase -.
type ShortenerUseCase struct {
	repo      UrlsRepo
	digitiser Digitiser
	blank     string
}

// NewShortener -.
func NewShortener(r UrlsRepo, d Digitiser, b string) *ShortenerUseCase {
	return &ShortenerUseCase{
		repo:      r,
		digitiser: d,
		blank:     b,
	}
}

// Shorten - shortens the URL, returns short URL.
func (uc *ShortenerUseCase) Shorten(ctx context.Context, url string, ttl int) (string, error) {
	// check in storage for already existed URL and return it
	short, err := uc.exist(ctx, url, ttl)
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

// Lengthen - returns the URL associated with the given short URL
func (uc *ShortenerUseCase) Lengthen(ctx context.Context, shortURL string) (string, error) {
	short, err := uc.parseShort(shortURL)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.parseShort: %w", err)
	}

	id, err := uc.digitiser.Digit(short)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.digitiser.Digit: %w", err)
	}

	URL, err := uc.repo.GetURL(ctx, id)
	if err != nil {
		if errors.Is(err, internal.ErrNotFoundURL) {
			return "", err
		}
		return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.repo.GetURL: %w", err)
	}

	return URL, nil
}

// exist - checks the repository for an already existing URL,
// if found, returns it short representation in URL.
func (uc *ShortenerUseCase) exist(ctx context.Context, url string, ttl int) (string, error) {
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

	err = uc.repo.Touch(ctx, id, ttl)
	if err != nil {
		return "", fmt.Errorf("exist - uc.repo.Touch: %w", err)
	}

	return short, nil
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
