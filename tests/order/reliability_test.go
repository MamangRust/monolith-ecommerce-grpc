package order_test

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

func (s *OrderServiceTestSuite) orderRequest(userID, merchID, prodID int) *requests.CreateOrderRequest {
	return &requests.CreateOrderRequest{
		UserID:     userID,
		MerchantID: merchID,
		TotalPrice: 20000,
		Items: []requests.CreateOrderItemRequest{
			{ProductID: prodID, Quantity: 1, Price: 10000},
		},
		ShippingAddress: requests.CreateShippingAddressRequest{
			Alamat:         "Test Address",
			Provinsi:       "West Java",
			Kota:           "Bandung",
			Courier:        "JNE",
			ShippingMethod: "REG",
			ShippingCost:   10000,
			Negara:         "Indonesia",
		},
	}
}

func (s *OrderServiceTestSuite) stockOf(ctx context.Context, prodID int) int32 {
	prodRes, err := pb.NewProductQueryServiceClient(s.Conns["product"]).FindById(ctx, &pb.FindByIdProductRequest{Id: int32(prodID)})
	s.Require().NoError(err)
	s.Require().NotNil(prodRes.Data)
	return prodRes.Data.CountInStock
}

// TestFailureInjectionInsufficientStock verifies that creating an order with a
// quantity beyond the available stock fails, rolls back any partial reservation,
// and leaves inventory untouched.
func (s *OrderServiceTestSuite) TestFailureInjectionInsufficientStock() {
	ctx := context.Background()

	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	s.Require().Equal(int32(100), s.stockOf(ctx, prodID))

	// Demand far more than the 100 units in stock.
	req := s.orderRequest(userID, merchID, prodID)
	req.Items[0].Quantity = 1000

	_, err := s.svc.OrderCommand.Create(ctx, req)
	s.Require().Error(err, "order creation beyond stock must fail")

	s.Equal(int32(100), s.stockOf(ctx, prodID), "stock must be untouched after a failed create")

	// No orphan reservation rows (or orders) may remain for the failed create.
	var reservations int64
	err = s.DBPool().QueryRow(ctx, `
		SELECT COUNT(*)
		FROM order_stock_reservations r
		JOIN orders o ON o.order_id = r.order_id
		WHERE o.user_id = $1`, int32(userID)).Scan(&reservations)
	s.Require().NoError(err)
	s.Zero(reservations, "no reservations may survive a rolled-back create")

	var orders int64
	err = s.DBPool().QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE user_id = $1`, int32(userID)).Scan(&orders)
	s.Require().NoError(err)
	s.Zero(orders, "no order rows may survive a rolled-back create")
}

// TestReconcileDriftAndCleanupIdempotency verifies durable reconciliation repairs
// reservation/order drift in both directions (with correct stock deltas) and that
// the retention policy purges only rows older than the cutoff.
func (s *OrderServiceTestSuite) TestReconcileDriftAndCleanupIdempotency() {
	ctx := context.Background()

	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	pool := s.DBPool()

	created, err := s.svc.OrderCommand.Create(ctx, s.orderRequest(userID, merchID, prodID))
	s.Require().NoError(err)
	orderID := int(created.OrderID)
	s.Equal(int32(99), s.stockOf(ctx, prodID), "create must reserve one unit of stock")

	reservationStatus := func() string {
		var status string
		err := pool.QueryRow(ctx,
			`SELECT status FROM order_stock_reservations WHERE order_id = $1 AND product_id = $2`,
			orderID, prodID).Scan(&status)
		s.Require().NoError(err)
		return status
	}
	s.Equal("reserved", reservationStatus(), "create must reserve stock")

	// Drift A: reservation wrongly marked released while the order is still active.
	// Durable reconciliation must re-reserve it (decrement stock again).
	_, err = pool.Exec(ctx,
		`UPDATE order_stock_reservations SET status = 'released' WHERE order_id = $1 AND product_id = $2`,
		orderID, prodID)
	s.Require().NoError(err)

	res, err := s.svc.OrderCommand.ReconcileStockReservations(ctx)
	s.Require().NoError(err)
	s.GreaterOrEqual(res.ReReserved, 1)
	s.Equal("reserved", reservationStatus(), "reconciliation must re-reserve stock for active orders")
	s.Equal(int32(98), s.stockOf(ctx, prodID), "re-reconciliation must decrement stock again")

	// Drift B: trash releases stock; simulate a reservation stuck as reserved.
	_, err = s.svc.OrderCommand.Trash(ctx, orderID)
	s.Require().NoError(err)
	s.Equal("released", reservationStatus(), "trash must release stock")
	s.Equal(int32(99), s.stockOf(ctx, prodID), "trash must return stock to inventory")

	_, err = pool.Exec(ctx,
		`UPDATE order_stock_reservations SET status = 'reserved' WHERE order_id = $1 AND product_id = $2`,
		orderID, prodID)
	s.Require().NoError(err)

	res2, err := s.svc.OrderCommand.ReconcileStockReservations(ctx)
	s.Require().NoError(err)
	s.GreaterOrEqual(res2.Released, 1)
	s.Equal("released", reservationStatus(), "reconciliation must release stock for trashed orders")
	s.Equal(int32(100), s.stockOf(ctx, prodID), "reconciliation must restore released stock")

	// Cleanup: age the released reservation and insert an old adjustment row,
	// then confirm the retention policy purges only rows older than the cutoff.
	_, err = pool.Exec(ctx,
		`UPDATE order_stock_reservations SET updated_at = NOW() - INTERVAL '10 days' WHERE order_id = $1 AND product_id = $2`,
		orderID, prodID)
	s.Require().NoError(err)

	_, err = pool.Exec(ctx,
		`INSERT INTO product_stock_adjustments (operation_id, product_id, delta, created_at)
		 VALUES ($1, $2, $3, $4)`,
		"test-old-adjustment", int32(prodID), 1, time.Now().AddDate(0, 0, -10))
	s.Require().NoError(err)

	cleanup, err := s.svc.OrderCommand.CleanupIdempotencyRecords(ctx, 1)
	s.Require().NoError(err)
	s.GreaterOrEqual(cleanup.ReleasedReservationsRemoved, int64(1))
	s.GreaterOrEqual(cleanup.AdjustmentsRemoved, int64(1))

	// The aged trashed reservation is gone.
	var stale int64
	err = pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM order_stock_reservations WHERE order_id = $1 AND product_id = $2`,
		orderID, prodID).Scan(&stale)
	s.Require().NoError(err)
	s.Zero(stale, "aged released reservation must be purged")

	// Fresh reservations created after cleanup must survive the retention policy.
	freshOrder, err := s.svc.OrderCommand.Create(ctx, s.orderRequest(userID, merchID, prodID))
	s.Require().NoError(err)
	freshOrderID := int(freshOrder.OrderID)

	var freshReservations int64
	err = pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM order_stock_reservations WHERE order_id = $1 AND product_id = $2 AND status = 'reserved'`,
		freshOrderID, prodID).Scan(&freshReservations)
	s.Require().NoError(err)
	s.Equal(int64(1), freshReservations, "fresh reservation must survive the retention cleanup")
}
