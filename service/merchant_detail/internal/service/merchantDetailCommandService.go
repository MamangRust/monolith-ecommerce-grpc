package service

import (
	"context"
	"os"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
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
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateMerchant")
	defer span.End()

	s.logger.Debug("Creating new merchant")

	merchant, err := s.merchantDetailCommandRepository.CreateMerchantDetail(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_MERCHANT")

		s.logger.Error("Failed to create merchant",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to create merchant")

		status = "failed_create_merchant_detail"

		return nil, merchantdetail_errors.ErrFailedCreateMerchantDetail
	}

	for _, social := range req.SocialLink {
		social.MerchantDetailID = &merchant.ID
		_, err := s.merchantSocialLinkRepository.CreateSocialLink(social)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_CREATE_SOCIAL_LINK")

			s.logger.Error("Failed to create social media link",
				zap.Error(err),
				zap.Any("social_link", social),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)

			span.SetStatus(codes.Error, "Failed to create social media link")

			status = "failed_create_social_link"

			return nil, merchantsociallink_errors.ErrFailedCreateMerchantSocialLink
		}
	}

	return s.mapping.ToMerchantDetailResponse(merchant), nil
}

func (s *merchantDetailCommandService) UpdateMerchant(req *requests.UpdateMerchantDetailRequest) (*response.MerchantDetailResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchantID", *req.MerchantDetailID),
	)

	s.logger.Debug("Updating merchant", zap.Int("merchantID", *req.MerchantDetailID))

	merchant, err := s.merchantDetailCommandRepository.UpdateMerchantDetail(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_MERCHANT")

		s.logger.Error("Failed to update merchant",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to update merchant")

		status = "failed_update_merchant_detail"

		return nil, merchantdetail_errors.ErrFailedUpdateMerchantDetail
	}

	for _, social := range req.SocialLink {
		social.MerchantDetailID = &merchant.ID
		_, err := s.merchantSocialLinkRepository.UpdateSocialLink(social)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SOCIAL_LINK")

			s.logger.Error("Failed to update social media link",
				zap.Error(err),
				zap.Any("social_link", social),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)

			span.SetStatus(codes.Error, "Failed to update social media link")

			status = "failed_update_social_link"

			return nil, merchantsociallink_errors.ErrFailedUpdateMerchantSocialLink
		}
	}

	return s.mapping.ToMerchantDetailResponse(merchant), nil
}

