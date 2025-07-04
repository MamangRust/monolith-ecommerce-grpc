package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	mencache                       mencache.ShippingAddressQueryCache
	errorhandler                   errorhandler.ShippingAddressQueryError
	trace                          trace.Tracer
	shippingAddressQueryRepository repository.ShippingAddressQueryRepository
	mapping                        response_service.ShippingAddressResponseMapper
	logger                         logger.LoggerInterface
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewShippingAddressQueryService(ctx context.Context,
	mencache mencache.ShippingAddressQueryCache,
	errorhandler errorhandler.ShippingAddressQueryError,
	shippingAddressQueryRepository repository.ShippingAddressQueryRepository, logger logger.LoggerInterface, mapping response_service.ShippingAddressResponseMapper) *shippingAddressQueryService {
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
		errorhandler:                   errorhandler,
		mencache:                       mencache,
		trace:                          otel.Tracer("shipping-address-query-service"),
		shippingAddressQueryRepository: shippingAddressQueryRepository,
		mapping:                        mapping,
		logger:                         logger,
		requestCounter:                 requestCounter,
		requestDuration:                requestDuration,
	}
}

func (s *shippingAddressQueryService) FindAll(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetShippingAddressAllCache(req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	shipping, totalRecords, err := s.shippingAddressQueryRepository.FindAllShippingAddress(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_SHIPPING_ADDRESS", span, &status, zap.Error(err))
	}

	shippingRes := s.mapping.ToShippingAddressesResponse(shipping)

	s.mencache.SetShippingAddressAllCache(req, shippingRes, totalRecords)

	logSuccess("Successfully fetched all shipping address", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return shippingRes, totalRecords, nil
}

func (s *shippingAddressQueryService) FindById(shipping_id int) (*response.ShippingAddressResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("shipping.id", shipping_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedShippingAddressCache(shipping_id); found {
		logSuccess("Successfully fetched shipping address from cache", zap.Int("shipping_id", shipping_id))

		return data, nil
	}

	shipping, err := s.shippingAddressQueryRepository.FindById(shipping_id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SHIPPING_ADDRESS_BY_ID", span, &status, shippingaddress_errors.ErrFailedFindShippingAddressByID, zap.Error(err))
	}

	so := s.mapping.ToShippingAddressResponse(shipping)

	s.mencache.SetCachedShippingAddressCache(so)

	logSuccess("Successfully fetched shipping address", zap.Int("shipping_id", shipping_id))

	return so, nil
}

func (s *shippingAddressQueryService) FindByOrder(order_id int) (*response.ShippingAddressResponse, *response.ErrorResponse) {
	const method = "FindByOrder"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("order.id", order_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedShippingAddressCache(order_id); found {
		logSuccess("Successfully fetched shipping address from cache", zap.Int("order.id", order_id))

		return data, nil
	}

	shipping, err := s.shippingAddressQueryRepository.FindByOrder(order_id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SHIPPING_ADDRESS_BY_ORDER", span, &status, shippingaddress_errors.ErrFailedFindShippingAddressByOrder, zap.Error(err))
	}

	so := s.mapping.ToShippingAddressResponse(shipping)

	s.mencache.SetCachedShippingAddressCache(so)

	logSuccess("Successfully fetched shipping address", zap.Int("order.id", order_id))

	return so, nil
}

func (s *shippingAddressQueryService) FindByActive(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetShippingAddressActiveCache(req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	cashiers, totalRecords, err := s.shippingAddressQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByActive", "FAILED_FIND_SHIPPING_ADDRESS", span, &status, shippingaddress_errors.ErrFailedFindActiveShippingAddresses, zap.Error(err))
	}

	so := s.mapping.ToShippingAddressesResponseDeleteAt(cashiers)

	s.mencache.SetShippingAddressActiveCache(req, so, totalRecords)

	logSuccess("Successfully fetched active shipping address", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *shippingAddressQueryService) FindByTrashed(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetShippingAddressTrashedCache(req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	shipping, totalRecords, err := s.shippingAddressQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_SHIPPING_ADDRESS", span, &status, shippingaddress_errors.ErrFailedFindTrashedShippingAddresses, zap.Error(err))
	}

	so := s.mapping.ToShippingAddressesResponseDeleteAt(shipping)
	s.mencache.SetShippingAddressTrashedCache(req, so, totalRecords)

	logSuccess("Successfully fetched trashed shipping address", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *shippingAddressQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *shippingAddressQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *shippingAddressQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
