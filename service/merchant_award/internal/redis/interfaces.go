package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantAwardQueryCache interface {
	GetCachedMerchantAwardAll(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponse, *int, bool)
	SetCachedMerchantAwardAll(req *requests.FindAllMerchant, data []*response.MerchantAwardResponse, totalRecords *int)

	GetCachedMerchantAwardActive(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, bool)
	SetCachedMerchantAwardActive(req *requests.FindAllMerchant, data []*response.MerchantAwardResponseDeleteAt, totalRecords *int)

	GetCachedMerchantAwardTrashed(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, bool)
	SetCachedMerchantAwardTrashed(req *requests.FindAllMerchant, data []*response.MerchantAwardResponseDeleteAt, totalRecords *int)

	GetCachedMerchantAward(id int) (*response.MerchantAwardResponse, bool)
	SetCachedMerchantAward(data *response.MerchantAwardResponse)
}

type MerchanrAwardCommandCache interface {
	DeleteMerchantAwardCache(id int)
}
