package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ProductQueryService interface {
	FindAll(req *requests.FindAllProduct) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByMerchant(req *requests.FindAllProductByMerchant) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByCategory(req *requests.FindAllProductByCategory) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(productID int) (*response.ProductResponse, *response.ErrorResponse)
}

type ProductCommandService interface {
	CreateProduct(req *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	UpdateProduct(req *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	TrashProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	RestoreProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	DeleteProductPermanent(productID int) (bool, *response.ErrorResponse)
	RestoreAllProducts() (bool, *response.ErrorResponse)
	DeleteAllProductsPermanent() (bool, *response.ErrorResponse)
}
