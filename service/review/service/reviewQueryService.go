package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-review/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type reviewQueryService struct {
	observability    observability.TraceLoggerObservability
	cache            cache.ReviewQueryCache
	reviewRepository repository.ReviewQueryRepository
	logger           logger.LoggerInterface
}

type ReviewQueryServiceDeps struct {
	Observability    observability.TraceLoggerObservability
	Cache            cache.ReviewQueryCache
	ReviewRepository repository.ReviewQueryRepository
	Logger           logger.LoggerInterface
}

func NewReviewQueryService(deps *ReviewQueryServiceDeps) ReviewQueryService {
	return &reviewQueryService{
		observability:    deps.Observability,
		cache:            deps.Cache,
		reviewRepository: deps.ReviewRepository,
		logger:           deps.Logger,
	}
}

func (s *reviewQueryService) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, *int, error) {
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

	if data, total, found := s.cache.GetReviewAllCache(ctx, req); found {
		logSuccess("Successfully retrieved all review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewsRow](
			s.logger,
			review_errors.ErrFailedFindAllReviews.WithInternal(err),
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewAllCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched all reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviews, &totalCount, nil
}

func (s *reviewQueryService) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, *int, error) {
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

	if data, total, found := s.cache.GetReviewActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewsActiveRow](
			s.logger,
			review_errors.ErrFailedFindActiveReviews.WithInternal(err),
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewActiveCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched active reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviews, &totalCount, nil
}

func (s *reviewQueryService) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, *int, error) {
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

	if data, total, found := s.cache.GetReviewTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewsTrashedRow](
			s.logger,
			review_errors.ErrFailedFindTrashedReviews.WithInternal(err),
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewTrashedCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched trashed reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return reviews, &totalCount, nil
}

func (s *reviewQueryService) FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, *int, error) {
	const method = "FindByProduct"

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
		attribute.String("search", search),
		attribute.Int("productID", req.ProductID))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetReviewByProductCache(ctx, req); found {
		logSuccess("Successfully retrieved product review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("productID", req.ProductID))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindByProduct(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewByProductIdRow](
			s.logger,
			review_errors.ErrFailedFindByProductReviews.WithInternal(err),
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("productID", req.ProductID),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewByProductCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched product reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.Int("productID", req.ProductID))

	return reviews, &totalCount, nil
}

func (s *reviewQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, *int, error) {
	const method = "FindByMerchant"

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
		attribute.String("search", search),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetReviewByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved merchant review records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("merchantID", req.MerchantID))
		return data, total, nil
	}

	reviews, err := s.reviewRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetReviewByMerchantIdRow](
			s.logger,
			review_errors.ErrFailedFindByMerchantReviews.WithInternal(err),
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	var totalCount int

	if len(reviews) > 0 {
		totalCount = int(reviews[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetReviewByMerchantCache(ctx, req, reviews, &totalCount)

	logSuccess("Successfully fetched merchant reviews",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.Int("merchantID", req.MerchantID))

	return reviews, &totalCount, nil
}

func (s *reviewQueryService) FindByID(ctx context.Context, id int) (*db.GetReviewByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetReviewByIdCache(ctx, id); found {
		logSuccess("Successfully retrieved review by ID from cache",
			zap.Int("id", id))
		return data, nil
	}

	review, err := s.reviewRepository.FindByID(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetReviewByIDRow](
			s.logger,
			review_errors.ErrReviewNotFound.WithInternal(err),
			method,
			span,
			zap.Int("id", id),
		)
	}

	s.cache.SetReviewByIdCache(ctx, review)

	logSuccess("Successfully fetched review by ID",
		zap.Int("id", id))

	return review, nil
}
