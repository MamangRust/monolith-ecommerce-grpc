package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type shippingQueryHandler struct {
	pb.UnimplementedShippingQueryServiceServer
	shippingQuery service.ShippingAddressQueryService
	logger        logger.LoggerInterface
}

func NewShippingQueryHandler(svc service.ShippingAddressQueryService, logger logger.LoggerInterface) pb.ShippingQueryServiceServer {
	return &shippingQueryHandler{
		shippingQuery: svc,
		logger:        logger,
	}
}

func (s *shippingQueryHandler) FindAll(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShipping, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	shippingAddresses, totalRecords, err := s.shippingQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbShippingAddresses := make([]*pb.ShippingResponse, len(shippingAddresses))
	for i, sh := range shippingAddresses {
		pbShippingAddresses[i] = mapToProtoShippingResponse(sh)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationShipping{
		Status:     "success",
		Message:    "Successfully fetched shipping addresses",
		Data:       pbShippingAddresses,
		Pagination: paginationMeta,
	}, nil
}

func (s *shippingQueryHandler) FindById(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShipping, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShipping{
		Status:  "success",
		Message: "Successfully fetched shipping address",
		Data:    mapToProtoShippingResponse(shipping),
	}, nil
}

func (s *shippingQueryHandler) FindByOrder(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShipping, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingQuery.FindByOrder(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShipping{
		Status:  "success",
		Message: "Successfully fetched shipping address by order ID",
		Data:    mapToProtoShippingResponse(shipping),
	}, nil
}

func (s *shippingQueryHandler) FindByActive(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShippingDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	shippingAddresses, totalRecords, err := s.shippingQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbShippingAddresses := make([]*pb.ShippingResponseDeleteAt, len(shippingAddresses))
	for i, sh := range shippingAddresses {
		pbShippingAddresses[i] = mapToProtoShippingResponseDeleteAt(sh)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationShippingDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active shipping addresses",
		Data:       pbShippingAddresses,
		Pagination: paginationMeta,
	}, nil
}

func (s *shippingQueryHandler) FindByTrashed(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShippingDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	shippingAddresses, totalRecords, err := s.shippingQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbShippingAddresses := make([]*pb.ShippingResponseDeleteAt, len(shippingAddresses))
	for i, sh := range shippingAddresses {
		pbShippingAddresses[i] = mapToProtoShippingResponseDeleteAt(sh)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationShippingDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed shipping addresses",
		Data:       pbShippingAddresses,
		Pagination: paginationMeta,
	}, nil
}
