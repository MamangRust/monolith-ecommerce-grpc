package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantAwardQueryError struct {
	logger logger.LoggerInterface
}

func NewMerchantAwardQueryError(logger logger.LoggerInterface) *merchantAwardQueryError {
	return &merchantAwardQueryError{
		logger: logger,
	}
}

func (e *merchantAwardQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantAwardResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantAwardResponse](e.logger, err, method, tracePrefix, span, status, merchantaward_errors.ErrFailedFindAllMerchantAwards, fields...)
}

func (e *merchantAwardQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantAwardResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (e *merchantAwardQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantAwardResponse](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
