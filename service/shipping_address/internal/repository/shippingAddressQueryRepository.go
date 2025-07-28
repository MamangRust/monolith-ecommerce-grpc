package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type shippingAddressQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.ShippingAddressMapping
}

func NewShippingAddressQueryRepository(db *db.Queries, mapping recordmapper.ShippingAddressMapping) *shippingAddressQueryRepository {
	return &shippingAddressQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *shippingAddressQueryRepository) FindAllShippingAddress(ctx context.Context, req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddress(ctx, reqDb)

	if err != nil {
		return nil, nil, shippingaddress_errors.ErrFindAllShippingAddress
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToShippingAddresssRecordPagination(res), &totalCount, nil
}

func (r *shippingAddressQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressActive(ctx, reqDb)

	if err != nil {
		return nil, nil, shippingaddress_errors.ErrFindActiveShippingAddress
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToShippingAddresssRecordActivePagination(res), &totalCount, nil
}

func (r *shippingAddressQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, shippingaddress_errors.ErrFindTrashedShippingAddress
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToShippingAddresssRecordTrashedPagination(res), &totalCount, nil
}

func (r *shippingAddressQueryRepository) FindById(ctx context.Context, shipping_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.GetShippingByID(ctx, int32(shipping_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByID
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}

func (r *shippingAddressQueryRepository) FindByOrder(ctx context.Context, order_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.GetShippingAddressByOrderID(ctx, int32(order_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}

	return r.mapping.ToShippingAddressRecord(res), nil
}
