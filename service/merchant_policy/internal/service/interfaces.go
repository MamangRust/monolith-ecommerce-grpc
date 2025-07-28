package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantPoliciesQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, merchantID int) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
}

type MerchantPoliciesCommandService interface {
	CreateMerchant(ctx context.Context, req *requests.CreateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
	UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
	TrashedMerchant(ctx context.Context, merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(ctx context.Context, merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
