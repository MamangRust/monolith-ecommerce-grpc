package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantPolicyQueryError struct {
	logger logger.LoggerInterface
}

func NewMerchantPolicyQueryError(logger logger.LoggerInterface) *merchantPolicyQueryError {
	return &merchantPolicyQueryError{
		logger: logger,
	}
}

func (e *merchantPolicyQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantPoliciesResponse, *int, *response.ErrorResponse) {
	return handleErrorPaginationTemplate[[]*response.MerchantPoliciesResponse](e.logger, err, method, tracePrefix, span, status, merchantpolicy_errors.ErrFailedFindAllMerchantPolicies, fields...)
}

func (e *merchantPolicyQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPaginationTemplate[[]*response.MerchantPoliciesResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (e *merchantPolicyQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantPoliciesResponse](e.logger, err, method, tracePrefix, span, status, merchantpolicy_errors.ErrFailedFindMerchantPolicyById, fields...)
}
