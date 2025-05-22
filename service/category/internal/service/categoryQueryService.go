package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
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
	trace                   trace.Tracer
	categoryQueryRepository repository.CategoryQueryRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.CategoryResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCategoryQueryService(ctx context.Context, categoryQueryRepository repository.CategoryQueryRepository, logger logger.LoggerInterface, mapping response_service.CategoryResponseMapper) *categoryQueryService {
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
		trace:                   otel.Tracer("category-query-service"),
		categoryQueryRepository: categoryQueryRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *categoryQueryService) FindAll(req *requests.FindAllCategory) ([]*response.CategoryResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, startTime)
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

	s.logger.Debug("Fetching all categories",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	category, totalRecords, err := s.categoryQueryRepository.FindAllCategory(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_CATEGORIES")

		s.logger.Error("Failed to fetch categories",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch categories")

		status = "failed_find_all_categories"

		return nil, nil, category_errors.ErrFailedFindAllCategories
	}

	categoriesResponse := s.mapping.ToCategorysResponse(category)

	s.logger.Debug("Successfully fetched category",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return categoriesResponse, totalRecords, nil
}

func (s *categoryQueryService) FindByActive(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, startTime)
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

	s.logger.Debug("Fetching categories active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	category, totalRecords, err := s.categoryQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_CATEGORIES")

		s.logger.Error("Failed to fetch categories",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch categories")

		status = "failed_find_active_categories"

		return nil, nil, category_errors.ErrFailedFindActiveCategories
	}

	s.logger.Debug("Successfully fetched categories",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToCategorysResponseDeleteAt(category), totalRecords, nil
}

func (s *categoryQueryService) FindByTrashed(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, startTime)
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

	s.logger.Debug("Fetching categories trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	categories, totalRecords, err := s.categoryQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_CATEGORIES")

		s.logger.Error("Failed to fetch categories",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch categories")

		status = "failed_find_trashed_categories"

		return nil, nil, category_errors.ErrFailedFindTrashedCategories
	}

	s.logger.Debug("Successfully fetched categories",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToCategorysResponseDeleteAt(categories), totalRecords, nil
}

func (s *categoryQueryService) FindById(category_id int) (*response.CategoryResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	s.logger.Debug("Fetching category by ID", zap.Int("category_id", category_id))

	category, err := s.categoryQueryRepository.FindById(category_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CATEGORY_BY_ID")

		s.logger.Error("Failed to fetch category by ID",
			zap.Error(err),
			zap.Int("category_id", category_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch category by ID")

		status = "failed_find_category_by_id"

		return nil, category_errors.ErrFailedFindCategoryById
	}

	return s.mapping.ToCategoryResponse(category), nil
}

func (s *categoryQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
