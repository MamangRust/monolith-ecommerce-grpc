package repository

import (
	"context"
	"log"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type cartQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CartRecordMapping
}

func NewCartQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CartRecordMapping) *cartQueryRepository {
	return &cartQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cartQueryRepository) FindCarts(req *requests.FindAllCarts) ([]*record.CartRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCartsParams{
		UserID:  int32(req.UserID),
		Column2: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCarts(r.ctx, reqDb)

	if err != nil {
		log.Fatal(err)
		return nil, nil, cart_errors.ErrFindAllCarts
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCartsRecordPagination(res), &totalCount, nil
}
