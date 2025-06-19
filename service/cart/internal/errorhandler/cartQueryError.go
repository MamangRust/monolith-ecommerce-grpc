package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cartQueryError struct {
	logger logger.LoggerInterface
}

func NewCartQueryError(logger logger.LoggerInterface) *cartQueryError {
	return &cartQueryError{
		logger: logger,
	}
}

func (c *cartQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CartResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.CartResponse](c.logger, err, method, tracePrefix, span, status, cart_errors.ErrCartNotFoundRes, fields...)
}
