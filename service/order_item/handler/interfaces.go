package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type OrderItemQueryHandler interface {
	pb.OrderItemQueryServiceServer
}

type OrderItemCommandHandler interface {
	pb.OrderItemCommandServiceServer
}
