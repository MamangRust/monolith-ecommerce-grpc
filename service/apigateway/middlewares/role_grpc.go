package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	apicache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RoleValidatorGRPC validates a user's roles via the role service gRPC
// (FindByUserId) and stores the role names in the Echo context.
//
// Ini menggantikan RoleValidator berbasis Kafka (request-response ke topic
// "request-role"/"response-role") yang tidak berfungsi di stack lokal sehingga
// semua route admin memakai middleware ini timeout 408. Dengan gRPC langsung,
// role diverifikasi secara sinkron dan RequireRoles dapat memutuskan 403.
func RoleValidatorGRPC(client pb.RoleQueryServiceClient, logger logger.LoggerInterface, cache apicache.RoleCache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDVal := c.Get("user_id")
			if userIDVal == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
			}

			userID, err := extractUserIDGeneric(userIDVal)
			if err != nil {
				logger.Error("Invalid User ID format", zap.Any("value", userIDVal), zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid User ID format")
			}

			ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
			defer cancel()

			if roles, found := cache.GetRoleCache(ctx, strconv.Itoa(userID)); found {
				c.Set("role_names", roles)
				return next(c)
			}

			res, err := client.FindByUserId(ctx, &pb.FindByIdUserRoleRequest{UserId: int32(userID)})
			if err != nil {
				logger.Error("Role validation via gRPC failed",
					zap.Int("user_id", userID), zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "Role validation failed")
			}

			roles := make([]string, 0, len(res.GetData()))
			for _, r := range res.GetData() {
				roles = append(roles, r.GetName())
			}
			if len(roles) == 0 {
				logger.Debug("Role validation failed (no roles)",
					zap.Int("user_id", userID))
				return echo.NewHTTPError(http.StatusUnauthorized, "Role validation failed")
			}

			cache.SetRoleCache(ctx, strconv.Itoa(userID), roles)
			c.Set("role_names", roles)
			return next(c)
		}
	}
}

func extractUserIDGeneric(userIDVal interface{}) (int, error) {
	switch val := userIDVal.(type) {
	case float64:
		return int(val), nil
	case int:
		return val, nil
	case int32:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, fmt.Errorf("unknown user ID type %T", userIDVal)
	}
}
