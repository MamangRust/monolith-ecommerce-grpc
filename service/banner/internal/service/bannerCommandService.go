package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type bannerCommandService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	bannerCommandRepository repository.BannerCommandRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.BannerResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewBannerCommandService(ctx context.Context, bannerCommand repository.BannerCommandRepository, logger logger.LoggerInterface, mapping response_service.BannerResponseMapper) *bannerCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "banner_command_service_request_total",
			Help: "Total number of requests to the BannerCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "banner_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the BannerCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &bannerCommandService{
		ctx:                     ctx,
		trace:                   otel.Tracer("banner-command-service"),
		bannerCommandRepository: bannerCommand,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *bannerCommandService) CreateBanner(req *requests.CreateBannerRequest) (*response.BannerResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateBanner", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateBanner")
	defer span.End()

	s.logger.Debug("Creating new Banner")

	Banner, err := s.bannerCommandRepository.CreateBanner(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_BANNER")

		s.logger.Error("Failed to create new Banner",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create new Banner")

		status = "failed_to_create_banner"

		return nil, banner_errors.ErrFailedCreateBanner
	}

	return s.mapping.ToBannerResponse(Banner), nil
}

func (s *bannerCommandService) UpdateBanner(req *requests.UpdateBannerRequest) (*response.BannerResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateBanner", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateBanner")
	defer span.End()

	span.SetAttributes(
		attribute.Int("BannerID", *req.BannerID),
	)

	s.logger.Debug("Updating Banner", zap.Int("BannerID", *req.BannerID))

	Banner, err := s.bannerCommandRepository.UpdateBanner(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_BANNER")

		s.logger.Error("Failed to update Banner",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update Banner")

		status = "failed_to_update_banner"

		return nil, banner_errors.ErrFailedUpdateBanner
	}

	return s.mapping.ToBannerResponse(Banner), nil
}

func (s *bannerCommandService) TrashedBanner(BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedBanner", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedBanner")
	defer span.End()

	span.SetAttributes(
		attribute.Int("BannerID", BannerID),
	)

	s.logger.Debug("Trashing Banner", zap.Int("BannerID", BannerID))

	Banner, err := s.bannerCommandRepository.TrashedBanner(BannerID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_BANNER")

		s.logger.Error("Failed to move Banner to trash",
			zap.Error(err),
			zap.Int("Banner_id", BannerID),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to move Banner to trash")

		status = "failed_to_move_banner_to_trash"

		return nil, banner_errors.ErrFailedTrashedBanner
	}

	return s.mapping.ToBannerResponseDeleteAt(Banner), nil
}

func (s *bannerCommandService) RestoreBanner(BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreBanner", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreBanner")
	defer span.End()

	s.logger.Debug("Restoring Banner", zap.Int("BannerID", BannerID))

	Banner, err := s.bannerCommandRepository.RestoreBanner(BannerID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_BANNER")

		s.logger.Error("Failed to restore Banner from trash",
			zap.Error(err),
			zap.Int("Banner_id", BannerID),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore Banner from trash")

		status = "failed_to_restore_banner_from_trash"

		return nil, banner_errors.ErrFailedRestoreBanner
	}

	return s.mapping.ToBannerResponseDeleteAt(Banner), nil
}

func (s *bannerCommandService) DeleteBannerPermanent(BannerID int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteBannerPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteBannerPermanent")
	defer span.End()

	s.logger.Debug("Deleting Banner permanently", zap.Int("BannerID", BannerID))

	success, err := s.bannerCommandRepository.DeleteBannerPermanent(BannerID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_BANNER_PERMANENT")

		s.logger.Error("Failed to permanently delete Banner",
			zap.Error(err),
			zap.Int("Banner_id", BannerID),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete Banner")

		status = "failed_to_permanently_delete_banner"

		return false, banner_errors.ErrFailedDeleteBanner
	}

	return success, nil
}

func (s *bannerCommandService) RestoreAllBanner() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllBanner", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllBanner")
	defer span.End()

	s.logger.Debug("Restoring all trashed Banners")

	success, err := s.bannerCommandRepository.RestoreAllBanner()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_BANNERS")

		s.logger.Error("Failed to restore all trashed Banners",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all trashed Banners")

		status = "failed_to_restore_all_trashed_banners"

		return false, banner_errors.ErrFailedRestoreAllBanners
	}

	return success, nil
}

func (s *bannerCommandService) DeleteAllBannerPermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllBannerPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllBannerPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all Banners")

	success, err := s.bannerCommandRepository.DeleteAllBannerPermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_BANNERS")

		s.logger.Error("Failed to permanently delete all trashed Banners",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all trashed Banners")

		status = "failed_to_permanently_delete_all_trashed_banners"

		return false, banner_errors.ErrFailedDeleteAllBanners
	}

	return success, nil
}

func (s *bannerCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
