package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type ReviewHandleGrpc interface {
	pb.ReviewQueryServiceServer
	pb.ReviewCommandServiceServer
}

type ReviewQueryHandler interface {
	pb.ReviewQueryServiceServer
}

type ReviewCommandHandler interface {
	pb.ReviewCommandServiceServer
}
