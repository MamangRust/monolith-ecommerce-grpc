package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type MerchantDetailQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantDetailResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse)
}

type MerchantDetailCommandError interface {
	HandleCreateMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantDetailResponse, *response.ErrorResponse)
	HandleUpdateMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantDetailResponse, *response.ErrorResponse)
	HandleTrashedMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllMerchantDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

type FileError interface {
	HandleErrorFileCover(
		logger logger.LoggerInterface,
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleErrorFileLogo(
		logger logger.LoggerInterface,
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}
