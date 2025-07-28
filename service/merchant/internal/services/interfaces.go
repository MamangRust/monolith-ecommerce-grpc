package services

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"golang.org/x/net/context"
)

type MerchantQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, merchantID int) (*response.MerchantResponse, *response.ErrorResponse)
}

type MerchantCommandService interface {
	CreateMerchant(ctx context.Context, req *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	UpdateMerchantStatus(ctx context.Context, req *requests.UpdateMerchantStatusRequest) (*response.MerchantResponse, *response.ErrorResponse)
	TrashedMerchant(ctx context.Context, merchantID int) (*response.MerchantResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(ctx context.Context, merchantID int) (*response.MerchantResponse, *response.ErrorResponse)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}

type MerchantDocumentQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
}

type MerchantDocumentCommandService interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	TrashedMerchantDocument(ctx context.Context, documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	RestoreMerchantDocument(ctx context.Context, documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, *response.ErrorResponse)
	RestoreAllMerchantDocument(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
