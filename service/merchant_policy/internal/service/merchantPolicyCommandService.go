package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
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
	errorhandler                    errorhandler.MerchantPolicyCommandError
	mencache                        mencache.MerchantPolicyCommandCache
	trace                           trace.Tracer
	logger                          logger.LoggerInterface
	merchantPolicyCommandRepository repository.MerchantPoliciesCommandRepository
	merchantQueryRepository         repository.MerchantQueryRepository
	mapping                         response_service.MerchantPolicyResponseMapper
	requestCounter                  *prometheus.CounterVec
	requestDuration                 *prometheus.HistogramVec
}

func NewMerchantPolicyCommandService(ctx context.Context,
	errorhandler errorhandler.MerchantPolicyCommandError,
	mencache mencache.MerchantPolicyCommandCache,
	logger logger.LoggerInterface, merchantPolicyCommandRepository repository.MerchantPoliciesCommandRepository, merchantQueryRepository repository.MerchantQueryRepository, mapping response_service.MerchantPolicyResponseMapper) *merchantPolicyCommandService {
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
		errorhandler:                    errorhandler,
		mencache:                        mencache,
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
	const method = "CreateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", req.MerchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantPolicyCommandRepository.CreateMerchantPolicy(req)

	if err != nil {
		return s.errorhandler.HandleCreateMerchantPolicyError(err, method, "FAILED_CREATE_MERCHANT_POLICY", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantPolicyResponse(merchant)

	logSuccess("Created merchant", zap.Int("merchant.id", req.MerchantID))

	return so, nil
}

func (s *merchantPolicyCommandService) UpdateMerchant(req *requests.UpdateMerchantPolicyRequest) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	const method = "UpdateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantPolicy.id", *req.MerchantPolicyID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantPolicyCommandRepository.UpdateMerchantPolicy(req)

	if err != nil {
		return s.errorhandler.HandleUpdateMerchantPolicyError(err, method, "FAILED_UPDATE_MERCHANT_POLICY", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantPolicyResponse(merchant)

	logSuccess("Updated merchant", zap.Int("merchantPolicy.id", *req.MerchantPolicyID))

	return so, nil
}

func (s *merchantPolicyCommandService) TrashedMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int(",merchantPolicy.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantPolicyCommandRepository.TrashedMerchantPolicy(merchantID)

	if err != nil {
		return s.errorhandler.HandleTrashedMerchantPolicyError(err, method, "FAILED_TRASH_MERCHANT_POLICY", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantPolicyResponseDeleteAt(merchant)

	s.mencache.DeleteMerchantPolicyCache(merchantID)

	logSuccess("Trashed merchant", zap.Int("merchantPolicy.id", merchantID))

	return so, nil
}

func (s *merchantPolicyCommandService) RestoreMerchant(merchantID int) (*response.MerchantPoliciesResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantPolicy.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantPolicyCommandRepository.RestoreMerchantPolicy(merchantID)

	if err != nil {
		return s.errorhandler.HandleRestoreMerchantPolicyError(err, method, "FAILED_RESTORE_MERCHANT_POLICY", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantPolicyResponseDeleteAt(merchant)

	logSuccess("Restored merchant", zap.Int("merchantPolicy.id", merchantID))

	return so, nil
}

func (s *merchantPolicyCommandService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	const method = "DeleteMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantPolicy.id", merchantID))

	defer func() {
		end(status)
	}()

	success, err := s.merchantPolicyCommandRepository.DeleteMerchantPolicyPermanent(merchantID)

	if err != nil {
		return s.errorhandler.HandleDeleteMerchantPolicyError(err, method, "FAILED_DELETE_MERCHANT_POLICY_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted merchant permanently", zap.Int("merchantPolicy.id", merchantID), zap.Bool("success", success))

	return success, nil
}

func (s *merchantPolicyCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	const method = "RestoreAllMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantPolicyCommandRepository.RestoreAllMerchantPolicy()

	if err != nil {
		return s.errorhandler.HandleRestoreAllMerchantPolicyError(err, method, "FAILED_RESTORE_ALL_MERCHANT_POLICY", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all merchants", zap.Bool("success", success))

	return success, nil
}

func (s *merchantPolicyCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantPolicyCommandRepository.DeleteAllMerchantPolicyPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllMerchantPolicyError(err, method, "FAILED_DELETE_ALL_MERCHANT_POLICY_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all merchants permanently", zap.Bool("success", success))

	return success, nil
}

func (s *merchantPolicyCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return span, end, status, logSuccess
}

func (s *merchantPolicyCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
