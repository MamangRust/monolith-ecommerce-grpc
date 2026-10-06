package apps

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-user/cache"
	"github.com/MamangRust/monolith-ecommerce-user/handler"
	"github.com/MamangRust/monolith-ecommerce-user/repository"
	"github.com/MamangRust/monolith-ecommerce-user/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	pbuserrole "github.com/MamangRust/monolith-ecommerce-pb/user_role"
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

	roleClient := pbrole.NewRoleQueryServiceClient(roleConn)
	// The user-role service piggybacks on the role service, so its client
	// shares the same connection/port.
	userRoleClient := pbuserrole.NewUserRoleServiceClient(roleConn)

	guardRole := resilience.NewDependencyGuard("role", 5, 30, 100, 3*time.Second, srv.Logger)
	guardUserRole := resilience.NewDependencyGuard("user_role", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(&repository.Deps{
		Db:              srv.DB,
		RoleQueryClient: roleClient,
		UserRoleClient:  userRoleClient,
		Guard: repository.GuardOptions{
			Role: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardRole),
			},
			UserRole: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardUserRole),
			},
		},
	})
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
		pbuser.RegisterUserQueryServiceServer(gs, h.UserQuery)
		pbuser.RegisterUserCommandServiceServer(gs, h.UserCommand)
	}

	return srv, nil
}
