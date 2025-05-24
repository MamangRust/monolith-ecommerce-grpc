package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type OrderItemHandlerGrpc interface {
	pb.OrderItemServiceServer
}
