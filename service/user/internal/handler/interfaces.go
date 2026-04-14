package handler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type UserQueryHandler interface {
	pb.UserQueryServiceServer
}

type UserCommandHandler interface {
	pb.UserCommandServiceServer
}
