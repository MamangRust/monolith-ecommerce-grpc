package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryDeleteAllServiceStub struct{}

func (categoryDeleteAllServiceStub) Create(context.Context, *requests.CreateCategoryRequest) (*db.CreateCategoryRow, error) {
	return nil, fmt.Errorf("unexpected Create call")
}

func (categoryDeleteAllServiceStub) Update(context.Context, *requests.UpdateCategoryRequest) (*db.UpdateCategoryRow, error) {
	return nil, fmt.Errorf("unexpected Update call")
}

func (categoryDeleteAllServiceStub) Trash(context.Context, int) (*db.Category, error) {
	return nil, fmt.Errorf("unexpected Trash call")
}

func (categoryDeleteAllServiceStub) Restore(context.Context, int) (*db.Category, error) {
	return nil, fmt.Errorf("unexpected Restore call")
}

func (categoryDeleteAllServiceStub) DeletePermanent(context.Context, int) (bool, error) {
	return false, fmt.Errorf("unexpected DeletePermanent call")
}

func (categoryDeleteAllServiceStub) RestoreAll(context.Context) (bool, error) {
	return false, fmt.Errorf("unexpected RestoreAll call")
}

func (categoryDeleteAllServiceStub) DeleteAll(context.Context) (bool, error) {
	return false, shared_errors.NewConflictError("cannot permanently delete categories while related products exist")
}

func TestCategoryCommandHandlerDeleteAllPreservesConflict(t *testing.T) {
	handler := NewCategoryCommandHandler(
		categoryDeleteAllServiceStub{},
		&logger.Logger{Log: zap.NewNop()},
	)

	_, err := handler.DeleteAllCategoryPermanent(context.Background(), &emptypb.Empty{})
	if status.Code(err) != codes.AlreadyExists {
		t.Fatalf("gRPC code = %s, want %s", status.Code(err), codes.AlreadyExists)
	}
	if status.Convert(err).Message() != "cannot permanently delete categories while related products exist" {
		t.Fatalf("gRPC message = %q, want conflict message", status.Convert(err).Message())
	}
}

var _ service.CategoryCommandService = categoryDeleteAllServiceStub{}
