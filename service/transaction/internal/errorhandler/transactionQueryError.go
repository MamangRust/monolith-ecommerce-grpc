package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionQueryError struct {
	logger logger.LoggerInterface
}

func NewTransactionQueryError(logger logger.LoggerInterface) *transactionQueryError {
	return &transactionQueryError{logger: logger}
}

func (t *transactionQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindAllTransactions, fields...)
}

func (t *transactionQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.TransactionResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (t *transactionQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
