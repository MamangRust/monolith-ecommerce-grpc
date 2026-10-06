// Package transaction adapts the Transaction service gRPC command API.
package transaction

import (
	"context"
	"fmt"

	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CommandRepository is the contract consumers depend on for transaction writes.
type CommandRepository interface {
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type repository struct {
	client pbtransaction.TransactionCommandServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a transaction command adapter. Passing zero options leaves the
// guard nil, which makes DependencyGuard.Call a plain passthrough.
func New(client pbtransaction.TransactionCommandServiceClient, opts ...adapter.GuardOption) *repository {
	r := &repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewCommandAdapter returns the adapter typed as the write-only contract.
func NewCommandAdapter(client pbtransaction.TransactionCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(client, opts...)
}

func (r *repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *repository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	_, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbtransaction.ApiResponseTransactionDelete, error) {
		return r.client.DeleteTransactionByOrderPermanent(callCtx, &pbtransaction.FindByIdTransactionRequest{
			Id: int32(orderID),
		})
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) DeleteAll(ctx context.Context) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pbtransaction.ApiResponseTransactionAll, error) {
		return r.client.DeleteAllTransactionPermanent(callCtx, &emptypb.Empty{})
	})
	if err != nil {
		return false, err
	}
	return resp != nil && resp.Status == "success", nil
}
