package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	mencache                    mencache.ReviewDetailQueryCache
	errorhandler                errorhandler.ReviewDetailQueryError
	trace                       trace.Tracer
	reviewDetailQueryRepository repository.ReviewDetailQueryRepository
	mapping                     response_service.ReviewDetailResponeMapper
	logger                      logger.LoggerInterface
	requestCounter              *prometheus.CounterVec
	requestDuration             *prometheus.HistogramVec
}

func NewReviewDetailQueryService(
	mencache mencache.ReviewDetailQueryCache,
	errorhandler errorhandler.ReviewDetailQueryError,
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
		mencache:                    mencache,
		errorhandler:                errorhandler,
		trace:                       otel.Tracer("review-detail-command-service"),
		reviewDetailQueryRepository: reviewDetailQueryRepository,
		mapping:                     mapping,
		logger:                      logger,
		requestCounter:              requestCounter,
		requestDuration:             requestDuration,
	}
}

func (s *reviewDetailQueryService) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetReviewDetailAllCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	res, totalRecords, err := s.reviewDetailQueryRepository.FindAllReviews(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_TO_FIND_REVIEW_DETAILS", span, &status, zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Error(err))
	}

	so := s.mapping.ToReviewsDetailsResponse(res)

	s.mencache.SetReviewDetailAllCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched all Review Details", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("totalRecords", *totalRecords))

	return so, totalRecords, nil
}

func (s *reviewDetailQueryService) FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetReviewDetailActiveCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	res, totalRecords, err := s.reviewDetailQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_REVIEW_DETAILS", span, &status, reviewdetail_errors.ErrFailedFindActiveReview, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailsResponseDeleteAt(res)

	s.mencache.SetReviewDetailActiveCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active Review Details", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("totalRecords", *totalRecords))

	return so, totalRecords, nil
}

func (s *reviewDetailQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetReviewDetailTrashedCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	res, totalRecords, err := s.reviewDetailQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FAILED_TO_FIND_REVIEW_DETAILS", span, &status, reviewdetail_errors.ErrFailedFindTrashedReview, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailsResponseDeleteAt(res)

	s.mencache.SetReviewDetailTrashedCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed Review Details", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("totalRecords", *totalRecords))

	return so, totalRecords, nil
}

func (s *reviewDetailQueryService) FindById(ctx context.Context, review_id int) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("reviewDetail.id", review_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedReviewDetailCache(ctx, review_id); found {
		logSuccess("Data found in cache", zap.Int("reviewDetail.id", review_id))

		return data, nil
	}

	res, err := s.reviewDetailQueryRepository.FindById(ctx, review_id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_TO_FIND_REVIEW_DETAILS", span, &status, reviewdetail_errors.ErrReviewDetailNotFoundRes, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailsResponse(res)

	s.mencache.SetCachedReviewDetailCache(ctx, so)

	logSuccess("Successfully fetched Review Detail", zap.Int("reviewDetail.id", review_id))

	return so, nil
}

func (s *reviewDetailQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

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

func (s *reviewDetailQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *reviewDetailQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
