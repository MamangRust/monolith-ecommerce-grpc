package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-user/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	roleAddr := viper.GetString("GRPC_ROLE_ADDR")
	
	roleConn, err := grpc.NewClient(
		roleAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to role service: %w", err)
	}

	roleClient := pb.NewRoleQueryServiceClient(roleConn)

	repos := repository.NewRepositories(srv.DB, roleClient)
	hashing := hash.NewHashingPassword()
	obs, _ := observability.NewObservability("user-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Repositories:  repos,
		Hash:          hashing,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterUserQueryServiceServer(gs, h.UserQuery)
		pb.RegisterUserCommandServiceServer(gs, h.UserCommand)
	}

	return srv, nil
}
