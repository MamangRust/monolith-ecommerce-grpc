package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardCommandService struct {
	errorhandler           errorhandler.CartCommandError
	trace                  trace.Tracer
	cartCommandRepository  repository.CartCommandRepository
	productQueryRepository repository.ProductQueryRepository
	userQueryRepository    repository.UserQueryRepository
	mapping                response_service.CartResponseMapper
	logger                 logger.LoggerInterface
	requestCounter         *prometheus.CounterVec
	requestDuration        *prometheus.HistogramVec
}

func NewCardCommandService(errorhandler errorhandler.CartCommandError, cardCommandRepository repository.CartCommandRepository,
	productQueryRepository repository.ProductQueryRepository, userQueryRepository repository.UserQueryRepository,
	logger logger.LoggerInterface, mapping response_service.CartResponseMapper) *cardCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cart_command_service_request_count",
			Help: "Total number of requests to the CartCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cart_command_service_request_duration",
			Help:    "Histogram of request durations for the CartCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardCommandService{
		errorhandler:           errorhandler,
		trace:                  otel.Tracer("cart-command-service"),
		cartCommandRepository:  cardCommandRepository,
		productQueryRepository: productQueryRepository,
		userQueryRepository:    userQueryRepository,
		mapping:                mapping,
		logger:                 logger,
		requestCounter:         requestCounter,
		requestDuration:        requestDuration,
	}
}

func (s *cardCommandService) CreateCart(ctx context.Context, req *requests.CreateCartRequest) (*response.CartResponse, *response.ErrorResponse) {
	const method = "CreateCart"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("product.id", req.ProductID))

	defer func() {
		end(status)
	}()

	product, err := s.productQueryRepository.FindById(ctx, req.ProductID)

	if err != nil {
		return s.errorhandler.HandleCreateCartError(err, method, "FAILED_FIND_PRODUCT", span, &status, zap.Int("product.id", req.ProductID), zap.Error(err))
	}

	_, err = s.userQueryRepository.FindById(ctx, req.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.CartResponse](s.logger, err, method, "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes, zap.Int("user.id", req.UserID), zap.Error(err))
	}

	cartRecord := &requests.CartCreateRecord{
		ProductID:    req.ProductID,
		UserID:       req.UserID,
		Name:         product.Name,
		Price:        product.Price,
		ImageProduct: product.ImageProduct,
		Quantity:     req.Quantity,
		Weight:       product.Weight,
	}

	res, err := s.cartCommandRepository.CreateCart(ctx, cartRecord)

	if err != nil {
		return s.errorhandler.HandleCreateCartError(err, method, "FAILED_CREATE_CART", span, &status, zap.Int("product.id", req.ProductID), zap.Error(err))
	}

	so := s.mapping.ToCartResponse(res)

	logSuccess("Successfully created cart", zap.Int("cart.id", so.ID))

	return so, nil
}

func (s *cardCommandService) DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, *response.ErrorResponse) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("cart.id", req.CartID), attribute.Int("user.id", req.UserID))

	defer func() {
		end(status)
	}()

	success, err := s.cartCommandRepository.DeletePermanent(ctx, req)

	if err != nil {
		return s.errorhandler.HandleDeletePermanentError(err, method, "FAILED_DELETE_CART_PERMANENTLY", span, &status, zap.Int("cart.id", req.CartID), zap.Error(err))
	}

	logSuccess("Successfully permanently deleted cart", zap.Int("cart.id", req.CartID), zap.Bool("success", success))

	return success, nil
}

func (s *cardCommandService) DeleteAllPermanently(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, *response.ErrorResponse) {
	const method = "DeleteAllPermanently"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.cartCommandRepository.DeleteAllPermanently(ctx, req)

	if err != nil {
		return s.errorhandler.HandleDeleteAllPermanentlyError(err, method, "FAILED_DELETE_ALL_CART_PERMANENTLY", span, &status, zap.Any("cart.id", req.CartIds), zap.Error(err))
	}

	logSuccess("Successfully permanently deleted all cart", zap.Any("cart.id", req.CartIds), zap.Bool("success", success))

	return success, nil
}

func (s *cardCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

	s.logger.Info("Start: " + method)

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

func (s *cardCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
