package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type ProductHandleGrpc interface {
	pb.ProductServiceServer
}
