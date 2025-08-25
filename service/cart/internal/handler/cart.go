package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
)

type cartHandleGrpc struct {
	pb.UnimplementedCartServiceServer
	cardQuery   service.CartQueryService
	cardCommand service.CartCommandService
	mapping     protomapper.CartProtoMapper
	logger      logger.LoggerInterface
}

func NewCartHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.CartServiceServer {
	return &cartHandleGrpc{
		cardQuery:   service.CartQuery,
		cardCommand: service.CartCommand,
		logger:      logger,
		mapping:     protomapper.NewCartProtoMapper(),
	}
}

func (s *cartHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllCartRequest) (*pb.ApiResponsePaginationCart, error) {
	userID := int(request.GetUserId())
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching cart items",
		zap.Int("user_id", userID),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllCarts{
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cartItems, totalRecords, err := s.cardQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch cart items",
			zap.Int("user_id", userID),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched cart items",
		zap.Int("user_id", userID),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationCart(paginationMeta, "success", "Successfully fetched cart items", cartItems)
	return so, nil
}

func (s *cartHandleGrpc) Create(ctx context.Context, request *pb.CreateCartRequest) (*pb.ApiResponseCart, error) {
	req := &requests.CreateCartRequest{
		Quantity:  int(request.GetQuantity()),
		ProductID: int(request.GetProductId()),
		UserID:    int(request.GetUserId()),
	}

	s.logger.Info("Adding item to cart",
		zap.Int("user_id", req.UserID),
		zap.Int("product_id", req.ProductID),
		zap.Int("quantity", req.Quantity),
	)

	if err := req.Validate(); err != nil {
		s.logger.Error("Cart item validation failed on create",
			zap.Int("user_id", req.UserID),
			zap.Int("product_id", req.ProductID),
			zap.Any("error", err),
		)
		return nil, cart_errors.ErrGrpcValidateCreateCart
	}

	cartItem, err := s.cardCommand.CreateCart(ctx, req)
	if err != nil {
		s.logger.Error("Failed to add item to cart",
			zap.Int("user_id", req.UserID),
			zap.Int("product_id", req.ProductID),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Item successfully added to cart",
		zap.Int("user_id", req.UserID),
		zap.Int("cart_item_id", int(cartItem.ID)),
		zap.Int("product_id", req.ProductID),
		zap.Int("quantity", req.Quantity),
	)

	so := s.mapping.ToProtoResponseCart("success", "Successfully added item to cart", cartItem)
	return so, nil
}

func (s *cartHandleGrpc) Delete(ctx context.Context, request *pb.DeleteCartRequest) (*pb.ApiResponseCartDelete, error) {
	userId := request.GetUserId()
	cartId := request.GetCartId()

	if userId == 0 || cartId == 0 {
		s.logger.Error("Invalid cart or user ID for deletion",
			zap.Int32("user_id", userId),
			zap.Int32("cart_id", cartId),
		)
		return nil, cart_errors.ErrGrpcCartInvalidId
	}

	_, err := s.cardCommand.DeletePermanent(ctx, &requests.DeleteCartRequest{
		UserID: int(userId),
		CartID: int(cartId),
	})
	if err != nil {
		s.logger.Error("Failed to remove item from cart",
			zap.Int32("user_id", userId),
			zap.Int32("cart_id", cartId),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Item removed from cart successfully",
		zap.Int32("user_id", userId),
		zap.Int32("cart_id", cartId),
	)

	so := s.mapping.ToProtoResponseCartDelete("success", "Successfully removed item from cart")
	return so, nil
}

func (s *cartHandleGrpc) DeleteAll(ctx context.Context, req *pb.DeleteAllCartRequest) (*pb.ApiResponseCartAll, error) {
	userId := req.GetUserId()

	if userId == 0 {
		s.logger.Error("Invalid user ID for clearing cart", zap.Int32("user_id", userId))
		return nil, cart_errors.ErrGrpcCartInvalidId
	}

	cartIDs := make([]int, len(req.GetCartIds()))
	for i, id := range req.GetCartIds() {
		cartIDs[i] = int(id)
	}

	s.logger.Info("Clearing multiple cart items",
		zap.Int32("user_id", userId),
		zap.Ints("cart_ids", cartIDs),
	)

	deleteRequest := &requests.DeleteAllCartRequest{
		UserID:  int(userId),
		CartIds: cartIDs,
	}

	_, err := s.cardCommand.DeleteAllPermanently(ctx, deleteRequest)
	if err != nil {
		s.logger.Error("Failed to clear cart items",
			zap.Int32("user_id", userId),
			zap.Ints("cart_ids", cartIDs),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully cleared cart items",
		zap.Int32("user_id", userId),
		zap.Int("items_cleared", len(cartIDs)),
	)

	so := s.mapping.ToProtoResponseCartAll("success", "Successfully cleared cart")
	return so, nil
}
