package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-merchant/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"

	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
)

type merchantQueryHandler struct {
	pbmerchant.UnimplementedMerchantQueryServiceServer
	merchantQuery service.MerchantQueryService
	logger        logger.LoggerInterface
}

func NewMerchantQueryHandler(svc service.MerchantQueryService, logger logger.LoggerInterface) pbmerchant.MerchantQueryServiceServer {
	return &merchantQueryHandler{
		merchantQuery: svc,
		logger:        logger,
	}
}

func (s *merchantQueryHandler) FindAll(ctx context.Context, request *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchant, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchants := make([]*pbmerchant.MerchantResponse, len(merchants))
	for i, m := range merchants {
		pbMerchants[i] = mapToProtoMerchantResponse(m)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbmerchant.ApiResponsePaginationMerchant{
		Status:     "success",
		Message:    "Successfully fetched merchants",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantQueryHandler) FindById(ctx context.Context, request *pbmerchant.FindByIdMerchantRequest) (*pbmerchant.ApiResponseMerchant, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    mapToProtoMerchantResponse(merchant),
	}, nil
}

func (s *merchantQueryHandler) FindByActive(ctx context.Context, request *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchantDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchants := make([]*pbmerchant.MerchantResponseDeleteAt, len(merchants))
	for i, m := range merchants {
		pbMerchants[i] = mapToProtoMerchantResponseDeleteAt(m)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbmerchant.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchants",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantQueryHandler) FindByTrashed(ctx context.Context, request *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchantDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbMerchants := make([]*pbmerchant.MerchantResponseDeleteAt, len(merchants))
	for i, m := range merchants {
		pbMerchants[i] = mapToProtoMerchantResponseTrashed(m)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbmerchant.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchants",
		Data:       pbMerchants,
		Pagination: paginationMeta,
	}, nil
}
