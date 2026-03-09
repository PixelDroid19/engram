package cloudserver

import (
	"strings"
	"unicode"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 50
	minPasswordLength = 8
	maxPasswordLength = 72
	maxSearchQueryLen = 500
	minSearchLimit    = 1
	maxSearchLimit    = 100
)

func ValidateUsername(username string) bool {
	username = strings.TrimSpace(username)
	runes := []rune(username)
	if len(runes) < minUsernameLength || len(runes) > maxUsernameLength {
		return false
	}

	for _, r := range runes {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.' {
			continue
		}
		return false
	}

	return true
}

func ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	at := strings.IndexByte(email, '@')
	if at <= 0 || at != strings.LastIndexByte(email, '@') || at == len(email)-1 {
		return false
	}

	domain := email[at+1:]
	if domain == "" {
		return false
	}
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	return true
}

func ValidatePassword(password string) bool {
	length := len(password)
	return length >= minPasswordLength && length <= maxPasswordLength
}

func ValidateSearchQuery(query string) bool {
	return len([]rune(strings.TrimSpace(query))) <= maxSearchQueryLen
}

func ValidateSearchLimit(limit int) bool {
	return limit >= minSearchLimit && limit <= maxSearchLimit
}
