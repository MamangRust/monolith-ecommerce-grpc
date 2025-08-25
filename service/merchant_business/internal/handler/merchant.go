package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantBusinessHandleGrpc struct {
	pb.UnimplementedMerchantBusinessServiceServer
	logger                         logger.LoggerInterface
	merchantBusinessQueryService   service.MerchantBusinessQueryService
	merchantBusinessCommandService service.MerchantBusinessCommandService
	mapping                        protomapper.MerchantBusinessProtoMapper
	mappingMerchant                protomapper.MerchantProtoMapper
}

func NewMerchantBusinessHandleGrpc(
	logger logger.LoggerInterface,
	service *service.Service,
	mapping protomapper.MerchantBusinessProtoMapper,
	mappingMerchant protomapper.MerchantProtoMapper,
) pb.MerchantBusinessServiceServer {
	return &merchantBusinessHandleGrpc{
		logger:                         logger,
		merchantBusinessQueryService:   service.MerchantBusinessQuery,
		merchantBusinessCommandService: service.MerchantBusinessCommand,
		mapping:                        mapping,
		mappingMerchant:                mappingMerchant,
	}
}

func (s *merchantBusinessHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusiness, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all merchant businesses",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	businesses, totalRecords, err := s.merchantBusinessQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all merchant businesses",
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

	s.logger.Info("Successfully fetched all merchant businesses",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantBusiness(paginationMeta, "success", "Successfully fetched merchant businesses", businesses)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant business ID provided", zap.Int("business_id", id))
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	s.logger.Info("Fetching merchant business by ID", zap.Int("business_id", id))

	business, err := s.merchantBusinessQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch merchant business by ID",
			zap.Int("business_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched merchant business by ID",
		zap.Int("business_id", id),
		zap.String("business_type", business.BusinessType),
		zap.Int("merchant_id", int(business.MerchantID)),
	)

	so := s.mapping.ToProtoResponseMerchantBusiness("success", "Successfully fetched merchant business", business)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusinessDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active merchant businesses",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	businesses, totalRecords, err := s.merchantBusinessQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active merchant businesses",
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

	s.logger.Info("Successfully fetched active merchant businesses",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantBusinessDeleteAt(paginationMeta, "success", "Successfully fetched active merchant businesses", businesses)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusinessDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed merchant businesses",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	businesses, totalRecords, err := s.merchantBusinessQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed merchant businesses",
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

	s.logger.Info("Successfully fetched trashed merchant businesses",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantBusinessDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant businesses", businesses)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	s.logger.Info("Creating merchant business information",
		zap.Int("merchant_id", int(request.GetMerchantId())),
		zap.String("business_type", request.GetBusinessType()),
		zap.Int("established_year", int(request.GetEstablishedYear())),
	)

	req := &requests.CreateMerchantBusinessInformationRequest{
		MerchantID:        int(request.GetMerchantId()),
		BusinessType:      request.GetBusinessType(),
		TaxID:             request.GetTaxId(),
		EstablishedYear:   int(request.GetEstablishedYear()),
		NumberOfEmployees: int(request.GetNumberOfEmployees()),
		WebsiteUrl:        request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant business creation",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("business_type", request.GetBusinessType()),
			zap.Error(err),
		)
		return nil, merchantbusiness_errors.ErrGrpcValidateCreateMerchantBusiness
	}

	business, err := s.merchantBusinessCommandService.CreateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create merchant business information",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("business_type", request.GetBusinessType()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant business information created successfully",
		zap.Int("business_info_id", int(business.ID)),
		zap.Int("merchant_id", int(business.MerchantID)),
		zap.String("business_type", business.BusinessType),
	)

	so := s.mapping.ToProtoResponseMerchantBusiness("success", "Successfully created merchant business information", business)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	id := int(request.GetMerchantBusinessInfoId())

	if id == 0 {
		s.logger.Error("Invalid business info ID provided for update", zap.Int("business_info_id", id))
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	s.logger.Info("Updating merchant business information", zap.Int("business_info_id", id))

	req := &requests.UpdateMerchantBusinessInformationRequest{
		MerchantBusinessInfoID: &id,
		BusinessType:           request.GetBusinessType(),
		TaxID:                  request.GetTaxId(),
		EstablishedYear:        int(request.GetEstablishedYear()),
		NumberOfEmployees:      int(request.GetNumberOfEmployees()),
		WebsiteUrl:             request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant business update",
			zap.Int("business_info_id", id),
			zap.String("business_type", request.GetBusinessType()),
			zap.Error(err),
		)
		return nil, merchantbusiness_errors.ErrGrpcValidateUpdateMerchantBusiness
	}

	business, err := s.merchantBusinessCommandService.UpdateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update merchant business information",
			zap.Int("business_info_id", id),
			zap.String("business_type", request.GetBusinessType()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant business information updated successfully",
		zap.Int("business_info_id", id),
		zap.String("business_type", business.BusinessType),
	)

	so := s.mapping.ToProtoResponseMerchantBusiness("success", "Successfully updated merchant business information", business)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) TrashedMerchantBusiness(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantBusinessDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid business info ID for trashing", zap.Int("business_info_id", id))
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	s.logger.Info("Moving merchant business info to trash", zap.Int("business_info_id", id))

	business, err := s.merchantBusinessCommandService.TrashedMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash merchant business information",
			zap.Int("business_info_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant business information moved to trash successfully",
		zap.Int("business_info_id", id),
		zap.Int("merchant_id", int(business.MerchantID)),
		zap.String("business_type", business.BusinessType),
	)

	so := s.mapping.ToProtoResponseMerchantBusinessDeleteAt("success", "Successfully trashed merchant business information", business)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) RestoreMerchantBusiness(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantBusinessDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid business info ID for restore", zap.Int("business_info_id", id))
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	s.logger.Info("Restoring merchant business info from trash", zap.Int("business_info_id", id))

	business, err := s.merchantBusinessCommandService.RestoreMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore merchant business information",
			zap.Int("business_info_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant business information restored successfully",
		zap.Int("business_info_id", id),
		zap.String("business_type", business.BusinessType),
	)

	so := s.mapping.ToProtoResponseMerchantBusinessDeleteAt("success", "Successfully restored merchant business information", business)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) DeleteMerchantBusinessPermanent(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid business info ID for permanent deletion", zap.Int("business_info_id", id))
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	s.logger.Info("Permanently deleting merchant business information", zap.Int("business_info_id", id))

	_, err := s.merchantBusinessCommandService.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete merchant business information",
			zap.Int("business_info_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant business information permanently deleted", zap.Int("business_info_id", id))

	so := s.mappingMerchant.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant business information permanently")
	return so, nil
}

func (s *merchantBusinessHandleGrpc) RestoreAllMerchantBusiness(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Restoring all trashed merchant business information")

	_, err := s.merchantBusinessCommandService.RestoreAllMerchant(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all merchant business information", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant business information restored successfully")

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully restored all merchant business information")
	return so, nil
}

func (s *merchantBusinessHandleGrpc) DeleteAllMerchantBusinessPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Permanently deleting all trashed merchant business information")

	_, err := s.merchantBusinessCommandService.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all merchant business information", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant business information permanently deleted")

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully deleted all merchant business information permanently")
	return so, nil
}
