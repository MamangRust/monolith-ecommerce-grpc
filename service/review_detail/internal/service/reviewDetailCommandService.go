package service

import (
	"context"
	"os"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	review_detail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type reviewDetailCommandService struct {
	observability               observability.TraceLoggerObservability
	cache                       cache.ReviewDetailCommandCache
	reviewDetailRepository      repository.ReviewDetailCommandRepository
	reviewDetailQueryRepository repository.ReviewDetailQueryRepository
	logger                      logger.LoggerInterface
}

type ReviewDetailCommandServiceDeps struct {
	Observability               observability.TraceLoggerObservability
	Cache                       cache.ReviewDetailCommandCache
	ReviewDetailRepository      repository.ReviewDetailCommandRepository
	ReviewDetailQueryRepository repository.ReviewDetailQueryRepository
	Logger                      logger.LoggerInterface
}

func NewReviewDetailCommandService(deps *ReviewDetailCommandServiceDeps) ReviewDetailCommandService {
	return &reviewDetailCommandService{
		observability:               deps.Observability,
		cache:                       deps.Cache,
		reviewDetailRepository:      deps.ReviewDetailRepository,
		reviewDetailQueryRepository: deps.ReviewDetailQueryRepository,
		logger:                      deps.Logger,
	}
}

func (s *reviewDetailCommandService) CreateReviewDetail(ctx context.Context, req *requests.CreateReviewDetailRequest) (*db.CreateReviewDetailRow, error) {
	const method = "CreateReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailRepository.CreateReviewDetail(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateReviewDetailRow](
			s.logger,
			review_detail_errors.ErrFailedCreateReviewDetail.WithInternal(err),
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, int(reviewDetail.ReviewDetailID))

	logSuccess("Successfully created review detail",
		zap.Int("review_detail_id", int(reviewDetail.ReviewDetailID)))

	return reviewDetail, nil
}

func (s *reviewDetailCommandService) UpdateReviewDetail(ctx context.Context, req *requests.UpdateReviewDetailRequest) (*db.UpdateReviewDetailRow, error) {
	const method = "UpdateReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_detail_id", *req.ReviewDetailID))

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailRepository.UpdateReviewDetail(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateReviewDetailRow](
			s.logger,
			review_detail_errors.ErrFailedUpdateReviewDetail.WithInternal(err),
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, int(reviewDetail.ReviewDetailID))

	logSuccess("Successfully updated review detail",
		zap.Int("review_detail_id", int(reviewDetail.ReviewDetailID)))

	return reviewDetail, nil
}

func (s *reviewDetailCommandService) TrashedReviewDetail(ctx context.Context, review_id int) (*db.ReviewDetail, error) {
	const method = "TrashedReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", review_id))

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailRepository.TrashedReviewDetail(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.ReviewDetail](
			s.logger,
			review_detail_errors.ErrFailedTrashedReviewDetail.WithInternal(err),
			method,
			span,

			zap.Int("review_id", review_id),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, int(reviewDetail.ReviewDetailID))

	logSuccess("Successfully trashed review detail",
		zap.Int("review_id", review_id))

	return reviewDetail, nil
}

func (s *reviewDetailCommandService) RestoreReviewDetail(ctx context.Context, review_id int) (*db.ReviewDetail, error) {
	const method = "RestoreReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", review_id))

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailRepository.RestoreReviewDetail(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.ReviewDetail](
			s.logger,
			review_detail_errors.ErrFailedRestoreReviewDetail.WithInternal(err),
			method,
			span,

			zap.Int("review_id", review_id),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, int(reviewDetail.ReviewDetailID))

	logSuccess("Successfully restored review detail",
		zap.Int("review_id", review_id))

	return reviewDetail, nil
}

func (s *reviewDetailCommandService) DeleteReviewDetailPermanent(ctx context.Context, review_id int) (bool, error) {
	const method = "DeleteReviewDetailPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", review_id))

	defer func() {
		end(status)
	}()

	reviewDetail, err := s.reviewDetailQueryRepository.FindByIdTrashed(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_detail_errors.ErrFailedFindReviewDetail.WithInternal(err),
			method,
			span,
			zap.Int("review_id", review_id),
		)
	}

	if reviewDetail.Url != "" {
		err := os.Remove(reviewDetail.Url)
		if err != nil {
			if os.IsNotExist(err) {
				status = "error"
				return errorhandler.HandleError[bool](
					s.logger,
					review_detail_errors.ErrFailedImageNotFound,
					method,
					span,
					zap.String("upload_path", reviewDetail.Url),
				)
			} else {
				status = "error"
				return errorhandler.HandleError[bool](
					s.logger,
					review_detail_errors.ErrFailedRemoveImage,
					method,
					span,
					zap.String("upload_path", reviewDetail.Url),
				)
			}
		}
	}

	success, err := s.reviewDetailRepository.DeleteReviewDetailPermanent(ctx, review_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_detail_errors.ErrFailedDeleteReviewDetailPermanent.WithInternal(err),
			method,
			span,
			zap.Int("review_id", review_id),
		)
	}

	s.cache.DeleteReviewDetailCache(ctx, review_id)

	logSuccess("Successfully permanently deleted review detail",
		zap.Int("review_id", review_id))

	return success, nil
}

func (s *reviewDetailCommandService) RestoreAllReviewDetail(ctx context.Context) (bool, error) {
	const method = "RestoreAllReviewDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewDetailRepository.RestoreAllReviewDetail(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_detail_errors.ErrFailedRestoreAllReviewDetails.WithInternal(err),
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed review details")

	return success, nil
}

func (s *reviewDetailCommandService) DeleteAllReviewDetailPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllReviewDetailPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewDetailRepository.DeleteAllReviewDetailPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_detail_errors.ErrFailedDeleteAllReviewDetails.WithInternal(err),
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed review details")

	return success, nil
}
