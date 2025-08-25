package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantsociallink_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantSocialLinkCommandError struct {
	logger logger.LoggerInterface
}

func NewMerchantSocialLinkCommandError(logger logger.LoggerInterface) *merchantSocialLinkCommandError {
	return &merchantSocialLinkCommandError{
		logger: logger,
	}
}

func (e *merchantSocialLinkCommandError) HandleCreateMerchantSocialLinkError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantSocialLinkResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantSocialLinkResponse](e.logger, err, method, tracePrefix, span, status, merchantsociallink_errors.ErrFailedCreateMerchantSocialLink, fields...)
}

func (e *merchantSocialLinkCommandError) HandleUpdateMerchantSocialLinkError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantSocialLinkResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantSocialLinkResponse](e.logger, err, method, tracePrefix, span, status, merchantsociallink_errors.ErrFailedUpdateMerchantSocialLink, fields...)
}
