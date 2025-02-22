package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if _, err := fmt.Println("Hello, World!"); err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("127.0.0.1", "8080"),
		Handler: NewServer(),
	}

	go func() {
		fmt.Printf("Listening on http://%s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			_, _ = fmt.Fprintf(os.Stderr, "Error listening and serving: %s\n", err)
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		_, _ = fmt.Println("Server shutting down gracefully, press Ctrl+C again to force.")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error shutting down: %s\n", err)
			os.Exit(1)
		}
	}()

	wg.Wait()

	return nil
}

func NewServer() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Processing /hello-world")
		w.Header().Set("Content-Type", "text/plain")
		if _, err := fmt.Fprintf(w, "Hello, World!"); err != nil {
			fmt.Printf("Error writing response: %s\n", err)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Processing /")
		http.NotFound(w, r)
	})

	return mux
}
