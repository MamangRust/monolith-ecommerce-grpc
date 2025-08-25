package service

import (
	"context"
	"os"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-category/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/utils"
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

type categoryCommandService struct {
	mencache                  mencache.CategoryCommandCache
	errorHandler              errorhandler.CategoryCommandError
	trace                     trace.Tracer
	categoryQueryRepository   repository.CategoryQueryRepository
	categoryCommandRepository repository.CategoryCommandRepository
	logger                    logger.LoggerInterface
	mapping                   response_service.CategoryResponseMapper
	requestCounter            *prometheus.CounterVec
	requestDuration           *prometheus.HistogramVec
}

func NewCategoryCommandService(
	mencache mencache.CategoryCommandCache,
	errorHandler errorhandler.CategoryCommandError,
	categoryCommandRepository repository.CategoryCommandRepository,
	categoryQueryRepository repository.CategoryQueryRepository,
	logger logger.LoggerInterface, mapping response_service.CategoryResponseMapper) *categoryCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "category_command_service_request_total",
			Help: "Total number of requests to the CategoryCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "category_command_service_request_duration_seconds",
			Help:    "Duration of requests to the CategoryCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &categoryCommandService{
		errorHandler:              errorHandler,
		mencache:                  mencache,
		trace:                     otel.Tracer("category-command-service"),
		categoryCommandRepository: categoryCommandRepository,
		categoryQueryRepository:   categoryQueryRepository,
		logger:                    logger,
		mapping:                   mapping,
		requestCounter:            requestCounter,
		requestDuration:           requestDuration,
	}
}

func (s *categoryCommandService) CreateCategory(ctx context.Context, req *requests.CreateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse) {
	const method = "CreateCategory"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("name", req.Name))

	defer func() {
		end(status)
	}()

	slug := utils.GenerateSlug(req.Name)

	req.SlugCategory = &slug

	cashier, err := s.categoryCommandRepository.CreateCategory(ctx, req)

	if err != nil {
		return s.errorHandler.HandleCreateCategoryError(err, method, "FAILED_CREATE_CATEGORY", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryResponse(cashier)

	logSuccess("Successfully created category", zap.Int("category.id", cashier.ID))

	return so, nil
}

func (s *categoryCommandService) UpdateCategory(ctx context.Context, req *requests.UpdateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse) {
	const method = "UpdateCategory"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("name", req.Name))

	defer func() {
		end(status)
	}()

	slug := utils.GenerateSlug(req.Name)

	req.SlugCategory = &slug

	category, err := s.categoryCommandRepository.UpdateCategory(ctx, req)

	if err != nil {
		return s.errorHandler.HandleUpdateCategoryError(err, method, "FAILED_UPDATE_CATEGORY", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryResponse(category)

	s.mencache.DeleteCachedCategoryCache(ctx, *req.CategoryID)

	logSuccess("Successfully updated category", zap.Int("category.id", category.ID))

	return so, nil
}

func (s *categoryCommandService) TrashedCategory(ctx context.Context, category_id int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedCategory"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("category.id", category_id))

	defer func() {
		end(status)
	}()

	category, err := s.categoryCommandRepository.TrashedCategory(ctx, category_id)

	if err != nil {
		return s.errorHandler.HandleTrashedCategoryError(err, method, "FAILED_TRASH_CATEGORY", span, &status, zap.Error(err))
	}
	so := s.mapping.ToCategoryResponseDeleteAt(category)

	s.mencache.DeleteCachedCategoryCache(ctx, category_id)

	logSuccess("Successfully trashed category", zap.Int("category.id", category_id))

	return so, nil
}

func (s *categoryCommandService) RestoreCategory(ctx context.Context, categoryID int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreCategory"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("category.id", categoryID))

	defer func() {
		end(status)
	}()

	category, err := s.categoryCommandRepository.RestoreCategory(ctx, categoryID)

	if err != nil {
		return s.errorHandler.HandleRestoreError(err, method, "FAILED_RESTORE_CATEGORY", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryResponseDeleteAt(category)

	s.mencache.DeleteCachedCategoryCache(ctx, categoryID)

	logSuccess("Successfully restored category", zap.Int("category.id", categoryID))

	return so, nil
}

func (s *categoryCommandService) DeleteCategoryPermanent(ctx context.Context, categoryID int) (bool, *response.ErrorResponse) {
	const method = "DeleteCategoryPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("category.id", categoryID))

	defer func() {
		end(status)
	}()

	res, err := s.categoryQueryRepository.FindByIdTrashed(ctx, categoryID)

	if err != nil {
		return s.errorHandler.HandleDeleteAllPermanentlyError(err, method, "FAILED_DELETE_CATEGORY_PERMANENT", span, &status, zap.Error(err))
	}

	if res.ImageCategory != "" {
		err := os.Remove(res.ImageCategory)
		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Debug("Image file does not exist, skipping deletion",
					zap.String("image_path", res.ImageCategory))

				span.SetAttributes(attribute.String("image_path", res.ImageCategory))
			} else {
				return errorhandler.HandleFiledError(s.logger, err, method, "FAILED_DELETE_CATEGORY_PERMANENT", res.ImageCategory, span, &status, category_errors.ErrFailedRemoveImageCategory, zap.Error(err))
			}
		} else {
			s.logger.Debug("Successfully deleted category image",
				zap.String("image_path", res.ImageCategory))
		}
	}

	success, err := s.categoryCommandRepository.DeleteCategoryPermanently(ctx, categoryID)

	if err != nil {
		return s.errorHandler.HandleDeleteError(err, method, "FAILED_DELETE_CATEGORY_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted category permanently", zap.Bool("success", success))

	return success, nil
}

func (s *categoryCommandService) RestoreAllCategories(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllCategories"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.categoryCommandRepository.RestoreAllCategories(ctx)

	if err != nil {
		return s.errorHandler.HandleRestoreAllError(err, method, "FAILED_RESTORE_ALL_TRASHED_CATEGORIES", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all trashed categories", zap.Bool("success", success))

	return success, nil
}

func (s *categoryCommandService) DeleteAllCategoriesPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllCategoriesPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.categoryCommandRepository.DeleteAllPermanentCategories(ctx)

	if err != nil {
		return s.errorHandler.HandleDeleteAllPermanentlyError(err, "DeleteAllCategoriesPermanent", "FAILED_DELETE_ALL_CATEGORIES_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all categories permanently", zap.Bool("success", success))

	return success, nil
}

func (s *categoryCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

	return ctx, span, end, status, logSuccess
}

func (s *categoryCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
