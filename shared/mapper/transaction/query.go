package transactionapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type transactionQueryResponseMapper struct{}

func NewTransactionQueryResponseMapper() TransactionQueryResponseMapper {
	return &transactionQueryResponseMapper{}
}

func (t *transactionQueryResponseMapper) ToResponseTransaction(transaction *pb.TransactionResponse) *response.TransactionResponse {
	if transaction == nil { return nil }
	return &response.TransactionResponse{
		ID:            int(transaction.Id),
		OrderID:       int(transaction.OrderId),
		MerchantID:    int(transaction.MerchantId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func (t *transactionQueryResponseMapper) ToResponsesTransaction(transactions []*pb.TransactionResponse) []*response.TransactionResponse {
	var mappedTransactions []*response.TransactionResponse
	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.ToResponseTransaction(transaction))
	}
	return mappedTransactions
}

func (t *transactionQueryResponseMapper) ToApiResponseTransaction(pbResponse *pb.ApiResponseTransaction) *response.ApiResponseTransaction {
	return &response.ApiResponseTransaction{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponseTransaction(pbResponse.Data),
	}
}

func (t *transactionQueryResponseMapper) ToApiResponsesTransaction(pbResponse *pb.ApiResponsesTransaction) *response.ApiResponsesTransaction {
	return &response.ApiResponsesTransaction{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponsesTransaction(pbResponse.Data),
	}
}

func (t *transactionQueryResponseMapper) ToApiResponsePaginationTransaction(pbResponse *pb.ApiResponsePaginationTransaction) *response.ApiResponsePaginationTransaction {
	return &response.ApiResponsePaginationTransaction{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       t.ToResponsesTransaction(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (t *transactionQueryResponseMapper) ToApiResponsePaginationTransactionDeleteAt(pbResponse *pb.ApiResponsePaginationTransactionDeleteAt) *response.ApiResponsePaginationTransactionDeleteAt {
	var mappedData []*response.TransactionResponseDeleteAt
	for _, tr := range pbResponse.Data {
		var deletedAt string
		if tr.DeletedAt != nil {
			deletedAt = tr.DeletedAt.Value
		}
		mappedData = append(mappedData, &response.TransactionResponseDeleteAt{
			ID:            int(tr.Id),
			OrderID:       int(tr.OrderId),
			MerchantID:    int(tr.MerchantId),
			PaymentMethod: tr.PaymentMethod,
			Amount:        int(tr.Amount),
			PaymentStatus: tr.PaymentStatus,
			CreatedAt:     tr.CreatedAt,
			UpdatedAt:     tr.UpdatedAt,
			DeletedAt:     &deletedAt,
		})
	}

	return &response.ApiResponsePaginationTransactionDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       mappedData,
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
