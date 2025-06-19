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
	ctx                           context.Context
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
	ctx context.Context,
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &reviewDetailCommandService{
		ctx:                           ctx,
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

func (s *reviewDetailCommandService) CreateReviewDetail(req *requests.CreateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	const method = "CreateReviewDetail"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("review.id", req.ReviewID))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailCommandRepository.CreateReviewDetail(req)

	if err != nil {
		return s.errorhandler.HandleCreateReviewDetailError(err, method, "FAILED_CREATE_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailsResponse(res)

	logSuccess("Successfully created Review Detail", zap.Int("reviewDetail.id", so.ID), zap.Int("review.id", req.ReviewID))

	return so, nil
}

func (s *reviewDetailCommandService) UpdateReviewDetail(req *requests.UpdateReviewDetailRequest) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	const method = "UpdateReviewDetail"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("reviewDetail.id", *req.ReviewDetailID))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailCommandRepository.UpdateReviewDetail(req)

	if err != nil {
		return s.errorhandler.HandleUpdateReviewDetailError(err, method, "FAILED_UPDATE_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailsResponse(res)

	logSuccess("Successfully updated Review Detail", zap.Int("reviewDetail.id", *req.ReviewDetailID))

	return so, nil
}

func (s *reviewDetailCommandService) TrashedReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedReviewDetail"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("reviewDetail.id", review_id))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailCommandRepository.TrashedReviewDetail(review_id)

	if err != nil {
		return s.errorhandler.HandleTrashedReviewDetailError(err, method, "FAILED_TRASH_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailResponseDeleteAt(res)

	logSuccess("Successfully trashed Review Detail", zap.Int("reviewDetail.id", review_id))

	return so, nil
}

func (s *reviewDetailCommandService) RestoreReviewDetail(review_id int) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreReviewDetail"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("reviewDetail.id", review_id))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailCommandRepository.RestoreReviewDetail(review_id)

	if err != nil {
		return s.errorhandler.HandleRestoreReviewDetailError(err, method, "FAILED_RESTORE_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	so := s.mapping.ToReviewDetailResponseDeleteAt(res)

	logSuccess("Successfully restored Review Detail", zap.Int("reviewDetail.id", review_id))

	return so, nil
}

func (s *reviewDetailCommandService) DeleteReviewDetailPermanent(review_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteReviewDetailPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("reviewDetail.id", review_id))

	defer func() {
		end(status)
	}()

	res, err := s.reviewDetailQueryRepository.FindByIdTrashed(review_id)

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

	success, err := s.reviewDetailCommandRepository.DeleteReviewDetailPermanent(review_id)

	if err != nil {
		return s.errorhandler.HandleDeleteReviewDetailError(err, method, "FAILED_DELETE_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted Review Detail", zap.Int("reviewDetail.id", review_id))

	return success, nil
}

func (s *reviewDetailCommandService) RestoreAllReviewDetail() (bool, *response.ErrorResponse) {
	const method = "RestoreAllReviewDetail"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewDetailCommandRepository.RestoreAllReviewDetail()

	if err != nil {
		return s.errorhandler.HandleRestoreAllReviewDetailError(err, method, "FAILED_RESTORE_ALL_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all Review Details", zap.Bool("success", success))

	return success, nil
}

func (s *reviewDetailCommandService) DeleteAllReviewDetailPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllReviewDetailPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.reviewDetailCommandRepository.DeleteAllReviewDetailPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllReviewDetailError(err, method, "FAILED_DELETE_ALL_PERMANENT_REVIEW_DETAIL", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all Review Details", zap.Bool("success", success))

	return success, nil
}

func (s *reviewDetailCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *reviewDetailCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
