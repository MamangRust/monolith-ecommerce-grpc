package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-product/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"

	"go.uber.org/zap"
)


type productQueryService struct {
	observability     observability.TraceLoggerObservability
	cache             cache.ProductQueryCache
	productRepository repository.ProductQueryRepository
	logger            logger.LoggerInterface
}

type ProductQueryServiceDeps struct {
	Observability     observability.TraceLoggerObservability
	Cache             cache.ProductQueryCache
	ProductRepository repository.ProductQueryRepository
	Logger            logger.LoggerInterface
}

func NewProductQueryService(deps *ProductQueryServiceDeps) ProductQueryService {
	return &productQueryService{
		observability:     deps.Observability,
		cache:             deps.Cache,
		productRepository: deps.ProductRepository,
		logger:            deps.Logger,
	}
}

func (s *productQueryService) FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, *int, error) {
	const method = "FindAll"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProducts(ctx, req); found {
		logSuccess("Successfully retrieved all product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	products, err := s.productRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsRow](
			s.logger,
			err,
			method,
			span,
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}


	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProducts(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched all products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return products, &totalCount, nil
}

func (s *productQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, *int, error) {
	const method = "FindByMerchant"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	merchantId := req.MerchantID

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.Int("merchant_id", merchantId))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProductsByMerchant(ctx, req); found {
		logSuccess("Successfully retrieved merchant product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("merchant_id", merchantId))
		return data, total, nil
	}

	products, err := s.productRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.Int("merchant_id", merchantId),
		)
	}


	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProductsByMerchant(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched merchant products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.Int("merchant_id", merchantId))

	return products, &totalCount, nil
}

func (s *productQueryService) FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, *int, error) {
	const method = "FindByCategory"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	category_name := req.CategoryName

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.String("category_name", category_name))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProductsByCategory(ctx, req); found {
		logSuccess("Successfully retrieved category product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("category_name", category_name))
		return data, total, nil
	}

	products, err := s.productRepository.FindByCategory(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsByCategoryNameRow](
			s.logger,
			err,
			method,
			span,
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("category_name", category_name),
		)
	}


	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProductsByCategory(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched category products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("category_name", category_name))

	return products, &totalCount, nil
}

func (s *productQueryService) FindByID(ctx context.Context, productID int) (*db.GetProductByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("productID", productID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedProduct(ctx, productID); found {
		logSuccess("Successfully retrieved product by ID from cache",
			zap.Int("productID", productID))
		return data, nil
	}

	product, err := s.productRepository.FindByID(ctx, productID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetProductByIDRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("productID", productID),
		)
	}


	s.cache.SetCachedProduct(ctx, product)

	logSuccess("Successfully fetched product by ID",
		zap.Int("productID", productID))

	return product, nil
}

func (s *productQueryService) FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, *int, error) {
	const method = "FindActive"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProductActive(ctx, req); found {
		logSuccess("Successfully retrieved active product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	products, err := s.productRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsActiveRow](
			s.logger,
			err,
			method,
			span,
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}


	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProductActive(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched active products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return products, &totalCount, nil
}

func (s *productQueryService) FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, *int, error) {
	const method = "FindTrashed"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedProductTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed product records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	products, err := s.productRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetProductsTrashedRow](
			s.logger,
			err,
			method,
			span,
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}


	var totalCount int
	if len(products) > 0 {
		totalCount = int(products[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedProductTrashed(ctx, req, products, &totalCount)

	logSuccess("Successfully fetched trashed products",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return products, &totalCount, nil
}
