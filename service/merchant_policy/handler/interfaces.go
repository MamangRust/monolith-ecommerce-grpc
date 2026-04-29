package handler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type MerchantPolicyQueryHandler interface {
	pb.MerchantPolicyQueryServiceServer
}

type MerchantPolicyCommandHandler interface {
	pb.MerchantPolicyCommandServiceServer
}
