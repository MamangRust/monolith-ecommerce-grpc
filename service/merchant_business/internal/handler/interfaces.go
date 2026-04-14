package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type MerchantBusinessQueryHandler interface {
	pb.MerchantBusinessQueryServiceServer
}

type MerchantBusinessCommandHandler interface {
	pb.MerchantBusinessCommandServiceServer
}
