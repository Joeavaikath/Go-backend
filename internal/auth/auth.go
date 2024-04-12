package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("No Authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed Authorization header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of Authorization header")
	}

	return vals[1], nil
}
