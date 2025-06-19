package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-category/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type categoryQueryService struct {
	ctx                     context.Context
	errorhandler            errorhandler.CategoryQueryError
	mencache                mencache.CategoryQueryCache
	trace                   trace.Tracer
	categoryQueryRepository repository.CategoryQueryRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.CategoryResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCategoryQueryService(ctx context.Context,
	errorhandler errorhandler.CategoryQueryError,
	mencache mencache.CategoryQueryCache,
	categoryQueryRepository repository.CategoryQueryRepository, logger logger.LoggerInterface, mapping response_service.CategoryResponseMapper) *categoryQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "category_query_service_request_total",
			Help: "Total number of requests to the CategoryQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "category_query_service_request_duration_seconds",
			Help:    "Duration of requests to the CategoryQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &categoryQueryService{
		ctx:                     ctx,
		errorhandler:            errorhandler,
		mencache:                mencache,
		trace:                   otel.Tracer("category-query-service"),
		categoryQueryRepository: categoryQueryRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *categoryQueryService) FindAll(req *requests.FindAllCategory) ([]*response.CategoryResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCategoriesCache(req); found {
		logSuccess("Successfully fetched categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	category, totalRecords, err := s.categoryQueryRepository.FindAllCategory(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_CATEGORY", span, &status, zap.Error(err))
	}

	categoriesResponse := s.mapping.ToCategorysResponse(category)

	s.mencache.SetCachedCategoriesCache(req, categoriesResponse, totalRecords)

	logSuccess("Successfully fetched categories", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return categoriesResponse, totalRecords, nil
}

func (s *categoryQueryService) FindByActive(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCategoryActiveCache(req); found {
		logSuccess("Successfully fetched categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	category, totalRecords, err := s.categoryQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_ACTIVE_CATEGORY", span, &status, category_errors.ErrFailedFindActiveCategories, zap.Error(err))
	}

	so := s.mapping.ToCategorysResponseDeleteAt(category)

	s.mencache.SetCachedCategoryActiveCache(req, so, totalRecords)

	logSuccess("Successfully fetched categories", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *categoryQueryService) FindByTrashed(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCategoryTrashedCache(req); found {
		logSuccess("Successfully fetched categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	categories, totalRecords, err := s.categoryQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_TRASHED_CATEGORY", span, &status, category_errors.ErrFailedFindTrashedCategories, zap.Error(err))
	}

	so := s.mapping.ToCategorysResponseDeleteAt(categories)

	s.mencache.SetCachedCategoryTrashedCache(req, so, totalRecords)

	logSuccess("Successfully fetched categories", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *categoryQueryService) FindById(category_id int) (*response.CategoryResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("category.id", category_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedCategoryCache(category_id); found {
		logSuccess("Successfully fetched category from cache", zap.Int("category.id", category_id))

		return data, nil
	}

	category, err := s.categoryQueryRepository.FindById(category_id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CATEGORY_BY_ID", span, &status, category_errors.ErrFailedFindCategoryById, zap.Error(err))
	}

	so := s.mapping.ToCategoryResponse(category)

	s.mencache.SetCachedCategoryCache(so)

	logSuccess("Successfully fetched category", zap.Int("category.id", category_id))

	return so, nil
}

func (s *categoryQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Info("Start: " + method)

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

	return span, end, status, logSuccess
}

func (s *categoryQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *categoryQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
