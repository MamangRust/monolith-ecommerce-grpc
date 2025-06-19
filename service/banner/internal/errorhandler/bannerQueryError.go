package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type bannerQueryError struct {
	logger logger.LoggerInterface
}

func NewBannerQueryError(logger logger.LoggerInterface) *bannerQueryError {
	return &bannerQueryError{
		logger: logger,
	}
}

func (b *bannerQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.BannerResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.BannerResponse](b.logger, err, method, tracePrefix, span, status, banner_errors.ErrBannerNotFoundRes, fields...)
}
func (b *bannerQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.BannerResponseDeleteAt](b.logger, err, method, tracePrefix, span, status, banner_errors.ErrBannerNotFoundRes, fields...)
}

func (b *bannerQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.BannerResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.BannerResponse](b.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
