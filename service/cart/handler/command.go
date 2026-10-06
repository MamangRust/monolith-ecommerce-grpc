package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-cart/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	"google.golang.org/protobuf/types/known/emptypb"

	pbcart "github.com/MamangRust/monolith-ecommerce-pb/cart"
)

type cartCommandHandler struct {
	pbcart.UnimplementedCartCommandServiceServer
	cartCommand service.CartCommandService
	logger      logger.LoggerInterface
}

func NewCartCommandHandler(cartCommand service.CartCommandService, logger logger.LoggerInterface) *cartCommandHandler {
	return &cartCommandHandler{
		cartCommand: cartCommand,
		logger:      logger,
	}
}

func (h *cartCommandHandler) Create(ctx context.Context, request *pbcart.CreateCartRequest) (*pbcart.ApiResponseCart, error) {
	req := &requests.CreateCartRequest{
		ProductID: int(request.GetProductId()),
		UserID:    int(request.GetUserId()),
		Quantity:  int(request.GetQuantity()),
	}

	if err := req.Validate(); err != nil {
		return nil, cart_errors.ErrGrpcValidateCreateCart
	}

	cart, err := h.cartCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbcart.ApiResponseCart{
		Status:  "success",
		Message: "Successfully created cart",
		Data:    mapToProtoCartResponse(cart),
	}, nil
}

func (h *cartCommandHandler) Delete(ctx context.Context, request *pbcart.DeleteCartRequest) (*pbcart.ApiResponseCartDelete, error) {
	req := &requests.DeleteCartRequest{
		CartID: int(request.GetCartId()),
		UserID: int(request.GetUserId()),
	}

	if err := req.Validate(); err != nil {
		return nil, cart_errors.ErrGrpcValidateDeleteCart
	}

	_, err := h.cartCommand.DeletePermanent(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbcart.ApiResponseCartDelete{
		Status:  "success",
		Message: "Successfully deleted cart item permanently",
	}, nil
}

func (h *cartCommandHandler) DeleteAll(ctx context.Context, request *pbcart.DeleteAllCartRequest) (*pbcart.ApiResponseCartAll, error) {
	req := &requests.DeleteAllCartRequest{
		UserID:  int(request.GetUserId()),
		CartIds: make([]int, len(request.GetCartIds())),
	}

	for i, id := range request.GetCartIds() {
		req.CartIds[i] = int(id)
	}

	if err := req.Validate(); err != nil {
		return nil, cart_errors.ErrGrpcValidateDeleteAllCart
	}

	_, err := h.cartCommand.DeleteAll(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbcart.ApiResponseCartAll{
		Status:  "success",
		Message: "Successfully deleted all cart items permanently",
	}, nil
}

func (h *cartCommandHandler) DeleteAllPermanentEmptypb(ctx context.Context, _ *emptypb.Empty) (*pbcart.ApiResponseCartAll, error) {
	_, err := h.cartCommand.DeleteAll(ctx, &requests.DeleteAllCartRequest{})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbcart.ApiResponseCartAll{
		Status:  "success",
		Message: "Successfully deleted all cart items permanently",
	}, nil
}
