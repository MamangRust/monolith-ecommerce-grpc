package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantBusinessQueryError struct {
	logger logger.LoggerInterface
}

func NewMerchantBusinessQueryError(logger logger.LoggerInterface) *merchantBusinessQueryError {
	return &merchantBusinessQueryError{
		logger: logger,
	}
}

func (e *merchantBusinessQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantBusinessResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantBusinessResponse](e.logger, err, method, tracePrefix, span, status, merchantbusiness_errors.ErrFailedFindAllMerchantBusiness, fields...)
}

func (e *merchantBusinessQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantBusinessResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (e *merchantBusinessQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantBusinessResponse](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
