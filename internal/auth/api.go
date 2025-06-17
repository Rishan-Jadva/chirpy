package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", errors.New("invalid authorization header format")
	}

	key := strings.TrimSpace(strings.TrimPrefix(authHeader, "ApiKey"))
	if key == "" {
		return "", errors.New("empty key string")
	}

	return key, nil
}
