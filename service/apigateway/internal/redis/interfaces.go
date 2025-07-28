package mencache

import "context"

type RoleCache interface {
	GetRoleCache(ctx context.Context, userID string) ([]string, bool)
	SetRoleCache(ctx context.Context, userID string, roles []string)
}
