package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderCommandService struct {
	observability             observability.TraceLoggerObservability
	cache                     cache.OrderCommandCache
	userQueryRepository       repository.UserQueryRepository
	productQueryRepository    repository.ProductQueryRepository
	productCommandRepository  repository.ProductCommandRepository
	orderQueryRepository      repository.OrderQueryRepository
	orderCommandRepository    repository.OrderCommandRepository
	orderItemQueryRepository  repository.OrderItemQueryRepository
	orderItemCommandRepos     repository.OrderItemCommandRepository
	merchantQueryRepository   repository.MerchantQueryRepository
	shippingAddressRepository repository.ShippingAddressCommandRepository
	logger                    logger.LoggerInterface
}

type OrderCommandServiceDeps struct {
	Observability             observability.TraceLoggerObservability
	Cache                     cache.OrderCommandCache
	UserQueryRepository       repository.UserQueryRepository
	ProductQueryRepository    repository.ProductQueryRepository
	ProductCommandRepository  repository.ProductCommandRepository
	OrderQueryRepository      repository.OrderQueryRepository
	OrderCommandRepository    repository.OrderCommandRepository
	OrderItemQueryRepository  repository.OrderItemQueryRepository
	OrderItemCommandRepository repository.OrderItemCommandRepository
	MerchantQueryRepository   repository.MerchantQueryRepository
	ShippingAddressRepository repository.ShippingAddressCommandRepository
	Logger                    logger.LoggerInterface
}

func NewOrderCommandService(deps *OrderCommandServiceDeps) OrderCommandService {
	return &orderCommandService{
		observability:             deps.Observability,
		cache:                     deps.Cache,
		userQueryRepository:       deps.UserQueryRepository,
		productQueryRepository:    deps.ProductQueryRepository,
		productCommandRepository:  deps.ProductCommandRepository,
		orderQueryRepository:      deps.OrderQueryRepository,
		orderCommandRepository:    deps.OrderCommandRepository,
		orderItemQueryRepository:  deps.OrderItemQueryRepository,
		orderItemCommandRepos:     deps.OrderItemCommandRepository,
		merchantQueryRepository:   deps.MerchantQueryRepository,
		shippingAddressRepository: deps.ShippingAddressRepository,
		logger:                    deps.Logger,
	}
}

func (s *orderCommandService) CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*db.CreateOrderRow, error) {
	const method = "CreateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("user.id", req.UserID))

	defer func() {
		end(status)
	}()

	order, err := s.orderCommandRepository.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		UserID:     req.UserID,
	})

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateOrderRow](s.logger, err, method, span)
	}

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindById(ctx, item.ProductID)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateOrderRow](s.logger, err, method, span)
		}

		if product.CountInStock < int32(item.Quantity) {
			status = "error"
			return errorhandler.HandleError[*db.CreateOrderRow](s.logger, order_errors.ErrInsufficientProductStock, method, span)
		}

		_, err = s.orderItemCommandRepos.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
			OrderID:   int(order.OrderID),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     int(product.Price),
		})

		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateOrderRow](s.logger, err, method, span)
		}

		product.CountInStock -= int32(item.Quantity)
		_, err = s.productCommandRepository.UpdateProductCountStock(ctx, int(product.ProductID), int(product.CountInStock))
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateOrderRow](s.logger, err, method, span)
		}
	}

	_, err = s.shippingAddressRepository.CreateShippingAddress(ctx, &requests.CreateShippingAddressRequest{
		OrderID:        pointerInt32ToInt(order.OrderID),
		Alamat:         req.ShippingAddress.Alamat,
		Provinsi:       req.ShippingAddress.Provinsi,
		Kota:           req.ShippingAddress.Kota,
		Courier:        req.ShippingAddress.Courier,
		ShippingMethod: req.ShippingAddress.ShippingMethod,
		ShippingCost:   req.ShippingAddress.ShippingCost,
		Negara:         req.ShippingAddress.Negara,
	})

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateOrderRow](s.logger, err, method, span)
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, int(order.OrderID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateOrderRow](s.logger, err, method, span)
	}

	res, err := s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    int(order.OrderID),
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})

	logSuccess("Successfully created order", zap.Int("order.id", int(order.OrderID)))

	return &db.CreateOrderRow{
		OrderID:    res.OrderID,
		UserID:     res.UserID,
		MerchantID: res.MerchantID,
		TotalPrice: res.TotalPrice,
		CreatedAt:  res.CreatedAt,
		UpdatedAt:  res.UpdatedAt,
	}, nil
}

