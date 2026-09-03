package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type cartCommandHandler struct {
	pb.UnimplementedCartCommandServiceServer
	cartCommand service.CartCommandService
	logger      logger.LoggerInterface
}

func NewCartCommandHandler(cartCommand service.CartCommandService, logger logger.LoggerInterface) *cartCommandHandler {
	return &cartCommandHandler{
		cartCommand: cartCommand,
		logger:      logger,
	}
}

func (h *cartCommandHandler) Create(ctx context.Context, request *pb.CreateCartRequest) (*pb.ApiResponseCart, error) {
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

	return &pb.ApiResponseCart{
		Status:  "success",
		Message: "Successfully created cart",
		Data:    mapToProtoCartResponse(cart),
	}, nil
}

// Delete implements the proto RPC (method name must match the proto rpc name).
func (h *cartCommandHandler) Delete(ctx context.Context, request *pb.DeleteCartRequest) (*pb.ApiResponseCartDelete, error) {
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

	return &pb.ApiResponseCartDelete{
		Status:  "success",
		Message: "Successfully deleted cart item permanently",
	}, nil
}

// DeleteAll implements the proto RPC (method name must match the proto rpc name).
func (h *cartCommandHandler) DeleteAll(ctx context.Context, request *pb.DeleteAllCartRequest) (*pb.ApiResponseCartAll, error) {
	cartIDs := make([]int, len(request.GetCartIds()))
	for i, id := range request.GetCartIds() {
		cartIDs[i] = int(id)
	}

	req := &requests.DeleteAllCartRequest{
		UserID:  int(request.GetUserId()),
		CartIds: cartIDs,
	}

	if err := req.Validate(); err != nil {
		return nil, cart_errors.ErrGrpcValidateDeleteAllCart
	}

	_, err := h.cartCommand.DeleteAll(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseCartAll{
		Status:  "success",
		Message: "Successfully deleted all cart items permanently",
	}, nil
}
