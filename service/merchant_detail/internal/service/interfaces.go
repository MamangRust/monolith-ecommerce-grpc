package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantDetailQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantDetailResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, merchantID int) (*response.MerchantDetailResponse, *response.ErrorResponse)
}

type MerchantDetailCommandService interface {
	CreateMerchant(ctx context.Context, req *requests.CreateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse)
	UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse)
	TrashedMerchant(ctx context.Context, merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(ctx context.Context, merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
