package cart_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CartQueryCache interface {
	GetCachedCarts(
		ctx context.Context,
		request *requests.FindAllCarts,
	) (*response.ApiResponseCartPagination, bool)

	SetCachedCarts(
		ctx context.Context,
		request *requests.FindAllCarts,
		response *response.ApiResponseCartPagination,
	)
}
