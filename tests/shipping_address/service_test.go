package shipping_address_test

import (
	"context"
	"testing"

	ship_cache "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ShippingAddressServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *ShippingAddressServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Shipping dependencies
	mencache := ship_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(queries)

	s.svc = service.NewService(&service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *ShippingAddressServiceTestSuite) TestShippingAddressLifecycle() {
	ctx := context.Background()

	// 1. Setup Dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderService()

	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	categoryID := s.SeedCategory(ctx)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)

	// 2. Create Shipping Address
	req := &requests.CreateShippingAddressRequest{
		OrderID:        &orderID,
		Alamat:         "Initial Address",
		Provinsi:       "West Java",
		Kota:           "Bandung",
		Courier:        "TIKI",
		ShippingMethod: "ONS",
		ShippingCost:   20000,
		Negara:         "Indonesia",
	}
	created, err := s.svc.ShippingAddressCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	shippingID := int(created.ShippingAddressID)

	// 3. FindByID
	found, err := s.svc.ShippingAddressQuery.FindByID(ctx, shippingID)
	s.Require().NoError(err)
	s.Equal(req.Alamat, found.Alamat)

	// 4. FindByOrder
	_, err = s.svc.ShippingAddressQuery.FindByOrder(ctx, orderID)
	s.Require().NoError(err)

	// 5. Update
	newAlamat := "Updated Address"
	updateReq := &requests.UpdateShippingAddressRequest{
		ShippingID:     &shippingID,
		OrderID:        &orderID,
		Alamat:         newAlamat,
		Provinsi:       req.Provinsi,
		Kota:           req.Kota,
		Courier:        req.Courier,
		ShippingMethod: req.ShippingMethod,
		ShippingCost:   req.ShippingCost,
		Negara:         req.Negara,
	}
	updated, err := s.svc.ShippingAddressCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newAlamat, updated.Alamat)

	// 6. FindAll
	_, total, err := s.svc.ShippingAddressQuery.FindAll(ctx, &requests.FindAllShippingAddress{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 7. Trash
	_, err = s.svc.ShippingAddressCommand.Trash(ctx, shippingID)
	s.Require().NoError(err)

	// 8. FindTrashed
	_, totalTrashed, err := s.svc.ShippingAddressQuery.FindTrashed(ctx, &requests.FindAllShippingAddress{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 9. FindActive
	active, _, err := s.svc.ShippingAddressQuery.FindActive(ctx, &requests.FindAllShippingAddress{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, a := range active {
		s.NotEqual(shippingID, int(a.ShippingAddressID))
	}

	// 10. Restore
	_, err = s.svc.ShippingAddressCommand.Restore(ctx, shippingID)
	s.Require().NoError(err)

	// 11. DeletePermanent
	_, err = s.svc.ShippingAddressCommand.Trash(ctx, shippingID)
	s.Require().NoError(err)
	success, err := s.svc.ShippingAddressCommand.DeletePermanent(ctx, shippingID)
	s.Require().NoError(err)
	s.True(success)

	// 12. RestoreAll & DeleteAll
	s1, _ := s.svc.ShippingAddressCommand.Create(ctx, &requests.CreateShippingAddressRequest{OrderID: &orderID, Alamat: "A1", Provinsi: "P1", Kota: "K1", Courier: "C1", ShippingMethod: "M1", ShippingCost: 1000, Negara: "N1"})
	s2, _ := s.svc.ShippingAddressCommand.Create(ctx, &requests.CreateShippingAddressRequest{OrderID: &orderID, Alamat: "A2", Provinsi: "P2", Kota: "K2", Courier: "C2", ShippingMethod: "M2", ShippingCost: 2000, Negara: "N2"})

	s.svc.ShippingAddressCommand.Trash(ctx, int(s1.ShippingAddressID))
	s.svc.ShippingAddressCommand.Trash(ctx, int(s2.ShippingAddressID))

	resRestoreAll, err := s.svc.ShippingAddressCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.ShippingAddressCommand.Trash(ctx, int(s1.ShippingAddressID))
	s.svc.ShippingAddressCommand.Trash(ctx, int(s2.ShippingAddressID))

	resDeleteAll, err := s.svc.ShippingAddressCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestShippingAddressServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ShippingAddressServiceTestSuite))
}
