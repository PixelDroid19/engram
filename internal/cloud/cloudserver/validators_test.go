package cloudserver

import (
	"strings"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		valid    bool
	}{
		{name: "simple", username: "alice_01", valid: true},
		{name: "dot accepted", username: "alan.b", valid: true},
		{name: "unicode accepted", username: "señor_dev", valid: true},
		{name: "too short", username: "ab", valid: false},
		{name: "too long", username: strings.Repeat("a", 51), valid: false},
		{name: "space rejected", username: "alan b", valid: false},
		{name: "at rejected", username: "alan@b", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateUsername(tt.username)
			if got != tt.valid {
				t.Fatalf("ValidateUsername(%q)=%v want %v", tt.username, got, tt.valid)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{name: "normal", email: "alice@example.com", valid: true},
		{name: "localhost", email: "dev@localhost", valid: true},
		{name: "empty", email: "", valid: false},
		{name: "no at", email: "alice.example.com", valid: false},
		{name: "empty domain", email: "alice@", valid: false},
		{name: "leading dot domain", email: "alice@.example.com", valid: false},
		{name: "trailing dot domain", email: "alice@example.com.", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateEmail(tt.email)
			if got != tt.valid {
				t.Fatalf("ValidateEmail(%q)=%v want %v", tt.email, got, tt.valid)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{name: "min boundary 8", password: strings.Repeat("a", 8), valid: true},
		{name: "max boundary 72", password: strings.Repeat("a", 72), valid: true},
		{name: "too short", password: strings.Repeat("a", 7), valid: false},
		{name: "too long", password: strings.Repeat("a", 73), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePassword(tt.password)
			if got != tt.valid {
				t.Fatalf("ValidatePassword(len=%d)=%v want %v", len(tt.password), got, tt.valid)
			}
		})
	}
}

func TestValidateSearch(t *testing.T) {
	if !ValidateSearchQuery(strings.Repeat("q", 500)) {
		t.Fatal("expected 500-char query to be valid")
	}
	if ValidateSearchQuery(strings.Repeat("q", 501)) {
		t.Fatal("expected 501-char query to be invalid")
	}

	if !ValidateSearchLimit(1) || !ValidateSearchLimit(100) {
		t.Fatal("expected search limits 1 and 100 to be valid")
	}
	if ValidateSearchLimit(0) || ValidateSearchLimit(101) {
		t.Fatal("expected search limits 0 and 101 to be invalid")
	}
}
