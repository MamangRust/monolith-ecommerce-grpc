package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/repository"
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

type bannerCommandService struct {
	errorhandler            errorhandler.BannerCommandError
	mencache                mencache.BannerCommandCache
	trace                   trace.Tracer
	bannerCommandRepository repository.BannerCommandRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.BannerResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewBannerCommandService(
	errorhandler errorhandler.BannerCommandError,
	mencache mencache.BannerCommandCache,
	bannerCommand repository.BannerCommandRepository, logger logger.LoggerInterface, mapping response_service.BannerResponseMapper) *bannerCommandService {
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
		errorhandler:            errorhandler,
		mencache:                mencache,
		trace:                   otel.Tracer("banner-command-service"),
		bannerCommandRepository: bannerCommand,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *bannerCommandService) CreateBanner(ctx context.Context, req *requests.CreateBannerRequest) (*response.BannerResponse, *response.ErrorResponse) {
	const method = "CreateBanner"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	Banner, err := s.bannerCommandRepository.CreateBanner(ctx, req)

	if err != nil {
		return s.errorhandler.HandleCreateBannerError(err, method, "FAILED_CREATE_BANNER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToBannerResponse(Banner)

	logSuccess("Successfully created Banner", zap.Int("banner.id", int(so.ID)))

	return so, nil
}

func (s *bannerCommandService) UpdateBanner(ctx context.Context, req *requests.UpdateBannerRequest) (*response.BannerResponse, *response.ErrorResponse) {
	const method = "UpdateBanner"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("banner.id", *req.BannerID))

	defer func() {
		end(status)
	}()

	Banner, err := s.bannerCommandRepository.UpdateBanner(ctx, req)

	if err != nil {
		return s.errorhandler.HandleUpdateBannerError(err, method, "FAILED_UPDATE_BANNER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToBannerResponse(Banner)

	s.mencache.DeleteBannerCache(ctx, *req.BannerID)

	logSuccess("Successfully updated Banner", zap.Int("banner.id", int(so.ID)))

	return so, nil
}

func (s *bannerCommandService) TrashedBanner(ctx context.Context, BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedBanner"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("banner.id", BannerID))

	defer func() {
		end(status)
	}()

	Banner, err := s.bannerCommandRepository.TrashedBanner(ctx, BannerID)

	if err != nil {
		return s.errorhandler.HandleTrashedBannerError(err, method, "FAILED_TRASH_BANNER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToBannerResponseDeleteAt(Banner)

	s.mencache.DeleteBannerCache(ctx, BannerID)

	logSuccess("Successfully trashed Banner", zap.Int("banner.id", int(so.ID)))

	return so, nil
}

func (s *bannerCommandService) RestoreBanner(ctx context.Context, BannerID int) (*response.BannerResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreBanner"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("banner.id", BannerID))

	defer func() {
		end(status)
	}()

	Banner, err := s.bannerCommandRepository.RestoreBanner(ctx, BannerID)

	if err != nil {
		return s.errorhandler.HandleRestoreBannerError(err, method, "FAILED_RESTORE_BANNER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToBannerResponseDeleteAt(Banner)

	s.mencache.DeleteBannerCache(ctx, BannerID)

	logSuccess("Successfully restored Banner", zap.Int("banner.id", int(so.ID)))

	return so, nil
}

func (s *bannerCommandService) DeleteBannerPermanent(ctx context.Context, BannerID int) (bool, *response.ErrorResponse) {
	const method = "DeleteBannerPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("banner.id", BannerID))

	defer func() {
		end(status)
	}()

	success, err := s.bannerCommandRepository.DeleteBannerPermanent(ctx, BannerID)

	if err != nil {
		return s.errorhandler.HandleDeleteBannerError(err, method, "FAILED_DELETE_BANNER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted Banner", zap.Int("banner.id", BannerID), zap.Bool("success", success))

	return success, nil
}

func (s *bannerCommandService) RestoreAllBanner(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllBanner"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.bannerCommandRepository.RestoreAllBanner(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllBannerError(err, method, "FAILED_RESTORE_ALL_TRASHED_BANNERS", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all trashed Banners", zap.Bool("success", success))

	return success, nil
}

func (s *bannerCommandService) DeleteAllBannerPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllBannerPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.bannerCommandRepository.DeleteAllBannerPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllBannerError(err, method, "FAILED_DELETE_ALL_BANNER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all Banners", zap.Bool("success", success))

	return success, nil
}

func (s *bannerCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *bannerCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
