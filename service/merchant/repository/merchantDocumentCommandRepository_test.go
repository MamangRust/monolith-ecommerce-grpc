package repository

import (
	"context"
	"strings"
	"testing"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type merchantDocumentUpdateDBTXStub struct {
	query string
	args  []interface{}
}

func (s *merchantDocumentUpdateDBTXStub) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (s *merchantDocumentUpdateDBTXStub) Query(context.Context, string, ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (s *merchantDocumentUpdateDBTXStub) QueryRow(_ context.Context, query string, args ...interface{}) pgx.Row {
	s.query = query
	s.args = args
	return merchantDocumentUpdateRowStub{}
}

type merchantDocumentUpdateRowStub struct{}

func (merchantDocumentUpdateRowStub) Scan(...interface{}) error {
	return nil
}

func TestMerchantDocumentCommandRepositoryUpdateTargetsDocumentID(t *testing.T) {
	const (
		documentID = 42
		merchantID = 7
	)

	dbtx := &merchantDocumentUpdateDBTXStub{}
	repo := NewMerchantDocumentCommandRepository(db.New(dbtx))
	note := "replace before review"

	_, err := repo.Update(context.Background(), &requests.UpdateMerchantDocumentRequest{
		DocumentID:   intPtr(documentID),
		MerchantID:   merchantID,
		DocumentType: "business_license",
		DocumentUrl:  "https://example.com/license.pdf",
		Status:       "pending",
		Note:         note,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	assertMerchantDocumentUpdateTarget(t, dbtx, documentID, merchantID)
}

func TestMerchantDocumentCommandRepositoryUpdateStatusTargetsDocumentID(t *testing.T) {
	const (
		documentID = 84
		merchantID = 13
	)

	dbtx := &merchantDocumentUpdateDBTXStub{}
	repo := NewMerchantDocumentCommandRepository(db.New(dbtx))
	note := "verified"

	_, err := repo.UpdateStatus(context.Background(), &requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: intPtr(documentID),
		MerchantID: merchantID,
		Status:     "approved",
		Note:       note,
	})
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}

	assertMerchantDocumentUpdateTarget(t, dbtx, documentID, merchantID)
}

func assertMerchantDocumentUpdateTarget(t *testing.T, dbtx *merchantDocumentUpdateDBTXStub, documentID, merchantID int) {
	t.Helper()

	if !strings.Contains(dbtx.query, "document_id = $1") {
		t.Fatalf("update query target = %q, want document_id", dbtx.query)
	}
	if len(dbtx.args) == 0 {
		t.Fatal("update query received no arguments")
	}
	gotDocumentID, ok := dbtx.args[0].(int32)
	if !ok {
		t.Fatalf("first update argument type = %T, want int32", dbtx.args[0])
	}
	if gotDocumentID != int32(documentID) {
		t.Fatalf("first update argument = %d, want DocumentID %d (MerchantID is %d)", gotDocumentID, documentID, merchantID)
	}
	if gotDocumentID == int32(merchantID) {
		t.Fatalf("first update argument incorrectly used MerchantID %d", merchantID)
	}
}

func intPtr(value int) *int {
	return &value
}
