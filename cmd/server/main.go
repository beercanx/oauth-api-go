package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"baconi.co.uk/oauth/internal/app/server"
	"github.com/codingconcepts/env"
)

func main() {

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false, // TODO - Decide if this is worth enabling
		Level:     slog.LevelDebug,
	})))

	if err := run(); err != nil {
		slog.Error("Main server error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := server.Config{}
	if configError := env.Set(&config); configError != nil {
		return configError
	}

	engine, engineError := server.Engine(ctx, config)
	if engineError != nil {
		return engineError
	}

	httpServer := &http.Server{
		Addr:              net.JoinHostPort(config.HttpHost, config.HttpPort),
		Handler:           engine,
		ReadHeaderTimeout: 1 * time.Second,
	}

	slog.Info(fmt.Sprintf("Listening and serving HTTP on http://%s", httpServer.Addr))

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below.
	listenErrCh := make(chan error, 1)
	go func() {
		if listenErr := httpServer.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			listenErrCh <- listenErr
		}
	}()

	// Listen for the interrupt signal.
	select {
	case <-ctx.Done():
	case listenErr := <-listenErrCh:
		return listenErr
	}

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	slog.Info("Server shutting down gracefully, press Ctrl+C again to force.")

	// The context is used to inform the server it has 5 seconds to finish the request it is currently handling.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Info("Server forced to shutdown")
		return err
	}

	slog.Info("Server exiting normally")
	return nil
}
