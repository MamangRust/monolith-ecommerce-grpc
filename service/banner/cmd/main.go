package main

import "github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/apps"

func main() {
	server, err := apps.NewServer()

	if err != nil {
		panic(err)
	}

	server.Run()
}
