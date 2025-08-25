package service

import (
	"context"
	"os"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-product/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-product/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/utils"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type productCommandService struct {
	errorhandler             errorhandler.ProductCommandError
	mencache                 mencache.ProductCommandCache
	trace                    trace.Tracer
	categoryRepository       repository.CategoryQueryRepository
	merchantRepository       repository.MerchantQueryRepository
	productQueryRepository   repository.ProductQueryRepository
	productCommandRepository repository.ProductCommandRepository
	mapping                  response_service.ProductResponseMapper
	logger                   logger.LoggerInterface
	requestCounter           *prometheus.CounterVec
	requestDuration          *prometheus.HistogramVec
}

func NewProductCommandService(
	errorhandler errorhandler.ProductCommandError,
	mencache mencache.ProductCommandCache,
	categoryRepository repository.CategoryQueryRepository,
	merchantRepository repository.MerchantQueryRepository,
	productQueryRepository repository.ProductQueryRepository,
	productCommandRepository repository.ProductCommandRepository,
	mapping response_service.ProductResponseMapper,
	logger logger.LoggerInterface,
) *productCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_command_service_requests_total",
			Help: "Total number of requests to the ProductCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the ProductCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &productCommandService{
		errorhandler:             errorhandler,
		mencache:                 mencache,
		trace:                    otel.Tracer("product-command-service"),
		categoryRepository:       categoryRepository,
		merchantRepository:       merchantRepository,
		productQueryRepository:   productQueryRepository,
		productCommandRepository: productCommandRepository,
		mapping:                  mapping,
		logger:                   logger,
		requestCounter:           requestCounter,
		requestDuration:          requestDuration,
	}
}

func (s *productCommandService) CreateProduct(ctx context.Context, req *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	const method = "CreateProduct"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("product.name", req.Name), attribute.Int("product.category_id", req.CategoryID), attribute.Int("product.merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	_, err := s.categoryRepository.FindById(ctx, req.CategoryID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.ProductResponse](s.logger, err, method, "FAILED_FIND_CATEGORY_BY_ID", span, &status, category_errors.ErrFailedFindCategoryById, zap.Error(err))
	}

	_, err = s.merchantRepository.FindById(ctx, req.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.ProductResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	slug := utils.GenerateSlug(req.Name)

	req.SlugProduct = &slug

	product, err := s.productCommandRepository.CreateProduct(ctx, req)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.ProductResponse](s.logger, err, method, "FAILED_CREATE_PRODUCT", span, &status, product_errors.ErrFailedCreateProduct, zap.Error(err))
	}

	so := s.mapping.ToProductResponse(product)

	logSuccess("Successfully created product", zap.Int("product.id", product.ID), zap.Bool("success", true))

	return so, nil
}

func (s *productCommandService) UpdateProduct(ctx context.Context, req *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	const method = "UpdateProduct"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.categoryRepository.FindById(ctx, req.CategoryID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.ProductResponse](s.logger, err, method, "FAILED_FIND_CATEGORY_BY_ID", span, &status, category_errors.ErrFailedFindCategoryById, zap.Error(err))
	}

	_, err = s.merchantRepository.FindById(ctx, req.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.ProductResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	slug := utils.GenerateSlug(req.Name)

	req.SlugProduct = &slug

	product, err := s.productCommandRepository.UpdateProduct(ctx, req)

	if err != nil {
		return s.errorhandler.HandleUpdateProductError(err, method, "FAILED_UPDATE_PRODUCT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToProductResponse(product)

	s.mencache.DeleteCachedProduct(ctx, *req.ProductID)

	logSuccess("Successfully updated product", zap.Int("product.id", *req.ProductID), zap.Bool("success", true))

	return so, nil
}

func (s *productCommandService) TrashProduct(ctx context.Context, productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashProduct"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	product, err := s.productCommandRepository.TrashedProduct(ctx, productID)

	if err != nil {
		return s.errorhandler.HandleTrashedProductError(err, method, "FAILED_TRASH_PRODUCT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToProductResponseDeleteAt(product)

	s.mencache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully trashed product", zap.Int("product.id", productID), zap.Bool("success", true))

	return so, nil
}

func (s *productCommandService) RestoreProduct(ctx context.Context, productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreProduct"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	product, err := s.productCommandRepository.RestoreProduct(ctx, productID)

	if err != nil {
		return s.errorhandler.HandleRestoreProductError(err, method, "FAILED_RESTORE_PRODUCT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToProductResponseDeleteAt(product)

	s.mencache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully restored product", zap.Int("product.id", productID), zap.Bool("success", true))

	return so, nil
}

func (s *productCommandService) DeleteProductPermanent(ctx context.Context, productID int) (bool, *response.ErrorResponse) {
	const method = "DeleteProductPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.productQueryRepository.FindByIdTrashed(ctx, productID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_FIND_PRODUCT_BY_ID_TRASHED", span, &status, product_errors.ErrFailedFindProductTrashedById, zap.Error(err))
	}

	if res.ImageProduct != "" {
		err := os.Remove(res.ImageProduct)

		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Error("Failed to delete product image",
					zap.String("image_path", res.ImageProduct),
					zap.Error(err))
			} else {
				return s.errorhandler.HandleFileError(err, method, "FAILED_DELETE_IMAGE_PRODUCT", res.ImageProduct, span, &status, zap.Error(err))
			}
		} else {
			s.logger.Debug("Successfully deleted category image",
				zap.String("image_path", res.ImageProduct))
		}
	}

	_, err = s.productCommandRepository.DeleteProductPermanent(ctx, productID)

	if err != nil {
		return s.errorhandler.HandleDeleteProductError(err, method, "FAILED_DELETE_PRODUCT_PERMANENT", span, &status, zap.Error(err))
	}

	msgSuccess := "Product deleted permanently successfully"

	logSuccess(msgSuccess, zap.Int("product.id", productID), zap.Bool("success", true))

	return true, nil
}

func (s *productCommandService) RestoreAllProducts(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllProducts"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.productCommandRepository.RestoreAllProducts(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllProductError(err, method, "FAILED_RESTORE_ALL_PRODUCTS", span, &status, zap.Error(err))
	}

	logSuccess("All trashed products restored successfully", zap.Bool("success", success))

	return success, nil
}

func (s *productCommandService) DeleteAllProductsPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllProductsPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.productCommandRepository.DeleteAllProductPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllProductError(err, method, "FAILED_DELETE_ALL_PRODUCTS_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("All trashed products deleted permanently successfully", zap.Bool("success", success))

	return success, nil
}

func (s *productCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *productCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
