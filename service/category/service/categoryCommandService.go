package service

import (
	"context"
	"os"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/utils"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type categoryCommandService struct {
	observability             observability.TraceLoggerObservability
	cache                     cache.CategoryCommandCache
	categoryQueryRepository   repository.CategoryQueryRepository
	categoryCommandRepository repository.CategoryCommandRepository
	logger                    logger.LoggerInterface
}

type CategoryCommandServiceDeps struct {
	Observability             observability.TraceLoggerObservability
	Cache                     cache.CategoryCommandCache
	CategoryQueryRepository   repository.CategoryQueryRepository
	CategoryCommandRepository repository.CategoryCommandRepository
	Logger                    logger.LoggerInterface
}

func NewCategoryCommandService(
	deps *CategoryCommandServiceDeps) *categoryCommandService {

	return &categoryCommandService{
		cache:                     deps.Cache,
		categoryCommandRepository: deps.CategoryCommandRepository,
		categoryQueryRepository:   deps.CategoryQueryRepository,
		logger:                    deps.Logger,
		observability:             deps.Observability,
	}
}

func (s *categoryCommandService) Create(ctx context.Context, req *requests.CreateCategoryRequest) (*db.CreateCategoryRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	if req.SlugCategory == nil || *req.SlugCategory == "" {
		generatedSlug := utils.GenerateSlug(req.Name)
		req.SlugCategory = &generatedSlug
	}

	category, err := s.categoryCommandRepository.Create(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateCategoryRow](
			s.logger,
			err,

			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteCachedCategoryCache(ctx, int(category.CategoryID))

	logSuccess("Successfully created category", zap.Int("categoryID", int(category.CategoryID)))
	return category, nil
}

func (s *categoryCommandService) Update(ctx context.Context, req *requests.UpdateCategoryRequest) (*db.UpdateCategoryRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("categoryID", *req.CategoryID))

	defer func() {
		end(status)
	}()

	if req.SlugCategory == nil || *req.SlugCategory == "" {
		generatedSlug := utils.GenerateSlug(req.Name)
		req.SlugCategory = &generatedSlug
	}

	category, err := s.categoryCommandRepository.Update(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateCategoryRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("category_id", *req.CategoryID),
		)
	}

	s.cache.DeleteCachedCategoryCache(ctx, *req.CategoryID)

	logSuccess("Successfully updated category", zap.Int("categoryID", *req.CategoryID))
	return category, nil
}

func (s *categoryCommandService) Trash(ctx context.Context, categoryID int) (*db.Category, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("categoryID", categoryID))

	defer func() {
		end(status)
	}()

	category, err := s.categoryCommandRepository.Trash(ctx, categoryID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Category](
			s.logger,
			err,

			method,
			span,

			zap.Int("category_id", categoryID),
		)
	}

	s.cache.DeleteCachedCategoryCache(ctx, categoryID)

	logSuccess("Successfully trashed category", zap.Int("categoryID", categoryID))
	return category, nil
}

func (s *categoryCommandService) Restore(ctx context.Context, categoryID int) (*db.Category, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("categoryID", categoryID))

	defer func() {
		end(status)
	}()

	category, err := s.categoryCommandRepository.Restore(ctx, categoryID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Category](
			s.logger,
			err,

			method,
			span,

			zap.Int("category_id", categoryID),
		)
	}

	s.cache.DeleteCachedCategoryCache(ctx, categoryID)

	logSuccess("Successfully restored category", zap.Int("categoryID", categoryID))
	return category, nil
}

func (s *categoryCommandService) DeletePermanent(ctx context.Context, categoryID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("categoryID", categoryID))

	defer func() {
		end(status)
	}()

	category, err := s.categoryQueryRepository.FindByIDTrashed(ctx, categoryID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,

			method,
			span,

			zap.Int("category_id", categoryID),
		)
	}

	if category.ImageCategory != nil && *category.ImageCategory != "" {
		err = os.Remove(*category.ImageCategory)
		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Debug(
					"Category image file not found, continuing with category deletion",
					zap.String("image_path", *category.ImageCategory),
				)
			} else {
				status = "error"
				return errorhandler.HandleError[bool](
					s.logger,
					category_errors.ErrFailedRemoveImageCategory,
					method,
					span,

					zap.String("image_path", *category.ImageCategory),
				)
			}
		} else {
			s.logger.Debug(
				"Successfully deleted category image",
				zap.String("image_path", *category.ImageCategory),
			)
		}
	}

	success, err := s.categoryCommandRepository.DeletePermanent(ctx, categoryID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,

			method,
			span,

			zap.Int("category_id", categoryID),
		)
	}

	s.cache.DeleteCachedCategoryCache(ctx, categoryID)

	logSuccess("Successfully deleted category permanently", zap.Int("categoryID", categoryID))
	return success, nil
}

func (s *categoryCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.categoryCommandRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,

			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed categories")
	return success, nil
}

func (s *categoryCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.categoryCommandRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,

			method,
			span,
		)
	}

	logSuccess("Successfully deleted all trashed categories permanently")
	return success, nil
}
