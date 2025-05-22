package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
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
	trace                          trace.Tracer
	merchantAwardCommandRepository repository.MerchantAwardCommandRepository
	logger                         logger.LoggerInterface
	mapping                        response_service.MerchantAwardResponseMapper
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewMerchantAwardCommandService(ctx context.Context, merchantAwardCommandRepositroy repository.MerchantAwardCommandRepository, logger logger.LoggerInterface, mapping response_service.MerchantAwardResponseMapper) *merchantAwardCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_award_Command_service_request_total",
			Help: "Total number of requests to the MerchantAwardCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_award_Command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantAwardCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantAwardCommandService{
		ctx:                            ctx,
		trace:                          otel.Tracer("merchant-award-Command-service"),
		merchantAwardCommandRepository: merchantAwardCommandRepositroy,
		logger:                         logger,
		mapping:                        mapping,
		requestCounter:                 requestCounter,
		requestDuration:                requestDuration,
	}
}

func (s *merchantAwardCommandService) CreateMerchant(req *requests.CreateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateMerchant")
	defer span.End()

	s.logger.Debug("Creating new merchant")

	merchant, err := s.merchantAwardCommandRepository.CreateMerchantAward(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_MERCHANT_AWARD")

		s.logger.Error("Failed to create merchant",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create merchant")

		status = "failed_create_merchant_award"

		return nil, merchantaward_errors.ErrFailedCreateMerchantAward
	}

	return s.mapping.ToMerchantAwardResponse(merchant), nil
}

func (s *merchantAwardCommandService) UpdateMerchant(req *requests.UpdateMerchantCertificationOrAwardRequest) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchantID", *req.MerchantCertificationID),
	)

	s.logger.Debug("Updating merchant", zap.Int("merchantID", *req.MerchantCertificationID))

	merchant, err := s.merchantAwardCommandRepository.UpdateMerchantAward(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_MERCHANT_AWARD")

		s.logger.Error("Failed to update merchant",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update merchant")

		status = "failed_update_merchant_award"

		return nil, merchantaward_errors.ErrFailedUpdateMerchantAward
	}

	return s.mapping.ToMerchantAwardResponse(merchant), nil
}

func (s *merchantAwardCommandService) TrashedMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse) {
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

	merchant, err := s.merchantAwardCommandRepository.TrashedMerchantAward(merchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_MERCHANT_AWARD")

		s.logger.Error("Failed to trash merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash merchant")

		status = "failed_trash_merchant_award"

		return nil, merchantaward_errors.ErrFailedTrashedMerchantAward
	}

	return s.mapping.ToMerchantAwardResponseDeleteAt(merchant), nil
}

func (s *merchantAwardCommandService) RestoreMerchant(merchantID int) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse) {
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

	merchant, err := s.merchantAwardCommandRepository.RestoreMerchantAward(merchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_MERCHANT_AWARD")

		s.logger.Error("Failed to restore merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore merchant")

		status = "failed_restore_merchant_award"

		return nil, merchantaward_errors.ErrFailedRestoreMerchantAward
	}

	return s.mapping.ToMerchantAwardResponseDeleteAt(merchant), nil
}

func (s *merchantAwardCommandService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
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

	success, err := s.merchantAwardCommandRepository.DeleteMerchantPermanent(merchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_MERCHANT_AWARD_PERMANENT")

		s.logger.Error("Failed to permanently delete merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete merchant")

		status = "failed_delete_merchant_award_permanent"

		return false, merchantaward_errors.ErrFailedDeleteMerchantAwardPermanent
	}

	return success, nil
}

func (s *merchantAwardCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllMerchant")
	defer span.End()

	s.logger.Debug("Restoring all trashed merchants")

	success, err := s.merchantAwardCommandRepository.RestoreAllMerchantAward()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_MERCHANT_AWARD")

		s.logger.Error("Failed to restore all trashed merchants",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all trashed merchants")

		status = "failed_restore_all_merchant_award"

		return false, merchantaward_errors.ErrFailedRestoreAllMerchantAwards
	}

	return success, nil
}

func (s *merchantAwardCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllMerchantPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllMerchantPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all merchants")

	success, err := s.merchantAwardCommandRepository.DeleteAllMerchantAwardPermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_MERCHANT_AWARD_PERMANENT")

		s.logger.Error("Failed to permanently delete all merchants",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all merchants")

		status = "failed_delete_all_merchant_award_permanent"

		return false, merchantaward_errors.ErrFailedDeleteAllMerchantAwardsPermanent
	}

	return success, nil
}

func (s *merchantAwardCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
