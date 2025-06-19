package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	errorhandler            errorhandler.ReviewCommandError
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

func NewReviewCommandService(ctx context.Context,
	errorhandler errorhandler.ReviewCommandError,
	productQueryRepository repository.ProductQueryRepository, userQueryRepository repository.UserQueryRepository,
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
		errorhandler:            errorhandler,
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
	const method = "CreateReview"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("user.id", req.UserID), attribute.Int("product.id", req.ProductID))

	defer func() {
		end(status)
	}()

	_, err := s.userQueryRepository.FindById(req.UserID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes)
	}

	_, err = s.productQueryRepository.FindById(req.ProductID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_PRODUCT_BY_ID", span, &status, product_errors.ErrFailedFindProductById, zap.Error(err))
	}

	review, err := s.reviewCommandRepository.CreateReview(req)

	if err != nil {
		return s.errorhandler.HandleCreateReviewError(err, method, "FAILED_CREATE_REVIEW", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewResponse(review)

	logSuccess("Successfully created review", zap.Int("review.id", review.ID), zap.Bool("success", true))

	return so, nil
}

func (s *reviewCommandService) UpdateReview(req *requests.UpdateReviewRequest) (*response.ReviewResponse, *response.ErrorResponse) {
	const method = "UpdateReview"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("review.id", *req.ReviewID))

	defer func() {
		end(status)
	}()

	_, err := s.reviewQueryRepository.FindById(*req.ReviewID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_REVIEW_BY_ID", span, &status, review_errors.ErrFailedReviewNotFound, zap.Error(err))
	}

	review, err := s.reviewCommandRepository.UpdateReview(req)

	if err != nil {
		return s.errorhandler.HandleUpdateReviewError(err, method, "FAILED_UPDATE_REVIEW", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewResponse(review)

	logSuccess("Successfully updated review", zap.Int("review.id", review.ID), zap.Bool("success", true))

	return so, nil
}

func (s *reviewCommandService) TrashedReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedReview"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("review.id", reviewID))

	defer func() {
		end(status)
	}()

	review, err := s.reviewCommandRepository.TrashReview(reviewID)

	if err != nil {
		return s.errorhandler.HandleTrashedReviewError(err, method, "FAILED_TRASH_REVIEW", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewResponseDeleteAt(review)

	msgSuccess := "Successfully trashed Review"

	logSuccess(msgSuccess, zap.Int("review.id", review.ID), zap.Bool("success", true))

	return so, nil
}

func (s *reviewCommandService) RestoreReview(reviewID int) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreReview"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("review.id", reviewID))

	defer func() {
		end(status)
	}()

	review, err := s.reviewCommandRepository.RestoreReview(reviewID)

	if err != nil {
		return s.errorhandler.HandleRestoreReviewError(err, method, "FAILED_RESTORE_REVIEW", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewResponseDeleteAt(review)

	logSuccess("Successfully restored review", zap.Int("review.id", review.ID), zap.Bool("success", true))

	return so, nil
}

func (s *reviewCommandService) DeleteReviewPermanent(reviewID int) (bool, *response.ErrorResponse) {
	const method = "DeleteReviewPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("review.id", reviewID))

	defer func() {
		end(status)
	}()

	success, err := s.reviewCommandRepository.DeleteReviewPermanently(reviewID)

	if err != nil {
		return s.errorhandler.HandleDeleteReviewError(err, method, "FAILED_DELETE_PERMANENT_REVIEW", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted review permanently", zap.Int("review.id", reviewID), zap.Bool("success", success))

	return success, nil
}

func (s *reviewCommandService) RestoreAllReviews() (bool, *response.ErrorResponse) {
	const method = "RestoreAllReviews"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewCommandRepository.RestoreAllReview()

	if err != nil {
		return s.errorhandler.HandleRestoreAllReviewError(err, method, "FAILED_RESTORE_ALL_REVIEW", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all reviews", zap.Bool("success", success))

	return success, nil
}

func (s *reviewCommandService) DeleteAllReviewsPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllReviewsPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewCommandRepository.DeleteAllPermanentReview()

	if err != nil {
		return s.errorhandler.HandleDeleteAllReviewError(err, method, "FAILED_DELETE_ALL_PERMANENT_REVIEW", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all reviews permanently", zap.Bool("success", success))

	return success, nil
}

func (s *reviewCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
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

	return span, end, status, logSuccess
}

func (s *reviewCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
