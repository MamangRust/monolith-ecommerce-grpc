package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderHandleGrpc struct {
	pb.UnimplementedOrderServiceServer
	orderQuery           service.OrderQueryService
	orderCommand         service.OrderCommandService
	orderStats           service.OrderStatsService
	orderStatsByMerchant service.OrderStatsByMerchantService
	logger               logger.LoggerInterface
	mapping              protomapper.OrderProtoMapper
}

func NewOrderHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.OrderServiceServer {
	return &orderHandleGrpc{
		orderQuery:           service.OrderQuery,
		orderCommand:         service.OrderCommand,
		orderStats:           service.OrderStats,
		orderStatsByMerchant: service.OrderStatsByMerchant,
		logger:               logger,
		mapping:              protomapper.NewOrderProtoMapper(),
	}
}

func (s *orderHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrder, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all orders",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all orders",
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

	s.logger.Info("Successfully fetched all orders",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationOrder(paginationMeta, "success", "Successfully fetched orders", orders)
	return so, nil
}

func (s *orderHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrder, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid order ID provided", zap.Int("order_id", id))
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	s.logger.Info("Fetching order by ID", zap.Int("order_id", id))

	order, err := s.orderQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch order by ID",
			zap.Int("order_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched order by ID",
		zap.Int("order_id", id),
		zap.Float64("total", float64(order.TotalPrice)),
		zap.Int("user_id", int(order.UserID)),
	)

	so := s.mapping.ToProtoResponseOrder("success", "Successfully fetched order", order)
	return so, nil
}

func (s *orderHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active orders",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active orders",
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

	s.logger.Info("Successfully fetched active orders",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationOrderDeleteAt(paginationMeta, "success", "Successfully fetched active orders", orders)
	return so, nil
}

func (s *orderHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed orders",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orders, totalRecords, err := s.orderQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed orders",
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

	s.logger.Info("Successfully fetched trashed orders",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationOrderDeleteAt(paginationMeta, "success", "Successfully fetched trashed orders", orders)
	return so, nil
}

func (s *orderHandleGrpc) FindMonthlyTotalRevenue(ctx context.Context, req *pb.FindYearMonthTotalRevenue) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly revenue", zap.Int("year", year))
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for monthly revenue", zap.Int("month", month))
		return nil, order_errors.ErrGrpcInvalidMonth
	}

	s.logger.Info("Fetching monthly total revenue for all merchants",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthTotalRevenue{
		Year:  year,
		Month: month,
	}

	revenueData, err := s.orderStats.FindMonthlyTotalRevenue(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly total revenue",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly total revenue",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("data_points", len(revenueData)),
	)

	return s.mapping.ToProtoResponseMonthlyTotalRevenue("success", "Monthly sales retrieved successfully", revenueData), nil
}

func (s *orderHandleGrpc) FindYearlyTotalRevenue(ctx context.Context, req *pb.FindYearTotalRevenue) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly revenue", zap.Int("year", year))
		return nil, order_errors.ErrGrpcInvalidYear
	}

	s.logger.Info("Fetching yearly total revenue for all merchants", zap.Int("year", year))

	revenueData, err := s.orderStats.FindYearlyTotalRevenue(ctx, year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly total revenue",
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly total revenue",
		zap.Int("year", year),
		zap.Int("data_points", len(revenueData)),
	)

	return s.mapping.ToProtoResponseYearlyTotalRevenue("success", "Yearly sales retrieved successfully", revenueData), nil
}

func (s *orderHandleGrpc) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearMonthTotalRevenueByMerchant) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	merchantID := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly merchant revenue", zap.Int("year", year))
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for monthly merchant revenue", zap.Int("month", month))
		return nil, order_errors.ErrGrpcInvalidMonth
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for monthly revenue", zap.Int("merchant_id", merchantID))
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	s.logger.Info("Fetching monthly total revenue by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthTotalRevenueMerchant{
		Year:       year,
		Month:      month,
		MerchantID: merchantID,
	}

	revenueData, err := s.orderStatsByMerchant.FindMonthlyTotalRevenueByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly total revenue by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly total revenue by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	return s.mapping.ToProtoResponseMonthlyTotalRevenue("success", "Monthly sales retrieved successfully", revenueData), nil
}

