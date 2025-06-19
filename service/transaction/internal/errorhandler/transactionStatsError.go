package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatsError struct {
	logger logger.LoggerInterface
}

func NewTransactionStatsError(logger logger.LoggerInterface) *transactionStatsError {
	return &transactionStatsError{logger: logger}
}

func (t *transactionStatsError) HandleMonthlyAmountSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthlyAmountSuccessResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyAmountSuccess, fields...)
}

func (t *transactionStatsError) HandleYearlyAmountSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyAmountSuccessResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyAmountSuccess, fields...)
}

func (t *transactionStatsError) HandleMonthlyAmountFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthlyAmountFailedResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyAmountFailed, fields...)
}

func (t *transactionStatsError) HandleYearlyAmountFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyAmountFailedResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyAmountFailed, fields...)
}

func (t *transactionStatsError) HandleMonthlyMethodSuccessError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthlyMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyMethod, fields...)
}

func (t *transactionStatsError) HandleYearlyMethodSuccessError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyMethod, fields...)
}

func (t *transactionStatsError) HandleMonthlyMethodFailedError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthlyMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyMethod, fields...)
}

func (t *transactionStatsError) HandleYearlyMethodFailedError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyMethod, fields...)
}
