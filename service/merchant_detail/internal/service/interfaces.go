package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantDetailQueryService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantDetailResponse, *response.ErrorResponse)
}

type MerchantDetailCommandService interface {
	CreateMerchant(req *requests.CreateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (
		*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}
