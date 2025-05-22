package service

import (
	"context"
	"os"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewDetailCommandService struct {
	ctx                           context.Context
	trace                         trace.Tracer
	reviewDetailQueryRepository   repository.ReviewDetailQueryRepository
	reviewDetailCommandRepository repository.ReviewDetailCommandRepository
	mapping                       response_service.ReviewDetailResponeMapper
	logger                        logger.LoggerInterface
	requestCounter                *prometheus.CounterVec
	requestDuration               *prometheus.HistogramVec
}

func NewReviewDetailCommandService(
	ctx context.Context,
	reviewDetailQueryRepository repository.ReviewDetailQueryRepository,
	reviewDetailCommandRepository repository.ReviewDetailCommandRepository,
	mapping response_service.ReviewDetailResponeMapper,
	logger logger.LoggerInterface,
) *reviewDetailCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "review_detail_command_service_request_count",
			Help: "Total number of requests to the ReviewDetailCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "review_detail_command_service_request_duration",
			Help:    "Histogram of request durations for the ReviewDetailCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &reviewDetailCommandService{
		ctx:                           ctx,
		trace:                         otel.Tracer("review-detail-command-service"),
		reviewDetailQueryRepository:   reviewDetailQueryRepository,
		reviewDetailCommandRepository: reviewDetailCommandRepository,
		mapping:                       mapping,
		logger:                        logger,
		requestCounter:                requestCounter,
		requestDuration:               requestDuration,
	}
}

func (s *reviewDetailCommandService) CreateReviewDetail(req *requests.CreateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new Review Detail")

	res, err := s.reviewDetailCommandRepository.CreateReviewDetail(req)

	if err != nil {
		s.logger.Error("Failed to create new Review Detail",
			zap.Error(err),
			zap.Any("request", req))

		return nil, reviewdetail_errors.ErrFailedCreateReviewDetail
	}

	return s.mapping.ToReviewDetailsResponse(res), nil
}

func (s *reviewDetailCommandService) UpdateReviewDetail(req *requests.UpdateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating Review Detail", zap.Int("Review DetailID", *req.ReviewDetailID))

	res, err := s.reviewDetailCommandRepository.UpdateReviewDetail(req)

	if err != nil {
		s.logger.Error("Failed to update Review Detail",
			zap.Error(err),
			zap.Any("request", req))

		return nil, reviewdetail_errors.ErrFailedUpdateReviewDetail
	}

	return s.mapping.ToReviewDetailsResponse(res), nil
}

func (s *reviewDetailCommandService) TrashedReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing Review Detail", zap.Int("Review DetailID", review_id))

	res, err := s.reviewDetailCommandRepository.TrashedReviewDetail(review_id)

	if err != nil {
		s.logger.Error("Failed to move Review Detail to trash",
			zap.Error(err),
			zap.Int("Review Detail_id", review_id))

		return nil, reviewdetail_errors.ErrFailedTrashedReviewDetail
	}

	return s.mapping.ToReviewDetailResponseDeleteAt(res), nil
}

func (s *reviewDetailCommandService) RestoreReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring Review Detail", zap.Int("Review DetailID", review_id))

	res, err := s.reviewDetailCommandRepository.RestoreReviewDetail(review_id)

	if err != nil {
		s.logger.Error("Failed to restore Review Detail from trash",
			zap.Error(err),
			zap.Int("Review Detail_id", review_id))

		return nil, reviewdetail_errors.ErrFailedRestoreReviewDetail
	}

	return s.mapping.ToReviewDetailResponseDeleteAt(res), nil
}

func (s *reviewDetailCommandService) DeleteReviewDetailPermanent(review_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting Review Detail permanently", zap.Int("Review DetailID", review_id))

	res, err := s.reviewDetailQueryRepository.FindByIdTrashed(review_id)

	if err != nil {
		s.logger.Error("Failed to find review detail",
			zap.Int("review_id", review_id),
			zap.Error(err))

		return false, reviewdetail_errors.ErrFailedDeletePermanentReview
	}

	if res.Url != "" {
		err := os.Remove(res.Url)
		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Debug("review detail upload path file not found, continuing with review detail deletion",
					zap.String("upload path", res.Url))

				return false, reviewdetail_errors.ErrFailedImageNotFound
			} else {
				s.logger.Debug("Failed to delete review detail upload path",
					zap.String("upload path", res.Url),
					zap.Error(err))

				return false, reviewdetail_errors.ErrFailedRemoveImage
			}
		} else {
			s.logger.Debug("Successfully deleted review detail upload path",
				zap.String("upload path", res.Url))
		}
	}

	success, err := s.reviewDetailCommandRepository.DeleteReviewDetailPermanent(review_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete Review Detail",
			zap.Error(err),
			zap.Int("Review Detail_id", review_id))

		return false, reviewdetail_errors.ErrFailedDeletePermanentReview
	}

	return success, nil
}

func (s *reviewDetailCommandService) RestoreAllReviewDetail() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all trashed Review Details")

	success, err := s.reviewDetailCommandRepository.RestoreAllReviewDetail()

	if err != nil {
		s.logger.Error("Failed to restore all trashed Review Details",
			zap.Error(err))

		return false, reviewdetail_errors.ErrFailedRestoreAllReviewDetail
	}

	return success, nil
}

func (s *reviewDetailCommandService) DeleteAllReviewDetailPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all Review Details")

	success, err := s.reviewDetailCommandRepository.DeleteAllReviewDetailPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed Review Details",
			zap.Error(err))

		return false, reviewdetail_errors.ErrFailedDeleteAllReviewDetail
	}

	return success, nil
}
