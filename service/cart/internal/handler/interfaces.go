package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type CartHandleGrpc interface {
	pb.CartServiceServer
}
