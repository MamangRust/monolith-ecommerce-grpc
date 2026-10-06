package repository

import (
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency. Callers build them with adapter.WithDependencyGuard (typically
// around resilience.NewDependencyGuard) so the repository owns the wiring.
type GuardOptions struct {
	User []adapter.GuardOption
}

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	UserQuery               UserQueryRepository
}

func NewRepositories(DB *db.Queries, userQueryClient pbuser.UserQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions

	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(DB),
		MerchantCommand:         NewMerchantCommandRepository(DB),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(DB),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(DB),
		UserQuery:               useradapter.NewQueryAdapter(userQueryClient, g.User...),
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
