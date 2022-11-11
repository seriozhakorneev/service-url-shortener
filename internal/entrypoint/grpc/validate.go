package grpc

import (
	"fmt"
	"regexp"
)

const Regexp = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)`

func validateURL(u string) error {
	if u == "" {
		return fmt.Errorf("url field can not be empty")
	}

	r := regexp.MustCompile(Regexp)
	if !r.MatchString(u) {
		return fmt.Errorf("required field 'URL'(%s) not matching a valid URL", u)

	}

	return nil
}
