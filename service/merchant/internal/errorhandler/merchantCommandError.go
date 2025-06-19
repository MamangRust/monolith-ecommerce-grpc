package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantCommandError struct {
	logger logger.LoggerInterface
}

func NewMerchantCommandError(logger logger.LoggerInterface) *merchantCommandError {
	return &merchantCommandError{
		logger: logger,
	}
}

func (e *merchantCommandError) HandleCreateMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedCreateMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleUpdateMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedUpdateMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleUpdateMerchantStatusError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedUpdateMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleTrashedMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedTrashedMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleRestoreMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedRestoreMerchant,
		fields...,
	)
}

func (e *merchantCommandError) HandleDeleteMerchantPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedDeleteMerchantPermanent,
		fields...,
	)
}

func (e *merchantCommandError) HandleRestoreAllMerchantError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedRestoreAllMerchants,
		fields...,
	)
}

func (e *merchantCommandError) HandleDeleteAllMerchantPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchant_errors.ErrFailedDeleteAllMerchantsPermanent,
		fields...,
	)
}
