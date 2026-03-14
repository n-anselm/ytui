package validator

import (
	"net/url"
	"regexp"
)

var urlPattern = regexp.MustCompile(`^https?://`)

func IsValidURL(text string) bool {
	if !urlPattern.MatchString(text) {
		return false
	}

	_, err := url.ParseRequestURI(text)
	return err == nil && len(text) > 10
}
