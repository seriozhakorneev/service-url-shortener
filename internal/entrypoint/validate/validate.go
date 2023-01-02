package validate

import (
	"fmt"
	"regexp"
	"time"
)

const (
	regexpURL = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)`

	minMinuteTTL = 1
	maxMinuteTTL = 525600

	minHourTTL = 1
	maxHourTTL = 8760
)

// URL validates provided url u with regexp
func URL(u string) error {
	if u == "" {
		return fmt.Errorf("url field can not be empty")
	}

	r := regexp.MustCompile(regexpURL)
	if !r.MatchString(u) {
		return fmt.Errorf("required url field: %s, not matching a valid URL", u)
	}

	return nil
}

// TTL validates ttl data, converts it to time.Duration
func TTL(unit string, value int32) (time.Duration, error) {
	rangeErr := "required TTL value field with unit: %s should be in range %d...%d"

	switch unit {
	case "min":
		if value < minMinuteTTL || value > maxMinuteTTL {
			return 0, fmt.Errorf(rangeErr, unit, minMinuteTTL, maxMinuteTTL)
		}
		return time.Duration(value) * time.Minute, nil

	case "hour":
		if value < minHourTTL || value > maxHourTTL {
			return 0, fmt.Errorf(rangeErr, unit, minHourTTL, maxHourTTL)
		}
		return time.Duration(value) * time.Hour, nil
	default:
		return 0, fmt.Errorf(
			"imposible TTL field unit: %s, expected: min, hour",
			unit,
		)
	}
}
