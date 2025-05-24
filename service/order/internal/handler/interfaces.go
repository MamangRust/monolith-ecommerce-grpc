package handler

import "github.com/MamangRust/monolith-point-of-sale-shared/pb"

type OrderHandleGrpc interface {
	pb.OrderServiceServer
}