func (s *orderHandleGrpc) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearTotalRevenueByMerchant) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	year := int(req.GetYear())
	merchantID := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly merchant revenue", zap.Int("year", year))
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for yearly revenue", zap.Int("merchant_id", merchantID))
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	s.logger.Info("Fetching yearly total revenue by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.YearTotalRevenueMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	revenueData, err := s.orderStatsByMerchant.FindYearlyTotalRevenueByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly total revenue by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly total revenue by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	return s.mapping.ToProtoResponseYearlyTotalRevenue("success", "Yearly sales retrieved successfully", revenueData), nil
}

func (s *orderHandleGrpc) FindMonthlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderMonthly, error) {
	year := int(request.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly revenue trend", zap.Int("year", year))
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	s.logger.Info("Fetching monthly revenue trend for all merchants", zap.Int("year", year))

	res, err := s.orderStats.FindMonthlyOrder(ctx, year)
	if err != nil {
		s.logger.Error("Failed to fetch monthly revenue trend",
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly revenue trend",
		zap.Int("year", year),
		zap.Int("data_points", len(res)),
	)

	so := s.mapping.ToProtoResponseMonthlyRevenue("success", "Monthly revenue data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) FindYearlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderYearly, error) {
	year := int(request.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly revenue trend", zap.Int("year", year))
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	s.logger.Info("Fetching yearly revenue trend for all merchants", zap.Int("year", year))

	res, err := s.orderStats.FindYearlyOrder(ctx, year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly revenue trend",
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly revenue trend",
		zap.Int("year", year),
		zap.Int("data_points", len(res)),
	)

	so := s.mapping.ToProtoResponseYearlyRevenue("success", "Yearly revenue data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) FindMonthlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderMonthly, error) {
	year := int(request.GetYear())
	merchantID := int(request.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly merchant revenue trend", zap.Int("year", year))
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for monthly revenue trend", zap.Int("merchant_id", merchantID))
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	s.logger.Info("Fetching monthly revenue trend by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.MonthOrderMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	res, err := s.orderStatsByMerchant.FindMonthlyOrderByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly revenue trend by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly revenue trend by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("data_points", len(res)),
	)

	so := s.mapping.ToProtoResponseMonthlyRevenue("success", "Monthly revenue by merchant data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) FindYearlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderYearly, error) {
	year := int(request.GetYear())
	merchantID := int(request.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly merchant revenue trend", zap.Int("year", year))
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for yearly revenue trend", zap.Int("merchant_id", merchantID))
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	s.logger.Info("Fetching yearly revenue trend by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.YearOrderMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	res, err := s.orderStatsByMerchant.FindYearlyOrderByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly revenue trend by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly revenue trend by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("data_points", len(res)),
	)

	so := s.mapping.ToProtoResponseYearlyRevenue("success", "Yearly revenue by merchant data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) Create(ctx context.Context, request *pb.CreateOrderRequest) (*pb.ApiResponseOrder, error) {
	s.logger.Info("Creating new order",
		zap.Int("merchant_id", int(request.GetMerchantId())),
		zap.Int("user_id", int(request.UserId)),
		zap.Int("total_price", int(request.GetTotalPrice())),
		zap.Int("items_count", len(request.GetItems())),
	)

	req := &requests.CreateOrderRequest{
		MerchantID: int(request.GetMerchantId()),
		UserID:     int(request.UserId),
		TotalPrice: int(request.GetTotalPrice()),
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.CreateOrderItemRequest{
			ProductID: int(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
			Price:     int(item.GetPrice()),
		})
	}

	if request.Shipping != nil {
		req.ShippingAddress = requests.CreateShippingAddressRequest{
			Alamat:         request.Shipping.GetAlamat(),
			Provinsi:       request.Shipping.GetProvinsi(),
			Kota:           request.Shipping.GetKota(),
			Courier:        request.Shipping.GetCourier(),
			ShippingMethod: request.Shipping.GetShippingMethod(),
			ShippingCost:   int(request.Shipping.GetShippingCost()),
			Negara:         request.Shipping.GetNegara(),
		}
	}

	s.logger.Debug("CreateOrder request payload", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on order creation",
			zap.Int("user_id", int(request.UserId)),
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.Error(err),
		)
		return nil, order_errors.ErrGrpcValidateCreateOrder
	}

	order, err := s.orderCommand.CreateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create order",
			zap.Int("user_id", int(request.UserId)),
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Order created successfully",
		zap.Int("order_id", int(order.ID)),
		zap.Int("total_price", order.TotalPrice),
	)

	so := s.mapping.ToProtoResponseOrder("success", "Successfully created order", order)
	return so, nil
}

func (s *orderHandleGrpc) Update(ctx context.Context, request *pb.UpdateOrderRequest) (*pb.ApiResponseOrder, error) {
	id := int(request.GetOrderId())
	idShipping := int(request.GetShipping().GetShippingId())

	if id == 0 {
		s.logger.Error("Invalid order ID provided for update", zap.Int("order_id", id))
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	s.logger.Info("Updating order", zap.Int("order_id", id))

	req := &requests.UpdateOrderRequest{
		OrderID:    &id,
		UserID:     int(request.GetUserId()),
		TotalPrice: int(request.GetTotalPrice()),
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.UpdateOrderItemRequest{
			OrderItemID: int(item.GetOrderItemId()),
			ProductID:   int(item.GetProductId()),
			Quantity:    int(item.GetQuantity()),
			Price:       int(item.GetPrice()),
		})
	}

	if request.Shipping != nil {
		req.ShippingAddress = requests.UpdateShippingAddressRequest{
			ShippingID:     &idShipping,
			Alamat:         request.Shipping.GetAlamat(),
			Provinsi:       request.Shipping.GetProvinsi(),
			Kota:           request.Shipping.GetKota(),
			Courier:        request.Shipping.GetCourier(),
			ShippingMethod: request.Shipping.GetShippingMethod(),
			ShippingCost:   int(request.Shipping.GetShippingCost()),
			Negara:         request.Shipping.GetNegara(),
		}
	}

	s.logger.Debug("UpdateOrder request payload", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on order update",
			zap.Int("order_id", id),
			zap.Error(err),
		)
		return nil, order_errors.ErrGrpcValidateUpdateOrder
	}

	order, err := s.orderCommand.UpdateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update order",
			zap.Int("order_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Order updated successfully",
		zap.Int("order_id", id),
	)

	so := s.mapping.ToProtoResponseOrder("success", "Successfully updated order", order)
	return so, nil
}

func (s *orderHandleGrpc) TrashedOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid order ID for trashing", zap.Int("order_id", id))
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	s.logger.Info("Moving order to trash", zap.Int("order_id", id))

	order, err := s.orderCommand.TrashedOrder(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash order",
			zap.Int("order_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Order moved to trash successfully",
		zap.Int("order_id", id),
	)

	so := s.mapping.ToProtoResponseOrderDeleteAt("success", "Successfully trashed order", order)
	return so, nil
}

func (s *orderHandleGrpc) RestoreOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid order ID for restore", zap.Int("order_id", id))
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	s.logger.Info("Restoring order from trash", zap.Int("order_id", id))

	order, err := s.orderCommand.RestoreOrder(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore order",
			zap.Int("order_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Order restored successfully",
		zap.Int("order_id", id),
	)

	so := s.mapping.ToProtoResponseOrderDeleteAt("success", "Successfully restored order", order)
	return so, nil
}

func (s *orderHandleGrpc) DeleteOrderPermanent(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid order ID for permanent deletion", zap.Int("order_id", id))
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	s.logger.Info("Permanently deleting order", zap.Int("order_id", id))

	_, err := s.orderCommand.DeleteOrderPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete order",
			zap.Int("order_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Order permanently deleted", zap.Int("order_id", id))

	so := s.mapping.ToProtoResponseOrderDelete("success", "Successfully deleted order permanently")
	return so, nil
}

func (s *orderHandleGrpc) RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	s.logger.Info("Restoring all trashed orders")

	_, err := s.orderCommand.RestoreAllOrder(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all orders", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All orders restored successfully")

	so := s.mapping.ToProtoResponseOrderAll("success", "Successfully restored all orders")
	return so, nil
}

func (s *orderHandleGrpc) DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	s.logger.Info("Permanently deleting all trashed orders")

	_, err := s.orderCommand.DeleteAllOrderPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all orders", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All orders permanently deleted")

	so := s.mapping.ToProtoResponseOrderAll("success", "Successfully deleted all orders permanently")
	return so, nil
}
