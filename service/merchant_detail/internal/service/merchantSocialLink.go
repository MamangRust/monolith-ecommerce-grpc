package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/repository"
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

type merchantSocialLinkService struct {
	ctx                          context.Context
	errorhandler                 errorhandler.MerchantSocialLinkCommandError
	trace                        trace.Tracer
	merchantSocialLinkRepository repository.MerchantSocialLinkCommandRepository
	mapping                      response_service.MerchantSocialLinkResponseMapper
	logger                       logger.LoggerInterface
	requestCounter               *prometheus.CounterVec
	requestDuration              *prometheus.HistogramVec
}

func NewMerchantSocialLinkService(ctx context.Context,
	errorhandler errorhandler.MerchantSocialLinkCommandError,
	merchantSocialLinkRepository repository.MerchantSocialLinkCommandRepository,
	mapping response_service.MerchantSocialLinkResponseMapper, logger logger.LoggerInterface) *merchantSocialLinkService {

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_social_link_service_requests_total",
			Help: "Total number of requests to the NewMerchantSocialLinkService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_social_link_service_request_duration_seconds",
			Help:    "Histogram of request durations for the NewMerchantSocialLinkService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantSocialLinkService{
		ctx:                          ctx,
		errorhandler:                 errorhandler,
		trace:                        otel.Tracer("merchant-social-link-service"),
		merchantSocialLinkRepository: merchantSocialLinkRepository,
		mapping:                      mapping,
		logger:                       logger,
		requestCounter:               requestCounter,
		requestDuration:              requestDuration,
	}
}

func (s *merchantSocialLinkService) CreateMerchant(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*response.MerchantSocialLinkResponse, *response.ErrorResponse) {
	const methd = "CreateMerchant"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, methd, attribute.Int("merchant_detail.id", *req.MerchantDetailID))

	defer end()

	merchant, err := s.merchantSocialLinkRepository.CreateSocialLink(ctx, req)
	if err != nil {
		return s.errorhandler.HandleCreateMerchantSocialLinkError(err, methd, "FAILED_CREATE_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantSocialLinkResponse(merchant)

	logSuccess("Merchant social link created", zap.Int("merchantID", merchant.ID))

	return so, nil
}

func (s *merchantSocialLinkService) UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*response.MerchantSocialLinkResponse, *response.ErrorResponse) {
	const methd = "UpdateMerchant"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, methd, attribute.Int("merchant_detail.id", *req.MerchantDetailID))

	defer end()

	merchant, err := s.merchantSocialLinkRepository.UpdateSocialLink(ctx, req)
	if err != nil {
		return s.errorhandler.HandleUpdateMerchantSocialLinkError(err, methd, "FAILED_UPDATE_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantSocialLinkResponse(merchant)

	logSuccess("Merchant social link updated", zap.Int("merchant.id", merchant.ID))

	return so, nil
}

func (s *merchantSocialLinkService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func() {
		s.recordMetrics(method, status, start)
		span.SetStatus(codes.Ok, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *merchantSocialLinkService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
