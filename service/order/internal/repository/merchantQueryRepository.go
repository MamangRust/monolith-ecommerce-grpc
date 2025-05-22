package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantQueryRepository {
	return &merchantQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantQueryRepository) FindById(user_id int) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByID(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchant_errors.ErrFindByIdMerchant
	}

	return r.mapping.ToMerchantRecord(res), nil
}
