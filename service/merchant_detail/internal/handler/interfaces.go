package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type MerchantDetailQueryHandler interface {
	pb.MerchantDetailQueryServiceServer
}

type MerchantDetailCommandHandler interface {
	pb.MerchantDetailCommandServiceServer
}

type MerchantSocialLinkCommandHandler interface {
	pb.MerchantSocialCommandServiceServer
}
