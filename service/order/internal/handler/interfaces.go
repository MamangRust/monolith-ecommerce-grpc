package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type OrderHandleGrpc interface {
	pb.OrderServiceServer
}
