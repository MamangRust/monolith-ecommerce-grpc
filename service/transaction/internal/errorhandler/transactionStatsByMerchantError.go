package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatsByMerchantError struct {
	logger logger.LoggerInterface
}

func NewTransactionStatsByMerchantError(logger logger.LoggerInterface) *transactionStatsByMerchantError {
	return &transactionStatsByMerchantError{logger: logger}
}

func (t *transactionStatsByMerchantError) HandleMonthlyAmountSuccessByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthlyAmountSuccessResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyAmountSuccess, fields...)
}

func (t *transactionStatsByMerchantError) HandleYearlyAmountSuccessByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyAmountSuccessResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyAmountSuccess, fields...)
}

func (t *transactionStatsByMerchantError) HandleMonthlyAmountFailedByMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthlyAmountFailedResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyAmountFailed, fields...)
}

func (t *transactionStatsByMerchantError) HandleYearlyAmountFailedByMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyAmountFailedResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyAmountFailed, fields...)
}

func (t *transactionStatsByMerchantError) HandleMonthlyMethodSuccessByMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthlyMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyMethod, fields...)
}

func (t *transactionStatsByMerchantError) HandleYearlyMethodSuccessByMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyMethod, fields...)
}

func (t *transactionStatsByMerchantError) HandleMonthlyMethodFailedByMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionMonthlyMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindMonthlyMethod, fields...)
}

func (t *transactionStatsByMerchantError) HandleYearlyMethodFailedByMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.TransactionYearlyMethodResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedFindYearlyMethod, fields...)
}
