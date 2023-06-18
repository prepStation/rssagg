package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts the api key from
// the request's header
// Example:
// Authorization:  ApiKey {api key value}
func GetAPIKey(header http.Header) (string, error) {
	apiKey := header.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("cannot find the api key")
	}

	api := strings.Split(apiKey, " ")
	if len(api) != 2 {
		return "", errors.New("malformed api key")
	}
	if api[0] != "ApiKey" {
		return "", errors.New("malformed the first part of api key")
	}

	return api[1], nil
}
