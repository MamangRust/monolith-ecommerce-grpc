package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantAwardQueryRepository interface {
	FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error)
	FindById(user_id int) (*record.MerchantAwardRecord, error)
}

type MerchantAwardCommandRepository interface {
	CreateMerchantAward(request *requests.CreateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error)
	UpdateMerchantAward(request *requests.UpdateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error)
	TrashedMerchantAward(merchant_id int) (*record.MerchantAwardRecord, error)
	RestoreMerchantAward(merchant_id int) (*record.MerchantAwardRecord, error)
	DeleteMerchantPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchantAward() (bool, error)
	DeleteAllMerchantAwardPermanent() (bool, error)
}

type MerchantQueryRepository interface {
	FindById(user_id int) (*record.MerchantRecord, error)
}
