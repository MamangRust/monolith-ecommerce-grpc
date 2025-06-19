package service

import (
	"context"
	"os"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	merchantsociallink_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDetailCommandService struct {
	ctx                             context.Context
	errorhandler                    errorhandler.MerchantDetailCommandError
	fileError                       errorhandler.FileError
	mencache                        mencache.MerchanrDetailCommandCache
	trace                           trace.Tracer
	merchantDetailQueryRepository   repository.MerchantDetailQueryRepository
	merchantDetailCommandRepository repository.MerchantDetailCommandRepository
	merchantSocialLinkRepository    repository.MerchantSocialLinkCommandRepository
	mapping                         response_service.MerchantDetailResponseMapper
	logger                          logger.LoggerInterface
	requestCounter                  *prometheus.CounterVec
	requestDuration                 *prometheus.HistogramVec
}

func NewMerchantDetailCommandService(ctx context.Context,
	errorhandler errorhandler.MerchantDetailCommandError,
	fileError errorhandler.FileError,
	mencache mencache.MerchanrDetailCommandCache,
	merchantDetailQueryRepository repository.MerchantDetailQueryRepository,
	merchantDetailCommandRepository repository.MerchantDetailCommandRepository,
	merchantSocialLinkRepository repository.MerchantSocialLinkCommandRepository,
	mapping response_service.MerchantDetailResponseMapper, logger logger.LoggerInterface) *merchantDetailCommandService {

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_detail_command_service_requests_total",
			Help: "Total number of requests to the MerchantDetailCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_detail_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantDetailCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantDetailCommandService{
		ctx:                             ctx,
		errorhandler:                    errorhandler,
		mencache:                        mencache,
		trace:                           otel.Tracer("merchant-detail-command-service"),
		merchantDetailCommandRepository: merchantDetailCommandRepository,
		merchantDetailQueryRepository:   merchantDetailQueryRepository,
		merchantSocialLinkRepository:    merchantSocialLinkRepository,
		mapping:                         mapping,
		logger:                          logger,
		requestCounter:                  requestCounter,
		requestDuration:                 requestDuration,
	}
}

func (s *merchantDetailCommandService) CreateMerchant(req *requests.CreateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse) {
	const methd = "CreateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(methd, attribute.Int("merchant.id", req.MerchantID))

	defer end()

	merchant, err := s.merchantDetailCommandRepository.CreateMerchantDetail(req)
	if err != nil {
		return s.errorhandler.HandleCreateMerchantDetailError(err, methd, "FAILED_CREATE_MERCHANT", span, &status, zap.Error(err))
	}

	for _, social := range req.SocialLink {
		social.MerchantDetailID = &merchant.ID
		_, err := s.merchantSocialLinkRepository.CreateSocialLink(social)
		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.MerchantDetailResponse](s.logger, err, methd, "FAILED_CREATE_SOCIAL_LINK", span, &status, merchantsociallink_errors.ErrFailedCreateMerchantSocialLink, zap.Error(err))
		}
	}

	so := s.mapping.ToMerchantDetailResponse(merchant)

	logSuccess("Merchant detail created", zap.Int("merchantID", merchant.ID))

	return so, nil
}

func (s *merchantDetailCommandService) UpdateMerchant(req *requests.UpdateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse) {
	const methd = "UpdateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(methd, attribute.Int("merchant.id", *req.MerchantDetailID))

	defer end()

	merchant, err := s.merchantDetailCommandRepository.UpdateMerchantDetail(req)
	if err != nil {
		return s.errorhandler.HandleUpdateMerchantDetailError(err, methd, "FAILED_UPDATE_MERCHANT", span, &status, zap.Error(err))
	}

	for _, social := range req.SocialLink {
		social.MerchantDetailID = &merchant.ID
		_, err := s.merchantSocialLinkRepository.UpdateSocialLink(social)
		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.MerchantDetailResponse](s.logger, err, methd, "FAILED_UPDATE_SOCIAL_LINK", span, &status, merchantsociallink_errors.ErrFailedUpdateMerchantSocialLink, zap.Error(err))
		}
	}

	so := s.mapping.ToMerchantDetailResponse(merchant)

	s.mencache.DeleteMerchantDetailCache(*req.MerchantDetailID)

	logSuccess("Merchant detail updated", zap.Int("merchant.id", merchant.ID))

	return so, nil
}

func (s *merchantDetailCommandService) TrashedMerchant(merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse) {
	const methd = "TrashedMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(methd, attribute.Int("merchant.id", merchantID))

	defer end()

	merchant, err := s.merchantDetailCommandRepository.TrashedMerchantDetail(merchantID)
	if err != nil {
		return s.errorhandler.HandleTrashedMerchantDetailError(err, methd, "FAILED_TRASH_MERCHANT", span, &status, zap.Error(err))
	}

	_, err = s.merchantSocialLinkRepository.TrashSocialLink(merchant.ID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDetailResponseDeleteAt](s.logger, err, methd, "FAILED_TRASH_SOCIAL_LINK", span, &status, merchantsociallink_errors.ErrFailedTrashMerchantSocialLink, zap.Error(err))
	}

	so := s.mapping.ToMerchantDetailResponseDeleteAt(merchant)

	s.mencache.DeleteMerchantDetailCache(merchantID)

	logSuccess("Merchant detail trashed", zap.Int("merchant.id", merchantID))

	return so, nil
}

