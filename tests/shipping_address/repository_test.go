package shipping_address_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ShippingAddressRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *ShippingAddressRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(queries)
}

func (s *ShippingAddressRepositoryTestSuite) TestAddressLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)

	req := &requests.CreateShippingAddressRequest{
		OrderID:        &orderID,
		Alamat:         "Home",
		Provinsi:       "Provinsi",
		Kota:           "Kota",
		Courier:        "Courier",
		ShippingMethod: "Method",
		ShippingCost:   1000,
		Negara:         "Indonesia",
	}

	created, err := s.repo.ShippingAddressCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.Alamat, created.Alamat)

	// 2. Find by ID
	found, err := s.repo.ShippingAddressQuery.FindByID(ctx, int(created.ShippingAddressID))
	s.NoError(err)
	s.Equal(created.Alamat, found.Alamat)

	// 3. Find by Order ID
	all, err := s.repo.ShippingAddressQuery.FindByOrder(ctx, orderID)
	s.NoError(err)
	s.NotEmpty(all)
}

func TestShippingAddressRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ShippingAddressRepositoryTestSuite))
}
