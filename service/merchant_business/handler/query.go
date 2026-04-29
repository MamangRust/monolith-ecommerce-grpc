package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type merchantBusinessQueryHandler struct {
	pb.UnimplementedMerchantBusinessQueryServiceServer
	merchantBusinessQuery service.MerchantBusinessQueryService
	logger                logger.LoggerInterface
}

func NewMerchantBusinessQueryHandler(svc service.MerchantBusinessQueryService, logger logger.LoggerInterface) MerchantBusinessQueryHandler {
	return &merchantBusinessQueryHandler{
		merchantBusinessQuery: svc,
		logger:                logger,
	}
}

func (s *merchantBusinessQueryHandler) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusiness, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantBusinessQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchants := make([]*pb.MerchantBusinessResponse, len(merchants))
	for i, merchant := range merchants {
		protoMerchants[i] = mapToProtoMerchantBusinessResponse(merchant)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantBusiness{
		Status:     "success",
		Message:    "Successfully fetched merchant",
		Data:       protoMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantBusinessQueryHandler) FindById(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	merchant, err := s.merchantBusinessQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantBusiness{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    mapToProtoMerchantBusinessResponse(merchant),
	}, nil
}

func (s *merchantBusinessQueryHandler) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusinessDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantBusinessQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchants := make([]*pb.MerchantBusinessResponseDeleteAt, len(merchants))
	for i, merchant := range merchants {
		protoMerchants[i] = mapToProtoMerchantBusinessResponseDeleteAt(merchant)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantBusinessDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchant",
		Data:       protoMerchants,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantBusinessQueryHandler) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusinessDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantBusinessQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoMerchants := make([]*pb.MerchantBusinessResponseDeleteAt, len(merchants))
	for i, merchant := range merchants {
		protoMerchants[i] = mapToProtoMerchantBusinessResponseDeleteAt(merchant)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantBusinessDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant",
		Data:       protoMerchants,
		Pagination: paginationMeta,
	}, nil
}
