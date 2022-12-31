package usecase

import (
	"context"
	"time"
)

type (
	// Shortener -.
	Shortener interface {
		Shorten(context.Context, string, time.Duration) (string, error)
		Lengthen(context.Context, string) (string, error)
	}

	// Digitiser -.
	Digitiser interface {
		Digit(string) (int, error)
		String(int) (string, error)
		Max() int
		Length() int
	}

	// UrlsRepo -.
	UrlsRepo interface {
		Create(context.Context, string, time.Duration) (int, error)
		Rewrite(context.Context, string, time.Duration) (int, error)
		GetURL(context.Context, int) (string, time.Time, error)
		GetID(context.Context, string) (int, error)
		Activate(context.Context, int, time.Duration) error
		Count(context.Context) (int, error)
	}

	// UrlsCache -.
	UrlsCache interface {
		Set(string, string, time.Duration) error
		Get(string) (*string, error)
	}
)
