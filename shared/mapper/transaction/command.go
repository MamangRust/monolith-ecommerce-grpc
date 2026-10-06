package transactionapimapper

import (
	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type transactionCommandResponseMapper struct{}

func NewTransactionCommandResponseMapper() TransactionCommandResponseMapper {
	return &transactionCommandResponseMapper{}
}

func (t *transactionCommandResponseMapper) ToResponseTransaction(transaction *pbtransaction.TransactionResponse) *response.TransactionResponse {
	if transaction == nil {
		return nil
	}
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

func (t *transactionCommandResponseMapper) ToResponsesTransaction(transactions []*pbtransaction.TransactionResponse) []*response.TransactionResponse {
	var mappedTransactions []*response.TransactionResponse
	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.ToResponseTransaction(transaction))
	}
	return mappedTransactions
}

func (t *transactionCommandResponseMapper) ToApiResponseTransaction(pbResponse *pbtransaction.ApiResponseTransaction) *response.ApiResponseTransaction {
	return &response.ApiResponseTransaction{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponseTransaction(pbResponse.Data),
	}
}

func (t *transactionCommandResponseMapper) ToResponseTransactionDeleteAt(transaction *pbtransaction.TransactionResponseDeleteAt) *response.TransactionResponseDeleteAt {
	if transaction == nil {
		return nil
	}
	var deletedAt string
	if transaction.DeletedAt != nil {
		deletedAt = transaction.DeletedAt.Value
	}

	return &response.TransactionResponseDeleteAt{
		ID:            int(transaction.Id),
		OrderID:       int(transaction.OrderId),
		MerchantID:    int(transaction.MerchantId),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int(transaction.Amount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		DeletedAt:     &deletedAt,
	}
}

func (t *transactionCommandResponseMapper) ToResponsesTransactionDeleteAt(transactions []*pbtransaction.TransactionResponseDeleteAt) []*response.TransactionResponseDeleteAt {
	var mappedTransactions []*response.TransactionResponseDeleteAt
	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.ToResponseTransactionDeleteAt(transaction))
	}
	return mappedTransactions
}

func (t *transactionCommandResponseMapper) ToApiResponseTransactionDeleteAt(pbResponse *pbtransaction.ApiResponseTransactionDeleteAt) *response.ApiResponseTransactionDeleteAt {
	return &response.ApiResponseTransactionDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    t.ToResponseTransactionDeleteAt(pbResponse.Data),
	}
}

func (t *transactionCommandResponseMapper) ToApiResponseTransactionDelete(pbResponse *pbtransaction.ApiResponseTransactionDelete) *response.ApiResponseTransactionDelete {
	return &response.ApiResponseTransactionDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (t *transactionCommandResponseMapper) ToApiResponseTransactionAll(pbResponse *pbtransaction.ApiResponseTransactionAll) *response.ApiResponseTransactionAll {
	return &response.ApiResponseTransactionAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (t *transactionCommandResponseMapper) ToApiResponsePaginationTransactionDeleteAt(pbResponse *pbtransaction.ApiResponsePaginationTransactionDeleteAt) *response.ApiResponsePaginationTransactionDeleteAt {
	return &response.ApiResponsePaginationTransactionDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       t.ToResponsesTransactionDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
