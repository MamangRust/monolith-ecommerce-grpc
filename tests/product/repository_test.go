package product_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ProductRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *ProductRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies (Role, User, Category, Merchant)
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(
		queries,
		pb.NewCategoryQueryServiceClient(s.Conns["category"]),
		pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
}

func (s *ProductRepositoryTestSuite) TestProductCreate() {
	ctx := context.Background()
	
	// Seed a category first
	catRes, err := pb.NewCategoryCommandServiceClient(s.Conns["category"]).Create(ctx, &pb.CreateCategoryRequest{
		Name:          "Product Test Category",
		Description:   "Test Description",
		SlugCategory:  "product-test-cat",
		ImageCategory: "seed.jpg",
	})
	s.NoError(err)
	catID := int32(catRes.Data.Id)

	// Seed User and Merchant
	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)

	rating := 5
	slug := "test-product"
	req := &requests.CreateProductRequest{
		Name:         "Test Product",
		Description:  "Product for repository test",
		Price:        100,
		CountInStock: 50,
		CategoryID:   int(catID),
		MerchantID:   merchantID,
		Brand:        "Test Brand",
		Weight:       1,
		Rating:       &rating,
		SlugProduct:  &slug,
		ImageProduct: "test.jpg",
	}

	created, err := s.repo.ProductCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.Name, created.Name)
}

func TestProductRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductRepositoryTestSuite))
}
