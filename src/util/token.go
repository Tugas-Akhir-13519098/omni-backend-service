package util

import (
	"errors"
	"strings"
)

var (
	AuthorizationHeaderError = errors.New("error extracting authorization header")
)

func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", AuthorizationHeaderError
	}

	token := strings.Split(header, " ")
	if len(token) != 2 {
		return "", AuthorizationHeaderError
	}

	return token[1], nil
}
