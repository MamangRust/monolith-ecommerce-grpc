package service

import (
	"context"
	"os"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-product/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/utils"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type productCommandService struct {
	observability      observability.TraceLoggerObservability
	cache              cache.ProductCommandCache
	categoryRepository repository.CategoryQueryRepository
	merchantRepository repository.MerchantQueryRepository
	productQueryRepo   repository.ProductQueryRepository
	productRepository  repository.ProductCommandRepository
	logger             logger.LoggerInterface
}

type ProductCommandServiceDeps struct {
	Observability      observability.TraceLoggerObservability
	Cache              cache.ProductCommandCache
	CategoryRepository repository.CategoryQueryRepository
	MerchantRepository repository.MerchantQueryRepository
	ProductQueryRepo   repository.ProductQueryRepository
	ProductRepository  repository.ProductCommandRepository
	Logger             logger.LoggerInterface
}

func NewProductCommandService(deps *ProductCommandServiceDeps) ProductCommandService {
	return &productCommandService{
		observability:      deps.Observability,
		cache:              deps.Cache,
		categoryRepository: deps.CategoryRepository,
		merchantRepository: deps.MerchantRepository,
		productQueryRepo:   deps.ProductQueryRepo,
		productRepository:  deps.ProductRepository,
		logger:             deps.Logger,
	}
}

func (s *productCommandService) Create(ctx context.Context, req *requests.CreateProductRequest) (*db.CreateProductRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("categoryID", req.CategoryID),
		attribute.Int("merchantID", req.MerchantID),
		attribute.String("name", req.Name))

	defer func() {
		end(status)
	}()

	_, err := s.categoryRepository.FindByID(ctx, req.CategoryID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateProductRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("categoryID", req.CategoryID),
		)
	}


	_, err = s.merchantRepository.FindByID(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateProductRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchantID", req.MerchantID),
		)
	}


	slug := utils.GenerateSlug(req.Name)
	req.SlugProduct = &slug

	product, err := s.productRepository.Create(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateProductRow](
			s.logger,
			err,
			method,
			span,
		)
	}


	s.cache.DeleteCachedProduct(ctx, int(product.ProductID))

	logSuccess("Successfully created product",
		zap.Int("productID", int(product.ProductID)),
		zap.String("slug", slug))

	return product, nil
}

func (s *productCommandService) Update(ctx context.Context, req *requests.UpdateProductRequest) (*db.UpdateProductRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("productID", *req.ProductID),
		attribute.Int("categoryID", req.CategoryID),
		attribute.Int("merchantID", req.MerchantID),
		attribute.String("name", req.Name))

	defer func() {
		end(status)
	}()

	_, err := s.categoryRepository.FindByID(ctx, req.CategoryID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateProductRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("categoryID", req.CategoryID),
		)
	}


	_, err = s.merchantRepository.FindByID(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateProductRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchantID", req.MerchantID),
		)
	}


	slug := utils.GenerateSlug(req.Name)
	req.SlugProduct = &slug

	product, err := s.productRepository.Update(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateProductRow](
			s.logger,
			err,
			method,
			span,
		)
	}


	s.cache.DeleteCachedProduct(ctx, int(product.ProductID))

	logSuccess("Successfully updated product",
		zap.Int("productID", int(product.ProductID)),
		zap.String("slug", slug))

	return product, nil
}

func (s *productCommandService) UpdateProductCountStock(ctx context.Context, productID int, stock int) (*db.UpdateProductCountStockRow, error) {
	const method = "UpdateProductCountStock"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("product_id", productID),
		attribute.Int("stock", stock))

	defer func() {
		end(status)
	}()

	product, err := s.productRepository.UpdateProductCountStock(ctx, productID, stock)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateProductCountStockRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("product_id", productID),
			zap.Int("stock", stock),
		)
	}


	s.cache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully updated product stock",
		zap.Int("product_id", productID),
		zap.Int("new_stock", stock))

	return product, nil
}

func (s *productCommandService) Trash(ctx context.Context, productID int) (interface{}, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("productID", productID))

	defer func() {
		end(status)
	}()

	product, err := s.productRepository.Trash(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Product](
			s.logger,
			err,
			method,
			span,
			zap.Int("product_id", productID),
		)
	}


	s.cache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully trashed product",
		zap.Int("productID", productID))

	return product, nil
}

func (s *productCommandService) Restore(ctx context.Context, productID int) (interface{}, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("productID", productID))

	defer func() {
		end(status)
	}()

	product, err := s.productRepository.Restore(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Product](
			s.logger,
			err,
			method,
			span,
			zap.Int("product_id", productID),
		)
	}


	s.cache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully restored product",
		zap.Int("productID", productID))

	return product, nil
}

func (s *productCommandService) DeletePermanent(ctx context.Context, productID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("productID", productID))

	defer func() {
		end(status)
	}()

	product, err := s.productQueryRepo.FindByIDTrashed(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Int("product_id", productID),
		)
	}


	if product.ImageProduct != nil && *product.ImageProduct != "" {
		if err := os.Remove(*product.ImageProduct); err != nil {
			if !os.IsNotExist(err) {
				status = "error"
				return errorhandler.HandleError[bool](
					s.logger,
					product_errors.ErrFailedDeleteImageProduct,
					method,
					span,
					zap.String("image_path", *product.ImageProduct),
				)
			}
		}
	}

	_, err = s.productRepository.DeletePermanent(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Int("product_id", productID),
		)
	}


	s.cache.DeleteCachedProduct(ctx, productID)

	logSuccess("Successfully permanently deleted product", zap.Int("productID", productID))

	return true, nil
}

func (s *productCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.productRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}


	logSuccess("Successfully restored all trashed products")

	return success, nil
}

func (s *productCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.productRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			product_errors.ErrFailedDeleteAllProductsPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed products")

	return success, nil
}
