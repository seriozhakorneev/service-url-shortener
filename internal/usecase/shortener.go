package usecase

import (
	"context"
	"fmt"
)

// ShortenerUseCase -.
type ShortenerUseCase struct {
	digitiser Digitiser
	// string getter from short url
	// repo
}

// New -.
func New(d Digitiser) *ShortenerUseCase {
	return &ShortenerUseCase{digitiser: d}
}

// Shorten - shortens the URI, makes URI entry in the database
func (uc *ShortenerUseCase) Shorten(ctx context.Context, URI string) (string, error) {

	// TODO check in db for existed one

	id, err := uc.digitiser.Digit(URI)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.digitiser.Digit: %w", err)
	}

	// TODO: set id, url to db
	// if no error, returns short_link

	shortURL, err := uc.digitiser.String(id)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.digitiser.String: %w", err)
	}

	return shortURL, nil
}

// Lengthen - returns the URI associated with the given short identifier
func (uc *ShortenerUseCase) Lengthen(ctx context.Context, short string) (string, error) {

	_, err := uc.digitiser.Digit(short)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.digitiser.Digit: %w", err)
	}

	//TODO: get n return URI by id from db
	return "URI", nil
}
