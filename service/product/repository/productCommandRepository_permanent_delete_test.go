package repository

import (
	"context"
	"errors"
	"testing"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	shared_errors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type productPermanentDeleteDBTXStub struct {
	err error
}

func (s *productPermanentDeleteDBTXStub) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, s.err
}

func (s *productPermanentDeleteDBTXStub) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (s *productPermanentDeleteDBTXStub) QueryRow(context.Context, string, ...interface{}) pgx.Row {
	return productPermanentDeleteRowStub{}
}

type productPermanentDeleteRowStub struct{}

func (productPermanentDeleteRowStub) Scan(...interface{}) error {
	return nil
}

func TestProductCommandRepositoryDeletePermanentReturnsConflictOnForeignKey(t *testing.T) {
	repo := NewProductCommandRepository(db.New(&productPermanentDeleteDBTXStub{
		err: &pgconn.PgError{Code: "23503", ConstraintName: "reviews_product_id_fkey"},
	}))

	deleted, err := repo.DeletePermanent(context.Background(), 7)
	assertProductPermanentDeleteConflict(t, deleted, err)
}

func TestProductCommandRepositoryDeleteAllReturnsConflictOnForeignKey(t *testing.T) {
	repo := NewProductCommandRepository(db.New(&productPermanentDeleteDBTXStub{
		err: &pgconn.PgError{Code: "23503", ConstraintName: "reviews_product_id_fkey"},
	}))

	deleted, err := repo.DeleteAll(context.Background())
	assertProductPermanentDeleteConflict(t, deleted, err)
}

func assertProductPermanentDeleteConflict(t *testing.T, deleted bool, err error) {
	t.Helper()

	if deleted {
		t.Fatal("delete reported success after a foreign-key violation")
	}

	var appErr *shared_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want AppError", err, err)
	}
	if appErr.Type != shared_errors.ErrorTypeConflict || appErr.Code != 409 {
		t.Fatalf("error type/code = %s/%d, want CONFLICT/409", appErr.Type, appErr.Code)
	}
}
