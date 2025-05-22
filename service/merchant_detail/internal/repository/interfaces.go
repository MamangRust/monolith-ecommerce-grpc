package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindById(user_id int) (*record.MerchantRecord, error)
}

type MerchantDetailQueryRepository interface {
	FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error)
	FindById(user_id int) (*record.MerchantDetailRecord, error)
	FindByIdTrashed(user_id int) (*record.MerchantDetailRecord, error)
}

type MerchantDetailCommandRepository interface {
	CreateMerchantDetail(request *requests.CreateMerchantDetailRequest) (*record.MerchantDetailRecord, error)
	UpdateMerchantDetail(request *requests.UpdateMerchantDetailRequest) (*record.MerchantDetailRecord, error)

	TrashedMerchantDetail(merchant_detail_id int) (*record.MerchantDetailRecord, error)
	RestoreMerchantDetail(merchant_detail_id int) (*record.MerchantDetailRecord, error)
	DeleteMerchantDetailPermanent(merchant_detail_id int) (bool, error)
	RestoreAllMerchantDetail() (bool, error)
	DeleteAllMerchantDetailPermanent() (bool, error)
}

type MerchantSocialLinkCommandRepository interface {
	CreateSocialLink(req *requests.CreateMerchantSocialRequest) (bool, error)
	UpdateSocialLink(req *requests.UpdateMerchantSocialRequest) (bool, error)
	TrashSocialLink(socialID int) (bool, error)
	RestoreSocialLink(socialID int) (bool, error)
	DeletePermanentSocialLink(socialID int) (bool, error)
	RestoreAllSocialLink() (bool, error)
	DeleteAllPermanentSocialLink() (bool, error)
}
