package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-grpc-order/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	traceunic "github.com/MamangRust/monolith-point-of-sale-pkg/trace_unic"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	shippingaddress_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/shipping_address_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderCommandService struct {
	ctx                              context.Context
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
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateOrder", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchantID", req.MerchantID),
		attribute.Int("userID", req.UserID),
	)

	s.logger.Debug("Creating new order with items", zap.Int("merchantID", req.MerchantID), zap.Int("userID", req.UserID))

	_, err := s.merchantQueryRepository.FindById(req.MerchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_BY_ID")

		s.logger.Error("Failed to find merchant by ID",
			zap.Error(err),
			zap.Int("merchant_id", req.MerchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("merchant_id", req.MerchantID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find merchant by ID")

		status = "failed_find_merchant_by_id"

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	_, err = s.userQueryRepository.FindById(req.UserID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("USER_NOT_FOUND")

		s.logger.Error("User not found",
			zap.Error(err),
			zap.Int("user_id", req.UserID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("user_id", req.UserID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")

		status = "user_not_found"

		return nil, user_errors.ErrUserNotFoundRes
	}

	order, err := s.orderCommandRepository.CreateOrder(&requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		UserID:     req.UserID,
	})

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_ORDER")

		s.logger.Error("Failed to create order",
			zap.Error(err),
			zap.Int("merchant_id", req.MerchantID),
			zap.Int("user_id", req.UserID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("merchant_id", req.MerchantID),
			attribute.Int("user_id", req.UserID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create order")

		status = "failed_create_order"

		return nil, order_errors.ErrFailedCreateOrder
	}

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindById(item.ProductID)

		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCT_BY_ID")

			s.logger.Error("Failed to find product by ID",
				zap.Error(err),
				zap.Int("product_id", item.ProductID),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.Int("product_id", item.ProductID),
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to find product by ID")

			status = "failed_find_product_by_id"

			return nil, product_errors.ErrFailedFindProductById
		}

		if product.CountInStock < item.Quantity {
			traceID := traceunic.GenerateTraceID("INSUFFICIENT_PRODUCT_STOCK")

			s.logger.Error("Insufficient product stock",
				zap.Error(err),
				zap.Int("product_id", item.ProductID),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.Int("product_id", item.ProductID),
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Insufficient product stock")

			status = "insufficient_product_stock"

			return nil, order_errors.ErrInsufficientProductStock
		}

		_, err = s.orderItemCommandRepositroy.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_CREATE_ORDER_ITEM")

			s.logger.Error("Failed to create order item",
				zap.Error(err),
				zap.Int("order_id", order.ID),
				zap.Int("product_id", item.ProductID),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.Int("order_id", order.ID),
				attribute.Int("product_id", item.ProductID),
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to create order item")

			status = "failed_create_order_item"

			return nil, orderitem_errors.ErrFailedCreateOrderItem
		}

		product.CountInStock -= item.Quantity
		_, err = s.productCommandRepository.UpdateProductCountStock(product.ID, product.CountInStock)

		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_COUNT_STOCK")

			s.logger.Error("Failed to count stock",
				zap.Error(err),
				zap.Int("product_id", product.ID),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.Int("product_id", product.ID),
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to count stock")

			status = "failed_count_stock"

			return nil, product_errors.ErrFailedCountStock
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
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_SHIPPING_ADDRESS")

		s.logger.Error("Failed to create shipping address",
			zap.Error(err),
			zap.Int("order_id", order.ID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order.ID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create shipping address")

		status = "failed_create_shipping_address"

		return nil, shippingaddress_errors.ErrFailedCreateShippingAddress
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(order.ID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CALCULATE_TOTAL_PRICE")

		s.logger.Error("Failed to calculate total price",
			zap.Error(err),
			zap.Int("order_id", order.ID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order.ID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to calculate total price")

		status = "failed_calculate_total_price"

		return nil, orderitem_errors.ErrFailedCalculateTotal
	}

	_, err = s.orderCommandRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    order.ID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_ORDER")

		s.logger.Error("Failed to update order",
			zap.Error(err),
			zap.Int("order_id", order.ID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order.ID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update order")

		status = "failed_update_order"

		return nil, order_errors.ErrFailedUpdateOrder
	}

	return s.mapping.ToOrderResponse(order), nil
}

func (s *orderCommandService) UpdateOrder(req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateOrder", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("orderID", *req.OrderID),
	)

	s.logger.Debug("Updating order with items", zap.Int("orderID", *req.OrderID))

	existingOrder, err := s.orderQueryRepository.FindById(*req.OrderID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_BY_ID")

		s.logger.Error("Failed to retrieve order details",
			zap.Error(err),
			zap.Int("order_id", *req.OrderID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", *req.OrderID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve order details")

		status = "failed_find_order_by_id"

		return nil, order_errors.ErrFailedFindOrderById
	}

	_, err = s.userQueryRepository.FindById(req.UserID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_USER_BY_ID")

		s.logger.Error("Failed to retrieve user details",
			zap.Error(err),
			zap.Int("user_id", req.UserID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("user_id", req.UserID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve user details")

		status = "failed_find_user_by_id"

		return nil, user_errors.ErrUserNotFoundRes
	}

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindById(item.ProductID)

		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCT_BY_ID")

			s.logger.Error("Failed to retrieve product details",
				zap.Error(err),
				zap.Int("product_id", item.ProductID),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.Int("product_id", item.ProductID),
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to retrieve product details")

			status = "failed_find_product_by_id"

			return nil, product_errors.ErrFailedFindProductById
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemCommandRepositroy.UpdateOrderItem(&requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       product.Price,
			})

			if err != nil {
				traceID := traceunic.GenerateTraceID("FAILED_UPDATE_ORDER_ITEM")

				s.logger.Error("Failed to update order item",
					zap.Error(err),
					zap.Int("order_item_id", item.OrderItemID),
					zap.String("traceID", traceID))

				span.SetAttributes(
					attribute.Int("order_item_id", item.OrderItemID),
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)
				span.SetStatus(codes.Error, "Failed to update order item")

				status = "failed_update_order_item"

				return nil, orderitem_errors.ErrFailedUpdateOrderItem
			}
		} else {
			if product.CountInStock < item.Quantity {
				traceID := traceunic.GenerateTraceID("INSUFFICIENT_STOCK")

				s.logger.Error("Insufficient product stock",
					zap.Error(err),
					zap.Int("product_id", item.ProductID),
					zap.String("traceID", traceID))

				span.SetAttributes(
					attribute.Int("product_id", item.ProductID),
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)
				span.SetStatus(codes.Error, "Insufficient product stock")

				status = "insufficient_stock"

				return nil, order_errors.ErrInsufficientProductStock
			}

			_, err := s.orderItemCommandRepositroy.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
				OrderID:   *req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})

			if err != nil {
				traceID := traceunic.GenerateTraceID("FAILED_CREATE_ORDER_ITEM")

				s.logger.Error("Failed to create order item",
					zap.Error(err),
					zap.Int("order_id", *req.OrderID),
					zap.String("traceID", traceID))

				span.SetAttributes(
					attribute.Int("order_id", *req.OrderID),
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)
				span.SetStatus(codes.Error, "Failed to create order item")

				status = "failed_create_order_item"

				return nil, orderitem_errors.ErrFailedCreateOrderItem
			}

			product.CountInStock -= item.Quantity
			_, err = s.productCommandRepository.UpdateProductCountStock(product.ID, product.CountInStock)

			if err != nil {
				traceID := traceunic.GenerateTraceID("FAILED_COUNT_STOCK")

				s.logger.Error("Failed to count stock",
					zap.Error(err),
					zap.Int("product_id", item.ProductID),
					zap.String("traceID", traceID))

				span.SetAttributes(
					attribute.Int("product_id", item.ProductID),
					attribute.String("traceID", traceID),
				)

				span.RecordError(err)
				span.SetStatus(codes.Error, "Failed to count stock")

				status = "failed_count_stock"

				return nil, product_errors.ErrFailedCountStock
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
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SHIPPING_ADDRESS")

		s.logger.Error("Failed to update shipping address",
			zap.Error(err),
			zap.Int("shipping_id", *req.ShippingAddress.ShippingID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("shipping_id", *req.ShippingAddress.ShippingID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update shipping address")

		status = "failed_update_shipping_address"

		return nil, shippingaddress_errors.ErrFailedUpdateShippingAddress
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(*req.OrderID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CALCULATE_TOTAL")

		s.logger.Error("Failed to calculate total price",
			zap.Error(err),
			zap.Int("order_id", *req.OrderID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", *req.OrderID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to calculate total price")

		status = "failed_calculate_total"

		return nil, orderitem_errors.ErrFailedCalculateTotal
	}

	_, err = s.orderCommandRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    *req.OrderID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_ORDER")

		s.logger.Error("Failed to update order",
			zap.Error(err),
			zap.Int("order_id", *req.OrderID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", *req.OrderID),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update order")

		status = "failed_update_order"

		return nil, order_errors.ErrFailedUpdateOrder
	}

	return s.mapping.ToOrderResponse(existingOrder), nil
}

func (s *orderCommandService) TrashedOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedOrder", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order_id", order_id),
	)

	s.logger.Debug("Moving order to trash",
		zap.Int("order_id", order_id))

	order, err := s.orderQueryRepository.FindById(order_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_BY_ID")

		s.logger.Error("Failed to retrieve order details",
			zap.Error(err),
			zap.Int("order_id", order_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order_id),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve order details")

		status = "failed_find_order_by_id"

		return nil, order_errors.ErrFailedFindOrderById
	}

	if order.DeletedAt != nil {
		traceID := traceunic.GenerateTraceID("FAILED_NOT_DELETE_AT_ORDER")

		s.logger.Error("Order already trashed",
			zap.Int("order_id", order_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order_id),
			attribute.String("traceID", traceID),
		)

		span.SetStatus(codes.Error, "Order already trashed")

		status = "failed_not_delete_at_order"

		return nil, order_errors.ErrFailedNotDeleteAtOrder
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEM_BY_ORDER")

		s.logger.Error("Failed to retrieve order items",
			zap.Error(err),
			zap.Int("order_id", order_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order_id),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve order items")

		status = "failed_find_order_item_by_order"

		return nil, order_errors.ErrFailedNotDeleteAtOrder
	}

	for _, item := range orderItems {
		if item.DeletedAt != nil {
			traceID := traceunic.GenerateTraceID("FAILED_NOT_DELETE_AT_ORDER_ITEM")

			s.logger.Error("Order item already trashed",
				zap.Int("order_item_id", item.ID),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.Int("order_item_id", item.ID),
				attribute.String("traceID", traceID),
			)

			span.SetStatus(codes.Error, "Order item already trashed")

			status = "failed_not_delete_at_order_item"

			return nil, orderitem_errors.ErrFailedNotDeleteAtOrderItem
		}

		trashedItem, err := s.orderItemCommandRepositroy.TrashedOrderItem(item.ID)

		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_TRASHED_ORDER_ITEM")

			s.logger.Error("Failed to move order item to trash",
				zap.Int("order_item_id", item.ID),
				zap.Error(err),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.Int("order_item_id", item.ID),
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to move order item to trash")

			status = "failed_trashed_order_item"

			return nil, orderitem_errors.ErrFailedTrashedOrderItem
		}

		s.logger.Debug("Order item trashed successfully",
			zap.Int("order_item_id", trashedItem.ID),
			zap.String("deleted_at", *trashedItem.DeletedAt))
	}

	trashedOrder, err := s.orderCommandRepository.TrashedOrder(order_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASHED_ORDER")

		s.logger.Error("Failed to move order to trash",
			zap.Int("order_id", order_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order_id),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to move order to trash")

		status = "failed_create_order"

		return nil, order_errors.ErrFailedCreateOrder
	}

	s.logger.Debug("Order moved to trash successfully",
		zap.Int("order_id", order_id),
		zap.String("deleted_at", *trashedOrder.DeletedAt))

	return s.mapping.ToOrderResponseDeleteAt(trashedOrder), nil
}

func (s *orderCommandService) RestoreOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreOrder", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order_id", order_id),
	)

	s.logger.Debug("Restoring order and related order items", zap.Int("order_id", order_id))

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEM_BY_ORDER")

		s.logger.Error("Failed to retrieve order items",
			zap.Error(err),
			zap.Int("order_id", order_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order_id),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve order items")

		status = "failed_find_order_item_by_order"

		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	for _, item := range orderItems {
		_, err := s.orderItemCommandRepositroy.RestoreOrderItem(item.ID)

		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ORDER_ITEM")

			s.logger.Error("Failed to restore order item from trash",
				zap.Int("order_item_id", item.ID),
				zap.Error(err),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.Int("order_item_id", item.ID),
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to restore order item from trash")

			status = "failed_restore_order_item"

			return nil, orderitem_errors.ErrFailedRestoreOrderItem
		}
	}

	order, err := s.orderCommandRepository.RestoreOrder(order_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ORDER")

		s.logger.Error("Failed to restore order from trash",
			zap.Int("order_id", order_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order_id),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore order from trash")

		status = "failed_restore_order"

		return nil, order_errors.ErrFailedRestoreOrder
	}

	return s.mapping.ToOrderResponseDeleteAt(order), nil
}

func (s *orderCommandService) DeleteOrderPermanent(order_id int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteOrderPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteOrderPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order_id", order_id),
	)

	s.logger.Debug("Permanently deleting order and related order items", zap.Int("order_id", order_id))

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(order_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEM_BY_ORDER")

		s.logger.Error("Failed to retrieve order items",
			zap.Error(err),
			zap.Int("order_id", order_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order_id),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve order items")

		status = "failed_find_order_item_by_order"

		return false, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	for _, item := range orderItems {
		_, err := s.orderItemCommandRepositroy.
			DeleteOrderItemPermanent(item.ID)

		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_DELETE_ORDER_ITEM_PERMANENT")

			s.logger.Error("Failed to permanently delete order item",
				zap.Int("order_item_id", item.ID),
				zap.Error(err),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.Int("order_item_id", item.ID),
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to permanently delete order item")

			status = "failed_delete_order_item_permanent"

			return false, orderitem_errors.ErrFailedDeleteOrderItem
		}
	}

	success, err := s.orderCommandRepository.DeleteOrderPermanent(order_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ORDER_PERMANENT")

		s.logger.Error("Failed to permanently delete order",
			zap.Int("order_id", order_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("order_id", order_id),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete order")

		status = "failed_delete_order_permanent"

		return false, order_errors.ErrFailedDeleteOrderPermanent
	}

	return success, nil
}

func (s *orderCommandService) RestoreAllOrder() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllOrder", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllOrder")
	defer span.End()

	s.logger.Debug("Restoring all trashed orders and related order items")

	successItems, err := s.orderItemCommandRepositroy.RestoreAllOrderItem()

	if err != nil || !successItems {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_ORDER_ITEM")

		s.logger.Error("Failed to restore all order items",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all order items")

		status = "failed_restore_all_order_item"

		return false, orderitem_errors.ErrFailedRestoreAllOrderItem
	}

	success, err := s.orderCommandRepository.RestoreAllOrder()
	if err != nil || !success {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_ORDER")

		s.logger.Error("Failed to restore all orders",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all orders")

		status = "failed_restore_all_order"

		return false, order_errors.ErrFailedRestoreAllOrder
	}

	return success, nil
}

func (s *orderCommandService) DeleteAllOrderPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllOrderPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllOrderPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all orders and related order items")

	successItems, err := s.orderItemCommandRepositroy.DeleteAllOrderPermanent()

	if err != nil || !successItems {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_ORDER_ITEM_PERMANENT")

		s.logger.Error("Failed to permanently delete all order items",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all order items")

		status = "failed_delete_all_order_item_permanent"

		return false, orderitem_errors.ErrFailedDeleteAllOrderItem
	}

	success, err := s.orderCommandRepository.DeleteAllOrderPermanent()

	if err != nil || !success {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_ORDER_PERMANENT")

		s.logger.Error("Failed to permanently delete all orders",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all orders")

		status = "failed_delete_all_order_permanent"

		return false, order_errors.ErrFailedDeleteAllOrderPermanent
	}

	return success, nil
}

func (s *orderCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
