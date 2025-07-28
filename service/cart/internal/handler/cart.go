package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type cartHandleGrpc struct {
	pb.UnimplementedCartServiceServer
	cardQuery   service.CartQueryService
	cardCommand service.CartCommandService
	mapping     protomapper.CartProtoMapper
}

func NewCartHandleGrpc(service *service.Service) *cartHandleGrpc {
	return &cartHandleGrpc{
		cardQuery:   service.CartQuery,
		cardCommand: service.CartCommand,
		mapping:     protomapper.NewCartProtoMapper(),
	}
}

func (s *cartHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllCartRequest) (*pb.ApiResponsePaginationCart, error) {
	user_id := int(request.GetUserId())
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCarts{
		UserID:   user_id,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cartItems, totalRecords, err := s.cardQuery.FindAll(ctx, &reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationCart(paginationMeta, "success", "Successfully fetched cart items", cartItems)
	return so, nil
}

func (s *cartHandleGrpc) Create(ctx context.Context, request *pb.CreateCartRequest) (*pb.ApiResponseCart, error) {
	req := &requests.CreateCartRequest{
		Quantity:  int(request.GetQuantity()),
		ProductID: int(request.GetProductId()),
		UserID:    int(request.GetUserId()),
	}

	if err := req.Validate(); err != nil {
		return nil, cart_errors.ErrGrpcValidateCreateCart
	}

	cartItem, err := s.cardCommand.CreateCart(ctx, req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCart("success", "Successfully added item to cart", cartItem)
	return so, nil
}

func (s *cartHandleGrpc) Delete(ctx context.Context, request *pb.DeleteCartRequest) (*pb.ApiResponseCartDelete, error) {
	userId := request.GetUserId()
	cartId := request.GetCartId()

	if userId == 0 || cartId == 0 {
		return nil, cart_errors.ErrGrpcCartInvalidId
	}

	_, err := s.cardCommand.DeletePermanent(ctx, &requests.DeleteCartRequest{
		UserID: int(userId),
		CartID: int(cartId),
	})

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCartDelete("success", "Successfully removed item from cart")

	return so, nil
}

func (s *cartHandleGrpc) DeleteAll(ctx context.Context, req *pb.DeleteAllCartRequest) (*pb.ApiResponseCartAll, error) {
	userId := req.GetUserId()

	if userId == 0 {
		return nil, cart_errors.ErrGrpcCartInvalidId
	}

	cartIDs := make([]int, len(req.GetCartIds()))
	for i, id := range req.GetCartIds() {
		cartIDs[i] = int(id)
	}

	deleteRequest := &requests.DeleteAllCartRequest{
		UserID:  int(userId),
		CartIds: cartIDs,
	}

	_, err := s.cardCommand.DeleteAllPermanently(ctx, deleteRequest)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCartAll("success", "Successfully cleared cart")
	return so, nil
}
