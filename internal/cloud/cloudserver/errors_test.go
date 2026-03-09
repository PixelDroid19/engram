package cloudserver

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gentleman-Programming/engram/internal/cloud/cloudstore"
)

func TestWriteStoreError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		fallback   string
		statusCode int
		errorBody  string
	}{
		{name: "not found", err: cloudstore.ErrNotFound, fallback: "chunk not found", statusCode: http.StatusNotFound, errorBody: "not found"},
		{name: "wrapped not found", err: fmt.Errorf("lookup failed: %w", cloudstore.ErrNotFound), fallback: "missing", statusCode: http.StatusNotFound, errorBody: "not found"},
		{name: "db down", err: errors.New("driver: bad connection"), fallback: "db", statusCode: http.StatusServiceUnavailable, errorBody: "database unavailable"},
		{name: "fallback internal", err: errors.New("boom"), fallback: "boom", statusCode: http.StatusInternalServerError, errorBody: "boom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeStoreError(rec, tt.err, tt.fallback)
			if rec.Code != tt.statusCode {
				t.Fatalf("status=%d want %d body=%s", rec.Code, tt.statusCode, rec.Body.String())
			}

			var body map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["error"] == "" {
				t.Fatal("expected error message in response")
			}
			if body["error"] != tt.errorBody {
				t.Fatalf("error body=%q want %q", body["error"], tt.errorBody)
			}
		})
	}
}

func TestIsDBConnectionError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "sql err conn done", err: sql.ErrConnDone, want: true},
		{name: "driver bad connection", err: errors.New("driver: bad connection"), want: true},
		{name: "connection refused", err: errors.New("dial tcp: connection refused"), want: true},
		{name: "generic validation error", err: errors.New("validation failed"), want: false},
		{name: "err not found", err: cloudstore.ErrNotFound, want: false},
		{name: "wrapped err not found", err: fmt.Errorf("wrapped: %w", cloudstore.ErrNotFound), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDBConnectionError(tt.err)
			if got != tt.want {
				t.Fatalf("isDBConnectionError(%v)=%v want %v", tt.err, got, tt.want)
			}
		})
	}
}
