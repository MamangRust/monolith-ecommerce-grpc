package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantBusinessQueryService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantBusinessResponse, *response.ErrorResponse)
}

type MerchantBusinessCommandService interface {
	CreateMerchant(req *requests.CreateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}
