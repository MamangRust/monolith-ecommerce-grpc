package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
)

func getStringValue(v interface{}) *wrapperspb.StringValue {
	if ts, ok := v.(pgtype.Timestamptz); ok && ts.Valid {
		return wrapperspb.String(ts.Time.Format("2006-01-02 15:04:05"))
	}
	if ts, ok := v.(pgtype.Timestamp); ok && ts.Valid {
		return wrapperspb.String(ts.Time.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func formatTimestamp(v interface{}) string {
	if ts, ok := v.(pgtype.Timestamptz); ok && ts.Valid {
		return ts.Time.Format("2006-01-02 15:04:05")
	}
	if ts, ok := v.(pgtype.Timestamp); ok && ts.Valid {
		return ts.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func createPaginationMeta(page, pageSize, totalRecords int) *pbcommon.PaginationMeta {
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoTransactionResponse(m interface{}) *pbtransaction.TransactionResponse {
	switch v := m.(type) {
	case *db.Transaction:
		return &pbtransaction.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionsRow:
		return &pbtransaction.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionsActiveRow:
		return &pbtransaction.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionByIDRow:
		return &pbtransaction.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionByOrderIDRow:
		return &pbtransaction.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionByMerchantRow:
		return &pbtransaction.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateTransactionRow:
		return &pbtransaction.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateTransactionRow:
		return &pbtransaction.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoTransactionResponseDeleteAt(m interface{}) *pbtransaction.TransactionResponseDeleteAt {
	switch v := m.(type) {
	case *db.Transaction:
		return &pbtransaction.TransactionResponseDeleteAt{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
			DeletedAt:     getStringValue(v.DeletedAt),
		}
	case *db.GetTransactionsActiveRow:
		return &pbtransaction.TransactionResponseDeleteAt{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
			DeletedAt:     getStringValue(v.DeletedAt),
		}
	case *db.GetTransactionsTrashedRow:
		return &pbtransaction.TransactionResponseDeleteAt{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
			DeletedAt:     getStringValue(v.DeletedAt),
		}
	default:
		return nil
	}
}
