package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type CartQueryHandler interface {
	pb.CartQueryServiceServer
}

type CartCommandHandler interface {
	pb.CartCommandServiceServer
}