func pointerInt32ToInt(v int32) *int {
	res := int(v)
	return &res
}

func (s *orderCommandService) UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*db.UpdateOrderRow, error) {
	const method = "UpdateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order.id", *req.OrderID),
		attribute.Int("user.id", req.UserID))

	defer func() {
		end(status)
	}()

	existingOrder, err := s.orderQueryRepository.FindById(ctx, *req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
	}

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindById(ctx, item.ProductID)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemCommandRepos.UpdateOrderItem(ctx, &requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       int(product.Price),
			})
			if err != nil {
				status = "error"
				return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
			}
		} else {
			if product.CountInStock < int32(item.Quantity) {
				status = "error"
				return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, order_errors.ErrInsufficientProductStock, method, span)
			}

			_, err := s.orderItemCommandRepos.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
				OrderID:   *req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     int(product.Price),
			})
			if err != nil {
				status = "error"
				return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
			}

			product.CountInStock -= int32(item.Quantity)
			_, err = s.productCommandRepository.UpdateProductCountStock(ctx, int(product.ProductID), int(product.CountInStock))
			if err != nil {
				status = "error"
				return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
			}
		}
	}

	_, err = s.shippingAddressRepository.UpdateShippingAddress(ctx, &requests.UpdateShippingAddressRequest{
		ShippingID:     req.ShippingAddress.ShippingID,
		OrderID:        pointerInt32ToInt(existingOrder.OrderID),
		Alamat:         req.ShippingAddress.Alamat,
		Provinsi:       req.ShippingAddress.Provinsi,
		Kota:           req.ShippingAddress.Kota,
		Courier:        req.ShippingAddress.Courier,
		ShippingMethod: req.ShippingAddress.ShippingMethod,
		ShippingCost:   req.ShippingAddress.ShippingCost,
		Negara:         req.ShippingAddress.Negara,
	})

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, *req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
	}

	res, err := s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    *req.OrderID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice),
	})

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
	}

	s.cache.DeleteOrderCache(ctx, *req.OrderID)

	logSuccess("Successfully updated order", zap.Int("order.id", *req.OrderID))

	return res, nil
}

func (s *orderCommandService) TrashedOrder(ctx context.Context, orderID int) (*db.Order, error) {
	const method = "TrashedOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID))

	defer func() {
		end(status)
	}()

	order, err := s.orderCommandRepository.TrashedOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Order](s.logger, err, method, span)
	}

	s.cache.DeleteOrderCache(ctx, orderID)

	logSuccess("Successfully trashed order", zap.Int("orderID", orderID))

	return order, nil
}

func (s *orderCommandService) RestoreOrder(ctx context.Context, orderID int) (*db.Order, error) {
	const method = "RestoreOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID))

	defer func() {
		end(status)
	}()

	order, err := s.orderCommandRepository.RestoreOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Order](s.logger, err, method, span)
	}

	s.cache.DeleteOrderCache(ctx, orderID)

	logSuccess("Successfully restored order", zap.Int("orderID", orderID))

	return order, nil
}

func (s *orderCommandService) DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	const method = "DeleteOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID))

	defer func() {
		end(status)
	}()

	success, err := s.orderCommandRepository.DeleteOrderPermanent(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.DeleteOrderCache(ctx, orderID)

	logSuccess("Successfully deleted order permanently", zap.Int("orderID", orderID))

	return success, nil
}

func (s *orderCommandService) RestoreAllOrder(ctx context.Context) (bool, error) {
	const method = "RestoreAllOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.orderCommandRepository.RestoreAllOrder(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	logSuccess("Successfully restored all orders")

	return success, nil
}

func (s *orderCommandService) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.orderCommandRepository.DeleteAllOrderPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	logSuccess("Successfully deleted all orders permanently")

	return success, nil
}
