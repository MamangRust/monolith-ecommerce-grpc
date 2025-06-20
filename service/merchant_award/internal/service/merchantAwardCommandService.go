package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/repository"
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

type merchantAwardCommandService struct {
	ctx                            context.Context
	errorhandler                   errorhandler.MerchantAwardCommandError
	mencache                       mencache.MerchanrAwardCommandCache
	trace                          trace.Tracer
	merchantAwardCommandRepository repository.MerchantAwardCommandRepository
	logger                         logger.LoggerInterface
	mapping                        response_service.MerchantAwardResponseMapper
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewMerchantAwardCommandService(ctx context.Context,
	errorhandler errorhandler.MerchantAwardCommandError,
	mencache mencache.MerchanrAwardCommandCache,
	merchantAwardCommandRepositroy repository.MerchantAwardCommandRepository, logger logger.LoggerInterface, mapping response_service.MerchantAwardResponseMapper) *merchantAwardCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_award_command_service_request_total",
			Help: "Total number of requests to the MerchantAwardCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_award_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantAwardCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantAwardCommandService{
		ctx:                            ctx,
		errorhandler:                   errorhandler,
		mencache:                       mencache,
		trace:                          otel.Tracer("merchant-award-Command-service"),
		merchantAwardCommandRepository: merchantAwardCommandRepositroy,
		logger:                         logger,
		mapping:                        mapping,
		requestCounter:                 requestCounter,
		requestDuration:                requestDuration,
	}
}

func (s *merchantAwardCommandService) CreateMerchant(req *requests.CreateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	const method = "CreateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", req.MerchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardCommandRepository.CreateMerchantAward(req)

	if err != nil {
		return s.errorhandler.HandleCreateMerchantAwardError(err, method, "FAILED_CREATE_MERCHANT_AWARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantAwardResponse(merchant)

	logSuccess("Merchant Award created", zap.Int("merchantAward.id", so.ID))

	return so, nil
}

func (s *merchantAwardCommandService) UpdateMerchant(req *requests.UpdateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	const method = "UpdateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantAward.id", *req.MerchantCertificationID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardCommandRepository.UpdateMerchantAward(req)

	if err != nil {
		return s.errorhandler.HandleUpdateMerchantAwardError(err, method, "FAILED_UPDATE_MERCHANT_AWARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantAwardResponse(merchant)

	s.mencache.DeleteMerchantAwardCache(*req.MerchantCertificationID)

	logSuccess("Merchant Award updated", zap.Int("merchantAward.id", *req.MerchantCertificationID))

	return so, nil
}

func (s *merchantAwardCommandService) TrashedMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardCommandRepository.TrashedMerchantAward(merchantID)

	if err != nil {
		return s.errorhandler.HandleTrashedMerchantAwardError(err, method, "FAILED_TRASH_MERCHANT_AWARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantAwardResponseDeleteAt(merchant)

	s.mencache.DeleteMerchantAwardCache(merchantID)

	logSuccess("Merchant Award trashed", zap.Int("merchantAward.id", merchantID))

	return so, nil
}

func (s *merchantAwardCommandService) RestoreMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardCommandRepository.RestoreMerchantAward(merchantID)

	if err != nil {
		return s.errorhandler.HandleRestoreMerchantAwardError(err, method, "FAILED_RESTORE_MERCHANT_AWARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantAwardResponseDeleteAt(merchant)

	logSuccess("Merchant Award restored", zap.Int("merchantAward.id", merchantID))

	return so, nil
}

func (s *merchantAwardCommandService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	const method = "DeleteMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardCommandRepository.DeleteMerchantPermanent(merchantID)

	if err != nil {
		return s.errorhandler.HandleDeleteMerchantAwardError(err, "DeleteMerchantPermanent", "FAILED_DELETE_MERCHANT_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted merchant permanently", zap.Int("merchantAward.id", merchantID), zap.Bool("success", success))

	return success, nil
}

func (s *merchantAwardCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	const method = "RestoreAllMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardCommandRepository.RestoreAllMerchantAward()

	if err != nil {
		return s.errorhandler.HandleRestoreAllMerchantAwardError(err, method, "FAILED_RESTORE_ALL_MERCHANT_AWARD", span, &status, zap.Error(err))
	}

	logSuccess("All trashed merchants restored", zap.Bool("success", success))

	return success, nil
}

func (s *merchantAwardCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardCommandRepository.DeleteAllMerchantAwardPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllMerchantAwardError(err, method, "FAILED_DELETE_ALL_MERCHANT_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all merchants permanently", zap.Bool("success", success))

	return success, nil
}

func (s *merchantAwardCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantAwardCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
