package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantDetailQueryCache interface {
	GetCachedMerchantDetailAll(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponse, *int, bool)
	SetCachedMerchantDetailAll(req *requests.FindAllMerchant, data []*response.MerchantDetailResponse, totalRecords *int)

	GetCachedMerchantDetailActive(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, bool)
	SetCachedMerchantDetailActive(req *requests.FindAllMerchant, data []*response.MerchantDetailResponseDeleteAt, totalRecords *int)

	GetCachedMerchantDetailTrashed(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, bool)
	SetCachedMerchantDetailTrashed(req *requests.FindAllMerchant, data []*response.MerchantDetailResponseDeleteAt, totalRecords *int)

	GetCachedMerchantDetail(id int) (*response.MerchantDetailResponse, bool)
	SetCachedMerchantDetail(data *response.MerchantDetailResponse)
}

type MerchanrDetailCommandCache interface {
	DeleteMerchantDetailCache(id int)
}
