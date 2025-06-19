package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ProductQueryCache interface {
	GetCachedProducts(req *requests.FindAllProduct) ([]*response.ProductResponse, *int, bool)
	SetCachedProducts(req *requests.FindAllProduct, data []*response.ProductResponse, total *int)

	GetCachedProductsByMerchant(req *requests.FindAllProductByMerchant) ([]*response.ProductResponse, *int, bool)
	SetCachedProductsByMerchant(req *requests.FindAllProductByMerchant, data []*response.ProductResponse, total *int)

	GetCachedProductsByCategory(req *requests.FindAllProductByCategory) ([]*response.ProductResponse, *int, bool)
	SetCachedProductsByCategory(req *requests.FindAllProductByCategory, data []*response.ProductResponse, total *int)

	GetCachedProductActive(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductActive(req *requests.FindAllProduct, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProductTrashed(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductTrashed(req *requests.FindAllProduct, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProduct(productID int) (*response.ProductResponse, bool)
	SetCachedProduct(data *response.ProductResponse)
}

type ProductCommandCache interface {
	DeleteCachedProduct(productID int)
}
