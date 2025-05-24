package handler

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderHandleGrpc struct {
	pb.UnimplementedOrderServiceServer
	orderQuery           service.OrderQueryService
	orderCommand         service.OrderCommandService
	orderStats           service.OrderStatsService
	orderStatsByMerchant service.OrderStatsByMerchantService
	mapping              protomapper.OrderProtoMapper
}

func NewOrderHandleGrpc(service service.Service) *orderHandleGrpc {
	return &orderHandleGrpc{
		orderQuery:           service.OrderQuery,
		orderCommand:         service.OrderCommand,
		orderStats:           service.OrderStats,
		orderStatsByMerchant: service.OrderStatsByMerchant,
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

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationOrder(paginationMeta, "success", "Successfully fetched order", merchant)
	return so, nil
}

func (s *orderHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrder, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderQuery.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseOrder("success", "Successfully fetched order", merchant)

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

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindByActive(&reqService)

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
	so := s.mapping.ToProtoResponsePaginationOrderDeleteAt(paginationMeta, "success", "Successfully fetched active order", merchant)

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

	reqService := requests.FindAllOrder{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.orderQuery.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationOrderDeleteAt(paginationMeta, "success", "Successfully fetched trashed order", users)

	return so, nil
}

func (s *orderHandleGrpc) FindMonthlyTotalRevenue(ctx context.Context, req *pb.FindYearMonthTotalRevenue) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		return nil, order_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthTotalRevenue{
		Year:  year,
		Month: month,
	}

	methods, err := s.orderStats.FindMonthlyTotalRevenue(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthlyTotalRevenue("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindYearlyTotalRevenue(ctx context.Context, req *pb.FindYearTotalRevenue) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	methods, err := s.orderStats.FindYearlyTotalRevenue(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearlyTotalRevenue("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearMonthTotalRevenueByMerchant) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		return nil, order_errors.ErrGrpcInvalidMonth
	}

	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthTotalRevenueMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	methods, err := s.orderStatsByMerchant.FindMonthlyTotalRevenueByMerchant(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthlyTotalRevenue("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearTotalRevenueByMerchant) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearTotalRevenueMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.orderStatsByMerchant.FindYearlyTotalRevenueByMerchant(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearlyTotalRevenue("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *orderHandleGrpc) FindMonthlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderMonthly, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindMonthlyOrder(year)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyRevenue("success", "Monthly revenue data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) FindYearlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderYearly, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindYearlyOrder(year)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyRevenue("success", "Yearly revenue data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) FindMonthlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderMonthly, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderStatsByMerchant.FindMonthlyOrderByMerchant(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyRevenue("success", "Monthly revenue by merchant data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) FindYearlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderYearly, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderStatsByMerchant.FindYearlyOrderByMerchant(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyRevenue("success", "Yearly revenue by merchant data retrieved", res)
	return so, nil
}

func (s *orderHandleGrpc) Create(ctx context.Context, request *pb.CreateOrderRequest) (*pb.ApiResponseOrder, error) {
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

	fmt.Println(req)

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateCreateOrder
	}

	order, err := s.orderCommand.CreateOrder(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseOrder("success", "Successfully created order", order)
	return so, nil
}

func (s *orderHandleGrpc) Update(ctx context.Context, request *pb.UpdateOrderRequest) (*pb.ApiResponseOrder, error) {
	id := int(request.GetOrderId())
	idShipping := int(request.GetShipping().GetShippingId())

	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

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

	if err := req.Validate(); err != nil {
		log.Fatal(err)
		return nil, order_errors.ErrGrpcValidateUpdateOrder
	}

	order, err := s.orderCommand.UpdateOrder(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseOrder("success", "Successfully updated order", order)
	return so, nil
}

func (s *orderHandleGrpc) TrashedOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.TrashedOrder(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseOrderDeleteAt("success", "Successfully trashed order", merchant)

	return so, nil
}

func (s *orderHandleGrpc) RestoreOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.RestoreOrder(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseOrderDeleteAt("success", "Successfully restored order", merchant)

	return so, nil
}

func (s *orderHandleGrpc) DeleteOrderPermanent(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.orderCommand.DeleteOrderPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseOrderDelete("success", "Successfully deleted order permanently")

	return so, nil
}

func (s *orderHandleGrpc) RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	_, err := s.orderCommand.RestoreAllOrder()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseOrderAll("success", "Successfully restore all order")

	return so, nil
}

func (s *orderHandleGrpc) DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	_, err := s.orderCommand.DeleteAllOrderPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseOrderAll("success", "Successfully delete order permanen")

	return so, nil
}
