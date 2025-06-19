package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TransactionStatsError interface {
	HandleMonthlyAmountSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse)
	HandleYearlyAmountSuccessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse)

	HandleMonthlyAmountFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse)
	HandleYearlyAmountFailedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse)

	HandleMonthlyMethodSuccessError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	HandleYearlyMethodSuccessError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)

	HandleMonthlyMethodFailedError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)

	HandleYearlyMethodFailedError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)
}

type TransactionStatsByMerchantError interface {
	HandleMonthlyAmountSuccessByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse)
	HandleYearlyAmountSuccessByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse)

	HandleMonthlyAmountFailedByMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse)
	HandleYearlyAmountFailedByMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse)

	HandleMonthlyMethodSuccessByMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)

	HandleYearlyMethodSuccessByMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)

	HandleMonthlyMethodFailedByMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)

	HandleYearlyMethodFailedByMerchantError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)
}

type TransactionQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TransactionResponse, *response.ErrorResponse)
}

type TransactionCommandError interface {
	HandleCannotModifiedStatus(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.TransactionResponse, *response.ErrorResponse)
	HandleInsufficientBalance(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)
	HandleInvalidOrderItem(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)
	HandleCreateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)
	HandleUpdateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse)
	HandleTrashedTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
