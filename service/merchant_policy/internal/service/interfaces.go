package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantPoliciesQueryService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
}

type MerchantPoliciesCommandService interface {
	CreateMerchant(req *requests.CreateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}
