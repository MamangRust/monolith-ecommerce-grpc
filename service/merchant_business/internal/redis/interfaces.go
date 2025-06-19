package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantBusinessQueryCache interface {
	GetCachedMerchantBusinessAll(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponse, *int, bool)
	SetCachedMerchantBusinessAll(req *requests.FindAllMerchant, data []*response.MerchantBusinessResponse, totalRecords *int)

	GetCachedMerchantBusinessActive(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, bool)
	SetCachedMerchantBusinessActive(req *requests.FindAllMerchant, data []*response.MerchantBusinessResponseDeleteAt, totalRecords *int)

	GetCachedMerchantBusinessTrashed(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, bool)
	SetCachedMerchantBusinessTrashed(req *requests.FindAllMerchant, data []*response.MerchantBusinessResponseDeleteAt, totalRecords *int)

	GetCachedMerchantBusiness(id int) (*response.MerchantBusinessResponse, bool)
	SetCachedMerchantBusiness(data *response.MerchantBusinessResponse)
}

type MerchanrBusinessCommandCache interface {
	DeleteMerchantBusinessCache(id int)
}
