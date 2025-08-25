package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewDetailHandleGrpc struct {
	pb.UnimplementedReviewDetailServiceServer
	reviewDetailQueryService   service.ReviewDetailQueryService
	reviewDetailCommandService service.ReviewDetailCommandService
	mapping                    protomapper.ReviewDetailProtoMapper
	mappingReview              protomapper.ReviewProtoMapper
	logger                     logger.LoggerInterface
}

func NewReviewDetailHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.ReviewDetailServiceServer {
	return &reviewDetailHandleGrpc{
		reviewDetailQueryService:   service.ReviewDetailQuery,
		reviewDetailCommandService: service.ReviewDetailCommand,
		mapping:                    protomapper.NewReviewDetailProtoMapper(),
		mappingReview:              protomapper.NewReviewProtoMapper(),
		logger:                     logger,
	}
}

func (s *reviewDetailHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetails, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all review details",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviews, totalRecords, err := s.reviewDetailQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all review details",
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

	s.logger.Info("Successfully fetched all review details",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_details_count", len(reviews)),
	)

	so := s.mapping.ToProtoResponsePaginationReviewDetail(paginationMeta, "success", "Successfully fetched review details", reviews)
	return so, nil
}

func (s *reviewDetailHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid review detail ID provided", zap.Int("detail_id", id))
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Fetching review detail by ID", zap.Int("detail_id", id))

	review, err := s.reviewDetailQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch review detail by ID",
			zap.Int("detail_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched review detail by ID",
		zap.Int("detail_id", id),
		zap.Int("review_id", int(review.ReviewID)),
	)

	so := s.mapping.ToProtoResponseReviewDetail("success", "Successfully fetched review detail", review)
	return so, nil
}

func (s *reviewDetailHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active review details",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviews, totalRecords, err := s.reviewDetailQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active review details",
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

	s.logger.Info("Successfully fetched active review details",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_details_count", len(reviews)),
	)

	so := s.mapping.ToProtoResponsePaginationReviewDetailDeleteAt(paginationMeta, "success", "Successfully fetched active review details", reviews)
	return so, nil
}

func (s *reviewDetailHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed review details",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviews, totalRecords, err := s.reviewDetailQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed review details",
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

	s.logger.Info("Successfully fetched trashed review details",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_details_count", len(reviews)),
	)

	so := s.mapping.ToProtoResponsePaginationReviewDetailDeleteAt(paginationMeta, "success", "Successfully fetched trashed review details", reviews)
	return so, nil
}

func (s *reviewDetailHandleGrpc) Create(ctx context.Context, request *pb.CreateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	s.logger.Info("Creating review detail",
		zap.Int("review_id", int(request.GetReviewId())),
		zap.String("type", request.GetType()),
		zap.String("url", request.GetUrl()),
	)

	req := &requests.CreateReviewDetailRequest{
		ReviewID: int(request.GetReviewId()),
		Type:     request.GetType(),
		Url:      request.GetUrl(),
		Caption:  request.GetCaption(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on review detail creation",
			zap.Int("review_id", int(request.GetReviewId())),
			zap.String("type", request.GetType()),
			zap.Error(err),
		)
		return nil, reviewdetail_errors.ErrGrpcValidateCreateReviewDetail
	}

	detail, err := s.reviewDetailCommandService.CreateReviewDetail(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create review detail",
			zap.Int("review_id", int(request.GetReviewId())),
			zap.String("type", request.GetType()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review detail created successfully",
		zap.Int("detail_id", int(detail.ID)),
		zap.Int("review_id", int(detail.ReviewID)),
		zap.String("type", detail.Type),
	)

	so := s.mapping.ToProtoResponseReviewDetail("success", "Successfully created review detail", detail)
	return so, nil
}

func (s *reviewDetailHandleGrpc) Update(ctx context.Context, request *pb.UpdateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	id := int(request.GetReviewDetailId())

	if id == 0 {
		s.logger.Error("Invalid review detail ID provided for update", zap.Int("detail_id", id))
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Updating review detail", zap.Int("detail_id", id))

	req := &requests.UpdateReviewDetailRequest{
		ReviewDetailID: &id,
		Type:           request.GetType(),
		Url:            request.GetUrl(),
		Caption:        request.GetCaption(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on review detail update",
			zap.Int("detail_id", id),
			zap.String("type", request.GetType()),
			zap.Error(err),
		)
		return nil, reviewdetail_errors.ErrGrpcValidateUpdateReviewDetail
	}

	detail, err := s.reviewDetailCommandService.UpdateReviewDetail(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update review detail",
			zap.Int("detail_id", id),
			zap.String("type", request.GetType()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review detail updated successfully",
		zap.Int("detail_id", id),
		zap.String("type", detail.Type),
	)

	so := s.mapping.ToProtoResponseReviewDetail("success", "Successfully updated review detail", detail)
	return so, nil
}

func (s *reviewDetailHandleGrpc) TrashedReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid review detail ID for trashing", zap.Int("detail_id", id))
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Moving review detail to trash", zap.Int("detail_id", id))

	detail, err := s.reviewDetailCommandService.TrashedReviewDetail(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash review detail",
			zap.Int("detail_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review detail moved to trash successfully",
		zap.Int("detail_id", id),
		zap.Int("review_id", int(detail.ReviewID)),
		zap.String("type", detail.Type),
	)

	so := s.mapping.ToProtoResponseReviewDetailDeleteAt("success", "Successfully trashed review detail", detail)
	return so, nil
}

func (s *reviewDetailHandleGrpc) RestoreReview(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid review detail ID for restore", zap.Int("detail_id", id))
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Restoring review detail from trash", zap.Int("detail_id", id))

	detail, err := s.reviewDetailCommandService.RestoreReviewDetail(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore review detail",
			zap.Int("detail_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review detail restored successfully",
		zap.Int("detail_id", id),
		zap.String("type", detail.Type),
	)

	so := s.mapping.ToProtoResponseReviewDetailDeleteAt("success", "Successfully restored review detail", detail)
	return so, nil
}

func (s *reviewDetailHandleGrpc) DeleteReviewPermanent(ctx context.Context, request *pb.FindByIdReviewRequest) (*pb.ApiResponseReviewDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid review detail ID for permanent deletion", zap.Int("detail_id", id))
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Permanently deleting review detail", zap.Int("detail_id", id))

	_, err := s.reviewDetailCommandService.DeleteReviewDetailPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete review detail",
			zap.Int("detail_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Review detail permanently deleted", zap.Int("detail_id", id))

	so := s.mappingReview.ToProtoResponseReviewDelete("success", "Successfully deleted review detail permanently")
	return so, nil
}

func (s *reviewDetailHandleGrpc) RestoreAllReview(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	s.logger.Info("Restoring all trashed review details")

	_, err := s.reviewDetailCommandService.RestoreAllReviewDetail(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all review details", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All review details restored successfully")

	so := s.mappingReview.ToProtoResponseReviewAll("success", "Successfully restored all review details")
	return so, nil
}

func (s *reviewDetailHandleGrpc) DeleteAllReviewPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	s.logger.Info("Permanently deleting all trashed review details")

	_, err := s.reviewDetailCommandService.DeleteAllReviewDetailPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all review details", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All review details permanently deleted")

	so := s.mappingReview.ToProtoResponseReviewAll("success", "Successfully deleted all review details permanently")
	return so, nil
}
