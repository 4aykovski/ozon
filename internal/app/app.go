// Package app contains the application lifecycle and transport wiring.
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const shutdownTimeout = 10 * time.Second

// NewHTTPServer creates the HTTP transport used by the service.
//
// The endpoint is intentionally small for now and serves as a readiness probe
// until the first Ozon API handlers are added.
func NewHTTPServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

// Run starts the HTTP server and shuts it down when ctx is canceled.
func Run(ctx context.Context, addr string) error {
	server := NewHTTPServer(addr)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		return nil
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
