package slider_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type SliderRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo     *repository.Repositories
	sliderID int
}

func (s *SliderRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(queries)
}

func (s *SliderRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *SliderRepositoryTestSuite) TestSliderLifecycle() {
	ctx := context.Background()

	pageReq := &requests.FindAllSlider{Search: "", Page: 1, PageSize: 10}

	// 1. Create
	req := &requests.CreateSliderRequest{
		Nama:     "Spring Collection",
		FilePath: "http://example.com/spring.jpg",
	}
	created, err := s.repo.SliderCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.Nama, created.Name)
	sliderID := int(created.SliderID)

	// 2. FindByID
	found, err := s.repo.SliderQuery.FindByID(ctx, sliderID)
	s.NoError(err)
	s.Equal(created.Image, found.Image)

	// 3. FindAll
	all, err := s.repo.SliderQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 4. FindActive (entity not yet trashed — should appear)
	active, err := s.repo.SliderQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 5. Update
	updateReq := &requests.UpdateSliderRequest{
		ID:       &sliderID,
		Nama:     "Updated Collection",
		FilePath: "http://example.com/updated.jpg",
	}
	updated, err := s.repo.SliderCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated Collection", updated.Name)

	// 6. Trash
	trashed, err := s.repo.SliderCommand.Trash(ctx, sliderID)
	s.NoError(err)
	s.NotNil(trashed)

	// 7. FindTrashed
	trashedItems, err := s.repo.SliderQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 8. FindActive after trash — should NOT include trashed entity
	activeAfterTrash, err := s.repo.SliderQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(sliderID, int(item.SliderID))
	}

	// 9. Restore
	restored, err := s.repo.SliderCommand.Restore(ctx, sliderID)
	s.NoError(err)
	s.NotNil(restored)

	// 10. Trash again then DeletePermanent
	_, err = s.repo.SliderCommand.Trash(ctx, sliderID)
	s.Require().NoError(err)

	success, err := s.repo.SliderCommand.DeletePermanent(ctx, sliderID)
	s.NoError(err)
	s.True(success)

	// 11. RestoreAll
	s.repo.SliderCommand.Create(ctx, &requests.CreateSliderRequest{Nama: "S1", FilePath: "http://example.com/s1.jpg"})
	second, _ := s.repo.SliderCommand.Create(ctx, &requests.CreateSliderRequest{Nama: "S2", FilePath: "http://example.com/s2.jpg"})
	s.repo.SliderCommand.Trash(ctx, int(second.SliderID))

	resRestore, err := s.repo.SliderCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 12. DeleteAll
	s.repo.SliderCommand.Trash(ctx, int(second.SliderID))
	resDelete, err := s.repo.SliderCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestSliderRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SliderRepositoryTestSuite))
}
