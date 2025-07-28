package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantBusinessQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, merchantID int) (*response.MerchantBusinessResponse, *response.ErrorResponse)
}

type MerchantBusinessCommandService interface {
	CreateMerchant(ctx context.Context, req *requests.CreateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse)
	UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse)
	TrashedMerchant(ctx context.Context, merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(ctx context.Context, merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
