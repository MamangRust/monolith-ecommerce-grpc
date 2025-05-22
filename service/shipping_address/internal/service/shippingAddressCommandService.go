package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
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
	trace                            trace.Tracer
	shippingAddressCommandRepository repository.ShippingAddressCommandRepository
	mapping                          response_service.ShippingAddressResponseMapper
	logger                           logger.LoggerInterface
	requestCounter                   *prometheus.CounterVec
	requestDuration                  *prometheus.HistogramVec
}

func NewShippingAddressCommandService(ctx context.Context, shippingAddressCommandRepository repository.ShippingAddressCommandRepository, logger logger.LoggerInterface, mapping response_service.ShippingAddressResponseMapper) *shippingAddressCommandService {
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
		trace:                            otel.Tracer("shipping-address-command-service"),
		shippingAddressCommandRepository: shippingAddressCommandRepository,
		mapping:                          mapping,
		logger:                           logger,
		requestCounter:                   requestCounter,
		requestDuration:                  requestDuration,
	}
}

func (s *shippingAddressCommandService) TrashShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashShippingAddress", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashShippingAddress")
	defer span.End()

	s.logger.Debug("Trashing shipping address", zap.Int("category", shipping_id))

	category, err := s.shippingAddressCommandRepository.TrashShippingAddress(shipping_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_SHIPPING_ADDRESS")

		s.logger.Error("Failed to trash shipping address",
			zap.Error(err),
			zap.Int("shipping_id", shipping_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("shipping_id", shipping_id),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash shipping address")
		status = "failed_trash_shipping_address"

		return nil, shippingaddress_errors.ErrFailedTrashShippingAddress
	}

	return s.mapping.ToShippingAddressResponseDeleteAt(category), nil
}

func (s *shippingAddressCommandService) RestoreShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreShippingAddress", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreShippingAddress")
	defer span.End()

	span.SetAttributes(
		attribute.Int("shipping_id", shipping_id),
	)

	s.logger.Debug("Restoring Shipping Address", zap.Int("shipping_id", shipping_id))

	shipping, err := s.shippingAddressCommandRepository.RestoreShippingAddress(shipping_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_SHIPPING_ADDRESS")

		s.logger.Error("Failed to restore shipping address",
			zap.Error(err),
			zap.Int("shipping_id", shipping_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("shipping_id", shipping_id),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore shipping address")
		status = "failed_restore_shipping_address"

		return nil, shippingaddress_errors.ErrFailedRestoreShippingAddress
	}

	return s.mapping.ToShippingAddressResponseDeleteAt(shipping), nil
}

func (s *shippingAddressCommandService) DeleteShippingAddressPermanently(shipping_id int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteShippingAddressPermanently", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteShippingAddressPermanently")
	defer span.End()

	span.SetAttributes(
		attribute.Int("shipping_id", shipping_id),
	)

	s.logger.Debug("Permanently deleting shipping address", zap.Int("shipping_id", shipping_id))

	success, err := s.shippingAddressCommandRepository.DeleteShippingAddressPermanently(shipping_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_SHIPPING_ADDRESS_PERMANENTLY")

		s.logger.Error("Failed to permanently delete shipping address",
			zap.Error(err),
			zap.Int("shipping_id", shipping_id),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.Int("shipping_id", shipping_id),
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete shipping address")
		status = "failed_delete_shipping_address_permanently"

		return false, shippingaddress_errors.ErrFailedDeleteShippingAddressPermanent
	}

	return success, nil
}

func (s *shippingAddressCommandService) RestoreAllShippingAddress() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllShippingAddress", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllShippingAddress")
	defer span.End()

	s.logger.Debug("Restoring all trashed shipping address")

	success, err := s.shippingAddressCommandRepository.RestoreAllShippingAddress()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_SHIPPING_ADDRESS")

		s.logger.Debug("Failed to restore all trashed shipping address", zap.Error(err), zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all trashed shipping address")
		status = "failed_restore_all_shipping_address"

		return false, shippingaddress_errors.ErrFailedRestoreAllShippingAddresses
	}

	return success, nil
}

func (s *shippingAddressCommandService) DeleteAllPermanentShippingAddress() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllPermanentShippingAddress", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllPermanentShippingAddress")
	defer span.End()

	s.logger.Debug("Permanently deleting all shipping address")

	success, err := s.shippingAddressCommandRepository.DeleteAllPermanentShippingAddress()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_SHIPPING_ADDRESS_PERMANENT")

		s.logger.Error("Failed to permanently delete all shipping address",
			zap.Error(err),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all shipping address")
		status = "failed_delete_all_shipping_address_permanent"

		return false, shippingaddress_errors.ErrFailedDeleteAllShippingAddressesPermanent
	}

	return success, nil
}

func (s *shippingAddressCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
