package merchant_detail_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	detail_cache "github.com/MamangRust/monolith-ecommerce-merchant_detail/cache"
	detail_handler "github.com/MamangRust/monolith-ecommerce-merchant_detail/handler"
	detail_repo "github.com/MamangRust/monolith-ecommerce-merchant_detail/repository"
	detail_service "github.com/MamangRust/monolith-ecommerce-merchant_detail/service"
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pbmerchant_detail "github.com/MamangRust/monolith-ecommerce-pb/merchant_detail"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
)

type MerchantDetailGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbmerchant_detail.MerchantDetailQueryServiceClient
	commandClient pbmerchant_detail.MerchantDetailCommandServiceClient
}

func (s *MerchantDetailGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Detail dependencies
	mencache := detail_cache.NewMencache(cacheStore)
	repos := detail_repo.NewRepositories(
		queries,
		pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
	svc := detail_service.NewService(&detail_service.Deps{
		Cache:         mencache,
		Repository:    repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := detail_handler.NewHandler(&detail_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbmerchant_detail.RegisterMerchantDetailQueryServiceServer(server, handler.MerchantDetailQuery)
	pbmerchant_detail.RegisterMerchantDetailCommandServiceServer(server, handler.MerchantDetailCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbmerchant_detail.NewMerchantDetailQueryServiceClient(conn)
	s.commandClient = pbmerchant_detail.NewMerchantDetailCommandServiceClient(conn)
}

func (s *MerchantDetailGapiTestSuite) TestMerchantDetailGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := int32(s.SeedMerchant(ctx, userID))

	// 2. Create
	createRes, err := s.commandClient.Create(ctx, &pbmerchant_detail.CreateMerchantDetailRequest{
		MerchantId:       merchID,
		DisplayName:      "GAPI Detail",
		ShortDescription: "GAPI Description",
		WebsiteUrl:       "https://gapi.com",
	})
	s.NoError(err)
	s.NotNil(createRes)
	detailID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbmerchant_detail.FindByIdMerchantDetailRequest{Id: detailID})
	s.NoError(err)
	s.Equal("GAPI Detail", getRes.Data.DisplayName)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbmerchant.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbmerchant.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.commandClient.Update(ctx, &pbmerchant_detail.UpdateMerchantDetailRequest{
		MerchantDetailId: detailID,
		DisplayName:      "GAPI Detail Updated",
		ShortDescription: "Updated Description",
		WebsiteUrl:       "https://gapi.com/updated",
	})
	s.NoError(err)
	s.Equal("GAPI Detail Updated", updateRes.Data.DisplayName)

	// 7. Trash
	_, err = s.commandClient.TrashedMerchantDetail(ctx, &pbmerchant_detail.FindByIdMerchantDetailRequest{Id: detailID})
	s.NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbmerchant.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreMerchantDetail(ctx, &pbmerchant_detail.FindByIdMerchantDetailRequest{Id: detailID})
	s.NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedMerchantDetail(ctx, &pbmerchant_detail.FindByIdMerchantDetailRequest{Id: detailID})
	_, err = s.commandClient.DeleteMerchantDetailPermanent(ctx, &pbmerchant_detail.FindByIdMerchantDetailRequest{Id: detailID})
	s.NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllMerchantDetail(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllMerchantDetailPermanent(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func TestMerchantDetailGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDetailGapiTestSuite))
}
