package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDetailCommandError struct {
	logger logger.LoggerInterface
}

func NewMerchantDetailCommandError(logger logger.LoggerInterface) *merchantDetailCommandError {
	return &merchantDetailCommandError{
		logger: logger,
	}
}

func (e *merchantDetailCommandError) HandleCreateMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantDetailResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDetailResponse](e.logger, err, method, tracePrefix, span, status, merchantdetail_errors.ErrFailedCreateMerchantDetail, fields...)
}

func (e *merchantDetailCommandError) HandleUpdateMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantDetailResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDetailResponse](e.logger, err, method, tracePrefix, span, status, merchantdetail_errors.ErrFailedUpdateMerchantDetail, fields...)
}

func (e *merchantDetailCommandError) HandleTrashedMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDetailResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchantdetail_errors.ErrFailedTrashedMerchantDetail, fields...)
}

func (e *merchantDetailCommandError) HandleRestoreMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDetailResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchantdetail_errors.ErrFailedRestoreMerchantDetail, fields...)
}

func (e *merchantDetailCommandError) HandleDeleteMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](e.logger, err, method, tracePrefix, span, status, merchantdetail_errors.ErrFailedDeleteMerchantDetailPermanent, fields...)
}

func (e *merchantDetailCommandError) HandleRestoreAllMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](e.logger, err, method, tracePrefix, span, status, merchantdetail_errors.ErrFailedRestoreAllMerchantDetail, fields...)
}

func (e *merchantDetailCommandError) HandleDeleteAllMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](e.logger, err, method, tracePrefix, span, status, merchantdetail_errors.ErrFailedDeleteAllMerchantDetailPermanent, fields...)
}
