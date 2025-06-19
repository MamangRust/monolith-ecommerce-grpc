package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type MerchantPolicyQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantPoliciesResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
}

type MerchantPolicyCommandError interface {
	HandleCreateMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
	HandleUpdateMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantPoliciesResponse, *response.ErrorResponse)
	HandleTrashedMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
