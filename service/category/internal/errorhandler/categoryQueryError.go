package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type categoryQueryError struct {
	logger logger.LoggerInterface
}

func NewCategoryQueryError(logger logger.LoggerInterface) *categoryQueryError {
	return &categoryQueryError{logger: logger}
}

func (c *categoryQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CategoryResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.CategoryResponse](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedFindAllCategories, fields...)
}

func (c *categoryQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.CategoryResponseDeleteAt](c.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (c *categoryQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.CategoryResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CategoryResponse](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedFindCategoryById, fields...)
}
