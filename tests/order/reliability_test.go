package order_test

import (
	"context"

	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
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
	prodRes, err := pbproduct.NewProductQueryServiceClient(s.Conns["product"]).FindById(ctx, &pbproduct.FindByIdProductRequest{Id: int32(prodID)})
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
