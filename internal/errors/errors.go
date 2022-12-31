package errors

import "errors"

var ErrNotFoundURL = errors.New("url not found")
var ErrImpossibleShortURL = errors.New("provided short URL too long or impossible with current configurations")
var ErrCaching = errors.New("caching error happens")
