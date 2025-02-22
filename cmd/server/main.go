package main

import (
	"baconi.co.uk/oauth/internal/app/server"
	"context"
	"github.com/codingconcepts/env"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config := server.Config{}
	if err := env.Set(&config); err != nil {
		return err
	}

	engine := server.Engine(&config)

	return engine.Run(net.JoinHostPort(config.HttpHost, config.HttpPort))
}
