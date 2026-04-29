package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type userQueryRepository struct {
	client pb.UserQueryServiceClient
}

func NewUserQueryRepository(client pb.UserQueryServiceClient) *userQueryRepository {
	return &userQueryRepository{
		client: client,
	}
}

func (r *userQueryRepository) FindByID(ctx context.Context, user_id int) (*db.GetUserByIDRow, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdUserRequest{Id: int32(user_id)})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &db.GetUserByIDRow{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		Password:  "", // Not provided by gRPC, but required by struct
	}, nil
}
