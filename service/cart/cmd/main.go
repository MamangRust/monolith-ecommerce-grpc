package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/apps"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		log.Println("Shutting down gracefully...")
		cancel()
	}()

	server, shutdown, err := apps.NewServer(ctx)

	if err != nil {
		server.Logger.Fatal("Failed to create server", zap.Error(err))
		panic(err)
	}

	defer func() {
		if err := shutdown(server.Ctx); err != nil {
			server.Logger.Error("Failed to shutdown tracer", zap.Error(err))
		}
	}()

	server.Run()
}
