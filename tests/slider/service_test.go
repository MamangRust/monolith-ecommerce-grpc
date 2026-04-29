package slider_test

import (
	"context"
	"testing"

	slider_cache "github.com/MamangRust/monolith-ecommerce-grpc-slider/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type SliderServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *SliderServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Slider dependencies
	mencache := slider_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(queries)

	s.svc = service.NewService(&service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *SliderServiceTestSuite) TestSliderLifecycle() {
	ctx := context.Background()

	// 1. Create
	req := &requests.CreateSliderRequest{
		Nama:     "Test Slider",
		FilePath: "slider.jpg",
	}
	created, err := s.svc.SliderCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	sliderID := int(created.SliderID)

	// 2. FindByID
	found, err := s.svc.SliderQuery.FindByID(ctx, sliderID)
	s.Require().NoError(err)
	s.Equal(req.Nama, found.Name)

	// 3. Update
	newName := "Updated Test Slider"
	updateReq := &requests.UpdateSliderRequest{
		ID:       &sliderID,
		Nama:     newName,
		FilePath: req.FilePath,
	}
	updated, err := s.svc.SliderCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newName, updated.Name)

	// 4. FindAll
	_, total, err := s.svc.SliderQuery.FindAll(ctx, &requests.FindAllSlider{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.SliderCommand.Trash(ctx, sliderID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.SliderQuery.FindTrashed(ctx, &requests.FindAllSlider{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.svc.SliderQuery.FindActive(ctx, &requests.FindAllSlider{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, sl := range active {
		s.NotEqual(sliderID, int(sl.SliderID))
	}

	// 8. Restore
	_, err = s.svc.SliderCommand.Restore(ctx, sliderID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.svc.SliderCommand.Trash(ctx, sliderID)
	s.Require().NoError(err)
	success, err := s.svc.SliderCommand.DeletePermanent(ctx, sliderID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	s1, _ := s.svc.SliderCommand.Create(ctx, &requests.CreateSliderRequest{Nama: "S1", FilePath: "G1"})
	s2, _ := s.svc.SliderCommand.Create(ctx, &requests.CreateSliderRequest{Nama: "S2", FilePath: "G2"})
	
	s.svc.SliderCommand.Trash(ctx, int(s1.SliderID))
	s.svc.SliderCommand.Trash(ctx, int(s2.SliderID))

	resRestoreAll, err := s.svc.SliderCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.SliderCommand.Trash(ctx, int(s1.SliderID))
	s.svc.SliderCommand.Trash(ctx, int(s2.SliderID))

	resDeleteAll, err := s.svc.SliderCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestSliderServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SliderServiceTestSuite))
}
