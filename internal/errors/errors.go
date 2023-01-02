package errors

import "errors"

var (
	ErrNotFoundURL        = errors.New("url not found")
	ErrImpossibleShortURL = errors.New("provided short url too long or impossible with current configurations")
	ErrCaching            = errors.New("caching error happens")
)
