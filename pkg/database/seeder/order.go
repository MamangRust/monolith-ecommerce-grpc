package seeder

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"go.uber.org/zap"
)

type orderSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewOrderSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *orderSeeder {
	return &orderSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *orderSeeder) Seed() error {
	for i := 1; i <= 8; i++ {
		order, err := r.db.CreateOrder(r.ctx, db.CreateOrderParams{
			MerchantID: int32(i),
			UserID:     int32(i),
			TotalPrice: int32(10000 * i),
		})
		if err != nil {
			r.logger.Error("failed to create order", zap.Error(err))
			return err
		}

		_, err = r.db.CreateOrderItem(r.ctx, db.CreateOrderItemParams{
			OrderID:   order.OrderID,
			ProductID: int32(i),
			Quantity:  int32(i),
			Price:     int32(10000),
		})
		if err != nil {
			r.logger.Error("failed to create order item", zap.Error(err))
			return err
		}
	}

	r.logger.Info("order & order-item successfully seeded")

	return nil
}
