package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type BannerQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.BannerResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.BannerResponse, *response.ErrorResponse)
}

type BannerCommandError interface {
	HandleCreateBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.BannerResponse, *response.ErrorResponse)
	HandleUpdateBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.BannerResponse, *response.ErrorResponse)
	HandleTrashedBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.BannerResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.BannerResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
