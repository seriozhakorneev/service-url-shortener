package usecase

import "context"

type (
	// Shortener -.
	Shortener interface {
		Shorten(context.Context, string) (string, error)
		Lengthen(context.Context, string) (string, error)
	}

	// Digitiser -.
	Digitiser interface {
		Digit(string) (int, error)
		String(int) (string, error)
		Max() int
	}

	// ShortenerRepo -.
	ShortenerRepo interface {
		Test(context.Context) error
	}
)
