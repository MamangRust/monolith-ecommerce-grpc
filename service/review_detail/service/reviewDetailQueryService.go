package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	review_detail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type reviewDetailQueryService struct {
	observability          observability.TraceLoggerObservability
	cache                  cache.ReviewDetailQueryCache
	reviewDetailRepository repository.ReviewDetailQueryRepository
	logger                 logger.LoggerInterface
}

type ReviewDetailQueryServiceDeps struct {
	Observability          observability.TraceLoggerObservability
	Cache                  cache.ReviewDetailQueryCache
	ReviewDetailRepository repository.ReviewDetailQueryRepository
	Logger                 logger.LoggerInterface
}

func NewReviewDetailQueryService(deps *ReviewDetailQueryServiceDeps) ReviewDetailQueryService {
	return &reviewDetailQueryService{
		observability:          deps.Observability,
		cache:                  deps.Cache,
		reviewDetailRepository: deps.ReviewDetailRepository,
		logger:                 deps.Logger,
	}
}

func (s *reviewDetailQueryService) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, *int, error) {
	const method = "FindAll"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetReviewDetailAllCache(ctx, req); found {
		logSuccess("Successfully retrieved all review detail records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviewDetails, err := s.reviewDetailRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewDetailsRow](
			s.logger,
			review_detail_errors.ErrFailedFindAllReviewDetails.WithInternal(err),
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviewDetails) > 0 {
		totalCount = int(reviewDetails[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewDetailAllCache(ctx, req, reviewDetails, &totalCount)

	logSuccess("Successfully fetched all review details",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviewDetails, &totalCount, nil
}

func (s *reviewDetailQueryService) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, *int, error) {
	const method = "FindActive"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetReviewDetailActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active review detail records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviewDetails, err := s.reviewDetailRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewDetailsActiveRow](
			s.logger,
			review_detail_errors.ErrFailedFindActiveReviewDetails.WithInternal(err),
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviewDetails) > 0 {
		totalCount = int(reviewDetails[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewDetailActiveCache(ctx, req, reviewDetails, &totalCount)

	logSuccess("Successfully fetched active review details",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviewDetails, &totalCount, nil
}

func (s *reviewDetailQueryService) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, *int, error) {
	const method = "FindTrashed"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetReviewDetailTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed review detail records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviewDetails, err := s.reviewDetailRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewDetailsTrashedRow](
			s.logger,
			review_detail_errors.ErrFailedFindTrashedReviewDetails.WithInternal(err),
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviewDetails) > 0 {
		totalCount = int(reviewDetails[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewDetailTrashedCache(ctx, req, reviewDetails, &totalCount)

	logSuccess("Successfully fetched trashed review details",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviewDetails, &totalCount, nil
}

func (s *reviewDetailQueryService) FindByID(ctx context.Context, review_id int) (*db.GetReviewDetailRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", review_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedReviewDetailCache(ctx, review_id); found {
		logSuccess("Successfully retrieved review detail by ID from cache",
			zap.Int("review_id", review_id))
		return data, nil
	}

	reviewDetail, err := s.reviewDetailRepository.FindByID(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetReviewDetailRow](
			s.logger,
			review_detail_errors.ErrFailedFindReviewDetail.WithInternal(err),
			method,
			span,

			zap.Int("review_id", review_id),
		)
	}

	s.cache.SetCachedReviewDetailCache(ctx, reviewDetail)

	logSuccess("Successfully fetched review detail by ID",
		zap.Int("review_id", review_id))

	return reviewDetail, nil
}
