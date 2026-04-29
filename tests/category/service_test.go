package category_test

import (
	"context"
	"testing"

	cat_cache "github.com/MamangRust/monolith-ecommerce-grpc-category/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CategoryServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *CategoryServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Category dependencies
	mencache := cat_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(queries)

	s.svc = service.NewService(&service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *CategoryServiceTestSuite) TestCategoryLifecycle() {
	ctx := context.Background()
	slug := "initial-category"
	req := &requests.CreateCategoryRequest{
		Name:          "Initial Category",
		Description:   "Testing service layer",
		SlugCategory:  &slug,
		ImageCategory: "test.jpg",
	}

	// 1. Create
	created, err := s.svc.CategoryCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	catID := int(created.CategoryID)

	// 2. FindByID
	found, err := s.svc.CategoryQuery.FindByID(ctx, catID)
	s.Require().NoError(err)
	s.Equal(req.Name, found.Name)

	// 3. Update
	updateReq := &requests.UpdateCategoryRequest{
		CategoryID:    &catID,
		Name:          "Updated Category Name",
		Description:   "Updated desc",
		SlugCategory:  &slug,
		ImageCategory: "updated.jpg",
	}
	updated, err := s.svc.CategoryCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(updateReq.Name, updated.Name)

	// 4. FindAll
	_, total, err := s.svc.CategoryQuery.FindAll(ctx, &requests.FindAllCategory{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.CategoryCommand.Trash(ctx, catID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.CategoryQuery.FindTrashed(ctx, &requests.FindAllCategory{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.svc.CategoryQuery.FindActive(ctx, &requests.FindAllCategory{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, c := range active {
		s.NotEqual(catID, int(c.CategoryID))
	}

	// 8. Restore
	_, err = s.svc.CategoryCommand.Restore(ctx, catID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.svc.CategoryCommand.Trash(ctx, catID)
	s.Require().NoError(err)
	success, err := s.svc.CategoryCommand.DeletePermanent(ctx, catID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	s1 := "s1"
	s2 := "s2"
	c1, _ := s.svc.CategoryCommand.Create(ctx, &requests.CreateCategoryRequest{Name: "C1", Description: "D1", SlugCategory: &s1, ImageCategory: "I1"})
	c2, _ := s.svc.CategoryCommand.Create(ctx, &requests.CreateCategoryRequest{Name: "C2", Description: "D2", SlugCategory: &s2, ImageCategory: "I2"})
	
	s.svc.CategoryCommand.Trash(ctx, int(c1.CategoryID))
	s.svc.CategoryCommand.Trash(ctx, int(c2.CategoryID))

	resRestoreAll, err := s.svc.CategoryCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.CategoryCommand.Trash(ctx, int(c1.CategoryID))
	s.svc.CategoryCommand.Trash(ctx, int(c2.CategoryID))

	resDeleteAll, err := s.svc.CategoryCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestCategoryServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryServiceTestSuite))
}
