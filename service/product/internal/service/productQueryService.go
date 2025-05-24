package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type productQueryService struct {
	ctx                    context.Context
	trace                  trace.Tracer
	productQueryRepository repository.ProductQueryRepository
	mapping                response_service.ProductResponseMapper
	logger                 logger.LoggerInterface
	requestCounter         *prometheus.CounterVec
	requestDuration        *prometheus.HistogramVec
}

func NewProductQueryService(
	ctx context.Context,
	productQueryRepository repository.ProductQueryRepository,
	mapping response_service.ProductResponseMapper,
	logger logger.LoggerInterface,
) *productQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_query_service_requests_total",
			Help: "Total number of requests to the ProductQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the ProductQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &productQueryService{
		ctx:                    ctx,
		trace:                  otel.Tracer("product-query-service"),
		productQueryRepository: productQueryRepository,
		mapping:                mapping,
		logger:                 logger,
		requestCounter:         requestCounter,
		requestDuration:        requestDuration,
	}
}

func (s *productQueryService) FindAll(req *requests.FindAllProduct) ([]*response.ProductResponse, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all products",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productQueryRepository.FindAllProducts(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_PRODUCTS")

		s.logger.Error("Failed to retrieve product list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("traceID", traceID),
		)

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve product list")

		status = "failed_find_all_products"

		return nil, nil, product_errors.ErrFailedFindAllProducts
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponse(products), totalRecords, nil
}

func (s *productQueryService) FindByMerchant(req *requests.FindAllProductByMerchant) ([]*response.ProductResponse, *int, *response.ErrorResponse) {
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
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.Int("merchant_id", merchantId),
	)

	s.logger.Debug("Fetching all products by merchant",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search),
		zap.Int("merchant_id", merchantId),
	)

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productQueryRepository.FindByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCTS_BY_MERCHANT")

		s.logger.Error("Failed to retrieve product list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("merchant_id", merchantId),
			zap.String("traceID", traceID),
		)

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve product list")

		status = "failed_find_products_by_merchant"

		return nil, nil, product_errors.ErrFailedFindProductsByMerchant
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponse(products), totalRecords, nil
}

func (s *productQueryService) FindByCategory(req *requests.FindAllProductByCategory) ([]*response.ProductResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByCategory", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByCategory")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	category_name := req.CategoryName

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.String("category_name", category_name),
	)

	s.logger.Debug("Fetching all products by category name",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search),
		zap.String("category_name", category_name),
	)

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productQueryRepository.FindByCategory(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCTS_BY_CATEGORY")

		s.logger.Error("Failed to retrieve product list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("category_name", category_name),
			zap.String("traceID", traceID),
		)

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve product list")

		status = "failed_find_products_by_category"

		return nil, nil, product_errors.ErrFailedFindProductsByCategory
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponse(products), totalRecords, nil
}

func (s *productQueryService) FindById(productID int) (*response.ProductResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("productID", productID),
	)

	s.logger.Debug("Fetching product by ID", zap.Int("productID", productID))

	product, err := s.productQueryRepository.FindById(productID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCT_BY_ID")

		s.logger.Error("Failed to retrieve product",
			zap.Int("product_id", productID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve product")

		status = "failed_find_product_by_id"

		return nil, product_errors.ErrFailedFindProductById
	}

	return s.mapping.ToProductResponse(product), nil
}

func (s *productQueryService) FindByActive(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all products active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCTS_BY_ACTIVE")

		s.logger.Error("Failed to retrieve product list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("traceID", traceID),
		)

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve product list")

		status = "failed_find_products_by_active"

		return nil, nil, product_errors.ErrFailedFindProductsByActive
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponseDeleteAt(products), totalRecords, nil
}

func (s *productQueryService) FindByTrashed(req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all products trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	products, totalRecords, err := s.productQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCTS_BY_TRASHED")

		s.logger.Error("Failed to retrieve product list",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("traceID", traceID),
		)

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve product list")

		status = "failed_find_products_by_trashed"

		return nil, nil, product_errors.ErrFailedFindProductsByTrashed
	}

	s.logger.Debug("Successfully fetched products",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToProductsResponseDeleteAt(products), totalRecords, nil
}

func (s *productQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
