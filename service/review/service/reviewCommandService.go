package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-review/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type reviewCommandService struct {
	observability    observability.TraceLoggerObservability
	cache            cache.ReviewCommandCache
	reviewRepository      repository.ReviewCommandRepository
	reviewQueryRepository repository.ReviewQueryRepository
	userRepository        repository.UserQueryRepository
	productRepository     repository.ProductQueryRepository
	logger                logger.LoggerInterface
}

type ReviewCommandServiceDeps struct {
	Observability         observability.TraceLoggerObservability
	Cache                 cache.ReviewCommandCache
	ReviewRepository      repository.ReviewCommandRepository
	ReviewQueryRepository repository.ReviewQueryRepository
	UserRepository        repository.UserQueryRepository
	ProductRepository     repository.ProductQueryRepository
	Logger                logger.LoggerInterface
}

func NewReviewCommandService(deps *ReviewCommandServiceDeps) ReviewCommandService {
	return &reviewCommandService{
		observability:         deps.Observability,
		cache:                 deps.Cache,
		reviewRepository:      deps.ReviewRepository,
		reviewQueryRepository: deps.ReviewQueryRepository,
		userRepository:        deps.UserRepository,
		productRepository:     deps.ProductRepository,
		logger:                deps.Logger,
	}
}

func (s *reviewCommandService) Create(ctx context.Context, req *requests.CreateReviewRequest) (*db.CreateReviewRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", req.UserID),
		attribute.Int("product_id", req.ProductID))

	defer func() {
		end(status)
	}()

	_, err := s.userRepository.FindByID(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateReviewRow](
			s.logger,
			user_errors.ErrUserNotFound,
			method,
			span,
			zap.Int("user_id", req.UserID),
		)
	}

	_, err = s.productRepository.FindByID(ctx, req.ProductID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateReviewRow](
			s.logger,
			product_errors.ErrProductNotFound,
			method,
			span,
			zap.Int("product_id", req.ProductID),
		)
	}

	review, err := s.reviewRepository.Create(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateReviewRow](
			s.logger,
			review_errors.ErrFailedCreateReview.WithInternal(err),
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteReviewCache(ctx, int(review.ReviewID))

	logSuccess("Successfully created review",
		zap.Int("review_id", int(review.ReviewID)),
		zap.Int("user_id", req.UserID),
		zap.Int("product_id", req.ProductID))

	return review, nil
}

func (s *reviewCommandService) Update(ctx context.Context, req *requests.UpdateReviewRequest) (*db.UpdateReviewRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("review_id", *req.ReviewID))

	defer func() {
		end(status)
	}()

	_, err := s.reviewQueryRepository.FindByID(ctx, *req.ReviewID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateReviewRow](
			s.logger,
			review_errors.ErrReviewNotFound,
			method,
			span,
			zap.Int("review_id", *req.ReviewID),
		)
	}

	review, err := s.reviewRepository.Update(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateReviewRow](
			s.logger,
			review_errors.ErrFailedUpdateReview.WithInternal(err),
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteReviewCache(ctx, int(review.ReviewID))

	logSuccess("Successfully updated review",
		zap.Int("review_id", int(review.ReviewID)))

	return review, nil
}

func (s *reviewCommandService) Trash(ctx context.Context, reviewID int) (*db.Review, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("reviewID", reviewID))

	defer func() {
		end(status)
	}()

	review, err := s.reviewRepository.Trash(ctx, reviewID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Review](
			s.logger,
			review_errors.ErrFailedTrashedReview.WithInternal(err),
			method,
			span,

			zap.Int("reviewID", reviewID),
		)
	}

	s.cache.DeleteReviewCache(ctx, int(review.ReviewID))

	logSuccess("Successfully trashed review",
		zap.Int("review_id", int(review.ReviewID)))

	return review, nil
}

func (s *reviewCommandService) Restore(ctx context.Context, reviewID int) (*db.Review, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("reviewID", reviewID))

	defer func() {
		end(status)
	}()

	review, err := s.reviewRepository.Restore(ctx, reviewID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Review](
			s.logger,
			review_errors.ErrFailedRestoreReview.WithInternal(err),
			method,
			span,

			zap.Int("reviewID", reviewID),
		)
	}

	s.cache.DeleteReviewCache(ctx, int(review.ReviewID))

	logSuccess("Successfully restored review",
		zap.Int("review_id", int(review.ReviewID)))

	return review, nil
}

func (s *reviewCommandService) DeletePermanent(ctx context.Context, reviewID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("reviewID", reviewID))

	defer func() {
		end(status)
	}()

	success, err := s.reviewRepository.DeletePermanent(ctx, reviewID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_errors.ErrFailedDeletePermanentReview.WithInternal(err),
			method,
			span,

			zap.Int("reviewID", reviewID),
		)
	}

	s.cache.DeleteReviewCache(ctx, reviewID)

	logSuccess("Successfully permanently deleted review",
		zap.Int("review_id", reviewID))

	return success, nil
}

func (s *reviewCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_errors.ErrFailedRestoreAllReviews.WithInternal(err),
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed reviews")

	return success, nil
}

func (s *reviewCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			review_errors.ErrFailedDeleteAllPermanentReviews.WithInternal(err),
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all reviews")

	return success, nil
}
