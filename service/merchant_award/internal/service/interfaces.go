package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantAwardQueryService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantAwardResponse, *response.ErrorResponse)
}

type MerchantAwardCommandService interface {
	CreateMerchant(req *requests.CreateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}
