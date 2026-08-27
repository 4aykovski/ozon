package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mktvision/ozon/internal/app"
)

func TestHealthHandler(t *testing.T) {
	t.Parallel()

	server := app.NewHTTPServer(":8080")
	recorder := httptest.NewRecorder()

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/healthz", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	server.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
	}

	if recorder.Body.String() != "ok\n" {
		t.Fatalf("body = %q, want %q", recorder.Body.String(), "ok\n")
	}
}

func TestHealthHandlerMethodNotAllowed(t *testing.T) {
	t.Parallel()

	server := app.NewHTTPServer(":8080")
	recorder := httptest.NewRecorder()

	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/healthz", nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	server.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}
