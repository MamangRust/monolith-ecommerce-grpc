package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type bannerCommandError struct {
	logger logger.LoggerInterface
}

func NewBannerCommandError(logger logger.LoggerInterface) *bannerCommandError {
	return &bannerCommandError{
		logger: logger,
	}
}

func (b *bannerCommandError) HandleCreateBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.BannerResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.BannerResponse](b.logger, err, method, tracePrefix, span, status, banner_errors.ErrFailedCreateBanner, fields...)
}

func (b *bannerCommandError) HandleUpdateBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.BannerResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.BannerResponse](b.logger, err, method, tracePrefix, span, status, banner_errors.ErrFailedUpdateBanner, fields...)
}

func (b *bannerCommandError) HandleTrashedBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.BannerResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.BannerResponseDeleteAt](b.logger, err, method, tracePrefix, span, status, banner_errors.ErrBannerNotFoundRes, fields...)
}

func (b *bannerCommandError) HandleRestoreBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.BannerResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.BannerResponseDeleteAt](b.logger, err, method, tracePrefix, span, status, banner_errors.ErrBannerNotFoundRes, fields...)
}

func (b *bannerCommandError) HandleDeleteBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](b.logger, err, method, tracePrefix, span, status, banner_errors.ErrBannerNotFoundRes, fields...)
}

func (b *bannerCommandError) HandleRestoreAllBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](b.logger, err, method, tracePrefix, span, status, banner_errors.ErrBannerNotFoundRes, fields...)
}

func (b *bannerCommandError) HandleDeleteAllBannerError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](b.logger, err, method, tracePrefix, span, status, banner_errors.ErrBannerNotFoundRes, fields...)
}
