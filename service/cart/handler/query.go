package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-cart/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	pbcart "github.com/MamangRust/monolith-ecommerce-pb/cart"
)

type cartQueryHandler struct {
	pbcart.UnimplementedCartQueryServiceServer
	cartQuery service.CartQueryService
	logger    logger.LoggerInterface
}

func NewCartQueryHandler(cartQuery service.CartQueryService, logger logger.LoggerInterface) *cartQueryHandler {
	return &cartQueryHandler{
		cartQuery: cartQuery,
		logger:    logger,
	}
}

func (h *cartQueryHandler) FindAll(ctx context.Context, request *pbcart.FindAllCartRequest) (*pbcart.ApiResponsePaginationCart, error) {
	userID := int(request.GetUserId())
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllCarts{
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cartItems, totalRecords, err := h.cartQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCartItems := make([]*pbcart.CartResponse, len(cartItems))
	for i, cartItem := range cartItems {
		protoCartItems[i] = mapToProtoCartResponse(cartItem)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbcart.ApiResponsePaginationCart{
		Status:     "success",
		Message:    "Successfully fetched cart items",
		Data:       protoCartItems,
		Pagination: paginationMeta,
	}, nil
}
