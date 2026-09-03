package cart_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CartRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *CartRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(
		queries,
		pb.NewUserQueryServiceClient(s.Conns["user"]),
		pb.NewProductQueryServiceClient(s.Conns["product"]),
	)
}

func (s *CartRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

// Cart has no soft-delete (Trash/Restore) — only Create, Find, Delete, DeleteAll.
func (s *CartRepositoryTestSuite) TestCartLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	childUserID := s.SeedUser(ctx)
	childCategoryID := s.SeedCategory(ctx)
	childMerchantID := s.SeedMerchant(ctx, childUserID)
	prodID := s.SeedProduct(ctx, childMerchantID, childCategoryID)
	prodID2 := s.SeedProduct(ctx, childMerchantID, childCategoryID)

	// Fetch product details for the cart record
	prodRes, err := pb.NewProductQueryServiceClient(s.Conns["product"]).FindById(ctx, &pb.FindByIdProductRequest{Id: int32(prodID)})
	s.Require().NoError(err)
	s.Require().NotNil(prodRes)
	s.Require().NotNil(prodRes.Data)

	// 1. Create first cart item
	req := &requests.CartCreateRecord{
		UserID:       childUserID,
		ProductID:    prodID,
		Name:         prodRes.Data.Name,
		Price:        int(prodRes.Data.Price),
		ImageProduct: prodRes.Data.ImageProduct,
		Quantity:     2,
		Weight:       int(prodRes.Data.Weight),
	}
	created, err := s.repo.CartCommand.CreateCart(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(childUserID), created.UserID)
	cart1ID := int(created.CartID)

	// 2. Create second cart item
	prodRes2, err := pb.NewProductQueryServiceClient(s.Conns["product"]).FindById(ctx, &pb.FindByIdProductRequest{Id: int32(prodID2)})
	s.Require().NoError(err)

	req2 := &requests.CartCreateRecord{
		UserID:       childUserID,
		ProductID:    prodID2,
		Name:         prodRes2.Data.Name,
		Price:        int(prodRes2.Data.Price),
		ImageProduct: prodRes2.Data.ImageProduct,
		Quantity:     1,
		Weight:       int(prodRes2.Data.Weight),
	}
	created2, err := s.repo.CartCommand.CreateCart(ctx, req2)
	s.NoError(err)
	s.NotNil(created2)
	cart2ID := int(created2.CartID)

	// 3. FindCarts — list all items for the user
	items, err := s.repo.CartQuery.FindCarts(ctx, &requests.FindAllCarts{
		UserID:   childUserID,
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.NotEmpty(items)

	// 4. DeletePermanent — remove first cart item
	delSuccess, err := s.repo.CartCommand.DeletePermanent(ctx, &requests.DeleteCartRequest{
		UserID: childUserID,
		CartID: cart1ID,
	})
	s.NoError(err)
	s.True(delSuccess)

	// 5. Verify only one item remains
	remainingItems, err := s.repo.CartQuery.FindCarts(ctx, &requests.FindAllCarts{
		UserID:   childUserID,
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.Len(remainingItems, 1)
	s.Equal(cart2ID, int(remainingItems[0].CartID))

	// 6. DeleteAllPermanently — remove all remaining items
	delAllSuccess, err := s.repo.CartCommand.DeleteAllPermanently(ctx, &requests.DeleteAllCartRequest{
		UserID:  childUserID,
		CartIds: []int{cart2ID},
	})
	s.NoError(err)
	s.True(delAllSuccess)

	// 7. Verify cart is empty
	emptyItems, err := s.repo.CartQuery.FindCarts(ctx, &requests.FindAllCarts{
		UserID:   childUserID,
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.Empty(emptyItems)
}

func TestCartRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartRepositoryTestSuite))
}
