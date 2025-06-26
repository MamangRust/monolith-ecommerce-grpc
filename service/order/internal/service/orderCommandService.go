package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-order/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderCommandService struct {
	ctx                              context.Context
	errorhandler                     errorhandler.OrderCommandError
	mencache                         mencache.OrderCommandCache
	trace                            trace.Tracer
	userQueryRepository              repository.UserQueryRepository
	productQueryRepository           repository.ProductQueryRepository
	productCommandRepository         repository.ProductCommandRepository
	orderQueryRepository             repository.OrderQueryRepository
	orderCommandRepository           repository.OrderCommandRepository
	orderItemQueryRepository         repository.OrderItemQueryRepository
	orderItemCommandRepositroy       repository.OrderItemCommandRepository
	merchantQueryRepository          repository.MerchantQueryRepository
	shippingAddressCommandRepository repository.ShippingAddressCommandRepository
	logger                           logger.LoggerInterface
	mapping                          response_service.OrderResponseMapper
	requestCounter                   *prometheus.CounterVec
	requestDuration                  *prometheus.HistogramVec
}

func NewOrderCommandService(
	ctx context.Context,
	errorhandler errorhandler.OrderCommandError,
	mencache mencache.OrderCommandCache,
	userQueryRepository repository.UserQueryRepository,
	productQueryRepository repository.ProductQueryRepository,
	productCommandRepository repository.ProductCommandRepository,
	orderQueryRepository repository.OrderQueryRepository,
	orderCommandRepository repository.OrderCommandRepository,
	orderItemQueryRepository repository.OrderItemQueryRepository,
	orderItemCommandRepositroy repository.OrderItemCommandRepository,
	merchantQueryRepository repository.MerchantQueryRepository,
	shippingAddressCommandRepository repository.ShippingAddressCommandRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderResponseMapper,

) *orderCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_command_service_request_count",
			Help: "Total number of requests to the OrderCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_command_service_request_duration",
			Help:    "Histogram of request durations for the OrderCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &orderCommandService{
		ctx:                              ctx,
		errorhandler:                     errorhandler,
		mencache:                         mencache,
		trace:                            otel.Tracer("order-command-service"),
		userQueryRepository:              userQueryRepository,
		productQueryRepository:           productQueryRepository,
		productCommandRepository:         productCommandRepository,
		orderQueryRepository:             orderQueryRepository,
		orderCommandRepository:           orderCommandRepository,
		orderItemQueryRepository:         orderItemQueryRepository,
		orderItemCommandRepositroy:       orderItemCommandRepositroy,
		merchantQueryRepository:          merchantQueryRepository,
		shippingAddressCommandRepository: shippingAddressCommandRepository,
		logger:                           logger,
		mapping:                          mapping,
		requestCounter:                   requestCounter,
		requestDuration:                  requestDuration,
	}
}

func (s *orderCommandService) CreateOrder(req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	const method = "CreateOrder"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", req.MerchantID), attribute.Int("user.id", req.UserID))

	defer func() {
		end(status)
	}()

	_, err := s.merchantQueryRepository.FindById(req.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	_, err = s.userQueryRepository.FindById(req.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	order, err := s.orderCommandRepository.CreateOrder(&requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		UserID:     req.UserID,
	})

	if err != nil {
		return s.errorhandler.HandleCreateOrderError(err, method, "FAILED_CREATE_ORDER", span, &status, zap.Error(err))
	}

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindById(item.ProductID)

		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_PRODUCT_BY_ID", span, &status, product_errors.ErrFailedFindProductById, zap.Error(err))
		}

		if product.CountInStock < item.Quantity {
			return s.errorhandler.HandleErrorInsufficientStockTemplate(err, method, "FAILED_INSUFFICIENT_STOCK", span, &status, order_errors.ErrInsufficientProductStock, zap.Error(err))
		}

		_, err = s.orderItemCommandRepositroy.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_CREATE_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedCreateOrderItem, zap.Error(err))
		}

		product.CountInStock -= item.Quantity
		_, err = s.productCommandRepository.UpdateProductCountStock(product.ID, product.CountInStock)

		if err != nil {
			return s.errorhandler.HandleErrorInvalidCountStockTemplate(err, method, "FAILED_UPDATE_PRODUCT_COUNT_STOCK", span, &status, product_errors.ErrFailedCountStock, zap.Error(err))
		}
	}

	_, err = s.shippingAddressCommandRepository.CreateShippingAddress(&requests.CreateShippingAddressRequest{
		OrderID:        &order.ID,
		Alamat:         req.ShippingAddress.Alamat,
		Provinsi:       req.ShippingAddress.Provinsi,
		Kota:           req.ShippingAddress.Kota,
		Courier:        req.ShippingAddress.Courier,
		ShippingMethod: req.ShippingAddress.ShippingMethod,
		ShippingCost:   req.ShippingAddress.ShippingCost,
		Negara:         req.ShippingAddress.Negara,
	})

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_CREATE_SHIPPING_ADDRESS", span, &status, shippingaddress_errors.ErrFailedCreateShippingAddress, zap.Error(err))
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(order.ID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_CALCULATE_TOTAL_PRICE", span, &status, orderitem_errors.ErrFailedCalculateTotal, zap.Error(err))
	}

	_, err = s.orderCommandRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    order.ID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_UPDATE_ORDER", span, &status, order_errors.ErrFailedUpdateOrder, zap.Error(err))
	}

	so := s.mapping.ToOrderResponse(order)

	logSuccess("Successfully create order", zap.Int("order.id", order.ID))

	return so, nil
}

