package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDetailQueryError struct {
	logger logger.LoggerInterface
}

func NewMerchantDetailQueryError(logger logger.LoggerInterface) *merchantDetailQueryError {
	return &merchantDetailQueryError{
		logger: logger,
	}
}

func (e *merchantDetailQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantDetailResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantDetailResponse](e.logger, err, method, tracePrefix, span, status, merchantdetail_errors.ErrFailedFindAllMerchantDetail, fields...)
}

func (e *merchantDetailQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantDetailResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
