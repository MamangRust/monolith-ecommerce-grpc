package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	shipping_address_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
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

func (s *shippingAddressCommandService) TrashShippingAddress(ctx context.Context, shipping_id int) (*db.ShippingAddress, error) {
	const method = "TrashShippingAddress"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	shippingAddress, err := s.shippingAddressRepository.TrashShippingAddress(ctx, shipping_id)
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

func (s *shippingAddressCommandService) RestoreShippingAddress(ctx context.Context, shipping_id int) (*db.ShippingAddress, error) {
	const method = "RestoreShippingAddress"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	shippingAddress, err := s.shippingAddressRepository.RestoreShippingAddress(ctx, shipping_id)
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

func (s *shippingAddressCommandService) DeleteShippingAddressPermanently(ctx context.Context, shipping_id int) (bool, error) {
	const method = "DeleteShippingAddressPermanently"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressRepository.DeleteShippingAddressPermanently(ctx, shipping_id)
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

func (s *shippingAddressCommandService) RestoreAllShippingAddress(ctx context.Context) (bool, error) {
	const method = "RestoreAllShippingAddress"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressRepository.RestoreAllShippingAddress(ctx)
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

func (s *shippingAddressCommandService) DeleteAllPermanentShippingAddress(ctx context.Context) (bool, error) {
	const method = "DeleteAllPermanentShippingAddress"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.shippingAddressRepository.DeleteAllPermanentShippingAddress(ctx)
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
