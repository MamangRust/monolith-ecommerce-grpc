package main

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/apps"
	"go.uber.org/zap"
)

func main() {
	server, shutdown, err := apps.NewServer()

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
