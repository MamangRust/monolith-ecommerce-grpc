package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-review/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	mencache              mencache.ReviewQueryCache
	errorhandler          errorhandler.ReviewQueryError
	trace                 trace.Tracer
	reviewQueryRepository repository.ReviewQueryRepository
	mapping               response_service.ReviewResponseMapper
	logger                logger.LoggerInterface
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewReviewQueryService(
	mencache mencache.ReviewQueryCache,
	errorhandler errorhandler.ReviewQueryError,
	reviewQueryRepository repository.ReviewQueryRepository, mapping response_service.ReviewResponseMapper, logger logger.LoggerInterface) *reviewQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "review_query_service_request_count",
			Help: "Total number of requests to the ReviewQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "review_query_service_request_duration",
			Help:    "Histogram of request durations for the ReviewQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &reviewQueryService{
		errorhandler:          errorhandler,
		mencache:              mencache,
		trace:                 otel.Tracer("review-query-service"),
		reviewQueryRepository: reviewQueryRepository,
		mapping:               mapping,
		logger:                logger,
		requestCounter:        requestCounter,
		requestDuration:       requestDuration,
	}
}

func (s *reviewQueryService) FindAllReviews(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewResponse, *int, *response.ErrorResponse) {
	const method = "FindAllReviews"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetReviewAllCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	Reviews, totalRecords, err := s.reviewQueryRepository.FindAllReview(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_TO_FIND_ALL_REVIEWS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewsResponse(Reviews)

	s.mencache.SetReviewAllCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched all reviews", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *reviewQueryService) FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetReviewActiveCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	Reviews, totalRecords, err := s.reviewQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_BY_ACTIVE", span, &status, review_errors.ErrFailedFindActiveReviews)
	}

	so := s.mapping.ToReviewsResponseDeleteAt(Reviews)

	s.mencache.SetReviewActiveCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active reviews", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *reviewQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetReviewTrashedCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	Reviews, totalRecords, err := s.reviewQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_BY_TRASHED", span, &status, review_errors.ErrFailedFindTrashedReviews)
	}

	so := s.mapping.ToReviewsResponseDeleteAt(Reviews)

	s.mencache.SetReviewTrashedCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed reviews", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *reviewQueryService) FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse) {
	const method = "FindByProduct"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	productId := req.ProductID
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.Int("productId", productId))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetReviewByProductCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	reviews, totalRecords, err := s.reviewQueryRepository.FindByProduct(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDetailError(err, method, "FAILED_TO_FIND_BY_PRODUCT", span, &status, review_errors.ErrFailedFindByProductReviews, zap.Error(err))
	}

	so := s.mapping.ToReviewsDetailResponse(reviews)
	s.mencache.SetReviewByProductCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched reviews by product", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *reviewQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse) {
	const method = "FindByMerchant"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetReviewByMerchantCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	reviews, totalRecords, err := s.reviewQueryRepository.FindByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDetailError(err, method, "FAILED_TO_FIND_BY_MERCHANT", span, &status, review_errors.ErrFailedFindByMerchantReviews, zap.Error(err))
	}

	so := s.mapping.ToReviewsDetailResponse(reviews)
	s.mencache.SetReviewByMerchantCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched reviews by merchant", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *reviewQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *reviewQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *reviewQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
