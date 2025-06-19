package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantPolicyCommandError struct {
	logger logger.LoggerInterface
}

func NewMerchantPolicyCommandError(logger logger.LoggerInterface) *merchantPolicyCommandError {
	return &merchantPolicyCommandError{
		logger: logger,
	}
}

func (e *merchantPolicyCommandError) HandleCreateMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantPoliciesResponse](e.logger, err, method, tracePrefix, span, status, merchantpolicy_errors.ErrFailedCreateMerchantPolicy, fields...)
}

func (e *merchantPolicyCommandError) HandleUpdateMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantPoliciesResponse](e.logger, err, method, tracePrefix, span, status, merchantpolicy_errors.ErrFailedUpdateMerchantPolicy, fields...)
}

func (e *merchantPolicyCommandError) HandleTrashedMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantPoliciesResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchantpolicy_errors.ErrFailedTrashedMerchantPolicy, fields...)
}

func (e *merchantPolicyCommandError) HandleRestoreMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorTemplate[*response.MerchantPoliciesResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, merchantpolicy_errors.ErrFailedRestoreMerchantPolicy, fields...)
}

func (e *merchantPolicyCommandError) HandleDeleteMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](e.logger, err, method, tracePrefix, span, status, merchantpolicy_errors.ErrFailedDeleteMerchantPolicyPermanent, fields...)
}

func (e *merchantPolicyCommandError) HandleRestoreAllMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](e.logger, err, method, tracePrefix, span, status, merchantpolicy_errors.ErrFailedRestoreAllMerchantPolicies, fields...)
}

func (e *merchantPolicyCommandError) HandleDeleteAllMerchantPolicyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](e.logger, err, method, tracePrefix, span, status, merchantpolicy_errors.ErrFailedDeleteAllMerchantPoliciesPermanent, fields...)
}
