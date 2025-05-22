package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewCommandService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	productQueryRepository  repository.ProductQueryRepository
	userQueryRepository     repository.UserQueryRepository
	reviewQueryRepository   repository.ReviewQueryRepository
	reviewCommandRepository repository.ReviewCommandRepository
	mapping                 response_service.ReviewResponseMapper
	logger                  logger.LoggerInterface
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewReviewCommandService(ctx context.Context, productQueryRepository repository.ProductQueryRepository, userQueryRepository repository.UserQueryRepository,
	reviewQueryRepository repository.ReviewQueryRepository,
	reviewCommandRepository repository.ReviewCommandRepository, mapping response_service.ReviewResponseMapper, logger logger.LoggerInterface) *reviewCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "review_command_service_request_count",
			Help: "Total number of requests to the ReviewCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "review_command_service_request_count",
			Help:    "Histogram of request durations for the ReviewCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &reviewCommandService{
		ctx:                     ctx,
		trace:                   otel.Tracer("review-command-service"),
		productQueryRepository:  productQueryRepository,
		userQueryRepository:     userQueryRepository,
		reviewQueryRepository:   reviewQueryRepository,
		reviewCommandRepository: reviewCommandRepository,
		mapping:                 mapping,
		logger:                  logger,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}
func (s *reviewCommandService) CreateReview(req *requests.CreateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateReview", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateReview")
	defer span.End()

	s.logger.Debug("Creating new cashier")

	_, err := s.userQueryRepository.FindById(req.UserID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_USER_BY_ID")

		s.logger.Error("User not found for review creation",
			zap.Int("user_id", req.UserID),
			zap.Error(err),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found for review creation")
		status = "failed_find_user_by_id"

		return nil, user_errors.ErrUserNotFoundRes
	}

	_, err = s.productQueryRepository.FindById(req.ProductID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCT_BY_ID")

		s.logger.Error("Product not found for review creation",
			zap.Int("product_id", req.ProductID),
			zap.Error(err),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Product not found for review creation")
		status = "failed_find_product_by_id"

		return nil, product_errors.ErrFailedFindProductById
	}

	review, err := s.reviewCommandRepository.CreateReview(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_REVIEW")

		s.logger.Error("Failed to create review",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create review")
		status = "failed_create_review"

		return nil, review_errors.ErrFailedCreateReview
	}

	return s.mapping.ToReviewResponse(review), nil
}

func (s *reviewCommandService) UpdateReview(req *requests.UpdateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateReview", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateReview")
	defer span.End()

	span.SetAttributes(
		attribute.Int("review_id", *req.ReviewID),
	)

	s.logger.Debug("Updating review", zap.Int("review_id", *req.ReviewID))

	_, err := s.reviewQueryRepository.FindById(*req.ReviewID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_REVIEW_BY_ID")

		s.logger.Error("Review not found for review update",
			zap.Int("review_id", *req.ReviewID),
			zap.Error(err),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Review not found for review update")
		status = "failed_find_review_by_id"

		return nil, review_errors.ErrFailedReviewNotFound
	}

	review, err := s.reviewCommandRepository.UpdateReview(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_REVIEW")

		s.logger.Error("Failed to update review",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update review")
		status = "failed_update_review"

		return nil, review_errors.ErrFailedUpdateReview
	}

	return s.mapping.ToReviewResponse(review), nil
}

func (s *reviewCommandService) TrashedReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedReview", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedReview")
	defer span.End()

	span.SetAttributes(
		attribute.Int("reviewID", reviewID),
	)

	s.logger.Debug("Trashing review", zap.Int("reviewID", reviewID))

	review, err := s.reviewCommandRepository.TrashReview(reviewID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_REVIEW")

		s.logger.Error("Failed to trash review",
			zap.Error(err),
			zap.Int("reviewID", reviewID),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.Int("reviewID", reviewID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash review")
		status = "failed_trash_review"

		return nil, review_errors.ErrFailedTrashedReview
	}

	return s.mapping.ToReviewResponseDeleteAt(review), nil
}

func (s *reviewCommandService) RestoreReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreReview", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreReview")
	defer span.End()

	s.logger.Debug("Restoring review", zap.Int("reviewID", reviewID))

	review, err := s.reviewCommandRepository.RestoreReview(reviewID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_REVIEW")

		s.logger.Error("Failed to restore review",
			zap.Error(err),
			zap.Int("reviewID", reviewID),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.Int("reviewID", reviewID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore review")
		status = "failed_restore_review"

		return nil, review_errors.ErrFailedRestoreReview
	}

	return s.mapping.ToReviewResponseDeleteAt(review), nil
}

func (s *reviewCommandService) DeleteReviewPermanent(reviewID int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteReviewPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteReviewPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("reviewID", reviewID),
	)

	s.logger.Debug("Permanently deleting review", zap.Int("reviewID", reviewID))

	success, err := s.reviewCommandRepository.DeleteReviewPermanently(reviewID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_PERMANENT_REVIEW")

		s.logger.Error("Failed to permanently delete review",
			zap.Error(err),
			zap.Int("reviewID", reviewID),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.Int("reviewID", reviewID),
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete review")
		status = "failed_delete_permanent_review"

		return false, review_errors.ErrFailedDeletePermanentReview
	}
	return success, nil
}

func (s *reviewCommandService) RestoreAllReviews() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllReviews", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllReviews")
	defer span.End()

	s.logger.Debug("Restoring all trashed reviews")

	success, err := s.reviewCommandRepository.RestoreAllReview()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_REVIEW")

		s.logger.Error("Failed to restore all trashed reviews",
			zap.Error(err),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all trashed reviews")
		status = "failed_restore_all_reviews"

		return false, review_errors.ErrFailedRestoreAllReviews
	}

	return success, nil
}

func (s *reviewCommandService) DeleteAllReviewsPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllReviewsPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllReviewsPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all reviews")

	success, err := s.reviewCommandRepository.DeleteAllPermanentReview()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_REVIEW_PERMANENT")

		s.logger.Error("Failed to permanently delete all reviews",
			zap.Error(err),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all reviews")
		status = "failed_delete_all_permanent_reviews"

		return false, review_errors.ErrFailedDeleteAllPermanentReviews
	}

	return success, nil
}

func (s *reviewCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
