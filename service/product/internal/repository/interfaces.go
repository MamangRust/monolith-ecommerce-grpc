package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryQueryRepository interface {
	FindById(ctx context.Context, categoryID int) (*record.CategoriesRecord, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, merchantID int) (*record.MerchantRecord, error)
}

type ProductQueryRepository interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*record.ProductRecord, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*record.ProductRecord, *int, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*record.ProductRecord, *int, error)

	FindByIdTrashed(ctx context.Context, productID int) (*record.ProductRecord, error)
	FindById(ctx context.Context, productID int) (*record.ProductRecord, error)
}

type ProductCommandRepository interface {
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*record.ProductRecord, error)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*record.ProductRecord, error)
	UpdateProductCountStock(ctx context.Context, productID int, stock int) (*record.ProductRecord, error)
	TrashedProduct(ctx context.Context, productID int) (*record.ProductRecord, error)
	RestoreProduct(ctx context.Context, productID int) (*record.ProductRecord, error)
	DeleteProductPermanent(ctx context.Context, productID int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}
