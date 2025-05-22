package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantPolicyCommandService struct {
	ctx                             context.Context
	trace                           trace.Tracer
	logger                          logger.LoggerInterface
	merchantPolicyCommandRepository repository.MerchantPoliciesCommandRepository
	merchantQueryRepository         repository.MerchantQueryRepository
	mapping                         response_service.MerchantPolicyResponseMapper
	requestCounter                  *prometheus.CounterVec
	requestDuration                 *prometheus.HistogramVec
}

func NewMerchantPolicyCommandService(ctx context.Context, logger logger.LoggerInterface, merchantPolicyCommandRepository repository.MerchantPoliciesCommandRepository, merchantQueryRepository repository.MerchantQueryRepository, mapping response_service.MerchantPolicyResponseMapper) *merchantPolicyCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_policy_command_service_requests_total",
			Help: "Total number of requests to the MerchantPolicyCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_policy_command_service_request_duration_seconds",
			Help:    "Histogram of request duration for the MerchantPolicyCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantPolicyCommandService{
		ctx:                             ctx,
		trace:                           otel.Tracer("merchant-policy-command-service"),
		logger:                          logger,
		mapping:                         mapping,
		merchantPolicyCommandRepository: merchantPolicyCommandRepository,
		merchantQueryRepository:         merchantQueryRepository,
		requestCounter:                  requestCounter,
		requestDuration:                 requestDuration,
	}
}

func (s *merchantPolicyCommandService) CreateMerchant(req *requests.CreateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateMerchant")
	defer span.End()

	s.logger.Debug("Creating new merchant")

	merchant, err := s.merchantPolicyCommandRepository.CreateMerchantPolicy(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_MERCHANT_POLICY")

		s.logger.Error("Failed to create merchant",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to create merchant")

		status = "failed_create_merchant_policy"

		return nil, merchantpolicy_errors.ErrFailedCreateMerchantPolicy
	}

	return s.mapping.ToMerchantPolicyResponse(merchant), nil
}

func (s *merchantPolicyCommandService) UpdateMerchant(req *requests.UpdateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchant")
	defer span.End()

	s.logger.Debug("Updating merchant", zap.Int("merchantPolicy", *req.MerchantPolicyID))

	merchant, err := s.merchantPolicyCommandRepository.UpdateMerchantPolicy(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_MERCHANT_POLICY")

		s.logger.Error("Failed to update merchant",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to update merchant")

		status = "failed_update_merchant_policy"

		return nil, merchantpolicy_errors.ErrFailedUpdateMerchantPolicy
	}

	return s.mapping.ToMerchantPolicyResponse(merchant), nil
}

func (s *merchantPolicyCommandService) TrashedMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchantID", merchantID),
	)

	s.logger.Debug("Trashing merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantPolicyCommandRepository.TrashedMerchantPolicy(merchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_MERCHANT_POLICY")

		s.logger.Error("Failed to trash merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to trash merchant")

		status = "failed_trash_merchant_policy"

		return nil, merchantpolicy_errors.ErrFailedTrashedMerchantPolicy
	}

	return s.mapping.ToMerchantPolicyResponseDeleteAt(merchant), nil
}

func (s *merchantPolicyCommandService) RestoreMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchantID", merchantID),
	)

	s.logger.Debug("Restoring merchant", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantPolicyCommandRepository.RestoreMerchantPolicy(merchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_MERCHANT_POLICY")

		s.logger.Error("Failed to restore merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to restore merchant")

		status = "failed_restore_merchant_policy"

		return nil, merchantpolicy_errors.ErrFailedRestoreMerchantPolicy
	}

	return s.mapping.ToMerchantPolicyResponseDeleteAt(merchant), nil
}

func (s *merchantPolicyCommandService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteMerchantPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteMerchantPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchantID", merchantID),
	)

	s.logger.Debug("Deleting merchant permanently", zap.Int("merchantID", merchantID))

	success, err := s.merchantPolicyCommandRepository.DeleteMerchantPolicyPermanent(merchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_MERCHANT_POLICY_PERMANENT")

		s.logger.Error("Failed to permanently delete merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to permanently delete merchant")

		status = "failed_delete_merchant_policy_permanent"

		return false, merchantpolicy_errors.ErrFailedDeleteMerchantPolicyPermanent
	}

	return success, nil
}

func (s *merchantPolicyCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllMerchant")
	defer span.End()

	s.logger.Debug("Restoring all trashed merchants")

	success, err := s.merchantPolicyCommandRepository.RestoreAllMerchantPolicy()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_MERCHANT_POLICY")

		s.logger.Error("Failed to restore all trashed merchants",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to restore all trashed merchants")

		status = "failed_restore_all_merchant_policy"

		return false, merchantpolicy_errors.ErrFailedRestoreAllMerchantPolicies
	}

	return success, nil
}

func (s *merchantPolicyCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllMerchantPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllMerchantPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all merchants")

	success, err := s.merchantPolicyCommandRepository.DeleteAllMerchantPolicyPermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_MERCHANT_POLICY")

		s.logger.Error("Failed to permanently delete all merchants",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to permanently delete all merchants")

		status = "failed_delete_all_merchant_policy"

		return false, merchantpolicy_errors.ErrFailedDeleteAllMerchantPoliciesPermanent
	}

	return success, nil
}

func (s *merchantPolicyCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
