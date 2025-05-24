package service

import (
	"context"
	"os"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-grpc-product/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	traceunic "github.com/MamangRust/monolith-point-of-sale-pkg/trace_unic"
	"github.com/MamangRust/monolith-point-of-sale-pkg/utils"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	merchant_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type productCommandService struct {
	ctx                      context.Context
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
	ctx context.Context,
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &productCommandService{
		ctx:                      ctx,
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

func (s *productCommandService) CreateProduct(req *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateProduct", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateProduct")
	defer span.End()

	s.logger.Debug("Creating new product")

	_, err := s.categoryRepository.FindById(req.CategoryID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CATEGORY_BY_ID")

		s.logger.Error("Category not found for product creation",
			zap.Int("categoryID", req.CategoryID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Category not found for product creation")

		status = "failed_find_category_by_id"

		return nil, category_errors.ErrFailedFindCategoryById
	}

	_, err = s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_BY_ID")

		s.logger.Error("Merchant not found for product creation",
			zap.Int("merchantID", req.MerchantID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Merchant not found for product creation")

		status = "failed_find_merchant_by_id"

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	slug := utils.GenerateSlug(req.Name)

	req.SlugProduct = &slug

	product, err := s.productCommandRepository.CreateProduct(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_PRODUCT")

		s.logger.Error("Failed to create product",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create product")

		status = "failed_create_product"

		return nil, product_errors.ErrFailedCreateProduct
	}

	s.logger.Debug("Product created successfully", zap.Int("productID", product.ID))

	return s.mapping.ToProductResponse(product), nil
}

func (s *productCommandService) UpdateProduct(req *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateProduct", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateProduct")
	defer span.End()

	span.SetAttributes(
		attribute.Int("productID", *req.ProductID),
	)

	s.logger.Debug("Updating product", zap.Int("productID", *req.ProductID))

	_, err := s.categoryRepository.FindById(req.CategoryID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CATEGORY_BY_ID")

		s.logger.Error("Category not found for product update",
			zap.Int("categoryID", req.CategoryID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Category not found for product update")

		status = "failed_find_category_by_id"

		return nil, category_errors.ErrFailedFindCategoryById
	}

	_, err = s.merchantRepository.FindById(req.MerchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_BY_ID")

		s.logger.Error("Merchant not found for product update",
			zap.Int("merchantID", req.MerchantID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Merchant not found for product update")

		status = "failed_find_merchant_by_id"

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	slug := utils.GenerateSlug(req.Name)

	req.SlugProduct = &slug

	product, err := s.productCommandRepository.UpdateProduct(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_PRODUCT")

		s.logger.Error("Failed to update product",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update product")

		status = "failed_update_product"

		return nil, product_errors.ErrFailedUpdateProduct
	}

	s.logger.Debug("Product updated successfully", zap.Int("productID", product.ID))
	return s.mapping.ToProductResponse(product), nil
}

func (s *productCommandService) TrashProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashProduct", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashProduct")
	defer span.End()

	span.SetAttributes(
		attribute.Int("productID", productID),
	)

	s.logger.Debug("Trashing product", zap.Int("productID", productID))

	product, err := s.productCommandRepository.TrashedProduct(productID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_PRODUCT")

		s.logger.Error("Failed to trash product",
			zap.Int("product_id", productID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash product")

		status = "failed_trash_product"

		return nil, product_errors.ErrFailedTrashProduct
	}

	return s.mapping.ToProductResponseDeleteAt(product), nil
}

func (s *productCommandService) RestoreProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreProduct", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreProduct")
	defer span.End()

	span.SetAttributes(
		attribute.Int("productID", productID),
	)

	s.logger.Debug("Restoring product", zap.Int("productID", productID))

	product, err := s.productCommandRepository.RestoreProduct(productID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_PRODUCT")

		s.logger.Error("Failed to restore product",
			zap.Int("product_id", productID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore product")

		status = "failed_restore_product"

		return nil, product_errors.ErrFailedRestoreProduct
	}

	return s.mapping.ToProductResponseDeleteAt(product), nil
}

func (s *productCommandService) DeleteProductPermanent(productID int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteProductPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteProductPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("productID", productID),
	)

	s.logger.Debug("Permanently deleting product", zap.Int("productID", productID))

	res, err := s.productQueryRepository.FindByIdTrashed(productID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCT_TRASHED_BY_ID")

		s.logger.Error("Failed to find product trashed by id",
			zap.Int("product_id", productID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find product trashed by id")

		status = "failed_find_product_trashed_by_id"

		return false, product_errors.ErrFailedFindProductTrashedById
	}

	if res.ImageProduct != "" {
		err := os.Remove(res.ImageProduct)
		if err != nil {
			if os.IsNotExist(err) {
				traceID := traceunic.GenerateTraceID("FAILED_DELETE_IMAGE_PRODUCT")

				s.logger.Error("Failed to delete product image",
					zap.String("image_path", res.ImageProduct),
					zap.Error(err),
					zap.String("traceID", traceID))

				span.SetAttributes(
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)
				span.SetStatus(codes.Error, "Failed to delete product image")

				status = "failed_delete_image_product"

			} else {
				traceID := traceunic.GenerateTraceID("FAILED_DELETE_IMAGE_PRODUCT")

				s.logger.Error("Failed to delete product image",
					zap.String("image_path", res.ImageProduct),
					zap.Error(err),
					zap.String("traceID", traceID))

				span.SetAttributes(
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)
				span.SetStatus(codes.Error, "Failed to delete product image")

				status = "failed_delete_image_product"

				return false, product_errors.ErrFailedDeleteImageProduct
			}
		} else {
			s.logger.Debug("Successfully deleted product image",
				zap.String("image_path", res.ImageProduct))
		}
	}

	_, err = s.productCommandRepository.DeleteProductPermanent(productID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_PRODUCT_PERMANENT")

		s.logger.Error("Failed to permanently delete product",
			zap.Int("product_id", productID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete product")

		status = "failed_delete_product_permanent"

		return false, product_errors.ErrFailedDeleteProductPermanent
	}

	s.logger.Debug("Product permanently deleted successfully",
		zap.Int("product_id", productID))

	return true, nil
}

func (s *productCommandService) RestoreAllProducts() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllProducts", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllProducts")
	defer span.End()

	s.logger.Debug("Restoring all trashed products")

	success, err := s.productCommandRepository.RestoreAllProducts()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_PRODUCTS")

		s.logger.Error("Failed to restore all trashed products",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all trashed products")

		status = "failed_restore_all_products"

		return false, product_errors.ErrFailedRestoreAllProducts
	}

	s.logger.Debug("All trashed products restored successfully",
		zap.Bool("success", success))

	return success, nil
}

func (s *productCommandService) DeleteAllProductsPermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllProductsPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllProductsPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.String("method", "DeleteAllProductsPermanent"),
	)

	s.logger.Debug("Permanently deleting all products")

	success, err := s.productCommandRepository.DeleteAllProductPermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_PRODUCTS_PERMANENT")

		s.logger.Error("Failed to permanently delete all trashed products",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all trashed products")

		status = "failed_delete_all_products_permanent"

		return false, product_errors.ErrFailedDeleteAllProductsPermanent
	}

	s.logger.Debug("All trashed products permanently deleted successfully",
		zap.Bool("success", success))

	return success, nil
}

func (s *productCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