func (s *merchantDetailCommandService) TrashedMerchant(merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse) {
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

	merchant, err := s.merchantDetailCommandRepository.TrashedMerchantDetail(merchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_MERCHANT")

		s.logger.Error("Failed to trash merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to trash merchant")

		status = "failed_trash_merchant_detail"

		return nil, merchantdetail_errors.ErrFailedTrashedMerchantDetail
	}

	_, err = s.merchantSocialLinkRepository.TrashSocialLink(merchant.ID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_SOCIAL_LINK")

		s.logger.Error("Failed to trash merchant social link",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to trash merchant social link")

		status = "failed_trash_social_link"

		return nil, merchantsociallink_errors.ErrFailedTrashMerchantSocialLink
	}

	return s.mapping.ToMerchantDetailResponseDeleteAt(merchant), nil
}

func (s *merchantDetailCommandService) RestoreMerchant(merchantID int) (*response.MerchantDetailResponseDeleteAt, *response.ErrorResponse) {
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

	merchant, err := s.merchantDetailCommandRepository.RestoreMerchantDetail(merchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_MERCHANT")

		s.logger.Error("Failed to restore merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to restore merchant")

		status = "failed_restore_merchant_detail"

		return nil, merchantdetail_errors.ErrFailedRestoreMerchantDetail
	}

	_, err = s.merchantSocialLinkRepository.RestoreSocialLink(merchant.ID)
	if err != nil {
		s.logger.Debug("Failed to restore merchant social link", zap.Error(err), zap.Int("merchant_id", merchantID))

		return nil, merchantsociallink_errors.ErrFailedRestoreMerchantSocialLink
	}

	return s.mapping.ToMerchantDetailResponseDeleteAt(merchant), nil
}

func (s *merchantDetailCommandService) DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse) {
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

	res, err := s.merchantDetailQueryRepository.FindByIdTrashed(merchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_BY_ID")

		s.logger.Error("Failed to retrieve merchant by ID",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve merchant by ID")

		status = "failed_find_merchant_detail_by_id"

		return false, merchantdetail_errors.ErrFailedFindMerchantDetailById
	}

	if res.CoverImageUrl != "" {
		err := os.Remove(res.CoverImageUrl)
		if err != nil {
			if os.IsNotExist(err) {
				traceID := traceunic.GenerateTraceID("FAILED_DELETE_COVER_IMAGE")

				s.logger.Debug("Cover image not found, skipping delete",
					zap.String("cover_image_path", res.CoverImageUrl),
					zap.String("traceID", traceID))
				span.SetAttributes(
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)

				span.SetStatus(codes.Error, "Cover image not found, skipping delete")

				status = "failed_delete_cover_image"

				return false, merchantdetail_errors.ErrFailedImageNotFound
			} else {
				traceID := traceunic.GenerateTraceID("FAILED_REMOVE_COVER_IMAGE")

				s.logger.Error("Failed to delete cover image",
					zap.String("cover_image_path", res.CoverImageUrl),
					zap.Error(err),
					zap.String("traceID", traceID))

				span.SetAttributes(
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)

				span.SetStatus(codes.Error, "Failed to delete cover image")

				status = "failed_remove_cover_image"

				return false, merchantdetail_errors.ErrFailedRemoveImageMerchantDetail
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
				traceID := traceunic.GenerateTraceID("FAILED_DELETE_LOGO_IMAGE")

				s.logger.Debug("Logo image not found, skipping delete",
					zap.String("logo_path", res.LogoUrl),
					zap.String("traceID", traceID))
				span.SetAttributes(
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)

				span.SetStatus(codes.Error, "Logo image not found, skipping delete")

				status = "failed_delete_logo_image"

				return false, merchantdetail_errors.ErrFailedLogoNotFound
			} else {
				traceID := traceunic.GenerateTraceID("FAILED_REMOVE_LOGO_IMAGE")

				s.logger.Error("Failed to delete logo image",
					zap.String("logo_path", res.LogoUrl),
					zap.Error(err),
					zap.String("traceID", traceID))

				span.SetAttributes(
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)

				span.SetStatus(codes.Error, "Failed to delete logo image")

				status = "failed_remove_logo_image"

				return false, merchantdetail_errors.ErrFailedRemoveImageMerchantDetail
			}
		} else {
			s.logger.Debug("Successfully deleted logo image",
				zap.String("logo_path", res.LogoUrl))
		}
	}

	success, err := s.merchantDetailCommandRepository.DeleteMerchantDetailPermanent(merchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_MERCHANT")

		s.logger.Error("Failed to permanently delete merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to permanently delete merchant")

		status = "failed_delete_merchant"

		return false, merchantdetail_errors.ErrFailedDeleteMerchantDetailPermanent
	}

	_, err = s.merchantSocialLinkRepository.DeletePermanentSocialLink(merchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_SOCIAL_LINK")

		s.logger.Error("Failed to permanently delete social link",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to permanently delete social link")

		status = "failed_delete_social_link"

		return false, merchantsociallink_errors.ErrFailedDeletePermanentMerchantSocialLink
	}

	return success, nil
}

func (s *merchantDetailCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllMerchant")
	defer span.End()

	s.logger.Debug("Restoring all trashed merchants")

	success, err := s.merchantDetailCommandRepository.RestoreAllMerchantDetail()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_MERCHANT")

		s.logger.Error("Failed to restore all trashed merchants",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to restore all trashed merchants")

		status = "failed_restore_all_merchant"

		return false, merchantdetail_errors.ErrFailedRestoreAllMerchantDetail
	}

	_, err = s.merchantSocialLinkRepository.RestoreAllSocialLink()
	if err != nil {
		s.logger.Debug("Failed to restore all social links", zap.Error(err))
		return false, merchantsociallink_errors.ErrFailedRestoreAllMerchantSocialLinks
	}

	return success, nil
}

func (s *merchantDetailCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllMerchantPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllMerchantPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all merchants")

	success, err := s.merchantDetailCommandRepository.DeleteAllMerchantDetailPermanent()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_MERCHANT")

		s.logger.Error("Failed to permanently delete all merchants",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to permanently delete all merchants")

		status = "failed_delete_all_merchant"

		return false, merchantdetail_errors.ErrFailedDeleteAllMerchantDetailPermanent
	}

	_, err = s.merchantSocialLinkRepository.DeleteAllPermanentSocialLink()
	if err != nil {
		s.logger.Debug("Failed to delete all social links permanently", zap.Error(err))

		return false, merchantsociallink_errors.ErrFailedDeleteAllPermanentMerchantSocialLinks
	}

	return success, nil
}

func (s *merchantDetailCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
