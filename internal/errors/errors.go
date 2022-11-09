package errors

import "errors"

var ErrNotFoundURL = errors.New("url not found")
var ErrLengthTooHigh = errors.New("the length of the short url exceeds the possible")
