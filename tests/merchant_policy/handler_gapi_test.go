package merchant_policy_test

import (
	"context"
	"testing"

	policy_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/cache"
	policy_handler "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/handler"
	policy_repo "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/repository"
	policy_service "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MerchantPolicyGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pb.MerchantPolicyQueryServiceClient
	commandClient pb.MerchantPolicyCommandServiceClient
}

func (s *MerchantPolicyGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies
	s.SetupUserService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Policy dependencies
	mencache := policy_cache.NewMencache(cacheStore)
	repos := policy_repo.NewRepositories(
		queries,
		pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
	svc := policy_service.NewService(&policy_service.Deps{
		Cache:         mencache,
		Repository:    repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := policy_handler.NewHandler(&policy_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pb.RegisterMerchantPolicyQueryServiceServer(server, handler.MerchantPolicyQuery)
	pb.RegisterMerchantPolicyCommandServiceServer(server, handler.MerchantPolicyCommand)
	
	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pb.NewMerchantPolicyQueryServiceClient(conn)
	s.commandClient = pb.NewMerchantPolicyCommandServiceClient(conn)
}

func (s *MerchantPolicyGapiTestSuite) TestMerchantPolicyGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := int32(s.SeedMerchant(ctx, userID))

	// 2. Create
	createRes, err := s.commandClient.Create(ctx, &pb.CreateMerchantPoliciesRequest{
		MerchantId:  merchID,
		PolicyType:  "Warranty",
		Title:       "Warranty Policy",
		Description: "1 year limited warranty",
	})
	s.NoError(err)
	s.NotNil(createRes)
	policyID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.queryClient.FindById(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)
	s.Equal("Warranty", getRes.Data.PolicyType)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pb.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pb.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.commandClient.Update(ctx, &pb.UpdateMerchantPoliciesRequest{
		MerchantPolicyId: policyID,
		PolicyType:       "Warranty Updated",
		Title:            "Warranty Policy Updated",
		Description:      "2 years limited warranty",
	})
	s.NoError(err)
	s.Equal("Warranty Updated", updateRes.Data.PolicyType)

	// 7. Trash
	_, err = s.commandClient.TrashedMerchantPolicies(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pb.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreMerchantPolicies(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedMerchantPolicies(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: policyID})
	_, err = s.commandClient.DeleteMerchantPoliciesPermanent(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllMerchantPolicies(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllMerchantPoliciesPermanent(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func TestMerchantPolicyGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantPolicyGapiTestSuite))
}
