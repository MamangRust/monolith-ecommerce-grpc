package main

import "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/apps"

func main() {
	server, err := apps.NewServer()

	if err != nil {
		panic(err)
	}

	server.Run()
}
