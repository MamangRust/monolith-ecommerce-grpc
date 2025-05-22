package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewQueryService struct {
	ctx                   context.Context
	trace                 trace.Tracer
	reviewQueryRepository repository.ReviewQueryRepository
	mapping               response_service.ReviewResponseMapper
	logger                logger.LoggerInterface
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewReviewQueryService(ctx context.Context, reviewQueryRepository repository.ReviewQueryRepository, mapping response_service.ReviewResponseMapper, logger logger.LoggerInterface) *reviewQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "review_query_service_request_count",
			Help: "Total number of requests to the ReviewQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "review_query_service_request_count",
			Help:    "Histogram of request durations for the ReviewQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &reviewQueryService{
		ctx:                   ctx,
		trace:                 otel.Tracer("review-query-service"),
		reviewQueryRepository: reviewQueryRepository,
		mapping:               mapping,
		logger:                logger,
		requestCounter:        requestCounter,
		requestDuration:       requestDuration,
	}
}

func (s *reviewQueryService) FindAllReviews(req *requests.FindAllReview) ([]*response.ReviewResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAllReviews", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllReviews")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Reviews, totalRecords, err := s.reviewQueryRepository.FindAllReview(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_REVIEWS")

		s.logger.Error("Failed to retrieve review list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve review list")
		span.End()

		status = "failed_find_all_reviews"

		return nil, nil, review_errors.ErrFailedFindAllReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponse(Reviews), totalRecords, nil
}

func (s *reviewQueryService) FindByActive(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActive")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Reviews, totalRecords, err := s.reviewQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_REVIEWS")

		s.logger.Error("Failed to retrieve review active list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve review active list")
		span.End()

		status = "failed_find_active_reviews"

		return nil, nil, review_errors.ErrFailedFindActiveReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponseDeleteAt(Reviews), totalRecords, nil
}

func (s *reviewQueryService) FindByTrashed(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashed")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Reviews, totalRecords, err := s.reviewQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_REVIEWS")

		s.logger.Error("Failed to retrieve review trashed list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve review trashed list")
		span.End()

		status = "failed_find_trashed_reviews"

		return nil, nil, review_errors.ErrFailedFindTrashedReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsResponseDeleteAt(Reviews), totalRecords, nil
}

func (s *reviewQueryService) FindByProduct(req *requests.FindAllReviewByProduct) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByProduct", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByProduct")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	reviews, totalRecords, err := s.reviewQueryRepository.FindByProduct(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BY_PRODUCT_REVIEWS")

		s.logger.Error("Failed to retrieve review product list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve review product list")
		span.End()

		status = "failed_find_by_product_reviews"

		return nil, nil, review_errors.ErrFailedFindByProductReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsDetailResponse(reviews), totalRecords, nil
}

func (s *reviewQueryService) FindByMerchant(req *requests.FindAllReviewByMerchant) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByMerchant")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching Reviews",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	reviews, totalRecords, err := s.reviewQueryRepository.FindByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BY_MERCHANT_REVIEWS")

		s.logger.Error("Failed to retrieve review merchant list",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve review merchant list")
		span.End()

		status = "failed_find_by_merchant_reviews"

		return nil, nil, review_errors.ErrFailedFindByMerchantReviews
	}

	s.logger.Debug("Successfully fetched Reviews",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewsDetailResponse(reviews), totalRecords, nil
}

func (s *reviewQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
