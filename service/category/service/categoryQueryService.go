package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"

	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type categoryQueryService struct {
	observability           observability.TraceLoggerObservability
	cache                   cache.CategoryQueryCache
	categoryQueryRepository repository.CategoryQueryRepository
	logger                  logger.LoggerInterface
}

type CategoryQueryServiceDeps struct {
	Observability           observability.TraceLoggerObservability
	Cache                   cache.CategoryQueryCache
	CategoryQueryRepository repository.CategoryQueryRepository
	Logger                  logger.LoggerInterface
}

func NewCategoryQueryService(
	deps *CategoryQueryServiceDeps) *categoryQueryService {

	return &categoryQueryService{
		cache:                   deps.Cache,
		categoryQueryRepository: deps.CategoryQueryRepository,
		logger:                  deps.Logger,
		observability:           deps.Observability,
	}
}

func (s *categoryQueryService) FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedCategoriesCache(ctx, req); found {
		logSuccess("Successfully retrieved all category records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	categories, err := s.categoryQueryRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetCategoriesRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(categories) > 0 {
		totalCount = int(categories[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedCategoriesCache(ctx, req, categories, &totalCount)

	logSuccess("Successfully fetched all categories",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return categories, &totalCount, nil
}

func (s *categoryQueryService) FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedCategoryActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active category records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	categories, err := s.categoryQueryRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetCategoriesActiveRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(categories) > 0 {
		totalCount = int(categories[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedCategoryActiveCache(ctx, req, categories, &totalCount)

	logSuccess("Successfully fetched active categories",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return categories, &totalCount, nil
}

func (s *categoryQueryService) FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedCategoryTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed category records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	categories, err := s.categoryQueryRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetCategoriesTrashedRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(categories) > 0 {
		totalCount = int(categories[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedCategoryTrashedCache(ctx, req, categories, &totalCount)

	logSuccess("Successfully fetched trashed categories",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return categories, &totalCount, nil
}

func (s *categoryQueryService) FindByID(ctx context.Context, category_id int) (*db.GetCategoryByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("category_id", category_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedCategoryCache(ctx, category_id); found {
		logSuccess("Successfully retrieved category from cache", zap.Int("category_id", category_id))
		return data, nil
	}

	category, err := s.categoryQueryRepository.FindByID(ctx, category_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetCategoryByIDRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("category_id", category_id),
		)
	}

	s.cache.SetCachedCategoryCache(ctx, category)

	logSuccess("Successfully fetched category", zap.Int("category_id", category_id))
	return category, nil
}
