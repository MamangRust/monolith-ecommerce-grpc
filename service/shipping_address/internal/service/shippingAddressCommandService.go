package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type shippingAddressCommandService struct {
	mencache                         mencache.ShippingAddressCommandCache
	errorhandler                     errorhandler.ShippingAddressCommandError
	trace                            trace.Tracer
	shippingAddressCommandRepository repository.ShippingAddressCommandRepository
	mapping                          response_service.ShippingAddressResponseMapper
	logger                           logger.LoggerInterface
	requestCounter                   *prometheus.CounterVec
	requestDuration                  *prometheus.HistogramVec
}

func NewShippingAddressCommandService(
	mencache mencache.ShippingAddressCommandCache,
	errorhandler errorhandler.ShippingAddressCommandError,
	shippingAddressCommandRepository repository.ShippingAddressCommandRepository, logger logger.LoggerInterface, mapping response_service.ShippingAddressResponseMapper) *shippingAddressCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipping_address_command_service_request_count",
			Help: "Total number of requests to the ShippingAddressCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shipping_address_command_service_request_duration",
			Help:    "Histogram of request durations for the ShippingAddressCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &shippingAddressCommandService{
		mencache:                         mencache,
		errorhandler:                     errorhandler,
		trace:                            otel.Tracer("shipping-address-command-service"),
		shippingAddressCommandRepository: shippingAddressCommandRepository,
		mapping:                          mapping,
		logger:                           logger,
		requestCounter:                   requestCounter,
		requestDuration:                  requestDuration,
	}
}

func (s *shippingAddressCommandService) TrashShippingAddress(ctx context.Context, shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashShippingAddress"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("shipping.id", shipping_id))

	defer func() {
		end(status)
	}()

	category, err := s.shippingAddressCommandRepository.TrashShippingAddress(ctx, shipping_id)

	if err != nil {
		return s.errorhandler.HandleTrashedShippingAddressError(err, method, "FAILED_TRASH_SHIPPING_ADDRESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToShippingAddressResponseDeleteAt(category)

	logSuccess("Successfully trashed shipping address", zap.Int("shipping.id", shipping_id), zap.Bool("success", true))

	return so, nil
}

func (s *shippingAddressCommandService) RestoreShippingAddress(ctx context.Context, shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreShippingAddress"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("shipping.id", shipping_id))

	defer func() {
		end(status)
	}()

	shipping, err := s.shippingAddressCommandRepository.RestoreShippingAddress(ctx, shipping_id)

	if err != nil {
		return s.errorhandler.HandleRestoreShippingAddressError(err, method, "FAILED_RESTORE_SHIPPING_ADDRESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToShippingAddressResponseDeleteAt(shipping)

	logSuccess("Successfully restored shipping address", zap.Int("shipping.id", shipping_id), zap.Bool("success", true))

	return so, nil
}

func (s *shippingAddressCommandService) DeleteShippingAddressPermanently(ctx context.Context, shipping_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteShippingAddressPermanently"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("shipping.id", shipping_id))

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressCommandRepository.DeleteShippingAddressPermanently(ctx, shipping_id)

	if err != nil {
		return s.errorhandler.HandleDeleteShippingAddressError(err, method, "FAILED_DELETE_SHIPPING_ADDRESS_PERMANENTLY", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted shipping address permanently", zap.Int("shipping.id", shipping_id), zap.Bool("success", success))

	return success, nil
}

func (s *shippingAddressCommandService) RestoreAllShippingAddress(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllShippingAddress"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressCommandRepository.RestoreAllShippingAddress(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllShippingAddressError(err, method, "FAILED_RESTORE_ALL_TRASHED_SHIPPING_ADDRESS", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all trashed shipping address", zap.Bool("success", success))

	return success, nil
}

func (s *shippingAddressCommandService) DeleteAllPermanentShippingAddress(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllPermanentShippingAddress"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressCommandRepository.DeleteAllPermanentShippingAddress(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllShippingAddressError(err, method, "FAILED_DELETE_ALL_SHIPPING_ADDRESS_PERMANENTLY", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all shipping address permanently", zap.Bool("success", success))

	return success, nil
}

func (s *shippingAddressCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
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

	return ctx, span, end, status, logSuccess
}

func (s *shippingAddressCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
