package handler

import (
	"context"
	"fmt"
	"testing"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
)

type merchantDocumentCommandServiceStub struct {
	update       func(context.Context, *requests.UpdateMerchantDocumentRequest) (*db.UpdateMerchantDocumentRow, error)
	updateStatus func(context.Context, *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error)
}

func (s merchantDocumentCommandServiceStub) Create(context.Context, *requests.CreateMerchantDocumentRequest) (*db.CreateMerchantDocumentRow, error) {
	return nil, fmt.Errorf("unexpected Create call")
}

func (s merchantDocumentCommandServiceStub) Update(ctx context.Context, req *requests.UpdateMerchantDocumentRequest) (*db.UpdateMerchantDocumentRow, error) {
	if s.update == nil {
		return nil, fmt.Errorf("unexpected Update call")
	}
	return s.update(ctx, req)
}

func (s merchantDocumentCommandServiceStub) UpdateStatus(ctx context.Context, req *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error) {
	if s.updateStatus == nil {
		return nil, fmt.Errorf("unexpected UpdateStatus call")
	}
	return s.updateStatus(ctx, req)
}

func (s merchantDocumentCommandServiceStub) Trash(context.Context, int) (*db.MerchantDocument, error) {
	return nil, fmt.Errorf("unexpected Trash call")
}

func (s merchantDocumentCommandServiceStub) Restore(context.Context, int) (*db.MerchantDocument, error) {
	return nil, fmt.Errorf("unexpected Restore call")
}

func (s merchantDocumentCommandServiceStub) DeletePermanent(context.Context, int) (bool, error) {
	return false, fmt.Errorf("unexpected DeletePermanent call")
}

func (s merchantDocumentCommandServiceStub) RestoreAll(context.Context) (bool, error) {
	return false, fmt.Errorf("unexpected RestoreAll call")
}

func (s merchantDocumentCommandServiceStub) DeleteAll(context.Context) (bool, error) {
	return false, fmt.Errorf("unexpected DeleteAll call")
}

func newMerchantDocumentCommandHandlerForTest(stub merchantDocumentCommandServiceStub) pb.MerchantDocumentCommandServiceServer {
	return NewMerchantDocumentCommandHandler(stub, &logger.Logger{Log: zap.NewNop()})
}

func TestMerchantDocumentCommandHandlerUpdateUsesDocumentID(t *testing.T) {
	const (
		documentID = 42
		merchantID = 7
	)

	var gotRequest *requests.UpdateMerchantDocumentRequest
	handler := newMerchantDocumentCommandHandlerForTest(merchantDocumentCommandServiceStub{
		update: func(_ context.Context, req *requests.UpdateMerchantDocumentRequest) (*db.UpdateMerchantDocumentRow, error) {
			gotRequest = req
			note := req.Note
			return &db.UpdateMerchantDocumentRow{
				DocumentID:   int32(*req.DocumentID),
				MerchantID:   int32(req.MerchantID),
				DocumentType: req.DocumentType,
				DocumentUrl:  req.DocumentUrl,
				Status:       req.Status,
				Note:         &note,
			}, nil
		},
	})

	got, err := handler.Update(context.Background(), &pb.UpdateMerchantDocumentRequest{
		DocumentId:   documentID,
		MerchantId:   merchantID,
		DocumentType: "business_license",
		DocumentUrl:  "https://example.com/license.pdf",
		Status:       "pending",
		Note:         "replace before review",
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if gotRequest == nil {
		t.Fatal("Update() did not call the command service")
	}
	if gotRequest.DocumentID == nil || *gotRequest.DocumentID != documentID {
		t.Fatalf("service received DocumentID %v, want %d", gotRequest.DocumentID, documentID)
	}
	if gotRequest.MerchantID != merchantID {
		t.Fatalf("service received MerchantID %d, want %d", gotRequest.MerchantID, merchantID)
	}
	if got.GetStatus() != "success" {
		t.Fatalf("response status = %q, want success", got.GetStatus())
	}
}

func TestMerchantDocumentCommandHandlerUpdateStatusUsesDocumentID(t *testing.T) {
	const (
		documentID = 84
		merchantID = 13
	)

	var gotRequest *requests.UpdateMerchantDocumentStatusRequest
	handler := newMerchantDocumentCommandHandlerForTest(merchantDocumentCommandServiceStub{
		updateStatus: func(_ context.Context, req *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error) {
			gotRequest = req
			note := req.Note
			return &db.UpdateMerchantDocumentStatusRow{
				DocumentID: int32(*req.DocumentID),
				MerchantID: int32(req.MerchantID),
				Status:     req.Status,
				Note:       &note,
			}, nil
		},
	})

	got, err := handler.UpdateStatus(context.Background(), &pb.UpdateMerchantDocumentStatusRequest{
		DocumentId: documentID,
		MerchantId: merchantID,
		Status:     "approved",
		Note:       "verified",
	})
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
	if gotRequest == nil {
		t.Fatal("UpdateStatus() did not call the command service")
	}
	if gotRequest.DocumentID == nil || *gotRequest.DocumentID != documentID {
		t.Fatalf("service received DocumentID %v, want %d", gotRequest.DocumentID, documentID)
	}
	if gotRequest.MerchantID != merchantID {
		t.Fatalf("service received MerchantID %d, want %d", gotRequest.MerchantID, merchantID)
	}
	if got.GetStatus() != "success" {
		t.Fatalf("response status = %q, want success", got.GetStatus())
	}
}
