package entity

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// TODO: Выпилить но оставить валидатор

// ShortenerData -.
type ShortenerData struct {
	URL string `json:"URL"`
}

const URLRegexp = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)`

// UnmarshalJSON - .
func (d *ShortenerData) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		return nil
	}

	tmp := struct {
		URL string `json:"URL"`
	}{}

	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}

	if tmp.URL == "" {
		return fmt.Errorf("required field 'URL' not found")
	}

	r := regexp.MustCompile(URLRegexp)
	if !r.MatchString(tmp.URL) {
		return fmt.Errorf("required field 'URL'(%s) not matching a valid URL", tmp.URL)
	}

	*d = ShortenerData{URL: tmp.URL}
	return nil
}
