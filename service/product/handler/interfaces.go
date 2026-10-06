package handler

import (
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
)

type ProductQueryHandler interface {
	pbproduct.ProductQueryServiceServer
}

type ProductCommandHandler interface {
	pbproduct.ProductCommandServiceServer
}
