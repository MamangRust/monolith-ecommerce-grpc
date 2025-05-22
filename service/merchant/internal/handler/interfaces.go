package handler

import "github.com/MamangRust/monolith-ecommerce-shared/pb"

type MerchantDocumentHandleGrpc interface {
	pb.MerchantDocumentServiceServer
}

type MerchantHandleGrpc interface {
	pb.MerchantServiceServer
}
