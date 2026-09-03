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

func (s *ShippingAddressRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *ShippingAddressRepositoryTestSuite) TestAddressLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)

	pageReq := &requests.FindAllShippingAddress{Search: "", Page: 1, PageSize: 10}

	// 1. Create
	req := &requests.CreateShippingAddressRequest{
		OrderID:        &orderID,
		Alamat:         "Home",
		Provinsi:       "DKI Jakarta",
		Kota:           "Jakarta",
		Courier:        "JNE",
		ShippingMethod: "REG",
		ShippingCost:   10000,
		Negara:         "Indonesia",
	}
	created, err := s.repo.ShippingAddressCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.Alamat, created.Alamat)
	addrID := int(created.ShippingAddressID)

	// 2. FindByID
	found, err := s.repo.ShippingAddressQuery.FindByID(ctx, addrID)
	s.NoError(err)
	s.Equal(created.Alamat, found.Alamat)

	// 3. FindByOrder
	foundByOrder, err := s.repo.ShippingAddressQuery.FindByOrder(ctx, orderID)
	s.NoError(err)
	s.NotEmpty(foundByOrder)

	// 4. FindAll
	all, err := s.repo.ShippingAddressQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 5. FindActive (before trash)
	active, err := s.repo.ShippingAddressQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 6. Update
	updateReq := &requests.UpdateShippingAddressRequest{
		ShippingID:     &addrID,
		OrderID:        &orderID,
		Alamat:         "Office",
		Provinsi:       "Jawa Barat",
		Kota:           "Bandung",
		Courier:        "SiCepat",
		ShippingMethod: "YES",
		ShippingCost:   15000,
		Negara:         "Indonesia",
	}
	updated, err := s.repo.ShippingAddressCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Office", updated.Alamat)

	// 7. Trash
	trashed, err := s.repo.ShippingAddressCommand.Trash(ctx, addrID)
	s.NoError(err)
	s.NotNil(trashed)

	// 8. FindTrashed
	trashedItems, err := s.repo.ShippingAddressQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 9. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.ShippingAddressQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(addrID, int(item.ShippingAddressID))
	}

	// 10. Restore
	restored, err := s.repo.ShippingAddressCommand.Restore(ctx, addrID)
	s.NoError(err)
	s.NotNil(restored)

	// 11. Trash again then DeletePermanent
	_, err = s.repo.ShippingAddressCommand.Trash(ctx, addrID)
	s.Require().NoError(err)

	success, err := s.repo.ShippingAddressCommand.DeletePermanent(ctx, addrID)
	s.NoError(err)
	s.True(success)

	// 12. RestoreAll
	created2, _ := s.repo.ShippingAddressCommand.Create(ctx, &requests.CreateShippingAddressRequest{
		OrderID: &orderID, Alamat: "A", Provinsi: "P", Kota: "K", Courier: "C",
		ShippingMethod: "SM", ShippingCost: 1000, Negara: "ID",
	})
	s.repo.ShippingAddressCommand.Trash(ctx, int(created2.ShippingAddressID))

	resRestore, err := s.repo.ShippingAddressCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 13. DeleteAll
	s.repo.ShippingAddressCommand.Trash(ctx, int(created2.ShippingAddressID))
	resDelete, err := s.repo.ShippingAddressCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestShippingAddressRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ShippingAddressRepositoryTestSuite))
}
