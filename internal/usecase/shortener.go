package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	internal "service-url-shortener/internal/errors"
)

// ShortenerUseCase -.
type ShortenerUseCase struct {
	repo      ShortenerRepo
	digitiser Digitiser
	blank     string
}

// New -.
func New(r ShortenerRepo, d Digitiser, b string) *ShortenerUseCase {
	return &ShortenerUseCase{
		repo:      r,
		digitiser: d,
		blank:     b,
	}
}

// exist - checks the repository for an already existing URL,
// if found, returns it short representation in URL.
func (uc *ShortenerUseCase) exist(ctx context.Context, URL string) (string, error) {
	id, err := uc.repo.GetID(ctx, URL)
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

	err = uc.repo.Touch(ctx, id)
	if err != nil {
		return "", fmt.Errorf("exist - uc.repo.Touch: %w", err)
	}

	return short, nil
}

// Shorten - shortens the URL, returns short URL.
func (uc *ShortenerUseCase) Shorten(ctx context.Context, URL string) (string, error) {
	// check in storage for already existed URL and return it
	short, err := uc.exist(ctx, URL)
	if err != nil {
		if !errors.Is(err, internal.ErrNotFoundURL) {
			return "", fmt.Errorf("ShortenerUseCase - Shorten - %w", err)
		}
	} else {
		return uc.blank + short, nil
	}

	// count of urls in storage
	count, err := uc.repo.Count(ctx)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.repo.Count: %w", err)
	}

	var id int
	// creating new URL or rewriting the oldest URL
	if count < uc.digitiser.Max() {
		id, err = uc.repo.Create(ctx, URL)
		if err != nil {
			return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.repo.Create: %w", err)
		}
	} else {
		id, err = uc.repo.Rewrite(ctx, URL)
		if err != nil {
			return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.repo.Rewrite: %w", err)
		}
	}

	short, err = uc.digitiser.String(id)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.digitiser.String: %w", err)
	}

	return uc.blank + short, nil
}

// Lengthen - returns the URL associated with the given short URL
func (uc *ShortenerUseCase) Lengthen(ctx context.Context, shortURL string) (string, error) {
	u, _ := url.Parse(shortURL)

	id, err := uc.digitiser.Digit(strings.TrimLeft(u.EscapedPath(), "/"))
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
