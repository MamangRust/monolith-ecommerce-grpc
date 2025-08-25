package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bannerHandleGrpc struct {
	pb.UnimplementedBannerServiceServer
	bannerQueryService   service.BannerQueryService
	logger               logger.LoggerInterface
	bannerCommandService service.BannerCommandService
	mapping              protomapper.BannerProtoMapper
}

func NewBannerHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.BannerServiceServer {
	return &bannerHandleGrpc{
		bannerQueryService:   service.BannerQuery,
		logger:               logger,
		bannerCommandService: service.BannerCommand,
		mapping:              protomapper.NewBannerProtoMaper(),
	}

}

func (s *bannerHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllBannerRequest) (*pb.ApiResponsePaginationBanner, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all banners",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	banners, totalRecords, err := s.bannerQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all banners",
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

	s.logger.Info("Successfully fetched all banners",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationBanner(paginationMeta, "success", "Successfully fetched banners", banners)
	return so, nil
}

func (s *bannerHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBanner, error) {
	id := int(request.GetId())

	s.logger.Info("Fetching banner by ID", zap.Int("banner_id", id))

	banner, err := s.bannerQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch banner by ID",
			zap.Int("banner_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched banner by ID", zap.Int("banner_id", id))

	so := s.mapping.ToProtoResponseBanner("success", "Successfully fetched banner", banner)
	return so, nil
}

func (s *bannerHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllBannerRequest) (*pb.ApiResponsePaginationBannerDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active banners",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	banners, totalRecords, err := s.bannerQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active banners",
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

	s.logger.Info("Successfully fetched active banners",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Int32("total_records", int32(*totalRecords)),
	)

	so := s.mapping.ToProtoResponsePaginationBannerDeleteAt(paginationMeta, "success", "Successfully fetched active banners", banners)
	return so, nil
}

func (s *bannerHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllBannerRequest) (*pb.ApiResponsePaginationBannerDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed banners",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllBanner{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	banners, totalRecords, err := s.bannerQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed banners",
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

	s.logger.Info("Successfully fetched trashed banners",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Int32("total_records", int32(*totalRecords)),
	)

	so := s.mapping.ToProtoResponsePaginationBannerDeleteAt(paginationMeta, "success", "Successfully fetched trashed banners", banners)
	return so, nil
}

func (s *bannerHandleGrpc) Create(ctx context.Context, request *pb.CreateBannerRequest) (*pb.ApiResponseBanner, error) {
	s.logger.Info("Creating new banner",
		zap.String("name", request.GetName()),
		zap.String("start_date", request.GetStartDate()),
		zap.String("end_date", request.GetEndDate()),
	)

	req := &requests.CreateBannerRequest{
		Name:      request.GetName(),
		StartDate: request.GetStartDate(),
		EndDate:   request.GetEndDate(),
		StartTime: request.GetStartTime(),
		EndTime:   request.GetEndTime(),
		IsActive:  request.GetIsActive(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Banner validation failed on create",
			zap.String("name", request.GetName()),
			zap.Error(err),
		)
		return nil, banner_errors.ErrGrpcValidateCreateBanner
	}

	banner, err := s.bannerCommandService.CreateBanner(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create banner",
			zap.String("name", request.GetName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Banner created successfully",
		zap.Int("banner_id", int(banner.ID)),
		zap.String("name", banner.Name),
	)

	so := s.mapping.ToProtoResponseBanner("success", "Successfully created banner", banner)
	return so, nil
}

func (s *bannerHandleGrpc) Update(ctx context.Context, request *pb.UpdateBannerRequest) (*pb.ApiResponseBanner, error) {
	id := int(request.GetBannerId())

	s.logger.Info("Updating banner", zap.Int("banner_id", id))

	req := &requests.UpdateBannerRequest{
		BannerID:  &id,
		Name:      request.GetName(),
		StartDate: request.GetStartDate(),
		EndDate:   request.GetEndDate(),
		StartTime: request.GetStartTime(),
		EndTime:   request.GetEndTime(),
		IsActive:  request.GetIsActive(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Banner validation failed on update",
			zap.Int("banner_id", id),
			zap.Error(err),
		)
		return nil, banner_errors.ErrGrpcValidateUpdateBanner
	}

	banner, err := s.bannerCommandService.UpdateBanner(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update banner",
			zap.Int("banner_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Banner updated successfully",
		zap.Int("banner_id", id),
		zap.String("name", banner.Name),
	)

	so := s.mapping.ToProtoResponseBanner("success", "Successfully updated banner", banner)
	return so, nil
}

func (s *bannerHandleGrpc) TrashedBanner(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDeleteAt, error) {
	id := int(request.GetId())

	s.logger.Info("Trashing banner", zap.Int("banner_id", id))

	banner, err := s.bannerCommandService.TrashedBanner(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash banner",
			zap.Int("banner_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Banner moved to trash", zap.Int("banner_id", id))

	so := s.mapping.ToProtoResponseBannerDeleteAt("success", "Successfully trashed banner", banner)
	return so, nil
}

func (s *bannerHandleGrpc) RestoreBanner(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDeleteAt, error) {
	id := int(request.GetId())

	s.logger.Info("Restoring banner from trash", zap.Int("banner_id", id))

	banner, err := s.bannerCommandService.RestoreBanner(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore banner",
			zap.Int("banner_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Banner restored successfully", zap.Int("banner_id", id))

	so := s.mapping.ToProtoResponseBannerDeleteAt("success", "Successfully restored banner", banner)
	return so, nil
}

func (s *bannerHandleGrpc) DeleteBannerPermanent(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDelete, error) {
	id := int(request.GetId())

	s.logger.Info("Permanently deleting banner", zap.Int("banner_id", id))

	_, err := s.bannerCommandService.DeleteBannerPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete banner",
			zap.Int("banner_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Banner permanently deleted", zap.Int("banner_id", id))

	so := s.mapping.ToProtoResponseBannerDelete("success", "Successfully deleted banner permanently")
	return so, nil
}

func (s *bannerHandleGrpc) RestoreAllBanner(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseBannerAll, error) {
	s.logger.Info("Restoring all trashed banners")

	_, err := s.bannerCommandService.RestoreAllBanner(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all banners", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All banners restored successfully")

	so := s.mapping.ToProtoResponseBannerAll("success", "Successfully restored all banners")
	return so, nil
}

func (s *bannerHandleGrpc) DeleteAllBannerPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseBannerAll, error) {
	s.logger.Info("Permanently deleting all trashed banners")

	_, err := s.bannerCommandService.DeleteAllBannerPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all banners", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All banners permanently deleted")

	so := s.mapping.ToProtoResponseBannerAll("success", "Successfully deleted all banners permanently")
	return so, nil
}
