package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearerRaw := headers.Get("Authorization")
	if bearerRaw == "" {
		return "", fmt.Errorf("authorization header came empty")
	}
	bearer, found := strings.CutPrefix(bearerRaw, "Bearer ")
	if !found {
		return "", fmt.Errorf("authorization token not found")
	}
	return bearer, nil
}