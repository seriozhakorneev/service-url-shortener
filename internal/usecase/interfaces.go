package usecase

import "context"

type (
	// Shortener -.
	Shortener interface {
		Shorten(context.Context, string, int) (string, error)
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
		Create(context.Context, string, int) (int, error)
		Rewrite(context.Context, string, int) (int, error)
		GetURL(context.Context, int) (string, error)
		GetID(context.Context, string) (int, error)
		Touch(context.Context, int, int) error
		Count(context.Context) (int, error)
	}
)
