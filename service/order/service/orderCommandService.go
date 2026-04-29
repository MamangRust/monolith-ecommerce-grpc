package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
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
	transactionCommandRepos   repository.TransactionCommandRepository
	shippingQueryRepository   pb.ShippingQueryServiceClient
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
	TransactionCommandRepository repository.TransactionCommandRepository
	ShippingQueryRepository   pb.ShippingQueryServiceClient
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
		transactionCommandRepos:   deps.TransactionCommandRepository,
		shippingQueryRepository:   deps.ShippingQueryRepository,
		logger:                    deps.Logger,
	}
}

func (s *orderCommandService) Create(ctx context.Context, req *requests.CreateOrderRequest) (*db.CreateOrderRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("user.id", req.UserID))

	defer func() {
		end(status)
	}()

	order, err := s.orderCommandRepository.Create(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		UserID:     req.UserID,
	})

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateOrderRow](s.logger, err, method, span)
	}

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindByID(ctx, item.ProductID)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateOrderRow](s.logger, err, method, span)
		}

		if product.CountInStock < int32(item.Quantity) {
			status = "error"
			return errorhandler.HandleError[*db.CreateOrderRow](s.logger, order_errors.ErrInsufficientProductStock, method, span)
		}

		_, err = s.orderItemCommandRepos.Create(ctx, &requests.CreateOrderItemRecordRequest{
			OrderID:   int(order.OrderID),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
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

	_, err = s.shippingAddressRepository.Create(ctx, &requests.CreateShippingAddressRequest{
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

	res, err := s.orderCommandRepository.Update(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    int(order.OrderID),
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice) + req.ShippingAddress.ShippingCost,
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

func (s *orderCommandService) Update(ctx context.Context, req *requests.UpdateOrderRequest) (*db.UpdateOrderRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order.id", *req.OrderID),
		attribute.Int("user.id", req.UserID))

	defer func() {
		end(status)
	}()

	existingOrder, err := s.orderQueryRepository.FindByID(ctx, *req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
	}

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindByID(ctx, item.ProductID)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemCommandRepos.Update(ctx, &requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       item.Price,
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

			_, err := s.orderItemCommandRepos.Create(ctx, &requests.CreateOrderItemRecordRequest{
				OrderID:   *req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
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

	shippingID := req.ShippingAddress.ShippingID
	if shippingID == nil {
		shippingRes, err := s.shippingQueryRepository.FindByOrder(ctx, &pb.FindByIdShippingRequest{
			Id: int32(*req.OrderID),
		})
		if err == nil && shippingRes != nil && shippingRes.Data != nil {
			id := int(shippingRes.Data.Id)
			shippingID = &id
		}
	}

	_, err = s.shippingAddressRepository.Update(ctx, &requests.UpdateShippingAddressRequest{
		ShippingID:     shippingID,
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

	res, err := s.orderCommandRepository.Update(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    *req.OrderID,
		UserID:     req.UserID,
		TotalPrice: int(*totalPrice) + req.ShippingAddress.ShippingCost,
	})

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderRow](s.logger, err, method, span)
	}

	s.cache.DeleteOrderCache(ctx, *req.OrderID)

	logSuccess("Successfully updated order", zap.Int("order.id", *req.OrderID))

	return res, nil
}

func (s *orderCommandService) Trash(ctx context.Context, orderID int) (*db.Order, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID))

	defer func() {
		end(status)
	}()

	order, err := s.orderCommandRepository.Trash(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Order](s.logger, err, method, span)
	}

	s.cache.DeleteOrderCache(ctx, orderID)

	logSuccess("Successfully trashed order", zap.Int("orderID", orderID))

	return order, nil
}

func (s *orderCommandService) Restore(ctx context.Context, orderID int) (*db.Order, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID))

	defer func() {
		end(status)
	}()

	order, err := s.orderCommandRepository.Restore(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Order](s.logger, err, method, span)
	}

	s.cache.DeleteOrderCache(ctx, orderID)

	logSuccess("Successfully restored order", zap.Int("orderID", orderID))

	return order, nil
}

func (s *orderCommandService) DeletePermanent(ctx context.Context, orderID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID))

	defer func() {
		end(status)
	}()

	_, err := s.orderItemCommandRepos.DeleteByOrderIDPermanent(ctx, orderID)
	if err != nil {
		return false, err
	}

	_, err = s.transactionCommandRepos.DeleteByOrderIDPermanent(ctx, orderID)
	if err != nil {
		return false, err
	}

	_, err = s.shippingAddressRepository.DeleteByOrderIDPermanent(ctx, orderID)
	if err != nil {
		return false, err
	}

	success, err := s.orderCommandRepository.DeletePermanent(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.DeleteOrderCache(ctx, orderID)

	logSuccess("Successfully deleted order permanently", zap.Int("orderID", orderID))

	return success, nil
}

func (s *orderCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.orderCommandRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	logSuccess("Successfully restored all orders")

	return success, nil
}

func (s *orderCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.orderItemCommandRepos.DeleteAll(ctx)
	if err != nil {
		return false, err
	}

	_, err = s.transactionCommandRepos.DeleteAll(ctx)
	if err != nil {
		return false, err
	}

	_, err = s.shippingAddressRepository.DeleteAll(ctx)
	if err != nil {
		return false, err
	}

	// For DeleteAllOrderPermanent, we might want to also delete all shipping addresses and transactions.
	// However, these methods don't exist yet for "all permanent".
	// For simplicity in this task, we focus on the order-specific permanent deletion.

	success, err := s.orderCommandRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	logSuccess("Successfully deleted all orders permanently")

	return success, nil
}
