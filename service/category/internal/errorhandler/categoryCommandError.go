package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type categoryCommandError struct {
	logger logger.LoggerInterface
}

func NewCategoryCommandError(logger logger.LoggerInterface) *categoryCommandError {
	return &categoryCommandError{logger: logger}
}

func (c *categoryCommandError) HandleCreateCategoryError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CategoryResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CategoryResponse](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedCreateCategory, fields...)
}

func (c *categoryCommandError) HandleUpdateCategoryError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CategoryResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CategoryResponse](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedUpdateCategory, fields...)
}

func (c *categoryCommandError) HandleTrashedCategoryError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.CategoryResponseDeleteAt](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedTrashedCategory, fields...)
}

func (c *categoryCommandError) HandleDeletePermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedDeleteCategoryPermanent, fields...)
}

func (c *categoryCommandError) HandleDeleteAllPermanentlyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedDeleteAllCategoriesPermanent, fields...)
}

func (c *categoryCommandError) HandleRestoreError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.CategoryResponseDeleteAt](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedRestoreCategory, fields...)
}

func (c *categoryCommandError) HandleRestoreAllError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedRestoreAllCategories, fields...)
}

func (c *categoryCommandError) HandleDeleteError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedDeleteCategoryPermanent, fields...)
}

func (c *categoryCommandError) HandleDeleteAllError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedDeleteAllCategoriesPermanent, fields...)
}
