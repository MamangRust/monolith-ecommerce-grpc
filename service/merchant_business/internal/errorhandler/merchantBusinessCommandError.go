package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantBusinessCommandError struct {
	logger logger.LoggerInterface
}

func NewMerchantBusinessCommandError(logger logger.LoggerInterface) *merchantBusinessCommandError {
	return &merchantBusinessCommandError{
		logger: logger,
	}
}

func (e *merchantBusinessCommandError) HandleCreateMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantBusinessResponse](e.logger, err, method, tracePrefix, span, status, merchantbusiness_errors.ErrFailedCreateMerchantBusiness, fields...)
}

func (e *merchantBusinessCommandError) HandleUpdateMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantBusinessResponse](e.logger, err, method, tracePrefix, span, status, merchantbusiness_errors.ErrFailedUpdateMerchantBusiness, fields...)
}

func (e *merchantBusinessCommandError) HandleTrashedMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantBusinessResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchantbusiness_errors.ErrFailedTrashedMerchantBusiness, fields...)
}

func (e *merchantBusinessCommandError) HandleRestoreMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantBusinessResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchantbusiness_errors.ErrFailedRestoreMerchantBusiness, fields...)
}

func (e *merchantBusinessCommandError) HandleDeleteMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](e.logger, err, method, tracePrefix, span, status, merchantbusiness_errors.ErrFailedDeleteMerchantBusinessPermanent, fields...)
}

func (e *merchantBusinessCommandError) HandleRestoreAllMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](e.logger, err, method, tracePrefix, span, status, merchantbusiness_errors.ErrFailedRestoreAllMerchantBusiness, fields...)
}

func (e *merchantBusinessCommandError) HandleDeleteAllMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](e.logger, err, method, tracePrefix, span, status, merchantbusiness_errors.ErrFailedDeleteAllMerchantBusinessPermanent, fields...)
}
