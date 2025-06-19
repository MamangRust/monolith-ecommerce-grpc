package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantAwardCommandError struct {
	logger logger.LoggerInterface
}

func NewMerchantAwardCommandError(logger logger.LoggerInterface) *merchantAwardCommandError {
	return &merchantAwardCommandError{
		logger: logger,
	}
}

func (e *merchantAwardCommandError) HandleCreateMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantAwardResponse](e.logger, err, method, tracePrefix, span, status, merchantaward_errors.ErrFailedCreateMerchantAward, fields...)
}

func (e *merchantAwardCommandError) HandleUpdateMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantAwardResponse](e.logger, err, method, tracePrefix, span, status, merchantaward_errors.ErrFailedUpdateMerchantAward, fields...)
}

func (e *merchantAwardCommandError) HandleTrashedMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantAwardResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchantaward_errors.ErrFailedTrashedMerchantAward, fields...)
}

func (e *merchantAwardCommandError) HandleRestoreMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantAwardResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchantaward_errors.ErrFailedRestoreMerchantAward, fields...)
}

func (e *merchantAwardCommandError) HandleDeleteMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](e.logger, err, method, tracePrefix, span, status, merchantaward_errors.ErrFailedDeleteMerchantAwardPermanent, fields...)
}

func (e *merchantAwardCommandError) HandleRestoreAllMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](e.logger, err, method, tracePrefix, span, status, merchantaward_errors.ErrFailedRestoreAllMerchantAwards, fields...)
}

func (e *merchantAwardCommandError) HandleDeleteAllMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](e.logger, err, method, tracePrefix, span, status, merchantaward_errors.ErrFailedDeleteAllMerchantAwardsPermanent, fields...)
}
