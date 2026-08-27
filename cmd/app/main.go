package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mktvision/ozon/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, httpAddr()); err != nil {
		log.Fatal(err)
	}
}

func httpAddr() string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}

	return ":8080"
}