func (s *orderCommandService) UpdateOrder(req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	const method = "UpdateOrder"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("order.id", *req.OrderID), attribute.Int("user.id", req.UserID))

	defer func() {
		end(status)
	}()

	existingOrder, err := s.orderQueryRepository.FindById(*req.OrderID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_ORDER_BY_ID", span, &status, order_errors.ErrFailedFindOrderById, zap.Error(err))
	}

	_, err = s.userQueryRepository.FindById(req.UserID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindById(item.ProductID)

		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_PRODUCT_BY_ID", span, &status, product_errors.ErrFailedFindProductById, zap.Error(err))
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemCommandRepositroy.UpdateOrderItem(&requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       product.Price,
			})

			if err != nil {
				return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_UPDATE_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedUpdateOrderItem, zap.Error(err))
			}
		} else {
			if product.CountInStock < item.Quantity {
				return s.errorhandler.HandleErrorInsufficientStockTemplate(err, method, "FAILED_INSUFFICIENT_STOCK", span, &status, order_errors.ErrInsufficientProductStock, zap.Error(err))
			}

			_, err := s.orderItemCommandRepositroy.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
				OrderID:   *req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})

			if err != nil {
				return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_CREATE_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedCreateOrderItem, zap.Error(err))
			}

			product.CountInStock -= item.Quantity
			_, err = s.productCommandRepository.UpdateProductCountStock(product.ID, product.CountInStock)

			if err != nil {
				return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_UPDATE_PRODUCT_COUNT_STOCK", span, &status, product_errors.ErrFailedCountStock, zap.Error(err))
			}
		}
	}

	_, err = s.shippingAddressCommandRepository.UpdateShippingAddress(&requests.UpdateShippingAddressRequest{
		ShippingID:     req.ShippingAddress.ShippingID,
		OrderID:        &existingOrder.ID,
		Alamat:         req.ShippingAddress.Alamat,
		Provinsi:       req.ShippingAddress.Provinsi,
		Kota:           req.ShippingAddress.Kota,
		Courier:        req.ShippingAddress.Courier,
		ShippingMethod: req.ShippingAddress.ShippingMethod,
		ShippingCost:   req.ShippingAddress.ShippingCost,
		Negara:         req.ShippingAddress.Negara,
	})

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_UPDATE_SHIPPING_ADDRESS", span, &status, shippingaddress_errors.ErrFailedUpdateShippingAddress, zap.Error(err))
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(*req.OrderID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_CALCULATE_TOTAL_PRICE", span, &status, orderitem_errors.ErrFailedCalculateTotal, zap.Error(err))
	}

	res, err := s.orderCommandRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    *req.OrderID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_UPDATE_ORDER", span, &status, order_errors.ErrFailedUpdateOrder, zap.Error(err))
	}

	so := s.mapping.ToOrderResponse(res)

	s.mencache.DeleteOrderCache(*req.OrderID)

	logSuccess("Successfully updated order", zap.Int("order.id", *req.OrderID))

	return so, nil
}

