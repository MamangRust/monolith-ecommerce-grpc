package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type merchantQueryRepository struct {
	client pb.MerchantQueryServiceClient
}

func NewMerchantQueryRepository(client pb.MerchantQueryServiceClient) *merchantQueryRepository {
	return &merchantQueryRepository{
		client: client,
	}
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, merchant_id int) (*db.GetMerchantByIDRow, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdMerchantRequest{Id: int32(merchant_id)})
	if err != nil {
		// pertahankan status gRPC dari dependency service (NotFound -> 404, dst)
		return nil, err
	}

	return &db.GetMerchantByIDRow{
		MerchantID:   res.Data.Id,
		UserID:       res.Data.UserId,
		Name:         res.Data.Name,
		Description:  &res.Data.Description,
		Address:      &res.Data.Address,
		ContactEmail: &res.Data.ContactEmail,
		ContactPhone: &res.Data.ContactPhone,
		Status:       res.Data.Status,
	}, nil
}
