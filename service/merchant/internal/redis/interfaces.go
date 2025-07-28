package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantQueryCache interface {
	GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponse, *int, bool)
	SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantResponse, total *int)

	GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, bool)
	SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantResponseDeleteAt, total *int)

	GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, bool)
	SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantResponseDeleteAt, total *int)

	GetCachedMerchant(ctx context.Context, id int) (*response.MerchantResponse, bool)
	SetCachedMerchant(ctx context.Context, data *response.MerchantResponse)

	GetCachedMerchantsByUserId(ctx context.Context, id int) ([]*response.MerchantResponse, bool)
	SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*response.MerchantResponse)
}

type MerchantDocumentQueryCache interface {
	GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, bool)
	SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponse, total *int)

	GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool)
	SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int)

	GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool)
	SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int)

	GetCachedMerchantDocument(ctx context.Context, id int) (*response.MerchantDocumentResponse, bool)
	SetCachedMerchantDocument(ctx context.Context, data *response.MerchantDocumentResponse)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(ctx context.Context, id int)
}

type MerchantDocumentCommandCache interface {
	DeleteCachedMerchantDocuments(ctx context.Context, id int)
}
