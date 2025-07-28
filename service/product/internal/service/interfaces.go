package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ProductQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, productID int) (*response.ProductResponse, *response.ErrorResponse)
}

type ProductCommandService interface {
	CreateProduct(ctx context.Context, req *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	UpdateProduct(ctx context.Context, req *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	TrashProduct(ctx context.Context, productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	RestoreProduct(ctx context.Context, productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	DeleteProductPermanent(ctx context.Context, productID int) (bool, *response.ErrorResponse)
	RestoreAllProducts(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllProductsPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
