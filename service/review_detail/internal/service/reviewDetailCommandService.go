package service

import (
	"context"
	"os"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/repository"
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

type reviewDetailCommandService struct {
	mencache                      mencache.ReviewDetailCommandCache
	errorhandler                  errorhandler.ReviewDetailCommandError
	trace                         trace.Tracer
	reviewDetailQueryRepository   repository.ReviewDetailQueryRepository
	reviewDetailCommandRepository repository.ReviewDetailCommandRepository
	mapping                       response_service.ReviewDetailResponeMapper
	logger                        logger.LoggerInterface
	requestCounter                *prometheus.CounterVec
	requestDuration               *prometheus.HistogramVec
}

func NewReviewDetailCommandService(
	mencache mencache.ReviewDetailCommandCache,
	errorhandler errorhandler.ReviewDetailCommandError,
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &reviewDetailCommandService{
		mencache:                      mencache,
		errorhandler:                  errorhandler,
		trace:                         otel.Tracer("review-detail-command-service"),
		reviewDetailQueryRepository:   reviewDetailQueryRepository,
		reviewDetailCommandRepository: reviewDetailCommandRepository,
		mapping:                       mapping,
		logger:                        logger,
		requestCounter:                requestCounter,
		requestDuration:               requestDuration,
	}
}

func (s *reviewDetailCommandService) CreateReviewDetail(ctx context.Context, req *requests.CreateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	const method = "CreateReviewDetail"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("review.id", req.ReviewID))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailCommandRepository.CreateReviewDetail(ctx, req)

	if err != nil {
		return s.errorhandler.HandleCreateReviewDetailError(err, method, "FAILED_CREATE_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailsResponse(res)

	logSuccess("Successfully created Review Detail", zap.Int("reviewDetail.id", so.ID), zap.Int("review.id", req.ReviewID))

	return so, nil
}

func (s *reviewDetailCommandService) UpdateReviewDetail(ctx context.Context, req *requests.UpdateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	const method = "UpdateReviewDetail"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("reviewDetail.id", *req.ReviewDetailID))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailCommandRepository.UpdateReviewDetail(ctx, req)

	if err != nil {
		return s.errorhandler.HandleUpdateReviewDetailError(err, method, "FAILED_UPDATE_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailsResponse(res)

	logSuccess("Successfully updated Review Detail", zap.Int("reviewDetail.id", *req.ReviewDetailID))

	return so, nil
}

func (s *reviewDetailCommandService) TrashedReviewDetail(ctx context.Context, review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedReviewDetail"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("reviewDetail.id", review_id))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailCommandRepository.TrashedReviewDetail(ctx, review_id)

	if err != nil {
		return s.errorhandler.HandleTrashedReviewDetailError(err, method, "FAILED_TRASH_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailResponseDeleteAt(res)

	logSuccess("Successfully trashed Review Detail", zap.Int("reviewDetail.id", review_id))

	return so, nil
}

func (s *reviewDetailCommandService) RestoreReviewDetail(ctx context.Context, review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreReviewDetail"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("reviewDetail.id", review_id))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailCommandRepository.RestoreReviewDetail(ctx, review_id)

	if err != nil {
		return s.errorhandler.HandleRestoreReviewDetailError(err, method, "FAILED_RESTORE_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailResponseDeleteAt(res)

	logSuccess("Successfully restored Review Detail", zap.Int("reviewDetail.id", review_id))

	return so, nil
}

func (s *reviewDetailCommandService) DeleteReviewDetailPermanent(ctx context.Context, review_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteReviewDetailPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("reviewDetail.id", review_id))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailQueryRepository.FindByIdTrashed(ctx, review_id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_DELETE_REVIEW_DETAIL", span, &status, nil, zap.Error(err))
	}

	if res.Url != "" {
		err := os.Remove(res.Url)
		if err != nil {
			if os.IsNotExist(err) {
				s.logger.Debug("File not found, but continuing with review detail deletion",
					zap.String("upload path", res.Url))
			} else {
				return s.errorhandler.HandleInvalidFileError(err, method, "FAILED_DELETE_REVIEW_DETAIL", res.Url, span, &status, zap.Error(err))
			}
		} else {
			s.logger.Debug("Successfully deleted review detail upload path",
				zap.String("upload path", res.Url))
		}
	}

	success, err := s.reviewDetailCommandRepository.DeleteReviewDetailPermanent(ctx, review_id)

	if err != nil {
		return s.errorhandler.HandleDeleteReviewDetailError(err, method, "FAILED_DELETE_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted Review Detail", zap.Int("reviewDetail.id", review_id))

	return success, nil
}

func (s *reviewDetailCommandService) RestoreAllReviewDetail(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllReviewDetail"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewDetailCommandRepository.RestoreAllReviewDetail(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllReviewDetailError(err, method, "FAILED_RESTORE_ALL_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all Review Details", zap.Bool("success", success))

	return success, nil
}

func (s *reviewDetailCommandService) DeleteAllReviewDetailPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllReviewDetailPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewDetailCommandRepository.DeleteAllReviewDetailPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllReviewDetailError(err, method, "FAILED_DELETE_ALL_PERMANENT_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all Review Details", zap.Bool("success", success))

	return success, nil
}

func (s *reviewDetailCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *reviewDetailCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
