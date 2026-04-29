package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type merchantAwardQueryHandler struct {
	pb.UnimplementedMerchantAwardQueryServiceServer
	merchantAwardQuery service.MerchantAwardQueryService
	logger             logger.LoggerInterface
}

func NewMerchantAwardQueryHandler(svc service.MerchantAwardQueryService, logger logger.LoggerInterface) MerchantAwardQueryHandler {
	return &merchantAwardQueryHandler{
		merchantAwardQuery: svc,
		logger:             logger,
	}
}

func (s *merchantAwardQueryHandler) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantAward, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantAwardQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchants := make([]*pb.MerchantAwardResponse, len(merchants))
	for i, merchant := range merchants {
		protoMerchants[i] = mapToProtoMerchantAwardResponse(merchant)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantAward{
		Status:     "success",
		Message:    "Successfully fetched merchant",
		Data:       protoMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantAwardQueryHandler) FindById(ctx context.Context, request *pb.FindByIdMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, errors.ToGrpcError(errors.ErrInternal) // Should use specific error if available
	}

	merchant, err := s.merchantAwardQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantAward{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    mapToProtoMerchantAwardResponse(merchant),
	}, nil
}

func (s *merchantAwardQueryHandler) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantAwardDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantAwardQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchants := make([]*pb.MerchantAwardResponseDeleteAt, len(merchants))
	for i, merchant := range merchants {
		protoMerchants[i] = mapToProtoMerchantAwardResponseDeleteAt(merchant)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantAwardDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchant",
		Data:       protoMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantAwardQueryHandler) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantAwardDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantAwardQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchants := make([]*pb.MerchantAwardResponseDeleteAt, len(merchants))
	for i, merchant := range merchants {
		protoMerchants[i] = mapToProtoMerchantAwardResponseDeleteAt(merchant)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantAwardDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant",
		Data:       protoMerchants,
		Pagination: paginationMeta,
	}, nil
}
