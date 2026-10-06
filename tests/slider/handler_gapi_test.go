package slider_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pbslider "github.com/MamangRust/monolith-ecommerce-pb/slider"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	slider_cache "github.com/MamangRust/monolith-ecommerce-slider/cache"
	slider_handler "github.com/MamangRust/monolith-ecommerce-slider/handler"
	slider_repo "github.com/MamangRust/monolith-ecommerce-slider/repository"
	slider_service "github.com/MamangRust/monolith-ecommerce-slider/service"
	tests "github.com/MamangRust/monolith-ecommerce-test"
)

type SliderGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbslider.SliderQueryServiceClient
	commandClient pbslider.SliderCommandServiceClient
}

func (s *SliderGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Slider dependencies
	mencache := slider_cache.NewMencache(cacheStore)
	repos := slider_repo.NewRepositories(queries)
	svc := slider_service.NewService(&slider_service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := slider_handler.NewHandler(&slider_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbslider.RegisterSliderQueryServiceServer(server, handler.SliderQuery)
	pbslider.RegisterSliderCommandServiceServer(server, handler.SliderCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbslider.NewSliderQueryServiceClient(conn)
	s.commandClient = pbslider.NewSliderCommandServiceClient(conn)
}

func (s *SliderGapiTestSuite) TestSliderGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createRes, err := s.commandClient.Create(ctx, &pbslider.CreateSliderRequest{
		Name:  "GAPI Slider",
		Image: "http://example.com/gapi-slider.jpg",
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	sliderID := createRes.Data.Id

	// 2. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbslider.FindByIdSliderRequest{Id: sliderID})
	s.Require().NoError(err)
	s.Equal("GAPI Slider", getRes.Data.Name)

	// 3. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbslider.FindAllSliderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbslider.FindAllSliderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.commandClient.Update(ctx, &pbslider.UpdateSliderRequest{
		Id:    sliderID,
		Name:  "GAPI Slider Updated",
		Image: "http://example.com/gapi-slider-updated.jpg",
	})
	s.Require().NoError(err)
	s.Equal("GAPI Slider Updated", updateRes.Data.Name)

	// 6. Trash
	_, err = s.commandClient.TrashedSlider(ctx, &pbslider.FindByIdSliderRequest{Id: sliderID})
	s.Require().NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbslider.FindAllSliderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.commandClient.RestoreSlider(ctx, &pbslider.FindByIdSliderRequest{Id: sliderID})
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, _ = s.commandClient.TrashedSlider(ctx, &pbslider.FindByIdSliderRequest{Id: sliderID})
	_, err = s.commandClient.DeleteSliderPermanent(ctx, &pbslider.FindByIdSliderRequest{Id: sliderID})
	s.Require().NoError(err)

	// 10. RestoreAll
	_, err = s.commandClient.RestoreAllSlider(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 11. DeleteAll
	_, err = s.commandClient.DeleteAllSliderPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestSliderGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SliderGapiTestSuite))
}
