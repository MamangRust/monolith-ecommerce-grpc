package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/repository"
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

type merchantBusinessCommandService struct {
	ctx                               context.Context
	errorhandler                      errorhandler.MerchantBusinessCommandError
	mencache                          mencache.MerchanrBusinessCommandCache
	trace                             trace.Tracer
	merchantQueryRepository           repository.MerchantQueryRepository
	merchantBusinessCommandRepository repository.MerchantBusinessCommandRepository
	logger                            logger.LoggerInterface
	mapping                           response_service.MerchantBusinessResponseMapper
	requestCounter                    *prometheus.CounterVec
	requestDuration                   *prometheus.HistogramVec
}

func NewMerchantBusinessCommandService(
	ctx context.Context,
	errorhandler errorhandler.MerchantBusinessCommandError,
	mencache mencache.MerchanrBusinessCommandCache,
	merchantQueryRepository repository.MerchantQueryRepository,
	merchantBusinessCommandRepository repository.MerchantBusinessCommandRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantBusinessResponseMapper,
) *merchantBusinessCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_business_command_service_requests_total",
			Help: "Total number of requests to the MerchantBusinessCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_business_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantBusinessCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantBusinessCommandService{
		ctx:                               ctx,
		errorhandler:                      errorhandler,
		mencache:                          mencache,
		trace:                             otel.Tracer("merchant-business-command-service"),
		merchantQueryRepository:           merchantQueryRepository,
		merchantBusinessCommandRepository: merchantBusinessCommandRepository,
		logger:                            logger,
		mapping:                           mapping,
		requestCounter:                    requestCounter,
		requestDuration:                   requestDuration,
	}
}

func (s *merchantBusinessCommandService) CreateMerchant(req *requests.CreateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	const method = "CreateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantBusiness.id", req.MerchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessCommandRepository.CreateMerchantBusiness(req)

	if err != nil {
		return s.errorhandler.HandleCreateMerchantBusinessError(err, method, "FAILED_CREATE_MERCHANT_BUSINESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantBusinessResponse(merchant)

	logSuccess("Created merchant", zap.Int("merchantBusiness.id", req.MerchantID))

	return so, nil
}

func (s *merchantBusinessCommandService) UpdateMerchant(req *requests.UpdateMerchantBusinessInformationRequest) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	const method = "UpdateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantBusiness.id", *req.MerchantBusinessInfoID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessCommandRepository.UpdateMerchantBusiness(req)

	if err != nil {
		return s.errorhandler.HandleUpdateMerchantBusinessError(err, method, "FAILED_UPDATE_MERCHANT_BUSINESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantBusinessResponse(merchant)

	s.mencache.DeleteMerchantBusinessCache(*req.MerchantBusinessInfoID)

	logSuccess("Updated merchant", zap.Int("merchantBusiness.id", *req.MerchantBusinessInfoID))

	return so, nil
}

func (s *merchantBusinessCommandService) TrashedMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantBusiness.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessCommandRepository.TrashedMerchantBusiness(merchantID)

	if err != nil {
		return s.errorhandler.HandleTrashedMerchantBusinessError(err, method, "FAILED_TRASH_MERCHANT_BUSINESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantBusinessResponseDeleteAt(merchant)

	s.mencache.DeleteMerchantBusinessCache(merchantID)

	logSuccess("Trashed merchant", zap.Int("merchantBusiness.id", merchantID))

	return so, nil
}

func (s *merchantBusinessCommandService) RestoreMerchant(merchantID int) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantBusiness.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessCommandRepository.RestoreMerchantBusiness(merchantID)

	if err != nil {
		return s.errorhandler.HandleRestoreMerchantBusinessError(err, method, "FAILED_RESTORE_MERCHANT_BUSINESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantBusinessResponseDeleteAt(merchant)

	logSuccess("Restored merchant", zap.Int("merchantBusiness.id", merchantID))

	return so, nil

}

func (s *merchantBusinessCommandService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	const method = "DeleteMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantBusiness.id", merchantID))

	defer func() {
		end(status)
	}()

	success, err := s.merchantBusinessCommandRepository.DeleteMerchantBusinessPermanent(merchantID)

	if err != nil {
		return s.errorhandler.HandleDeleteMerchantBusinessError(err, method, "FAILED_DELETE_MERCHANT_BUSINESS_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted merchant permanently", zap.Int("merchantBusiness.id", merchantID), zap.Bool("success", success))

	return success, nil
}

func (s *merchantBusinessCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	const method = "RestoreAllMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantBusinessCommandRepository.RestoreAllMerchantBusiness()

	if err != nil {
		return s.errorhandler.HandleRestoreAllMerchantBusinessError(err, method, "FAILED_RESTORE_ALL_MERCHANT_BUSINESS", span, &status, zap.Error(err))
	}

	logSuccess("All trashed merchants restored", zap.Bool("success", success))

	return success, nil
}

func (s *merchantBusinessCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantBusinessCommandRepository.DeleteAllMerchantBusinessPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllMerchantBusinessError(err, method, "FAILED_DELETE_ALL_MERCHANT_BUSINESS_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all trashed merchants permanently", zap.Bool("success", success))

	return success, nil
}

func (s *merchantBusinessCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantBusinessCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
