package grpc

import (
	"fmt"
	"regexp"
)

const (
	regexpURL = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)`

	minTTL = 1   // hour
	maxTTL = 720 // month
)

func validateURL(u string) error {
	if u == "" {
		return fmt.Errorf("url field can not be empty")
	}

	r := regexp.MustCompile(regexpURL)
	if !r.MatchString(u) {
		return fmt.Errorf("required url field: %s, not matching a valid URL", u)
	}

	return nil
}

func validateTTL(ttl int32) error {
	if ttl <= minTTL || ttl > maxTTL {
		return fmt.Errorf(
			"required TTL field shoul be in range %d...%d",
			minTTL,
			maxTTL,
		)
	}

	return nil
}
