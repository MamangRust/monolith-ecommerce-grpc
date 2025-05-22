package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shippingAddressHandleGrpc struct {
	pb.UnimplementedShippingServiceServer
	shippingAddressQueryService   service.ShippingAddressQueryService
	shippingAddressCommandService service.ShippingAddressCommandService
	mapping                       protomapper.ShippingAddresProtoMapper
}

func NewShippingAddressHandleGrpc(service service.Service) *shippingAddressHandleGrpc {
	return &shippingAddressHandleGrpc{
		shippingAddressQueryService:   service.ShippingAddressQuery,
		shippingAddressCommandService: service.ShippingAddressCommand,
		mapping:                       protomapper.NewShippingAddressProtoMapper(),
	}
}

func (s *shippingAddressHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShipping, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	Shipping, totalRecords, err := s.shippingAddressQueryService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationShippingAddress(paginationMeta, "success", "Successfully fetched categories", Shipping)
	return so, nil
}

func (s *shippingAddressHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShipping, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingAddressQueryService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddress("success", "Successfully fetched shipping address", shipping)

	return so, nil

}

func (s *shippingAddressHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShippingDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.shippingAddressQueryService.FindByActive(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationShippingAddressDeleteAt(paginationMeta, "success", "Successfully fetched active categories", users)

	return so, nil
}

func (s *shippingAddressHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShippingDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.shippingAddressQueryService.FindByTrashed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(0),
	}

	so := s.mapping.ToProtoResponsePaginationShippingAddressDeleteAt(paginationMeta, "success", "Successfully fetched trashed categories", users)

	return so, nil
}

func (s *shippingAddressHandleGrpc) TrashedShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	Shipping, err := s.shippingAddressCommandService.TrashShippingAddress(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressDeleteAt("success", "Successfully trashed Shipping", Shipping)

	return so, nil
}

func (s *shippingAddressHandleGrpc) RestoreShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	Shipping, err := s.shippingAddressCommandService.RestoreShippingAddress(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressDeleteAt("success", "Successfully restored Shipping", Shipping)

	return so, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingPermanent(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	_, err := s.shippingAddressCommandService.DeleteShippingAddressPermanently(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressDelete("success", "Successfully deleted Shipping permanently")

	return so, nil
}

func (s *shippingAddressHandleGrpc) RestoreAllShipping(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingAddressCommandService.RestoreAllShippingAddress()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressAll("success", "Successfully restore all Shipping")

	return so, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingAddressPermanently(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingAddressCommandService.DeleteAllPermanentShippingAddress()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseShippingAddressAll("success", "Successfully delete Shipping permanen")

	return so, nil
}
