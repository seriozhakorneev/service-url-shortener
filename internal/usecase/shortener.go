package usecase

import (
	"context"
	"fmt"
)

// ShortenerUseCase -.
type ShortenerUseCase struct {
	digitiser Digitiser
	blank     string
	// string getter from short url
	// repo
}

// New -.
func New(d Digitiser, b string) *ShortenerUseCase {
	return &ShortenerUseCase{digitiser: d, blank: b}
}

// Shorten - shortens the URL, makes URL entry in the database
func (uc *ShortenerUseCase) Shorten(ctx context.Context, URL string) (string, error) {

	//TODO check in db for existed one, if there is no
	//TODO check the count of links in db
	count := len(URL)
	//TODO if count < uc.digitiser.Max() - create new
	//TODO else  - rewrite oldest(time)

	short, err := uc.digitiser.String(count)
	if err != nil {
		return "", fmt.Errorf("ShortenerUseCase - Shorten - uc.digitiser.String: %w", err)
	}

	return uc.blank + short, nil
}

// Lengthen - returns the URL associated with the given short URL
func (uc *ShortenerUseCase) Lengthen(ctx context.Context, short string) (string, error) {
	// TODO : get short string from url

	//_, err := uc.digitiser.Digit(short)
	//if err != nil {
	//	return "", fmt.Errorf("ShortenerUseCase - Lengthen - uc.digitiser.Digit: %w", err)
	//}

	//TODO: get n return URL by id from db
	return short, nil
}
