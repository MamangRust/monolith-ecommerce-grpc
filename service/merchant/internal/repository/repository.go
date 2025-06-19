package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	UserQuery               UserQueryRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		MerchantCommand:         NewMerchantCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantRecordMapper),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantDocumentRecordMapper),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.MerchantDocumentRecordMapper),
		UserQuery:               NewUserQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRecordMapper),
	}
}
