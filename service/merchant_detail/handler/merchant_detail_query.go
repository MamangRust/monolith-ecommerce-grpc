package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type merchantDetailQueryHandler struct {
	pb.UnimplementedMerchantDetailQueryServiceServer
	MerchantDetailQuery service.MerchantDetailQueryService
	logger              logger.LoggerInterface
}

func NewMerchantDetailQueryHandler(svc service.MerchantDetailQueryService, logger logger.LoggerInterface) MerchantDetailQueryHandler {
	return &merchantDetailQueryHandler{
		MerchantDetailQuery: svc,
		logger:              logger,
	}
}

func (s *merchantDetailQueryHandler) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetail, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	details, totalRecords, err := s.MerchantDetailQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoDetails := make([]*pb.MerchantDetailResponse, len(details))
	for i, detail := range details {
		protoDetails[i] = mapToProtoMerchantDetailResponse(detail)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantDetail{
		Status:     "success",
		Message:    "Successfully fetched merchant details",
		Data:       protoDetails,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDetailQueryHandler) FindById(ctx context.Context, request *pb.FindByIdMerchantDetailRequest) (*pb.ApiResponseMerchantDetail, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	detail, err := s.MerchantDetailQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantDetail{
		Status:  "success",
		Message: "Successfully fetched merchant detail",
		Data:    mapToProtoMerchantDetailResponse(detail),
	}, nil
}

func (s *merchantDetailQueryHandler) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetailDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	details, totalRecords, err := s.MerchantDetailQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoDetails := make([]*pb.MerchantDetailResponseDeleteAt, len(details))
	for i, detail := range details {
		protoDetails[i] = mapToProtoMerchantDetailResponseDeleteAt(detail)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantDetailDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active merchant details",
		Data:       protoDetails,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDetailQueryHandler) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetailDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	details, totalRecords, err := s.MerchantDetailQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoDetails := make([]*pb.MerchantDetailResponseDeleteAt, len(details))
	for i, detail := range details {
		protoDetails[i] = mapToProtoMerchantDetailResponseDeleteAt(detail)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantDetailDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant details",
		Data:       protoDetails,
		Pagination: paginationMeta,
	}, nil
}
