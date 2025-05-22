package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
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
	ctx                    context.Context
	trace                  trace.Tracer
	cartCommandRepository  repository.CartCommandRepository
	productQueryRepository repository.ProductQueryRepository
	userQueryRepository    repository.UserQueryRepository
	mapping                response_service.CartResponseMapper
	logger                 logger.LoggerInterface
	requestCounter         *prometheus.CounterVec
	requestDuration        *prometheus.HistogramVec
}

func NewCardCommandService(ctx context.Context, cardCommandRepository repository.CartCommandRepository,
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
		ctx:                    ctx,
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

func (s *cardCommandService) CreateCart(req *requests.CreateCartRequest) (*response.CartResponse, *response.ErrorResponse) {
	start := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("CreateCart", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateCart")
	defer span.End()

	product, err := s.productQueryRepository.FindById(req.ProductID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCT_BY_ID")

		s.logger.Error("Failed to retrieve product details",
			zap.Error(err),
			zap.Int("product_id", req.ProductID),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve product details")
		status = "failed_to_retrieve_product_details"

		return nil, product_errors.ErrFailedFindProductById
	}

	_, err = s.userQueryRepository.FindById(req.UserID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("USER_NOT_FOUND")

		s.logger.Error("Failed to find user by id", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		status = "user_not_found"

		return nil, user_errors.ErrUserNotFoundRes
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

	res, err := s.cartCommandRepository.CreateCart(cartRecord)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_CART")

		s.logger.Error("Failed to create cart",
			zap.Error(err),
			zap.String("trace_id", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create cart")
		status = "failed_to_create_cart"

		return nil, cart_errors.ErrFailedCreateCart
	}

	so := s.mapping.ToCartResponse(res)

	return so, nil
}

func (s *cardCommandService) DeletePermanent(cart_id int) (bool, *response.ErrorResponse) {
	start := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("DeletePermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeletePermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting cart", zap.Int("cart_id", cart_id))

	success, err := s.cartCommandRepository.DeletePermanent(cart_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_CART")

		s.logger.Error("Failed to permanently delete cart",
			zap.Error(err),
			zap.String("trace_id", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete cart")
		status = "failed_to_permanently_delete_cart"

		return false, cart_errors.ErrFailedDeleteCart
	}

	return success, nil
}

func (s *cardCommandService) DeleteAllPermanently(req *requests.DeleteCartRequest) (bool, *response.ErrorResponse) {
	start := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllPermanently", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllPermanently")
	defer span.End()

	s.logger.Debug("Permanently deleting cart all", zap.Any("cart_id", req.CartIds))

	success, err := s.cartCommandRepository.DeleteAllPermanently(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_CARTS")

		s.logger.Error("Failed to permanently delete all carts",
			zap.Error(err),
			zap.String("trace_id", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all carts")
		status = "failed_to_permanently_delete_all_carts"

		return false, cart_errors.ErrFailedDeleteAllCarts
	}

	return success, nil
}

func (s *cardCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
