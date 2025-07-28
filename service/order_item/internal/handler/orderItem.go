package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type orderItemHandleGrpc struct {
	pb.UnimplementedOrderItemServiceServer
	orderItemService service.OrderItemQueryService
	mapping          protomapper.OrderItemProtoMapper
}

func NewOrderItemHandleGrpc(
	orderItemService service.OrderItemQueryService,
	mapping protomapper.OrderItemProtoMapper,
) *orderItemHandleGrpc {
	return &orderItemHandleGrpc{
		orderItemService: orderItemService,
		mapping:          mapping,
	}
}

func (s *orderItemHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItem, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindAllOrderItems(ctx, &reqService)

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

	so := s.mapping.ToProtoResponsePaginationOrderItem(paginationMeta, "success", "Successfully fetched order items", orderItems)
	return so, nil
}

func (s *orderItemHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItemDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByActive(ctx, &reqService)

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

	so := s.mapping.ToProtoResponsePaginationOrderItemDeleteAt(paginationMeta, "success", "Successfully fetched active order items", orderItems)
	return so, nil
}

func (s *orderItemHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllOrderItemRequest) (*pb.ApiResponsePaginationOrderItemDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByTrashed(ctx, &reqService)

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

	so := s.mapping.ToProtoResponsePaginationOrderItemDeleteAt(paginationMeta, "success", "Successfully fetched trashed order items", orderItems)
	return so, nil
}

func (s *orderItemHandleGrpc) FindOrderItemByOrder(ctx context.Context, request *pb.FindByIdOrderItemRequest) (*pb.ApiResponsesOrderItem, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	orderItems, err := s.orderItemService.FindOrderItemByOrder(ctx, id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponsesOrderItem("success", "Successfully fetched order items by order", orderItems)
	return so, nil
}
