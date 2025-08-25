package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
)

type orderItemHandleGrpc struct {
	pb.UnimplementedOrderItemServiceServer
	orderItemService service.OrderItemQueryService
	logger           logger.LoggerInterface
	mapping          protomapper.OrderItemProtoMapper
}

func NewOrderItemHandleGrpc(
	orderItemService service.OrderItemQueryService,
	mapping protomapper.OrderItemProtoMapper,
	logger logger.LoggerInterface,
) pb.OrderItemServiceServer {
	return &orderItemHandleGrpc{
		orderItemService: orderItemService,
		logger:           logger,
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

	s.logger.Info("Fetching all order items",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindAllOrderItems(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all order items",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
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

	s.logger.Info("Successfully fetched all order items",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_items_count", len(orderItems)),
	)

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

	s.logger.Info("Fetching active order items",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active order items",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
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

	s.logger.Info("Successfully fetched active order items",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_items_count", len(orderItems)),
	)

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

	s.logger.Info("Fetching trashed order items",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed order items",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
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

	s.logger.Info("Successfully fetched trashed order items",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_items_count", len(orderItems)),
	)

	so := s.mapping.ToProtoResponsePaginationOrderItemDeleteAt(paginationMeta, "success", "Successfully fetched trashed order items", orderItems)
	return so, nil
}

func (s *orderItemHandleGrpc) FindOrderItemByOrder(ctx context.Context, request *pb.FindByIdOrderItemRequest) (*pb.ApiResponsesOrderItem, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid order ID provided for fetching order items", zap.Int("order_id", id))
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Fetching order items by order ID", zap.Int("order_id", id))

	orderItems, err := s.orderItemService.FindOrderItemByOrder(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch order items by order ID",
			zap.Int("order_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched order items by order ID",
		zap.Int("order_id", id),
		zap.Int("items_count", len(orderItems)),
	)

	so := s.mapping.ToProtoResponsesOrderItem("success", "Successfully fetched order items by order", orderItems)
	return so, nil
}
