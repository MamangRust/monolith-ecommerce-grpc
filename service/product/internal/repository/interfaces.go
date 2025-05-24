package repository

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CategoryQueryRepository interface {
	FindById(category_id int) (*record.CategoriesRecord, error)
}

type MerchantQueryRepository interface {
	FindById(merchant_id int) (*record.MerchantRecord, error)
}

type ProductQueryRepository interface {
	FindAllProducts(req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByActive(req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByTrashed(req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByMerchant(req *requests.FindAllProductByMerchant) ([]*record.ProductRecord, *int, error)
	FindByCategory(req *requests.FindAllProductByCategory) ([]*record.ProductRecord, *int, error)

	FindByIdTrashed(product_id int) (*record.ProductRecord, error)
	FindById(product_id int) (*record.ProductRecord, error)
}

type ProductCommandRepository interface {
	CreateProduct(request *requests.CreateProductRequest) (*record.ProductRecord, error)
	UpdateProduct(request *requests.UpdateProductRequest) (*record.ProductRecord, error)
	UpdateProductCountStock(product_id int, stock int) (*record.ProductRecord, error)
	TrashedProduct(product_id int) (*record.ProductRecord, error)
	RestoreProduct(product_id int) (*record.ProductRecord, error)
	DeleteProductPermanent(product_id int) (bool, error)
	RestoreAllProducts() (bool, error)
	DeleteAllProductPermanent() (bool, error)
}
