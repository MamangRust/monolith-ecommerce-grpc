package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantAwardQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantAwardResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, merchantID int) (*response.MerchantAwardResponse, *response.ErrorResponse)
}

type MerchantAwardCommandService interface {
	CreateMerchant(ctx context.Context, req *requests.CreateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse)
	UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse)
	TrashedMerchant(ctx context.Context, merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(ctx context.Context, merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
