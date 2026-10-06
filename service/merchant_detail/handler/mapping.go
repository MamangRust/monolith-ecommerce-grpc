package handler

import (
	"math"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbmerchant_detail "github.com/MamangRust/monolith-ecommerce-pb/merchant_detail"
)

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
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func formatTimestamp(v interface{}) string {
	if ts, ok := v.(pgtype.Timestamp); ok && ts.Valid {
		return ts.Time.Format("2006-01-02 15:04:05.000")
	}
	if ts, ok := v.(pgtype.Timestamptz); ok && ts.Valid {
		return ts.Time.Format("2006-01-02 15:04:05.000")
	}
	return ""
}

func mapToProtoMerchantDetailResponse(m interface{}) *pbmerchant_detail.MerchantDetailResponse {
	switch v := m.(type) {
	case *db.MerchantDetail:
		return &pbmerchant_detail.MerchantDetailResponse{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantDetailRow:
		return &pbmerchant_detail.MerchantDetailResponse{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateMerchantDetailRow:
		return &pbmerchant_detail.MerchantDetailResponse{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateMerchantDetailRow:
		return &pbmerchant_detail.MerchantDetailResponse{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantDetailsRow:
		return &pbmerchant_detail.MerchantDetailResponse{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantDetailResponseDeleteAt(m interface{}) *pbmerchant_detail.MerchantDetailResponseDeleteAt {
	var res *pbmerchant_detail.MerchantDetailResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *db.MerchantDetail:
		res = &pbmerchant_detail.MerchantDetailResponseDeleteAt{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetMerchantDetailsActiveRow:
		res = &pbmerchant_detail.MerchantDetailResponseDeleteAt{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetMerchantDetailsTrashedRow:
		res = &pbmerchant_detail.MerchantDetailResponseDeleteAt{
			Id:               v.MerchantDetailID,
			MerchantId:       v.MerchantID,
			DisplayName:      getString(v.DisplayName),
			CoverImageUrl:    getString(v.CoverImageUrl),
			LogoUrl:          getString(v.LogoUrl),
			ShortDescription: getString(v.ShortDescription),
			WebsiteUrl:       getString(v.WebsiteUrl),
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	default:
		return nil
	}

	if val := formatTimestamp(deletedAt); val != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: val}
	}

	return res
}

func mapToProtoMerchantSocialLinkResponse(m interface{}) *pbmerchant_detail.MerchantSocialMediaLinkResponse {
	switch v := m.(type) {
	case *db.MerchantSocialMediaLink:
		return &pbmerchant_detail.MerchantSocialMediaLinkResponse{
			Id:               v.MerchantSocialID,
			MerchantDetailId: v.MerchantDetailID,
			Platform:         v.Platform,
			Url:              v.Url,
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateMerchantSocialMediaLinkRow:
		return &pbmerchant_detail.MerchantSocialMediaLinkResponse{
			Id:               v.MerchantSocialID,
			MerchantDetailId: v.MerchantDetailID,
			Platform:         v.Platform,
			Url:              v.Url,
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateMerchantSocialMediaLinkRow:
		return &pbmerchant_detail.MerchantSocialMediaLinkResponse{
			Id:               v.MerchantSocialID,
			MerchantDetailId: v.MerchantDetailID,
			Platform:         v.Platform,
			Url:              v.Url,
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}
