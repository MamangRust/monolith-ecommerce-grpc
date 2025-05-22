package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindById(user_id int) (*record.MerchantRecord, error)
}

type MerchantBusinessQueryRepository interface {
	FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error)
	FindById(user_id int) (*record.MerchantBusinessRecord, error)
}

type MerchantBusinessCommandRepository interface {
	CreateMerchantBusiness(request *requests.CreateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error)
	UpdateMerchantBusiness(request *requests.UpdateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error)
	TrashedMerchantBusiness(merchant_id int) (*record.MerchantBusinessRecord, error)
	RestoreMerchantBusiness(merchant_id int) (*record.MerchantBusinessRecord, error)
	DeleteMerchantBusinessPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchantBusiness() (bool, error)
	DeleteAllMerchantBusinessPermanent() (bool, error)
}
