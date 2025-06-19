package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type fileError struct {
	logger logger.LoggerInterface
}

func NewFileError(logger logger.LoggerInterface) *fileError {
	return &fileError{
		logger: logger,
	}
}

func (f *fileError) HandleErrorFileCover(
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorFileError[bool](logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (f *fileError) HandleErrorFileLogo(
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorFileError[bool](logger, err, method, tracePrefix, span, status, errResp, fields...)
}
