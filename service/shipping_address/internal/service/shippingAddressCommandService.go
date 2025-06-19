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
	ctx                              context.Context
	mencache                         mencache.ShippingAddressCommandCache
	errorhandler                     errorhandler.ShippingAddressCommandError
	trace                            trace.Tracer
	shippingAddressCommandRepository repository.ShippingAddressCommandRepository
	mapping                          response_service.ShippingAddressResponseMapper
	logger                           logger.LoggerInterface
	requestCounter                   *prometheus.CounterVec
	requestDuration                  *prometheus.HistogramVec
}

func NewShippingAddressCommandService(ctx context.Context,
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &shippingAddressCommandService{
		ctx:                              ctx,
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

func (s *shippingAddressCommandService) TrashShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashShippingAddress"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("shipping.id", shipping_id))

	defer func() {
		end(status)
	}()

	category, err := s.shippingAddressCommandRepository.TrashShippingAddress(shipping_id)

	if err != nil {
		return s.errorhandler.HandleTrashedShippingAddressError(err, method, "FAILED_TRASH_SHIPPING_ADDRESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToShippingAddressResponseDeleteAt(category)

	logSuccess("Successfully trashed shipping address", zap.Int("shipping.id", shipping_id), zap.Bool("success", true))

	return so, nil
}

func (s *shippingAddressCommandService) RestoreShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreShippingAddress"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("shipping.id", shipping_id))

	defer func() {
		end(status)
	}()

	shipping, err := s.shippingAddressCommandRepository.RestoreShippingAddress(shipping_id)

	if err != nil {
		return s.errorhandler.HandleRestoreShippingAddressError(err, method, "FAILED_RESTORE_SHIPPING_ADDRESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToShippingAddressResponseDeleteAt(shipping)

	logSuccess("Successfully restored shipping address", zap.Int("shipping.id", shipping_id), zap.Bool("success", true))

	return so, nil
}

func (s *shippingAddressCommandService) DeleteShippingAddressPermanently(shipping_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteShippingAddressPermanently"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("shipping.id", shipping_id))

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressCommandRepository.DeleteShippingAddressPermanently(shipping_id)

	if err != nil {
		return s.errorhandler.HandleDeleteShippingAddressError(err, method, "FAILED_DELETE_SHIPPING_ADDRESS_PERMANENTLY", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted shipping address permanently", zap.Int("shipping.id", shipping_id), zap.Bool("success", success))

	return success, nil
}

func (s *shippingAddressCommandService) RestoreAllShippingAddress() (bool, *response.ErrorResponse) {
	const method = "RestoreAllShippingAddress"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressCommandRepository.RestoreAllShippingAddress()

	if err != nil {
		return s.errorhandler.HandleRestoreAllShippingAddressError(err, method, "FAILED_RESTORE_ALL_TRASHED_SHIPPING_ADDRESS", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all trashed shipping address", zap.Bool("success", success))

	return success, nil
}

func (s *shippingAddressCommandService) DeleteAllPermanentShippingAddress() (bool, *response.ErrorResponse) {
	const method = "DeleteAllPermanentShippingAddress"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressCommandRepository.DeleteAllPermanentShippingAddress()

	if err != nil {
		return s.errorhandler.HandleDeleteAllShippingAddressError(err, method, "FAILED_DELETE_ALL_SHIPPING_ADDRESS_PERMANENTLY", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all shipping address permanently", zap.Bool("success", success))

	return success, nil
}

func (s *shippingAddressCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *shippingAddressCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
