package review_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type ReviewRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *ReviewRepositoryTestSuite) SetupSuite() {
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

func (s *ReviewRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *ReviewRepositoryTestSuite) TestReviewLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	pageReq := &requests.FindAllReview{Search: "", Page: 1, PageSize: 10}

	// 1. Create Review
	req := &requests.CreateReviewRequest{
		UserID:    userID,
		ProductID: prodID,
		Rating:    5,
		Comment:   "Excellent product!",
	}
	created, err := s.repo.ReviewCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(5), created.Rating)
	reviewID := int(created.ReviewID)

	// 2. FindByID
	found, err := s.repo.ReviewQuery.FindByID(ctx, reviewID)
	s.NoError(err)
	s.Equal(created.Comment, found.Comment)

	// 3. FindByProduct
	reviewsByProd, err := s.repo.ReviewQuery.FindByProduct(ctx, &requests.FindAllReviewByProduct{
		ProductID: prodID,
		Rating:    5,
		Page:      1,
		PageSize:  10,
		Search:    "",
	})
	s.NoError(err)
	s.NotEmpty(reviewsByProd)

	// 4. FindAll
	all, err := s.repo.ReviewQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 5. FindActive (before trash)
	active, err := s.repo.ReviewQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 6. Update
	updateReq := &requests.UpdateReviewRequest{
		ReviewID: &reviewID,
		Name:     "Updated",
		Rating:   4,
		Comment:  "Updated comment",
	}
	updated, err := s.repo.ReviewCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(4), updated.Rating)

	// 7. Trash
	trashed, err := s.repo.ReviewCommand.Trash(ctx, reviewID)
	s.NoError(err)
	s.NotNil(trashed)

	// 8. FindTrashed
	trashedItems, err := s.repo.ReviewQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 9. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.ReviewQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(reviewID, int(item.ReviewID))
	}

	// 10. Restore
	restored, err := s.repo.ReviewCommand.Restore(ctx, reviewID)
	s.NoError(err)
	s.NotNil(restored)

	// 11. Trash again then DeletePermanent
	_, err = s.repo.ReviewCommand.Trash(ctx, reviewID)
	s.Require().NoError(err)

	success, err := s.repo.ReviewCommand.DeletePermanent(ctx, reviewID)
	s.NoError(err)
	s.True(success)

	// 12. RestoreAll
	s.repo.ReviewCommand.Create(ctx, &requests.CreateReviewRequest{
		UserID: userID, ProductID: prodID, Rating: 3, Comment: "OK",
	})
	second, _ := s.repo.ReviewCommand.Create(ctx, &requests.CreateReviewRequest{
		UserID: userID, ProductID: prodID, Rating: 4, Comment: "Good",
	})
	s.repo.ReviewCommand.Trash(ctx, int(second.ReviewID))

	resRestore, err := s.repo.ReviewCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 13. DeleteAll
	s.repo.ReviewCommand.Trash(ctx, int(second.ReviewID))
	resDelete, err := s.repo.ReviewCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestReviewRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ReviewRepositoryTestSuite))
}
