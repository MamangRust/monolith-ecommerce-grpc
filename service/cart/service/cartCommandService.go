package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cartCommandService struct {
	cartCommandRepository  repository.CartCommandRepository
	productQueryRepository repository.ProductQueryRepository
	userQueryRepository    repository.UserQueryRepository
	observability          observability.TraceLoggerObservability
	logger                 logger.LoggerInterface
}

type CartCommandServiceDeps struct {
	CartCommandRepository  repository.CartCommandRepository
	ProductQueryRepository repository.ProductQueryRepository
	UserQueryRepository    repository.UserQueryRepository
	Observability          observability.TraceLoggerObservability
	Logger                 logger.LoggerInterface
}

func NewCartCommandService(deps *CartCommandServiceDeps) *cartCommandService {
	return &cartCommandService{
		cartCommandRepository:  deps.CartCommandRepository,
		productQueryRepository: deps.ProductQueryRepository,
		userQueryRepository:    deps.UserQueryRepository,
		logger:                 deps.Logger,
		observability:          deps.Observability,
	}
}

func (s *cartCommandService) Create(ctx context.Context, req *requests.CreateCartRequest) (*db.Cart, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	product, err := s.productQueryRepository.FindById(ctx, req.ProductID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Cart](
			s.logger,
			err,
			method,
			span,
			zap.Int("product_id", req.ProductID),
		)
	}

	_, err = s.userQueryRepository.FindById(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Cart](
			s.logger,
			err,
			method,
			span,
			zap.Int("user_id", req.UserID),
		)
	}

	var imageProduct string
	if product.ImageProduct != nil {
		imageProduct = *product.ImageProduct
	}

	var weight int
	if product.Weight != nil {
		weight = int(*product.Weight)
	}

	cartRecord := &requests.CartCreateRecord{
		ProductID:    req.ProductID,
		UserID:       req.UserID,
		Name:         product.Name,
		Price:        int(product.Price),
		ImageProduct: imageProduct,
		Quantity:     req.Quantity,
		Weight:       weight,
	}

	res, err := s.cartCommandRepository.CreateCart(ctx, cartRecord)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Cart](
			s.logger,
			err,
			method,
			span,
			zap.Any("request", req),
		)
	}

	logSuccess("Successfully created cart", zap.Int("cartID", int(res.CartID)))
	return &db.Cart{
		CartID:    res.CartID,
		UserID:    res.UserID,
		ProductID: res.ProductID,
		Name:      res.Name,
		Price:     res.Price,
		Image:     res.Image,
		Quantity:  res.Quantity,
		Weight:    res.Weight,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func (s *cartCommandService) DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error) {
	const method = "DeletePermanentCart"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("cartID", req.CartID))

	defer func() {
		end(status)
	}()

	success, err := s.cartCommandRepository.DeletePermanent(ctx, &requests.DeleteCartRequest{
		CartID: req.CartID,
		UserID: req.UserID,
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Int("cart_id", req.CartID),
		)
	}

	logSuccess("Successfully deleted cart permanently", zap.Int("cartID", req.CartID))
	return success, nil
}

func (s *cartCommandService) DeleteAll(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.cartCommandRepository.DeleteAllPermanently(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all carts permanently")
	return success, nil
}
