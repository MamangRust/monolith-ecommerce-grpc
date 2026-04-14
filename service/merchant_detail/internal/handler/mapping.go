package handler

import (
	"math"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func createPaginationMeta(page, pageSize, totalRecords int) *pb.PaginationMeta {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pb.PaginationMeta{
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

func mapToProtoMerchantDetailResponse(m interface{}) *pb.MerchantDetailResponse {
	switch v := m.(type) {
	case *db.MerchantDetail:
		return &pb.MerchantDetailResponse{
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
		return &pb.MerchantDetailResponse{
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
		return &pb.MerchantDetailResponse{
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
		return &pb.MerchantDetailResponse{
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
		return &pb.MerchantDetailResponse{
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

func mapToProtoMerchantDetailResponseDeleteAt(m interface{}) *pb.MerchantDetailResponseDeleteAt {
	var res *pb.MerchantDetailResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *db.MerchantDetail:
		res = &pb.MerchantDetailResponseDeleteAt{
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
		res = &pb.MerchantDetailResponseDeleteAt{
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
		res = &pb.MerchantDetailResponseDeleteAt{
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

func mapToProtoMerchantSocialLinkResponse(m interface{}) *pb.MerchantSocialMediaLinkResponse {
	switch v := m.(type) {
	case *db.MerchantSocialMediaLink:
		return &pb.MerchantSocialMediaLinkResponse{
			Id:               v.MerchantSocialID,
			MerchantDetailId: v.MerchantDetailID,
			Platform:         v.Platform,
			Url:              v.Url,
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateMerchantSocialMediaLinkRow:
		return &pb.MerchantSocialMediaLinkResponse{
			Id:               v.MerchantSocialID,
			MerchantDetailId: v.MerchantDetailID,
			Platform:         v.Platform,
			Url:              v.Url,
			CreatedAt:        formatTimestamp(v.CreatedAt),
			UpdatedAt:        formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateMerchantSocialMediaLinkRow:
		return &pb.MerchantSocialMediaLinkResponse{
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