func (s *orderCommandService) TrashedOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedOrder"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("order.id", order_id))

	defer func() {
		end(status)
	}()

	order, err := s.orderQueryRepository.FindById(order_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_FIND_ORDER_BY_ID", span, &status, order_errors.ErrFailedFindOrderById, zap.Error(err))
	}

	if order.DeletedAt != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_NOT_DELETE_AT_ORDER", span, &status, order_errors.ErrFailedNotDeleteAtOrder, zap.Error(err))
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_FIND_ORDER_ITEM_BY_ORDER", span, &status, orderitem_errors.ErrFailedFindOrderItemByOrder, zap.Error(err))
	}

	for _, item := range orderItems {
		if item.DeletedAt != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_NOT_DELETE_AT_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedNotDeleteAtOrderItem, zap.Error(err))
		}

		trashedItem, err := s.orderItemCommandRepositroy.TrashedOrderItem(item.ID)

		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_TRASH_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedTrashedOrderItem, zap.Error(err))
		}

		s.logger.Debug("Order item trashed successfully",
			zap.Int("order_item_id", trashedItem.ID),
			zap.String("deleted_at", *trashedItem.DeletedAt))
	}

	trashedOrder, err := s.orderCommandRepository.TrashedOrder(order_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_TRASH_ORDER", span, &status, order_errors.ErrFailedTrashOrder, zap.Error(err))
	}

	so := s.mapping.ToOrderResponseDeleteAt(trashedOrder)

	s.mencache.DeleteOrderCache(order_id)

	logSuccess("Successfully trashed order", zap.Int("order.id", order_id))

	return so, nil
}

func (s *orderCommandService) RestoreOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreOrder"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("order.id", order_id))

	defer func() {
		end(status)
	}()

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_FIND_ORDER_ITEM_BY_ORDER", span, &status, orderitem_errors.ErrFailedFindOrderItemByOrder, zap.Error(err))
	}

	for _, item := range orderItems {
		_, err := s.orderItemCommandRepositroy.RestoreOrderItem(item.ID)

		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_RESTORE_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedRestoreOrderItem, zap.Error(err))
		}
	}

	order, err := s.orderCommandRepository.RestoreOrder(order_id)

	if err != nil {
		return s.errorhandler.HandleRestoreOrderError(err, method, "FAILED_RESTORE_ORDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderResponseDeleteAt(order)

	logSuccess("Successfully restored order", zap.Int("order.id", order_id))

	return so, nil
}

func (s *orderCommandService) DeleteOrderPermanent(order_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteOrderPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("order.id", order_id))

	defer func() {
		end(status)
	}()

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_FIND_ORDER_ITEM_BY_ORDER", span, &status, orderitem_errors.ErrFailedFindOrderItemByOrder, zap.Error(err))
	}

	for _, item := range orderItems {
		_, err := s.orderItemCommandRepositroy.
			DeleteOrderItemPermanent(item.ID)

		if err != nil {
			return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_DELETE_ORDER_ITEM_PERMANENT", span, &status, orderitem_errors.ErrFailedDeleteOrderItem, zap.Error(err))
		}
	}

	success, err := s.orderCommandRepository.DeleteOrderPermanent(order_id)

	if err != nil {
		return s.errorhandler.HandleDeleteOrderError(err, method, "FAILED_DELETE_ORDER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted order permanently", zap.Int("order.id", order_id))

	return success, nil
}

func (s *orderCommandService) RestoreAllOrder() (bool, *response.ErrorResponse) {
	const method = "RestoreAllOrder"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	successItems, err := s.orderItemCommandRepositroy.RestoreAllOrderItem()

	if err != nil || !successItems {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_RESTORE_ALL_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedRestoreAllOrderItem, zap.Error(err))
	}

	success, err := s.orderCommandRepository.RestoreAllOrder()
	if err != nil || !success {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_RESTORE_ALL_ORDER", span, &status, order_errors.ErrFailedRestoreAllOrder, zap.Error(err))
	}

	logSuccess("Successfully restored all orders", zap.Bool("success", success))

	return success, nil
}

func (s *orderCommandService) DeleteAllOrderPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllOrderPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	successItems, err := s.orderItemCommandRepositroy.DeleteAllOrderPermanent()

	if err != nil || !successItems {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_DELETE_ALL_ORDER_ITEM_PERMANENT", span, &status, orderitem_errors.ErrFailedDeleteAllOrderItem, zap.Error(err))
	}

	success, err := s.orderCommandRepository.DeleteAllOrderPermanent()

	if err != nil || !success {
		return s.errorhandler.HandleDeleteAllOrderError(err, method, "FAILED_DELETE_ALL_ORDER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all orders permanently", zap.Bool("success", success))

	return success, nil
}

func (s *orderCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *orderCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
