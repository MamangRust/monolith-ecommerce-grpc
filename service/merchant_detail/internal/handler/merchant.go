package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDetailHandleGrpc struct {
	pb.UnimplementedMerchantDetailServiceServer
	merchantDetailQueryService   service.MerchantDetailQueryService
	merchantDetailCommandService service.MerchantDetailCommandService
	mapping                      protomapper.MerchantDetailProtoMapper
	mappingMerchant              protomapper.MerchantProtoMapper
	logger                       logger.LoggerInterface
}

func NewMerchantDetailHandleGrpc(
	service *service.Service,
	mapping protomapper.MerchantDetailProtoMapper,
	mappingMerchant protomapper.MerchantProtoMapper,
	logger logger.LoggerInterface,
) pb.MerchantDetailServiceServer {
	return &merchantDetailHandleGrpc{
		merchantDetailQueryService:   service.MerchantDetailQuery,
		merchantDetailCommandService: service.MerchantDetailCommand,
		mapping:                      mapping,
		mappingMerchant:              mappingMerchant,
		logger:                       logger,
	}
}

func (s *merchantDetailHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetail, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all merchant details",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	details, totalRecords, err := s.merchantDetailQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all merchant details",
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

	s.logger.Info("Successfully fetched all merchant details",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantDetail(paginationMeta, "success", "Successfully fetched merchant details", details)
	return so, nil
}

func (s *merchantDetailHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantDetailRequest) (*pb.ApiResponseMerchantDetail, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant detail ID provided", zap.Int("detail_id", id))
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	s.logger.Info("Fetching merchant detail by ID", zap.Int("detail_id", id))

	detail, err := s.merchantDetailQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch merchant detail by ID",
			zap.Int("detail_id", id),
			zap.Any("erroro", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched merchant detail by ID",
		zap.Int("detail_id", id),
		zap.Int("merchant_id", int(detail.MerchantID)),
	)

	so := s.mapping.ToProtoResponseMerchantDetailRelation("success", "Successfully fetched merchant detail", detail)
	return so, nil
}

func (s *merchantDetailHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetailDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active merchant details",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	details, totalRecords, err := s.merchantDetailQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active merchant details",
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

	s.logger.Info("Successfully fetched active merchant details",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantDetailDeleteAt(paginationMeta, "success", "Successfully fetched active merchant details", details)
	return so, nil
}

func (s *merchantDetailHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetailDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed merchant details",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	details, totalRecords, err := s.merchantDetailQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed merchant details",
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

	s.logger.Info("Successfully fetched trashed merchant details",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantDetailDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant details", details)
	return so, nil
}

func (s *merchantDetailHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantDetailRequest) (*pb.ApiResponseMerchantDetail, error) {
	s.logger.Info("Creating merchant detail",
		zap.Int("merchant_id", int(request.GetMerchantId())),
		zap.String("display_name", request.GetDisplayName()),
	)

	req := &requests.CreateMerchantDetailRequest{
		MerchantID:       int(request.GetMerchantId()),
		DisplayName:      request.GetDisplayName(),
		CoverImageUrl:    request.GetCoverImageUrl(),
		LogoUrl:          request.GetLogoUrl(),
		ShortDescription: request.GetShortDescription(),
		WebsiteUrl:       request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant detail creation",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("display_name", request.GetDisplayName()),
			zap.Error(err),
		)
		return nil, merchantdetail_errors.ErrGrpcValidateCreateMerchantDetail
	}

	detail, err := s.merchantDetailCommandService.CreateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create merchant detail",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("display_name", request.GetDisplayName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant detail created successfully",
		zap.Int("detail_id", int(detail.ID)),
		zap.Int("merchant_id", int(detail.MerchantID)),
		zap.String("display_name", detail.DisplayName),
		zap.Int("social_links_count", len(detail.SocialMediaLinks)),
	)

	so := s.mapping.ToProtoResponseMerchantDetail("success", "Successfully created merchant detail", detail)
	return so, nil
}

func (s *merchantDetailHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantDetailRequest) (*pb.ApiResponseMerchantDetail, error) {
	id := int(request.GetMerchantDetailId())

	if id == 0 {
		s.logger.Error("Invalid merchant detail ID provided for update", zap.Int("detail_id", id))
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	s.logger.Info("Updating merchant detail", zap.Int("detail_id", id))

	req := &requests.UpdateMerchantDetailRequest{
		MerchantDetailID: &id,
		DisplayName:      request.GetDisplayName(),
		CoverImageUrl:    request.GetCoverImageUrl(),
		LogoUrl:          request.GetLogoUrl(),
		ShortDescription: request.GetShortDescription(),
		WebsiteUrl:       request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant detail update",
			zap.Int("detail_id", id),
			zap.String("display_name", request.GetDisplayName()),
			zap.Error(err),
		)
		return nil, merchantdetail_errors.ErrGrpcValidateUpdateMerchantDetail
	}

	detail, err := s.merchantDetailCommandService.UpdateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update merchant detail",
			zap.Int("detail_id", id),
			zap.String("display_name", request.GetDisplayName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant detail updated successfully",
		zap.Int("detail_id", id),
		zap.String("display_name", detail.DisplayName),
		zap.Int("social_links_count", len(detail.SocialMediaLinks)),
	)

	so := s.mapping.ToProtoResponseMerchantDetail("success", "Successfully updated merchant detail", detail)
	return so, nil
}

func (s *merchantDetailHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant detail ID for trashing", zap.Int("detail_id", id))
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	s.logger.Info("Moving merchant detail to trash", zap.Int("detail_id", id))

	detail, err := s.merchantDetailCommandService.TrashedMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash merchant detail",
			zap.Int("detail_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant detail moved to trash successfully",
		zap.Int("detail_id", id),
		zap.Int("merchant_id", int(detail.MerchantID)),
		zap.String("display_name", detail.DisplayName),
	)

	so := s.mapping.ToProtoResponseMerchantDetailDeleteAt("success", "Successfully trashed merchant detail", detail)
	return so, nil
}

func (s *merchantDetailHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant detail ID for restore", zap.Int("detail_id", id))
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	s.logger.Info("Restoring merchant detail from trash", zap.Int("detail_id", id))

	detail, err := s.merchantDetailCommandService.RestoreMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore merchant detail",
			zap.Int("detail_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant detail restored successfully",
		zap.Int("detail_id", id),
		zap.String("display_name", detail.DisplayName),
	)

	so := s.mapping.ToProtoResponseMerchantDetailDeleteAt("success", "Successfully restored merchant detail", detail)
	return so, nil
}

func (s *merchantDetailHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant detail ID for permanent deletion", zap.Int("detail_id", id))
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	s.logger.Info("Permanently deleting merchant detail", zap.Int("detail_id", id))

	_, err := s.merchantDetailCommandService.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete merchant detail",
			zap.Int("detail_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant detail permanently deleted", zap.Int("detail_id", id))

	so := s.mappingMerchant.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant detail permanently")
	return so, nil
}

func (s *merchantDetailHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Restoring all trashed merchant details")

	_, err := s.merchantDetailCommandService.RestoreAllMerchant(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all merchant details", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant details restored successfully")

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully restored all merchant details")
	return so, nil
}

func (s *merchantDetailHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Permanently deleting all trashed merchant details")

	_, err := s.merchantDetailCommandService.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all merchant details", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant details permanently deleted")

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully deleted all merchant details permanently")
	return so, nil
}
