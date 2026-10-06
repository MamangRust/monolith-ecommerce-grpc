package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pbmerchant_document "github.com/MamangRust/monolith-ecommerce-pb/merchant_document"
)

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

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

func mapToProtoMerchantResponse(m interface{}) *pbmerchant.MerchantResponse {
	switch v := m.(type) {
	case *db.Merchant:
		return &pbmerchant.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantsRow:
		return &pbmerchant.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantByIDRow:
		return &pbmerchant.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateMerchantRow:
		return &pbmerchant.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateMerchantRow:
		return &pbmerchant.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateMerchantStatusRow:
		return &pbmerchant.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantResponseDeleteAt(m interface{}) *pbmerchant.MerchantResponseDeleteAt {
	switch v := m.(type) {
	case *db.Merchant:
		return &pbmerchant.MerchantResponseDeleteAt{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
			DeletedAt:    getStringValue(v.DeletedAt),
		}
	case *db.GetMerchantsActiveRow:
		return &pbmerchant.MerchantResponseDeleteAt{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
			DeletedAt:    getStringValue(v.DeletedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantResponseTrashed(m interface{}) *pbmerchant.MerchantResponseDeleteAt {
	switch v := m.(type) {
	case *db.GetMerchantsTrashedRow:
		return &pbmerchant.MerchantResponseDeleteAt{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
			DeletedAt:    getStringValue(v.DeletedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantDocumentResponse(m interface{}) *pbmerchant_document.MerchantDocument {
	switch v := m.(type) {
	case *db.MerchantDocument:
		return &pbmerchant_document.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantDocumentsRow:
		return &pbmerchant_document.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetActiveMerchantDocumentsRow:
		return &pbmerchant_document.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantDocumentRow:
		return &pbmerchant_document.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateMerchantDocumentRow:
		return &pbmerchant_document.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantDocumentResponseAt(m interface{}) *pbmerchant_document.MerchantDocumentDeleteAt {
	switch v := m.(type) {
	case *db.GetTrashedMerchantDocumentsRow:
		return &pbmerchant_document.MerchantDocumentDeleteAt{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
			DeletedAt:    getStringValue(v.DeletedAt),
		}
	default:
		return nil
	}
}
