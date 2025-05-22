package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type shippingAddressQueryService struct {
	ctx                            context.Context
	trace                          trace.Tracer
	shippingAddressQueryRepository repository.ShippingAddressQueryRepository
	mapping                        response_service.ShippingAddressResponseMapper
	logger                         logger.LoggerInterface
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewShippingAddressQueryService(ctx context.Context, shippingAddressQueryRepository repository.ShippingAddressQueryRepository, logger logger.LoggerInterface, mapping response_service.ShippingAddressResponseMapper) *shippingAddressQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipping_address_query_service_request_count",
			Help: "Total number of requests to the ShippingAddressQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shipping_address_query_service_request_duration",
			Help:    "Histogram of request durations for the ShippingAddressQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &shippingAddressQueryService{
		ctx:                            ctx,
		trace:                          otel.Tracer("shipping-address-query-service"),
		shippingAddressQueryRepository: shippingAddressQueryRepository,
		mapping:                        mapping,
		logger:                         logger,
		requestCounter:                 requestCounter,
		requestDuration:                requestDuration,
	}
}

func (s *shippingAddressQueryService) FindAll(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAll")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching shipping address",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	shipping, totalRecords, err := s.shippingAddressQueryRepository.FindAllShippingAddress(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_SHIPPING_ADDRESS")

		s.logger.Error("Failed to retrieve shipping address",
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve shipping address")

		status = "failed_find_all_shipping_address"

		return nil, nil, shippingaddress_errors.ErrFailedFindAllShippingAddresses
	}

	shippingRes := s.mapping.ToShippingAddressesResponse(shipping)

	s.logger.Debug("Successfully fetched shipping address",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return shippingRes, totalRecords, nil
}

func (s *shippingAddressQueryService) FindById(shipping_id int) (*response.ShippingAddressResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("shipping_id", shipping_id),
	)

	s.logger.Debug("Fetching Shipping Address by ID", zap.Int("shipping_id", shipping_id))

	shipping, err := s.shippingAddressQueryRepository.FindById(shipping_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SHIPPING_ADDRESS_BY_ID")

		s.logger.Error("Failed to retrieve Shipping Address details",
			zap.Int("shipping_id", shipping_id),
			zap.Error(err))

		span.SetAttributes(
			attribute.Int("shipping_id", shipping_id),
			attribute.String("traceID", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve Shipping Address details")

		status = "failed_find_shipping_address_by_id"

		return nil, shippingaddress_errors.ErrFailedFindShippingAddressByID
	}

	return s.mapping.ToShippingAddressResponse(shipping), nil
}

func (s *shippingAddressQueryService) FindByOrder(order_id int) (*response.ShippingAddressResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByOrder", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order_id", order_id),
	)

	s.logger.Debug("Fetching shipping address by order id", zap.Int("shipping_id", order_id))

	shipping, err := s.shippingAddressQueryRepository.FindByOrder(order_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SHIPPING_ADDRESS_BY_ORDER")

		s.logger.Error("Failed to retrieve shipping address by order id",
			zap.Int("order_id", order_id),
			zap.Error(err))

		span.SetAttributes(
			attribute.Int("order_id", order_id),
			attribute.String("traceID", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve shipping address by order id")

		status = "failed_find_shipping_address_by_order"

		return nil, shippingaddress_errors.ErrFailedFindShippingAddressByOrder
	}

	return s.mapping.ToShippingAddressResponse(shipping), nil
}

func (s *shippingAddressQueryService) FindByActive(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching categories",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	cashiers, totalRecords, err := s.shippingAddressQueryRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active Shipping Address",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, shippingaddress_errors.ErrFailedFindActiveShippingAddresses
	}

	s.logger.Debug("Successfully fetched shipping address",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToShippingAddressesResponseDeleteAt(cashiers), totalRecords, nil
}

func (s *shippingAddressQueryService) FindByTrashed(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching Shipping Address",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	shipping, totalRecords, err := s.shippingAddressQueryRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed shipping address",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, shippingaddress_errors.ErrFailedFindTrashedShippingAddresses
	}

	s.logger.Debug("Successfully fetched shipping address",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToShippingAddressesResponseDeleteAt(shipping), totalRecords, nil
}

func (s *shippingAddressQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
