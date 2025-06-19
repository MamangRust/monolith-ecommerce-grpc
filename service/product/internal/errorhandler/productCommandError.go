package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type productCommandError struct {
	logger logger.LoggerInterface
}

func NewProductCommandError(logger logger.LoggerInterface) *productCommandError {
	return &productCommandError{
		logger: logger,
	}
}

func (o *productCommandError) HandleFileError(
	err error,
	method, tracePrefix, imagePath string,
	span trace.Span,
	status *string,
	fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleFiledError(o.logger, err, method, tracePrefix, imagePath, span, status, fields...)
}

func (p *productCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.ProductResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponse](p.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (p *productCommandError) HandleCreateProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponse](p.logger, err, method, tracePrefix, span, status, product_errors.ErrFailedCreateProduct, fields...)
}

func (p *productCommandError) HandleUpdateProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponse](p.logger, err, method, tracePrefix, span, status, product_errors.ErrFailedUpdateProduct, fields...)
}

func (p *productCommandError) HandleTrashedProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponseDeleteAt](p.logger, err, method, tracePrefix, span, status, product_errors.ErrFailedTrashProduct, fields...)
}

func (p *productCommandError) HandleRestoreProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponseDeleteAt](p.logger, err, method, tracePrefix, span, status, product_errors.ErrFailedRestoreProduct, fields...)
}

func (p *productCommandError) HandleDeleteProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](p.logger, err, method, tracePrefix, span, status, product_errors.ErrFailedDeleteProductPermanent, fields...)
}

func (p *productCommandError) HandleRestoreAllProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](p.logger, err, method, tracePrefix, span, status, product_errors.ErrFailedRestoreAllProducts, fields...)
}

func (p *productCommandError) HandleDeleteAllProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](p.logger, err, method, tracePrefix, span, status, product_errors.ErrFailedDeleteAllProductsPermanent, fields...)
}
