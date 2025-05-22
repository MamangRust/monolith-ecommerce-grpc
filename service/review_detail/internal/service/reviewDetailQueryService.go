package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewDetailQueryService struct {
	ctx                         context.Context
	trace                       trace.Tracer
	reviewDetailQueryRepository repository.ReviewDetailQueryRepository
	mapping                     response_service.ReviewDetailResponeMapper
	logger                      logger.LoggerInterface
	requestCounter              *prometheus.CounterVec
	requestDuration             *prometheus.HistogramVec
}

func NewReviewDetailQueryService(
	ctx context.Context,
	reviewDetailQueryRepository repository.ReviewDetailQueryRepository,
	mapping response_service.ReviewDetailResponeMapper,
	logger logger.LoggerInterface,
) *reviewDetailQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "review_detail_query_service_request_count",
			Help: "Total number of requests to the ReviewDetailQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "review_detail_query_service_request_duration",
			Help:    "Histogram of request durations for the ReviewDetailQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &reviewDetailQueryService{
		ctx:                         ctx,
		trace:                       otel.Tracer("review-detail-command-service"),
		reviewDetailQueryRepository: reviewDetailQueryRepository,
		mapping:                     mapping,
		logger:                      logger,
		requestCounter:              requestCounter,
		requestDuration:             requestDuration,
	}
}

func (s *reviewDetailQueryService) FindAll(req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAll")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all Review Details",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.reviewDetailQueryRepository.FindAllReviews(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_REVIEW_DETAIL")

		s.logger.Error("Failed to retrieve Review Details",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve Review Details")
		status = "failed_find_all_review_detail"

		return nil, nil, reviewdetail_errors.ErrFailedFindAllReview
	}

	s.logger.Debug("Successfully fetched Review Details",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToReviewsDetailsResponse(res), totalRecords, nil
}

func (s *reviewDetailQueryService) FindByActive(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all Review Details active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.reviewDetailQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_REVIEW_DETAIL")

		s.logger.Error("Failed to retrieve active Review Details",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to retrieve active Review Details")

		status = "failed_find_active_review_detail"

		return nil, nil, reviewdetail_errors.ErrFailedFindActiveReview
	}

	s.logger.Debug("Successfully fetched active Review Detail",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewDetailsResponseDeleteAt(res), totalRecords, nil
}

func (s *reviewDetailQueryService) FindByTrashed(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all Review Details trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.reviewDetailQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_REVIEW_DETAIL")

		s.logger.Error("Failed to retrieve trashed Review Details",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve trashed Review Details")
		status = "failed_find_trashed_review_detail"

		return nil, nil, reviewdetail_errors.ErrFailedFindTrashedReview
	}

	s.logger.Debug("Successfully fetched trashed Review Detail",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToReviewDetailsResponseDeleteAt(res), totalRecords, nil
}

func (s *reviewDetailQueryService) FindById(review_id int) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	s.logger.Debug("Fetching Review Detail by ID", zap.Int("Review DetailID", review_id))

	res, err := s.reviewDetailQueryRepository.FindById(review_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_REVIEW_DETAIL_BY_ID")

		s.logger.Error("Failed to retrieve Review Detail by ID",
			zap.Error(err),
			zap.Int("Review DetailID", review_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve Review Detail by ID")
		status = "failed_find_review_detail_by_id"

		return nil, reviewdetail_errors.ErrReviewDetailNotFoundRes
	}

	return s.mapping.ToReviewDetailsResponse(res), nil
}

func (s *reviewDetailQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
