package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"github.com/MamangRust/monolith-ecommerce-shipping-address/service"

	pbshipping_address "github.com/MamangRust/monolith-ecommerce-pb/shipping_address"
)

type shippingQueryHandler struct {
	pbshipping_address.UnimplementedShippingQueryServiceServer
	shippingQuery service.ShippingAddressQueryService
	logger        logger.LoggerInterface
}

func NewShippingQueryHandler(svc service.ShippingAddressQueryService, logger logger.LoggerInterface) pbshipping_address.ShippingQueryServiceServer {
	return &shippingQueryHandler{
		shippingQuery: svc,
		logger:        logger,
	}
}

func (s *shippingQueryHandler) FindAll(ctx context.Context, request *pbshipping_address.FindAllShippingRequest) (*pbshipping_address.ApiResponsePaginationShipping, error) {
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

	pbShippingAddresses := make([]*pbshipping_address.ShippingResponse, len(shippingAddresses))
	for i, sh := range shippingAddresses {
		pbShippingAddresses[i] = mapToProtoShippingResponse(sh)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbshipping_address.ApiResponsePaginationShipping{
		Status:     "success",
		Message:    "Successfully fetched shipping addresses",
		Data:       pbShippingAddresses,
		Pagination: paginationMeta,
	}, nil
}

func (s *shippingQueryHandler) FindById(ctx context.Context, request *pbshipping_address.FindByIdShippingRequest) (*pbshipping_address.ApiResponseShipping, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbshipping_address.ApiResponseShipping{
		Status:  "success",
		Message: "Successfully fetched shipping address",
		Data:    mapToProtoShippingResponse(shipping),
	}, nil
}

func (s *shippingQueryHandler) FindByOrder(ctx context.Context, request *pbshipping_address.FindByIdShippingRequest) (*pbshipping_address.ApiResponseShipping, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingQuery.FindByOrder(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbshipping_address.ApiResponseShipping{
		Status:  "success",
		Message: "Successfully fetched shipping address by order ID",
		Data:    mapToProtoShippingResponse(shipping),
	}, nil
}

func (s *shippingQueryHandler) FindByActive(ctx context.Context, request *pbshipping_address.FindAllShippingRequest) (*pbshipping_address.ApiResponsePaginationShippingDeleteAt, error) {
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

	pbShippingAddresses := make([]*pbshipping_address.ShippingResponseDeleteAt, len(shippingAddresses))
	for i, sh := range shippingAddresses {
		pbShippingAddresses[i] = mapToProtoShippingResponseDeleteAt(sh)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbshipping_address.ApiResponsePaginationShippingDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active shipping addresses",
		Data:       pbShippingAddresses,
		Pagination: paginationMeta,
	}, nil
}

func (s *shippingQueryHandler) FindByTrashed(ctx context.Context, request *pbshipping_address.FindAllShippingRequest) (*pbshipping_address.ApiResponsePaginationShippingDeleteAt, error) {
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

	pbShippingAddresses := make([]*pbshipping_address.ShippingResponseDeleteAt, len(shippingAddresses))
	for i, sh := range shippingAddresses {
		pbShippingAddresses[i] = mapToProtoShippingResponseDeleteAt(sh)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbshipping_address.ApiResponsePaginationShippingDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed shipping addresses",
		Data:       pbShippingAddresses,
		Pagination: paginationMeta,
	}, nil
}
