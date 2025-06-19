package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactonCommandError struct {
	logger logger.LoggerInterface
}

func NewTransactionCommandError(logger logger.LoggerInterface) *transactonCommandError {
	return &transactonCommandError{logger: logger}
}

func (t *transactonCommandError) HandleInsufficientBalance(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorInsufficientBalance[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedPaymentInsufficientBalance, fields...)
}

func (t *transactonCommandError) HandleCannotModifiedStatus(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorCannotModified[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedPaymentStatusCannotBeModified, fields...)
}

func (t *transactonCommandError) HandleInvalidOrderItem(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorInvalidQuantityOrderItem[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, orderitem_errors.ErrFailedFindOrderItemByOrder, fields...)
}

func (t *transactonCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (t *transactonCommandError) HandleCreateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedCreateTransaction, fields...)
}
func (t *transactonCommandError) HandleUpdateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedUpdateTransaction, fields...)
}

func (t *transactonCommandError) HandleTrashedTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedTrashedTransaction, fields...)
}

func (t *transactonCommandError) HandleRestoreTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedRestoreTransaction, fields...)
}

func (t *transactonCommandError) HandleDeleteTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedDeleteTransactionPermanently, fields...)
}

func (t *transactonCommandError) HandleRestoreAllTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedRestoreAllTransactions, fields...)
}

func (t *transactonCommandError) HandleDeleteAllTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedDeleteAllTransactionPermanent, fields...)
}
