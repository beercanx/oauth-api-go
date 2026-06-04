package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"baconi.co.uk/oauth/internal/app/server"
	"github.com/codingconcepts/env"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln("Error initializing server: ", err)
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

	log.Printf("Listening and serving HTTP on http://%s\n", httpServer.Addr)

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
	log.Println("Server shutting down gracefully, press Ctrl+C again to force.")

	// The context is used to inform the server it has 5 seconds to finish the request it is currently handling.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown: ", err)
		return err
	}

	log.Println("Server exiting.")
	return nil
}
