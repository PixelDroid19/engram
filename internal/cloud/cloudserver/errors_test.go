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
	}{
		{name: "not found", err: cloudstore.ErrNotFound, fallback: "chunk not found", statusCode: http.StatusNotFound},
		{name: "wrapped not found", err: fmt.Errorf("lookup failed: %w", cloudstore.ErrNotFound), fallback: "missing", statusCode: http.StatusNotFound},
		{name: "db down", err: errors.New("driver: bad connection"), fallback: "db", statusCode: http.StatusServiceUnavailable},
		{name: "fallback internal", err: errors.New("boom"), fallback: "boom", statusCode: http.StatusInternalServerError},
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
		})
	}
}

func TestIsDBConnectionError(t *testing.T) {
	if !isDBConnectionError(sql.ErrConnDone) {
		t.Fatal("expected sql.ErrConnDone to be db connection error")
	}
	if !isDBConnectionError(errors.New("driver: bad connection")) {
		t.Fatal("expected bad connection to be db connection error")
	}
	if isDBConnectionError(errors.New("validation failed")) {
		t.Fatal("did not expect generic validation error to be db connection error")
	}
	if isDBConnectionError(fmt.Errorf("wrapped: %w", cloudstore.ErrNotFound)) {
		t.Fatal("did not expect wrapped ErrNotFound to be db connection error")
	}
}
