package repository

import (
	"context"
	"fmt"

	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionCommandRepository struct {
	client pb.TransactionCommandServiceClient
}

func NewTransactionCommandRepository(client pb.TransactionCommandServiceClient) TransactionCommandRepository {
	return &transactionCommandRepository{client: client}
}

func (r *transactionCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, order_id int) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	_, err := r.client.DeleteTransactionByOrderPermanent(ctx, &pb.FindByIdTransactionRequest{
		Id: int32(order_id),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *transactionCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	res, err := r.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return res.Status == "success", nil
}
