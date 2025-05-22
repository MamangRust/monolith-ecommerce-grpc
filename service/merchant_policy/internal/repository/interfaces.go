package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryRepository interface {
	FindAllMerchantPolicy(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error)
	FindById(user_id int) (*record.MerchantPoliciesRecord, error)
}

type MerchantPoliciesCommandRepository interface {
	CreateMerchantPolicy(request *requests.CreateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error)
	UpdateMerchantPolicy(request *requests.UpdateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error)
	TrashedMerchantPolicy(merchant_id int) (*record.MerchantPoliciesRecord, error)
	RestoreMerchantPolicy(merchant_id int) (*record.MerchantPoliciesRecord, error)
	DeleteMerchantPolicyPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchantPolicy() (bool, error)
	DeleteAllMerchantPolicyPermanent() (bool, error)
}

type MerchantQueryRepository interface {
	FindById(user_id int) (*record.MerchantRecord, error)
}
