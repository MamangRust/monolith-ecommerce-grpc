package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	shipping_address_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type shippingAddressCommandService struct {
	observability             observability.TraceLoggerObservability
	cache                     cache.ShippingAddressCommandCache
	shippingAddressRepository repository.ShippingAddressCommandRepository
	logger                    logger.LoggerInterface
}

type ShippingAddressCommandServiceDeps struct {
	Observability             observability.TraceLoggerObservability
	Cache                     cache.ShippingAddressCommandCache
	ShippingAddressRepository repository.ShippingAddressCommandRepository
	Logger                    logger.LoggerInterface
}

func NewShippingAddressCommandService(deps *ShippingAddressCommandServiceDeps) ShippingAddressCommandService {
	return &shippingAddressCommandService{
		observability:             deps.Observability,
		cache:                     deps.Cache,
		shippingAddressRepository: deps.ShippingAddressRepository,
		logger:                    deps.Logger,
	}
}

func (s *shippingAddressCommandService) Create(ctx context.Context, request *requests.CreateShippingAddressRequest) (*db.CreateShippingAddressRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("alamat", request.Alamat),
		attribute.String("provinsi", request.Provinsi),
		attribute.String("kota", request.Kota))

	defer func() {
		end(status)
	}()

	shippingAddress, err := s.shippingAddressRepository.Create(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateShippingAddressRow](
			s.logger,
			shipping_address_errors.ErrFailedCreateShippingAddress,
			method,
			span,
			zap.Error(err),
		)
	}

	s.cache.InvalidateShippingAddressCache(ctx)

	logSuccess("Successfully created shipping address")

	return shippingAddress, nil
}

func (s *shippingAddressCommandService) Update(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*db.UpdateShippingAddressRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", *request.ShippingID))

	defer func() {
		end(status)
	}()

	shippingAddress, err := s.shippingAddressRepository.Update(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateShippingAddressRow](
			s.logger,
			shipping_address_errors.ErrFailedUpdateShippingAddress,
			method,
			span,
			zap.Error(err),
		)
	}

	s.cache.InvalidateShippingAddressCache(ctx)

	logSuccess("Successfully updated shipping address",
		zap.Int("shipping_id", *request.ShippingID))

	return shippingAddress, nil
}

func (s *shippingAddressCommandService) Trash(ctx context.Context, shipping_id int) (*db.ShippingAddress, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	shippingAddress, err := s.shippingAddressRepository.Trash(ctx, shipping_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.ShippingAddress](
			s.logger,
			shipping_address_errors.ErrFailedTrashShippingAddress,
			method,
			span,

			zap.Int("shipping_id", shipping_id),
		)
	}

	s.cache.InvalidateShippingAddressCache(ctx)

	logSuccess("Successfully trashed shipping address",
		zap.Int("shipping_id", shipping_id))

	return shippingAddress, nil
}

func (s *shippingAddressCommandService) Restore(ctx context.Context, shipping_id int) (*db.ShippingAddress, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	shippingAddress, err := s.shippingAddressRepository.Restore(ctx, shipping_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.ShippingAddress](
			s.logger,
			shipping_address_errors.ErrFailedRestoreShippingAddress,
			method,
			span,

			zap.Int("shipping_id", shipping_id),
		)
	}

	s.cache.InvalidateShippingAddressCache(ctx)

	logSuccess("Successfully restored shipping address",
		zap.Int("shipping_id", shipping_id))

	return shippingAddress, nil
}

func (s *shippingAddressCommandService) DeletePermanent(ctx context.Context, shipping_id int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressRepository.DeletePermanent(ctx, shipping_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			shipping_address_errors.ErrFailedDeleteShippingAddressPermanent,
			method,
			span,

			zap.Int("shipping_id", shipping_id),
		)
	}

	s.cache.InvalidateShippingAddressCache(ctx)

	logSuccess("Successfully permanently deleted shipping address",
		zap.Int("shipping_id", shipping_id))

	return success, nil
}

func (s *shippingAddressCommandService) DeleteShippingAddressByOrderPermanent(ctx context.Context, order_id int) (bool, error) {
	const method = "DeleteShippingAddressByOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", order_id))

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressRepository.DeleteByOrderIDPermanent(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			shipping_address_errors.ErrFailedDeleteShippingAddressPermanent,
			method,
			span,
			zap.Error(err),
			zap.Int("order_id", order_id),
		)
	}

	s.cache.InvalidateShippingAddressCache(ctx)

	logSuccess("Successfully permanently deleted shipping addresses by order",
		zap.Int("order_id", order_id))

	return success, nil
}

func (s *shippingAddressCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			shipping_address_errors.ErrFailedRestoreAllShippingAddresses,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed shipping addresses")

	return success, nil
}

func (s *shippingAddressCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			shipping_address_errors.ErrFailedDeleteAllShippingAddressesPermanent,
			method,
			span,
		)
	}

	s.cache.InvalidateShippingAddressCache(ctx)
	logSuccess("Successfully permanently deleted all trashed shipping addresses")

	return success, nil
}
