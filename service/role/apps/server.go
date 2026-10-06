package apps

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-role/cache"
	"github.com/MamangRust/monolith-ecommerce-role/handler"
	"github.com/MamangRust/monolith-ecommerce-role/repository"
	"github.com/MamangRust/monolith-ecommerce-role/service"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"google.golang.org/grpc"

	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	pbuserrole "github.com/MamangRust/monolith-ecommerce-pb/user_role"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	obs, _ := observability.NewObservability("role-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbrole.RegisterRoleQueryServiceServer(gs, h.RoleQuery)
		pbrole.RegisterRoleCommandServiceServer(gs, h.RoleCommand)
		pbuserrole.RegisterUserRoleServiceServer(gs, h.UserRole)
	}

	return srv, nil
}
