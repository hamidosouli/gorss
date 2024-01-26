package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetApiKey extracts an API key from
// the headers of HTTP request
// Example:
// Authorization: ApiKey {insert apikey here}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}
	values := strings.Split(val, " ")
	if len(values) != 2 {
		return "", errors.New("malformed auth header")
	}
	if values[0] != "ApiKey" {
		return "", errors.New("malformed first part of header")
	}
	return values[1], nil

}