func (s *merchantDetailCommandService) RestoreMerchant(merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse) {
	const methd = "RestoreMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(methd, attribute.Int("merchant.id", merchantID))

	defer end()

	merchant, err := s.merchantDetailCommandRepository.RestoreMerchantDetail(merchantID)
	if err != nil {
		return s.errorhandler.HandleRestoreMerchantDetailError(err, methd, "FAILED_RESTORE_MERCHANT", span, &status, zap.Error(err))
	}

	_, err = s.merchantSocialLinkRepository.RestoreSocialLink(merchant.ID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDetailResponseDeleteAt](s.logger, err, methd, "FAILED_RESTORE_SOCIAL_LINK", span, &status, merchantsociallink_errors.ErrFailedRestoreMerchantSocialLink, zap.Error(err))
	}

	logSuccess("Merchant detail restored", zap.Int("merchant.id", merchantID))

	so := s.mapping.ToMerchantDetailResponseDeleteAt(merchant)

	return so, nil
}

func (s *merchantDetailCommandService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
	const methd = "DeleteMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(methd, attribute.Int("merchant.id", merchantID))

	defer end()

	res, err := s.merchantDetailQueryRepository.FindByIdTrashed(merchantID)
	if err != nil {
		return s.errorhandler.HandleDeleteMerchantDetailError(err, methd, "FAILED_DELETE_MERCHANT_PERMANENT", span, &status, zap.Error(err))
	}

	if res.CoverImageUrl != "" {
		err := os.Remove(res.CoverImageUrl)
		if err != nil {
			if os.IsNotExist(err) {
				return s.fileError.HandleErrorFileCover(s.logger, err, methd, "FAILED_DELETE_COVER_IMAGE", span, &status, merchantdetail_errors.ErrFailedImageNotFound, zap.Error(err))
			} else {
				return s.fileError.HandleErrorFileCover(s.logger, err, "DeleteMerchantPermanent", "FAILED_REMOVE_COVER_IMAGE", span, &status, merchantdetail_errors.ErrFailedRemoveImageMerchantDetail, zap.Error(err))
			}
		} else {
			s.logger.Debug("Successfully deleted cover image",
				zap.String("cover_image_path", res.CoverImageUrl))
		}
	}

	if res.LogoUrl != "" {
		err := os.Remove(res.LogoUrl)
		if err != nil {
			if os.IsNotExist(err) {
				return s.fileError.HandleErrorFileCover(s.logger, err, methd, "FAILED_DELETE_LOGO_IMAGE", span, &status, merchantdetail_errors.ErrFailedLogoNotFound, zap.Error(err))
			} else {
				return s.fileError.HandleErrorFileCover(s.logger, err, methd, "FAILED_REMOVE_LOGO_IMAGE", span, &status, merchantdetail_errors.ErrFailedRemoveLogoMerchantDetail, zap.Error(err))
			}
		} else {
			s.logger.Debug("Successfully deleted logo image",
				zap.String("logo_path", res.LogoUrl))
		}
	}

	success, err := s.merchantDetailCommandRepository.DeleteMerchantDetailPermanent(merchantID)
	if err != nil {
		return s.errorhandler.HandleDeleteMerchantDetailError(err, methd, "FAILED_DELETE_MERCHANT_PERMANENT", span, &status, zap.Error(err))
	}

	_, err = s.merchantSocialLinkRepository.DeletePermanentSocialLink(merchantID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, methd, "FAILED_DELETE_SOCIAL_LINK", span, &status, merchantsociallink_errors.ErrFailedDeletePermanentMerchantSocialLink, zap.Error(err))
	}

	logSuccess("Merchant detail deleted permanently", zap.Int("merchant.id", merchantID))

	return success, nil
}

func (s *merchantDetailCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	const method = "RestoreAllMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer end()

	success, err := s.merchantDetailCommandRepository.RestoreAllMerchantDetail()
	if err != nil {
		return s.errorhandler.HandleRestoreAllMerchantDetailError(err, method, "FAILED_RESTORE_ALL_MERCHANT", span, &status, zap.Error(err))
	}

	_, err = s.merchantSocialLinkRepository.RestoreAllSocialLink()
	if err != nil {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_RESTORE_ALL_SOCIAL_LINK", span, &status, merchantsociallink_errors.ErrFailedRestoreAllMerchantSocialLinks, zap.Error(err))
	}

	logSuccess("All merchants restored", zap.Bool("success", success))

	return success, nil
}

func (s *merchantDetailCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer end()

	success, err := s.merchantDetailCommandRepository.DeleteAllMerchantDetailPermanent()
	if err != nil {
		return s.errorhandler.HandleDeleteAllMerchantDetailError(err, method, "FAILED_DELETE_ALL_MERCHANT_PERMANENT", span, &status, zap.Error(err))
	}

	_, err = s.merchantSocialLinkRepository.DeleteAllPermanentSocialLink()
	if err != nil {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_DELETE_ALL_SOCIAL_LINK", span, &status, merchantsociallink_errors.ErrFailedDeleteAllPermanentMerchantSocialLinks, zap.Error(err))
	}

	logSuccess("Successfully deleted all merchants permanently", zap.Bool("success", success))

	return success, nil
}

func (s *merchantDetailCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(),
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

	end := func() {
		s.recordMetrics(method, status, start)
		span.SetStatus(codes.Ok, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return span, end, status, logSuccess
}

func (s *merchantDetailCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
