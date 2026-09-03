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

type categoryPermanentDeleteDBTXStub struct {
	err error
}

func (s *categoryPermanentDeleteDBTXStub) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, s.err
}

func (s *categoryPermanentDeleteDBTXStub) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (s *categoryPermanentDeleteDBTXStub) QueryRow(context.Context, string, ...interface{}) pgx.Row {
	return categoryPermanentDeleteRowStub{}
}

type categoryPermanentDeleteRowStub struct{}

func (categoryPermanentDeleteRowStub) Scan(...interface{}) error {
	return nil
}

func TestCategoryCommandRepositoryDeletePermanentReturnsConflictOnForeignKey(t *testing.T) {
	repo := NewCategoryCommandRepository(db.New(&categoryPermanentDeleteDBTXStub{
		err: &pgconn.PgError{Code: "23503", ConstraintName: "products_category_id_fkey"},
	}))

	deleted, err := repo.DeletePermanent(context.Background(), 7)
	assertCategoryPermanentDeleteConflict(t, deleted, err)
}

func TestCategoryCommandRepositoryDeleteAllReturnsConflictOnForeignKey(t *testing.T) {
	repo := NewCategoryCommandRepository(db.New(&categoryPermanentDeleteDBTXStub{
		err: &pgconn.PgError{Code: "23503", ConstraintName: "products_category_id_fkey"},
	}))

	deleted, err := repo.DeleteAll(context.Background())
	assertCategoryPermanentDeleteConflict(t, deleted, err)
}

func assertCategoryPermanentDeleteConflict(t *testing.T, deleted bool, err error) {
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
